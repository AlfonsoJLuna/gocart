package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

// Base struct to embed in all models
// If a model embeds the Base struct, it satisfies the Model interface
type Base struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *Base) GetBase() *Base {
	return b
}

type Model interface {
	GetBase() *Base
}

// Error types
var (
	ErrNotFound      = errors.New("not found")
	ErrBadInput      = errors.New("bad input")
	ErrAlreadyExists = errors.New("already exists")
)

// Functions to handle CRUD operations of main records

// Adds a new record to a bucket
func create[T Model](tx *bbolt.Tx, bucketName string, model T) error {
	base := model.GetBase()

	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generating uuid: %w", err)
	}
	base.ID = id

	now := time.Now()
	base.CreatedAt = now
	base.UpdatedAt = now

	value, err := json.Marshal(model)
	if err != nil {
		return err
	}

	bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		return err
	}

	return bucket.Put(base.ID[:], value)
}

// Retrieves one record from a bucket by ID
func readTx[T any](tx *bbolt.Tx, bucketName string, id uuid.UUID) (*T, error) {
	var model T

	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, fmt.Errorf("bucket %q: %w", bucketName, ErrNotFound)
	}

	value := bucket.Get(id[:])
	if value == nil {
		return nil, fmt.Errorf("id %q in bucket %q: %w", id, bucketName, ErrNotFound)
	}

	if err := json.Unmarshal(value, &model); err != nil {
		return nil, err
	}

	return &model, nil
}

// Wrapper for readTx when no other operations are needed in the same Tx
func read[T any](db *bbolt.DB, bucketName string, id uuid.UUID) (*T, error) {
	var result *T
	err := db.View(func(tx *bbolt.Tx) error {
		var err error
		result, err = readTx[T](tx, bucketName, id)
		return err
	})
	return result, err
}

// Updates an existing record in a bucket
func update[T Model](tx *bbolt.Tx, bucketName string, model T) error {
	base := model.GetBase()

	if base.ID == uuid.Nil {
		return fmt.Errorf("model ID is nil: %w", ErrBadInput)
	}

	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return fmt.Errorf("bucket %q: %w", bucketName, ErrNotFound)
	}

	if bucket.Get(base.ID[:]) == nil {
		return fmt.Errorf("id %q in bucket %q: %w", base.ID, bucketName, ErrNotFound)
	}

	base.UpdatedAt = time.Now()

	value, err := json.Marshal(model)
	if err != nil {
		return err
	}

	return bucket.Put(base.ID[:], value)
}

// Removes an existing record from a bucket
func delete(tx *bbolt.Tx, bucketName string, id uuid.UUID) error {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return fmt.Errorf("bucket %q: %w", bucketName, ErrNotFound)
	}

	if bucket.Get(id[:]) == nil {
		return fmt.Errorf("id %q in bucket %q: %w", id, bucketName, ErrNotFound)
	}

	return bucket.Delete(id[:])
}

// Retrieves all records in a bucket. Supports pagination, use limit=0 to return all values
func list[T any](db *bbolt.DB, bucketName string, offset, limit int, desc bool) ([]*T, error) {
	if offset < 0 || limit < 0 {
		return nil, fmt.Errorf("offset and limit must be non-negative: %w", ErrBadInput)
	}

	var results []*T

	if err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %q: %w", bucketName, ErrNotFound)
		}

		c := bucket.Cursor()
		var k, v []byte

		// Position the cursor at first or last entry depending on desc
		if desc {
			k, v = c.Last()
		} else {
			k, v = c.First()
		}

		// Skip offset entries
		for i := 0; i < offset && k != nil; i++ {
			if desc {
				k, v = c.Prev()
			} else {
				k, v = c.Next()
			}
		}

		// Collect up to limit entries
		for k != nil && (limit == 0 || len(results) < limit) {
			var model T

			if err := json.Unmarshal(v, &model); err != nil {
				return err
			}

			results = append(results, &model)

			if desc {
				k, v = c.Prev()
			} else {
				k, v = c.Next()
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

// Functions to manage 1->1 (isUnique==true) or 1->N (isUnique==false) indexes

// Adds an index entry. Must be called within the same Tx as the delete function
// Within the Tx, must be called after create so create has already generated the ID
func createIndex(tx *bbolt.Tx, bucketName, indexName, key string, id uuid.UUID, isUnique bool) error {
	indexBucketName := bucketName + ":idx:" + indexName

	bucket, err := tx.CreateBucketIfNotExists([]byte(indexBucketName))
	if err != nil {
		return err
	}

	indexKey := key + ":" + id.String()

	if isUnique {
		// Prevent duplicate key:* entry (same key is not allowed even if id is different)
		prefix := []byte(key + ":")

		c := bucket.Cursor()
		k, _ := c.Seek(prefix)

		if k != nil && bytes.HasPrefix(k, prefix) {
			return fmt.Errorf("unique key %q in bucket %q: %w", key, indexBucketName, ErrAlreadyExists)
		}
	} else {
		// Prevent duplicate key:id entry (same key is allowed if id is different)
		if bucket.Get([]byte(indexKey)) != nil {
			return fmt.Errorf("entry %q in bucket %q: %w", indexKey, indexBucketName, ErrAlreadyExists)
		}
	}

	return bucket.Put([]byte(indexKey), id[:])
}

// Removes an index entry. Must be called within the same Tx as the delete function
func deleteIndex(tx *bbolt.Tx, bucketName, indexName, key string, id uuid.UUID) error {
	indexBucketName := bucketName + ":idx:" + indexName

	bucket := tx.Bucket([]byte(indexBucketName))
	if bucket == nil {
		return fmt.Errorf("bucket %q: %w", indexBucketName, ErrNotFound)
	}

	indexKey := key + ":" + id.String()

	if bucket.Get([]byte(indexKey)) == nil {
		return fmt.Errorf("entry %q in bucket %q: %w", indexKey, indexBucketName, ErrNotFound)
	}

	return bucket.Delete([]byte(indexKey))
}

// Updates an index entry. Must be called within the same Tx as the update function
func updateIndex(tx *bbolt.Tx, bucketName, indexName string, oldKey, newKey string, id uuid.UUID, isUnique bool) error {
	if oldKey == newKey {
		return nil // Nothing to update
	}

	if err := deleteIndex(tx, bucketName, indexName, oldKey, id); err != nil {
		return err // Couldn't delete, aborts the whole update transaction
	}

	if err := createIndex(tx, bucketName, indexName, newKey, id, isUnique); err != nil {
		return err // Couldn't create, aborts the whole update transaction
	}

	return nil // Updated ok
}

func createNullableIndex(tx *bbolt.Tx, bucketName, indexName string, id *uuid.UUID, recordID uuid.UUID) error {
    if id == nil {
        return nil
    }
    return createIndex(tx, bucketName, indexName, id.String(), recordID, false)
}

func deleteNullableIndex(tx *bbolt.Tx, bucketName, indexName string, id *uuid.UUID, recordID uuid.UUID) error {
    if id == nil {
        return nil
    }
    return deleteIndex(tx, bucketName, indexName, id.String(), recordID)
}

func updateNullableIndex(tx *bbolt.Tx, bucketName, indexName string, oldID, newID *uuid.UUID, recordID uuid.UUID) error {
    if err := deleteNullableIndex(tx, bucketName, indexName, oldID, recordID); err != nil {
        return err
    }
    return createNullableIndex(tx, bucketName, indexName, newID, recordID)
}

// Functions to retrieve records using indexes

// Retrieves all records of a bucket by index key. Supports pagination, use limit=0 to return all values
func listByIndex[T any](db *bbolt.DB, bucketName, indexName, key string, offset, limit int, desc bool) ([]*T, error) {
	if offset < 0 || limit < 0 {
		return nil, fmt.Errorf("offset and limit must be non-negative: %w", ErrBadInput)
	}

	var results []*T

	indexBucketName := bucketName + ":idx:" + indexName
	prefix := []byte(key + ":")

	err := db.View(func(tx *bbolt.Tx) error {
		indexBucket := tx.Bucket([]byte(indexBucketName))
		if indexBucket == nil {
			return fmt.Errorf("bucket %q: %w", indexBucketName, ErrNotFound)
		}

		dataBucket := tx.Bucket([]byte(bucketName))
		if dataBucket == nil {
			return fmt.Errorf("bucket %q: %w", bucketName, ErrNotFound)
		}

		c := indexBucket.Cursor()
		var k, v []byte

		// Position the cursor at first entry matching the key:* prefix
		k, v = c.Seek(prefix)

		if desc {
			// If desc, we need to position the cursor at last entry matching the key:* prefix

			// Advance until we fall off the end of matching keys
			for k != nil && bytes.HasPrefix(k, prefix) {
				k, v = c.Next()
			}

			if k != nil {
				k, v = c.Prev() // Landed past prefix, step back
			} else {
				k, v = c.Last() // Exhausted all keys, jump to last
			}
		}

		// Skip offset entries
		for i := 0; i < offset && k != nil && bytes.HasPrefix(k, prefix); i++ {
			if desc {
				k, v = c.Prev()
			} else {
				k, v = c.Next()
			}
		}

		// Collect up to limit entries
		for k != nil && bytes.HasPrefix(k, prefix) && (limit == 0 || len(results) < limit) {
			value := dataBucket.Get(v)
			if value != nil {
				var model T
				if err := json.Unmarshal(value, &model); err != nil {
					return err
				}
				results = append(results, &model)
			}

			if desc {
				k, v = c.Prev()
			} else {
				k, v = c.Next()
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

// Retrieves one record from a bucket by index key
func readByIndex[T any](db *bbolt.DB, bucketName, indexName, key string) (*T, error) {

	results, err := listByIndex[T](
		db,
		bucketName,
		indexName,
		key,
		0, // offset
		1, // limit
		false,
	)

	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("key %q in bucket %q: %w", key, bucketName, ErrNotFound)
	}

	return results[0], nil
}
