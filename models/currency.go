package models

import (
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

const currencyBucket = "currencies"

type Currency struct {
	Base
	ISOCode   string  `json:"iso_code"`		// 3-letter ISO 4217 code
	Name      string  `json:"name"`			// Currency name in English
	NameAlt   string  `json:"name_alt"`		// Currency name in local language
	Decimals  int     `json:"decimals"`		// Number of decimals to show in the UI
	FXRate    float64 `json:"fx_rate"`		// Latest exchange rate against the EUR
	IsEnabled bool    `json:"is_enabled"`
}

func CurrencyCreate(db *bbolt.DB, c *Currency) (uuid.UUID, error) {
	if err := db.Update(func(tx *bbolt.Tx) error {
		if err := create(tx, currencyBucket, c); err != nil {
			return err
		}
		if err := createIndex(tx, currencyBucket, "iso_code", c.ISOCode, c.ID, true); err != nil {
			return err
		}
		if err := createIndex(tx, currencyBucket, "name", c.Name, c.ID, true); err != nil {
			return err
		}
		if err := createIndex(tx, currencyBucket, "name_alt", c.NameAlt, c.ID, true); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return uuid.Nil, err
	}
	return c.ID, nil
}

func CurrencyListAll(db *bbolt.DB, offset, limit int, desc bool) ([]*Currency, error) {
	return list[Currency](db, currencyBucket, offset, limit, desc)
}

func CurrencyReadByID(db *bbolt.DB, id uuid.UUID) (*Currency, error) {
	return read[Currency](db, currencyBucket, id)
}

func CurrencyReadByISOCode(db *bbolt.DB, isoCode string) (*Currency, error) {
	return readByIndex[Currency](db, currencyBucket, "iso_code", isoCode)
}

func CurrencyReadByName(db *bbolt.DB, name string) (*Currency, error) {
	return readByIndex[Currency](db, currencyBucket, "name", name)
}

func CurrencyReadByNameAlt(db *bbolt.DB, nameAlt string) (*Currency, error) {
	return readByIndex[Currency](db, currencyBucket, "name_alt", nameAlt)
}

func CurrencyUpdate(db *bbolt.DB, c *Currency) error {
	return db.Update(func(tx *bbolt.Tx) error {
		old, err := readTx[Currency](tx, currencyBucket, c.ID)
		if err != nil {
			return err
		}
		if err := updateIndex(tx, currencyBucket, "iso_code", old.ISOCode, c.ISOCode, c.ID, true); err != nil {
			return err
		}
		if err := updateIndex(tx, currencyBucket, "name", old.Name, c.Name, c.ID, true); err != nil {
			return err
		}
		if err := updateIndex(tx, currencyBucket, "name_alt", old.NameAlt, c.NameAlt, c.ID, true); err != nil {
			return err
		}
		return update(tx, currencyBucket, c)
	})
}

func CurrencyDelete(db *bbolt.DB, id uuid.UUID) error {
	return db.Update(func(tx *bbolt.Tx) error {
		c, err := readTx[Currency](tx, currencyBucket, id)
		if err != nil {
			return err
		}
		if err := deleteIndex(tx, currencyBucket, "iso_code", c.ISOCode, id); err != nil {
			return err
		}
		if err := deleteIndex(tx, currencyBucket, "name", c.Name, id); err != nil {
			return err
		}
		if err := deleteIndex(tx, currencyBucket, "name_alt", c.NameAlt, id); err != nil {
			return err
		}
		return delete(tx, currencyBucket, id)
	})
}
