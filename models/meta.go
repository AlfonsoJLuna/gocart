package models

import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

const metaBucket = "meta"
const metaKey = "meta"

type Meta struct {
	SchemaVersion	int  	`json:"schema_version"`
	Initialized   	bool	`json:"initialized"`
}

func MetaRead(db *bbolt.DB) (*Meta, error) {
	var meta Meta
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(metaBucket))
		if bucket == nil {
			return fmt.Errorf("bucket %q: %w", metaBucket, ErrNotFound)
		}
		value := bucket.Get([]byte(metaKey))
		if value == nil {
			return fmt.Errorf("key %q in bucket %q: %w", metaKey, metaBucket, ErrNotFound)
		}
		return json.Unmarshal(value, &meta)
	})
	return &meta, err
}

func MetaWrite(db *bbolt.DB, meta *Meta) error {
	return db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(metaBucket))
		if err != nil {
			return fmt.Errorf("bucket %q: %w", metaBucket, err)
		}
		data, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(metaKey), data)
	})
}
