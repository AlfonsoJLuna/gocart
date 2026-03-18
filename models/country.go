package models

import (
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

const countryBucket = "countries"

type Region struct {
    Name      		string  	`json:"name"`				// Region name in English
    NameAlt   		string  	`json:"name_alt"`			// Region name in local language
    IsEU      		bool    	`json:"is_eu"`				// This region is part of the EU
    VATRate   		float64 	`json:"vat_rate"`			// Specific VAT rate for this region
    IsEnabled		bool    	`json:"is_enabled"`
}

type Country struct {
	Base
	ISOCode    		string  	`json:"iso_code"`			// 2-letter ISO 3166-1 alpha-2 code
	Name       		string  	`json:"name"`				// Country name in English
	NameAlt    		string  	`json:"name_alt"`			// Country name in local language
	CurrencyISOCode string		`json:"currency_iso_code"`	// Default currency for this country
	IsEU       		bool    	`json:"is_eu"`				// This country is part of the EU
	VATRate    		float64		`json:"vat_rate"`			// Default VAT rate for this country
	IsEnabled  		bool    	`json:"is_enabled"`
	Regions			[]*Region	`json:"regions,omitempty"`	// Countries can optionally have regions
}

func CountryCreate(db *bbolt.DB, c *Country) (uuid.UUID, error) {
	if err := db.Update(func(tx *bbolt.Tx) error {
		if err := create(tx, countryBucket, c); err != nil {
			return err
		}
		if err := createIndex(tx, countryBucket, "iso_code", c.ISOCode, c.ID, true); err != nil {
			return err
		}
		if err := createIndex(tx, countryBucket, "name", c.Name, c.ID, true); err != nil {
			return err
		}
		if err := createIndex(tx, countryBucket, "name_alt", c.NameAlt, c.ID, true); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return uuid.Nil, err
	}
	return c.ID, nil
}

func CountryListAll(db *bbolt.DB, offset, limit int, desc bool) ([]*Country, error) {
	return list[Country](db, countryBucket, offset, limit, desc)
}

func CountryReadByID(db *bbolt.DB, id uuid.UUID) (*Country, error) {
	return read[Country](db, countryBucket, id)
}

func CountryReadByISOCode(db *bbolt.DB, isoCode string) (*Country, error) {
	return readByIndex[Country](db, countryBucket, "iso_code", isoCode)
}

func CountryReadByName(db *bbolt.DB, name string) (*Country, error) {
	return readByIndex[Country](db, countryBucket, "name", name)
}

func CountryReadByNameAlt(db *bbolt.DB, nameAlt string) (*Country, error) {
	return readByIndex[Country](db, countryBucket, "name_alt", nameAlt)
}

func CountryUpdate(db *bbolt.DB, c *Country) error {
	return db.Update(func(tx *bbolt.Tx) error {
		old, err := readTx[Country](tx, countryBucket, c.ID)
		if err != nil {
			return err
		}
		if err := updateIndex(tx, countryBucket, "iso_code", old.ISOCode, c.ISOCode, c.ID, true); err != nil {
			return err
		}
		if err := updateIndex(tx, countryBucket, "name", old.Name, c.Name, c.ID, true); err != nil {
			return err
		}
		if err := updateIndex(tx, countryBucket, "name_alt", old.NameAlt, c.NameAlt, c.ID, true); err != nil {
			return err
		}
		return update(tx, countryBucket, c)
	})
}

func CountryDelete(db *bbolt.DB, id uuid.UUID) error {
	return db.Update(func(tx *bbolt.Tx) error {
		c, err := readTx[Country](tx, countryBucket, id)
		if err != nil {
			return err
		}
		if err := deleteIndex(tx, countryBucket, "iso_code", c.ISOCode, id); err != nil {
			return err
		}
		if err := deleteIndex(tx, countryBucket, "name", c.Name, id); err != nil {
			return err
		}
		if err := deleteIndex(tx, countryBucket, "name_alt", c.NameAlt, id); err != nil {
			return err
		}
		return delete(tx, countryBucket, id)
	})
}
