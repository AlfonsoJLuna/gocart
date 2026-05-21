package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gocart/models"
	
	"gocart/utils"
)

const regionColumns = `id, created_at, updated_at, is_enabled,
	country_id, name, name_alt, is_eu, vat_rate`

func scanRegion(row interface{ Scan(...any) error }) (*models.Region, error) {
	var r models.Region

	var createdAt, updatedAt string

	if err := row.Scan(
		&r.ID, &createdAt, &updatedAt, &r.IsEnabled,
		&r.CountryID, &r.Name, &r.NameAlt, &r.IsEU, &r.VATRate,
	); err != nil {
		return nil, err
	}

	var err error
	
	r.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse created_at: %w", err)
	}

	r.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse updated_at: %w", err)
	}

	return &r, nil
}

func RegionCreate(db *sql.DB, r *models.Region) error {
	if err := r.Init(); err != nil {
		return fmt.Errorf("initializing region: %w", err)
	}

	_, err := db.Exec(`
		INSERT INTO regions (
			id,
			created_at,
			updated_at,
			is_enabled,
			country_id,
			name,
			name_alt,
			is_eu,
			vat_rate
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.ID.String(),
		utils.TimeToString(r.CreatedAt),
		utils.TimeToString(r.UpdatedAt),
		r.IsEnabled,
		r.CountryID,
		r.Name,
		r.NameAlt,
		r.IsEU,
		r.VATRate,
	)
	if err != nil {
		return fmt.Errorf("RegionCreate: %w", err)
	}

	return nil
}

func RegionListByCountryID(db *sql.DB, countryID uuid.UUID) ([]*models.Region, error) {
	rows, err := db.Query(`
		SELECT `+regionColumns+` FROM regions
		WHERE country_id = ? ORDER BY name ASC`, countryID)
	if err != nil {
		return nil, fmt.Errorf("RegionListByCountryID: %w", err)
	}
	defer rows.Close()

	var regions []*models.Region
	for rows.Next() {
		r, err := scanRegion(rows)
		if err != nil {
			return nil, fmt.Errorf("RegionListByCountryID scan: %w", err)
		}
		regions = append(regions, r)
	}
	return regions, rows.Err()
}

func RegionReadByID(db *sql.DB, id uuid.UUID) (*models.Region, error) {
	row := db.QueryRow(`SELECT `+regionColumns+` FROM regions WHERE id = ?`, id)
	r, err := scanRegion(row)
	if err != nil {
		return nil, fmt.Errorf("RegionReadByID: %w", err)
	}
	return r, nil
}

func RegionUpdate(db *sql.DB, r *models.Region) error {
	now := time.Now().UTC()
	now_str := now.Format(time.RFC3339Nano)

	_, err := db.Exec(`
		UPDATE regions
		SET updated_at = ?, is_enabled = ?,
		    name = ?, name_alt = ?, is_eu = ?, vat_rate = ?
		WHERE id = ?`,
		now_str, r.IsEnabled,
		r.Name, r.NameAlt, r.IsEU, r.VATRate,
		r.ID,
	)
	if err != nil {
		return fmt.Errorf("RegionUpdate: %w", err)
	}
	r.UpdatedAt = now
	return nil
}

func RegionDelete(db *sql.DB, id uuid.UUID) error {
	if _, err := db.Exec(`DELETE FROM regions WHERE id = ?`, id); err != nil {
		return fmt.Errorf("RegionDelete: %w", err)
	}
	return nil
}
