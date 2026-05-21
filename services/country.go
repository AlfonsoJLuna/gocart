package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	
	"gocart/models"
	"gocart/utils"
)

const countryColumns = `id, created_at, updated_at, is_enabled,
	iso_code, name, name_alt, currency_id, is_eu, vat_rate`

func scanCountry(row interface{ Scan(...any) error }) (*models.Country, error) {
	var c models.Country

	var createdAt, updatedAt string

	if err := row.Scan(
		&c.ID, &createdAt, &updatedAt, &c.IsEnabled,
		&c.ISOCode, &c.Name, &c.NameAlt,
		&c.CurrencyID, &c.IsEU, &c.VATRate,
	); err != nil {
		return nil, err
	}

	var err error
	
	c.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse created_at: %w", err)
	}

	c.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse updated_at: %w", err)
	}

	return &c, nil
}

func CountryCreate(db *sql.DB, c *models.Country) error {
	if err := c.Init(); err != nil {
		return fmt.Errorf("initializing country: %w", err)
	}
	
	_, err := db.Exec(`
		INSERT INTO countries (
			id,
			created_at,
			updated_at,
			is_enabled,
			iso_code,
			name,
			name_alt,
			currency_id,
			is_eu,
			vat_rate
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID.String(),
		utils.TimeToString(c.CreatedAt),
		utils.TimeToString(c.UpdatedAt),
		c.IsEnabled,
		c.ISOCode,
		c.Name,
		c.NameAlt,
		c.CurrencyID,
		c.IsEU,
		c.VATRate,
	)
	if err != nil {
		return fmt.Errorf("CountryCreate: %w", err)
	}

	return nil
}

func CountryListAll(db *sql.DB, offset, limit int, desc bool) ([]*models.Country, error) {
	order := "ASC"
	if desc {
		order = "DESC"
	}
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM countries
		ORDER BY name %s LIMIT %d OFFSET %d`,
		countryColumns, order, limit, offset,
	))
	if err != nil {
		return nil, fmt.Errorf("CountryListAll: %w", err)
	}
	defer rows.Close()

	var countries []*models.Country
	for rows.Next() {
		c, err := scanCountry(rows)
		if err != nil {
			return nil, fmt.Errorf("CountryListAll scan: %w", err)
		}
		countries = append(countries, c)
	}
	return countries, rows.Err()
}

func CountryReadByID(db *sql.DB, id uuid.UUID) (*models.Country, error) {
	row := db.QueryRow(`SELECT `+countryColumns+` FROM countries WHERE id = ?`, id)
	c, err := scanCountry(row)
	if err != nil {
		return nil, fmt.Errorf("CountryReadByID: %w", err)
	}
	return c, nil
}

func CountryReadByISOCode(db *sql.DB, isoCode string) (*models.Country, error) {
	row := db.QueryRow(`SELECT `+countryColumns+` FROM countries WHERE iso_code = ?`, isoCode)
	c, err := scanCountry(row)
	if err != nil {
		return nil, fmt.Errorf("CountryReadByISOCode: %w", err)
	}
	return c, nil
}

func CountryUpdate(db *sql.DB, c *models.Country) error {
	now := time.Now().UTC()
	now_str := now.Format(time.RFC3339Nano)

	_, err := db.Exec(`
		UPDATE countries
		SET updated_at = ?, is_enabled = ?,
		    iso_code = ?, name = ?, name_alt = ?,
		    currency_id = ?, is_eu = ?, vat_rate = ?
		WHERE id = ?`,
		now_str, c.IsEnabled,
		c.ISOCode, c.Name, c.NameAlt,
		c.CurrencyID, c.IsEU, c.VATRate,
		c.ID,
	)
	if err != nil {
		return fmt.Errorf("CountryUpdate: %w", err)
	}
	c.UpdatedAt = now
	return nil
}

// CountryDelete permanently deletes a country and cascades to its regions.
func CountryDelete(db *sql.DB, id uuid.UUID) error {
	if _, err := db.Exec(`DELETE FROM countries WHERE id = ?`, id); err != nil {
		return fmt.Errorf("CountryDelete: %w", err)
	}
	return nil
}