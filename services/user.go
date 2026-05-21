package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"gocart/models"
	"gocart/utils"
)

const userColumns = `id, created_at, updated_at, is_enabled,
	email, username, password_hash, phone,
	currency_id, country_id, region_id, vat_number,
	billing_full_name, billing_company_name,
	billing_address_1, billing_address_2,
	billing_postal_code, billing_city,
	shipping_full_name, shipping_company_name,
	shipping_address_1, shipping_address_2,
	shipping_postal_code, shipping_city,
	is_verified, is_business,
	is_subscribed_info, is_subscribed_promos, is_admin,
	last_login_at, last_user_change_at, last_pass_change_at, last_email_change_at`

func scanUser(row interface{ Scan(...any) error }) (*models.User, error) {
	var u models.User

	var createdAt, updatedAt string

	var lastLoginAt, lastUserChangeAt, lastPassChangeAt, lastEmailChangeAt sql.NullInt64
	if err := row.Scan(
		&u.ID, &createdAt, &updatedAt, &u.IsEnabled,
		&u.Email, &u.Username, &u.PasswordHash, &u.Phone,
		&u.CurrencyID, &u.CountryID, &u.RegionID, &u.VATNumber,
		&u.BillingAddress.FullName, &u.BillingAddress.CompanyName,
		&u.BillingAddress.Address1, &u.BillingAddress.Address2,
		&u.BillingAddress.PostalCode, &u.BillingAddress.City,
		&u.ShippingAddress.FullName, &u.ShippingAddress.CompanyName,
		&u.ShippingAddress.Address1, &u.ShippingAddress.Address2,
		&u.ShippingAddress.PostalCode, &u.ShippingAddress.City,
		&u.IsVerified, &u.IsBusiness,
		&u.IsSubscribedInfo, &u.IsSubscribedPromos, &u.IsAdmin,
		&lastLoginAt, &lastUserChangeAt, &lastPassChangeAt, &lastEmailChangeAt,
	); err != nil {
		return nil, err
	}

	var err error

	u.CreatedAt, err = time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse created_at: %w", err)
	}

	u.UpdatedAt, err = time.Parse(time.RFC3339Nano, updatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse updated_at: %w", err)
	}

	if lastLoginAt.Valid {
		t := time.Unix(lastLoginAt.Int64, 0).UTC()
		u.LastLoginAt = &t
	}
	if lastUserChangeAt.Valid {
		t := time.Unix(lastUserChangeAt.Int64, 0).UTC()
		u.LastUserChangeAt = &t
	}
	if lastPassChangeAt.Valid {
		t := time.Unix(lastPassChangeAt.Int64, 0).UTC()
		u.LastPassChangeAt = &t
	}
	if lastEmailChangeAt.Valid {
		t := time.Unix(lastEmailChangeAt.Int64, 0).UTC()
		u.LastEmailChangeAt = &t
	}
	return &u, nil
}

func UserCreate(db *sql.DB, u *models.User, password string) error {
	if err := u.Init(); err != nil {
		return fmt.Errorf("initializing user: %w", err)
	}

	if err := u.SetPassword(password); err != nil {
		return fmt.Errorf("setting password: %w", err)
	}

	_, err := db.Exec(`
		INSERT INTO users (
			id,
			created_at,
			updated_at,
			is_enabled,
			email,
			username,
			password_hash,
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		u.ID.String(),
		utils.TimeToString(u.CreatedAt),
		utils.TimeToString(u.UpdatedAt),
		u.IsEnabled,
		u.Email,
		u.Username,
		u.PasswordHash,
	)
	if err != nil {
		return fmt.Errorf("UserCreate: %w", err)
	}

	return nil
}

func UserListAll(db *sql.DB, offset, limit int, desc bool) ([]*models.User, error) {
	order := "ASC"
	if desc {
		order = "DESC"
	}
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM users
		ORDER BY username %s LIMIT %d OFFSET %d`,
		userColumns, order, limit, offset,
	))
	if err != nil {
		return nil, fmt.Errorf("UserListAll: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("UserListAll scan: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func UserListByCountryID(db *sql.DB, countryID uuid.UUID, offset, limit int, desc bool) ([]*models.User, error) {
	order := "ASC"
	if desc {
		order = "DESC"
	}
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM users WHERE country_id = ?
		ORDER BY username %s LIMIT %d OFFSET %d`,
		userColumns, order, limit, offset,
	), countryID)
	if err != nil {
		return nil, fmt.Errorf("UserListByCountryID: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("UserListByCountryID scan: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func UserListByCurrencyID(db *sql.DB, currencyID uuid.UUID, offset, limit int, desc bool) ([]*models.User, error) {
	order := "ASC"
	if desc {
		order = "DESC"
	}
	rows, err := db.Query(fmt.Sprintf(`
		SELECT %s FROM users WHERE currency_id = ?
		ORDER BY username %s LIMIT %d OFFSET %d`,
		userColumns, order, limit, offset,
	), currencyID)
	if err != nil {
		return nil, fmt.Errorf("UserListByCurrencyID: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("UserListByCurrencyID scan: %w", err)
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func UserReadByID(db *sql.DB, id uuid.UUID) (*models.User, error) {
	row := db.QueryRow(`SELECT `+userColumns+` FROM users WHERE id = ?`, id)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("UserReadByID: %w", err)
	}
	return u, nil
}

func UserReadByEmail(db *sql.DB, email string) (*models.User, error) {
	row := db.QueryRow(`SELECT `+userColumns+` FROM users WHERE email = ?`, email)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("UserReadByEmail: %w", err)
	}
	return u, nil
}

func UserReadByUsername(db *sql.DB, username string) (*models.User, error) {
	row := db.QueryRow(`SELECT `+userColumns+` FROM users WHERE username = ?`, username)
	u, err := scanUser(row)
	if err != nil {
		return nil, fmt.Errorf("UserReadByUsername: %w", err)
	}
	return u, nil
}

func UserUpdate(db *sql.DB, u *models.User) error {
	now := time.Now().UTC()
	now_str := now.Format(time.RFC3339Nano)

	_, err := db.Exec(`
		UPDATE users
		SET updated_at = ?, is_enabled = ?,
		    email = ?, username = ?, phone = ?,
		    currency_id = ?, country_id = ?, region_id = ?, vat_number = ?,
		    billing_full_name = ?, billing_company_name = ?,
		    billing_address_1 = ?, billing_address_2 = ?,
		    billing_postal_code = ?, billing_city = ?,
		    shipping_full_name = ?, shipping_company_name = ?,
		    shipping_address_1 = ?, shipping_address_2 = ?,
		    shipping_postal_code = ?, shipping_city = ?,
		    is_verified = ?, is_business = ?,
		    is_subscribed_info = ?, is_subscribed_promos = ?, is_admin = ?
		WHERE id = ?`,
		now_str, u.IsEnabled,
		u.Email, u.Username, u.Phone,
		u.CurrencyID, u.CountryID, u.RegionID, u.VATNumber,
		u.BillingAddress.FullName, u.BillingAddress.CompanyName,
		u.BillingAddress.Address1, u.BillingAddress.Address2,
		u.BillingAddress.PostalCode, u.BillingAddress.City,
		u.ShippingAddress.FullName, u.ShippingAddress.CompanyName,
		u.ShippingAddress.Address1, u.ShippingAddress.Address2,
		u.ShippingAddress.PostalCode, u.ShippingAddress.City,
		u.IsVerified, u.IsBusiness,
		u.IsSubscribedInfo, u.IsSubscribedPromos, u.IsAdmin,
		u.ID,
	)
	if err != nil {
		return fmt.Errorf("UserUpdate: %w", err)
	}

	u.UpdatedAt = now
	
	return nil
}

func UserUpdatePassword(db *sql.DB, id uuid.UUID, newPassword string) error {
	u, err := UserReadByID(db, id)
	if err != nil {
		return err
	}

	if err := u.SetPassword(newPassword); err != nil {
		return fmt.Errorf("setting password: %w", err)
	}

	now := time.Now().UTC()
	u.UpdatedAt = now
	u.LastPassChangeAt = &now

	_, err = db.Exec(`
		UPDATE users
		SET
			password_hash = ?,	
			updated_at = ?,
			last_pass_change_at = ?
		WHERE id = ?`,
		u.PasswordHash,
		utils.TimeToString(u.UpdatedAt),
		utils.TimeToString(*u.LastPassChangeAt),
		u.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("UserUpdatePassword: %w", err)
	}

	return nil
}

func UserValidatePassword(db *sql.DB, id uuid.UUID, password string) error {
	u, err := UserReadByID(db, id)
	if err != nil {
		return err
	}

	return u.CheckPassword(password)
}

func UserDelete(db *sql.DB, id uuid.UUID) error {
	if _, err := db.Exec(`DELETE FROM users WHERE id = ?`, id); err != nil {
		return fmt.Errorf("UserDelete: %w", err)
	}
	return nil
}
