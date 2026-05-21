package seeds

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	
	"gocart/models"
	"gocart/services"
)

//go:embed currencies.json
var currenciesJSON []byte

//go:embed countries.json
var countriesJSON []byte

// countryJSON is a temporary struct used only for unmarshalling countries.json.
// It includes currency_iso_code (used to look up the currency UUID after
// seeding) and regions, which are stored in a separate table.
type countryJSON struct {
	models.Country
	CurrencyISOCode string       `json:"currency_iso_code"`
	Regions         []regionJSON `json:"regions"`
}
 
type regionJSON struct {
	Name      string  `json:"name"`
	NameAlt   string  `json:"name_alt"`
	IsEU      bool    `json:"is_eu"`
	VATRate   float64 `json:"vat_rate"`
}

func SeedCurrencies(db *sql.DB) error {
	var currencies []*models.Currency
	if err := json.Unmarshal(currenciesJSON, &currencies); err != nil {
		return fmt.Errorf("parsing currencies.json: %w", err)
	}

	existing, err := services.CurrencyListAll(db, 0, -1, false)
	if err != nil {
		return fmt.Errorf("listing existing currencies: %w", err)
	}
	for _, c := range existing {
		if err := services.CurrencyDelete(db, c.ID); err != nil {
			return fmt.Errorf("deleting currency %q: %w", c.ISOCode, err)
		}
	}

	for _, c := range currencies {
		c.IsEnabled = true
		if err := services.CurrencyCreate(db, c); err != nil {
			return fmt.Errorf("creating currency %q: %w", c.ISOCode, err)
		}
	}
	return nil
}

func SeedCountries(db *sql.DB) error {
	var countries []countryJSON
	if err := json.Unmarshal(countriesJSON, &countries); err != nil {
		return fmt.Errorf("parsing countries.json: %w", err)
	}
 
	existing, err := services.CountryListAll(db, 0, -1, false)
	if err != nil {
		return fmt.Errorf("listing existing countries: %w", err)
	}
	for _, c := range existing {
		// CountryDelete cascades to regions via ON DELETE CASCADE.
		if err := services.CountryDelete(db, c.ID); err != nil {
			return fmt.Errorf("deleting country %q: %w", c.ISOCode, err)
		}
	}
 
	for _, entry := range countries {
		entry.Country.IsEnabled = true
 
		// Look up the currency UUID by ISO code and assign it.
		if entry.CurrencyISOCode != "" {
			currency, err := services.CurrencyReadByISOCode(db, entry.CurrencyISOCode)
			if err != nil {
				return fmt.Errorf("looking up currency %q for country %q: %w", entry.CurrencyISOCode, entry.Country.ISOCode, err)
			}
			entry.Country.CurrencyID = &currency.ID
		}
 
		if err := services.CountryCreate(db, &entry.Country); err != nil {
			return fmt.Errorf("creating country %q: %w", entry.Country.ISOCode, err)
		}
 
		for _, r := range entry.Regions {
			region := &models.Region{
				CountryID: entry.Country.ID,
				Name:      r.Name,
				NameAlt:   r.NameAlt,
				IsEU:      r.IsEU,
				VATRate:   r.VATRate,
			}
			region.IsEnabled = true
			if err := services.RegionCreate(db, region); err != nil {
				return fmt.Errorf("creating region %q for country %q: %w", r.Name, entry.Country.ISOCode, err)
			}
		}
	}
	return nil
}

func SeedAll(db *sql.DB) error {
	log.Println("Seeding currencies...")
	if err := SeedCurrencies(db); err != nil {
		return fmt.Errorf("seeding currencies: %w", err)
	}
	
	log.Println("Seeding countries...")
	if err := SeedCountries(db); err != nil {
		return fmt.Errorf("seeding countries: %w", err)
	}

	return nil
}
