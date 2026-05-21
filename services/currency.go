package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	
	"gocart/models"
	"gocart/utils"
)

const currencyColumns = `id, created_at, updated_at, is_enabled,
	iso_code, name, name_alt, decimals, fx_rate`

func scanCurrency(row interface{ Scan(...any) error }) (*models.Currency, error) {
	var c models.Currency

	var createdAt, updatedAt string

	if err := row.Scan(
		&c.ID, &createdAt, &updatedAt, &c.IsEnabled,
		&c.ISOCode, &c.Name, &c.NameAlt, &c.Decimals, &c.FXRate,
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

func CurrencyCreate(db *sql.DB, c *models.Currency) error {
	if err := c.Init(); err != nil {
		return fmt.Errorf("initializing currency: %w", err)
	}

	_, err := db.Exec(`
		INSERT INTO currencies (
			id,
			created_at,
			updated_at,
			is_enabled,
			iso_code,
			name,
			name_alt,
			decimals,
			fx_rate
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		c.ID.String(),
		utils.TimeToString(c.CreatedAt),
		utils.TimeToString(c.UpdatedAt),
		c.IsEnabled,
		c.ISOCode,
		c.Name,
		c.NameAlt,
		c.Decimals,
		c.FXRate,
	)
	if err != nil {
		return fmt.Errorf("CurrencyCreate: %w", err)
	}

	return nil
}

func CurrencyListAll(db *sql.DB, offset, limit int, desc bool) ([]*models.Currency, error) {
	order := "ASC"
	if desc {
		order = "DESC"
	}
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM currencies
		ORDER BY name %s LIMIT %d OFFSET %d`,
		currencyColumns, order, limit, offset,
	))
	if err != nil {
		return nil, fmt.Errorf("CurrencyListAll: %w", err)
	}
	defer rows.Close()

	var currencies []*models.Currency
	for rows.Next() {
		c, err := scanCurrency(rows)
		if err != nil {
			return nil, fmt.Errorf("CurrencyListAll scan: %w", err)
		}
		currencies = append(currencies, c)
	}
	return currencies, rows.Err()
}

func CurrencyReadByID(db *sql.DB, id uuid.UUID) (*models.Currency, error) {
	row := db.QueryRow(`SELECT `+currencyColumns+` FROM currencies WHERE id = ?`, id)
	c, err := scanCurrency(row)
	if err != nil {
		return nil, fmt.Errorf("CurrencyReadByID: %w", err)
	}
	return c, nil
}

func CurrencyReadByISOCode(db *sql.DB, isoCode string) (*models.Currency, error) {
	row := db.QueryRow(`SELECT `+currencyColumns+` FROM currencies WHERE iso_code = ?`, isoCode)
	c, err := scanCurrency(row)
	if err != nil {
		return nil, fmt.Errorf("CurrencyReadByISOCode: %w", err)
	}
	return c, nil
}

func CurrencyUpdate(db *sql.DB, c *models.Currency) error {
	now := time.Now().UTC()
	now_str := now.Format(time.RFC3339Nano)
	
	_, err := db.Exec(`
		UPDATE currencies
		SET updated_at = ?, is_enabled = ?,
		    iso_code = ?, name = ?, name_alt = ?,
		    decimals = ?, fx_rate = ?
		WHERE id = ?`,
		now_str, c.IsEnabled,
		c.ISOCode, c.Name, c.NameAlt,
		c.Decimals, c.FXRate,
		c.ID,
	)
	if err != nil {
		return fmt.Errorf("CurrencyUpdate: %w", err)
	}
	c.UpdatedAt = now
	return nil
}

func CurrencyDelete(db *sql.DB, id uuid.UUID) error {
	if _, err := db.Exec(`DELETE FROM currencies WHERE id = ?`, id); err != nil {
		return fmt.Errorf("CurrencyDelete: %w", err)
	}
	return nil
}
