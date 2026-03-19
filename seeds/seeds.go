package seeds

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"errors"
	"log"

	"go.etcd.io/bbolt"

	"gocart/models"
)

//go:embed currencies.json
var currenciesJSON []byte

//go:embed countries.json
var countriesJSON []byte

func SeedCurrencies(db *bbolt.DB) error {
    var currencies []*models.Currency
    if err := json.Unmarshal(currenciesJSON, &currencies); err != nil {
        return fmt.Errorf("parsing currencies.json: %w", err)
    }

    existing, err := models.CurrencyListAll(db, 0, 0, false)
    if err != nil && !errors.Is(err, models.ErrNotFound) {
        return fmt.Errorf("listing existing currencies: %w", err)
    }

    for _, c := range existing {
        if err := models.CurrencyDelete(db, c.ID); err != nil {
            return fmt.Errorf("deleting currency %q: %w", c.ISOCode, err)
        }
    }

    for _, c := range currencies {
        c.IsEnabled = true
        if _, err := models.CurrencyCreate(db, c); err != nil {
            return fmt.Errorf("creating currency %q: %w", c.ISOCode, err)
        }
    }

    return nil
}

func SeedCountries(db *bbolt.DB) error {
	var countries []*models.Country

	if err := json.Unmarshal(countriesJSON, &countries); err != nil {
		return fmt.Errorf("parsing countries.json: %w", err)
	}

	existing, err := models.CountryListAll(db, 0, 0, false)

	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return fmt.Errorf("listing existing countries: %w", err)
	}
	
	for _, c := range existing {
		if err := models.CountryDelete(db, c.ID); err != nil {
			return fmt.Errorf("deleting country %q: %w", c.ISOCode, err)
		}
	}

	for _, c := range countries {
		c.IsEnabled = true

		for _, r := range c.Regions {
            r.IsEnabled = true
        }

		if _, err := models.CountryCreate(db, c); err != nil {
			return fmt.Errorf("creating country %q: %w", c.ISOCode, err)
		}
	}

	return nil
}

func SeedAll(db *bbolt.DB) error {
	log.Println("Seeding currencies...")

	if err := SeedCurrencies(db); err != nil {
		return fmt.Errorf("seeding currencies: %w", err)
	}

	log.Println("Seeding countries...")

	if err := SeedCountries(db); err != nil {
		return fmt.Errorf("seeding countries: %w", err)
	}

	log.Println("Database seeded successfully.")

	return nil
}
