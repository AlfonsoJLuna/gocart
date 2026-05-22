package models

import (
    "log"
    "time"
    "errors"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID                  uuid.UUID   `db:"id"                    json:"id"`
    Username            string      `db:"username"              json:"username"`
    PasswordHash        string      `db:"password_hash"         json:"password_hash"`
    Email               string      `db:"email"                 json:"email"`
    Phone               string      `db:"phone"                 json:"phone"`
    CountryID           *uuid.UUID  `db:"country_id"            json:"country_id"`
    CurrencyID          *uuid.UUID  `db:"currency_id"           json:"currency_id"`
    VATNumber           string      `db:"vat_number"            json:"vat_number"`
    BillingFullName     string      `db:"billing_full_name"     json:"billing_full_name"`
    BillingCompanyName  string      `db:"billing_company_name"  json:"billing_company_name"`
    BillingAddress1     string      `db:"billing_address_1"     json:"billing_address_1"`
    BillingAddress2     string      `db:"billing_address_2"     json:"billing_address_2"`
    BillingPostalCode   string      `db:"billing_postalcode"    json:"billing_postalcode"`
    BillingCity         string      `db:"billing_city"          json:"billing_city"`
    ShippingFullName    string      `db:"shipping_full_name"    json:"shipping_full_name"`
    ShippingCompanyName string      `db:"shipping_company_name" json:"shipping_company_name"`
    ShippingAddress1    string      `db:"shipping_address_1"    json:"shipping_address_1"`
    ShippingAddress2    string      `db:"shipping_address_2"    json:"shipping_address_2"`
    ShippingPostalCode  string      `db:"shipping_postalcode"   json:"shipping_postalcode"`
    ShippingCity        string      `db:"shipping_city"         json:"shipping_city"`
    IsVerified          bool        `db:"is_verified"           json:"is_verified"`
    IsBusiness          bool        `db:"is_business"           json:"is_business"`
    IsEnabled           bool        `db:"is_enabled"            json:"is_enabled"`
    IsSubscribedInfo    bool        `db:"is_subscribed_info"    json:"is_subscribed_info"`
    IsSubscribedPromos  bool        `db:"is_subscribed_promos"  json:"is_subscribed_promos"`
    IsAdmin             bool        `db:"is_admin"              json:"is_admin"`
    LastLoginAt         *time.Time  `db:"last_login_at"         json:"last_login_at"`
    LastUserChangeAt    *time.Time  `db:"last_user_change_at"   json:"last_user_change_at"`
    LastPassChangeAt    *time.Time  `db:"last_pass_change_at"   json:"last_pass_change_at"`
    LastEmailChangeAt   *time.Time  `db:"last_email_change_at"  json:"last_email_change_at"`
    CreatedAt           time.Time   `db:"created_at"            json:"created_at"`
    UpdatedAt           time.Time   `db:"updated_at"            json:"updated_at"`
}

func CreateTableUsers(db *sqlx.DB) {
    query := `
        CREATE EXTENSION IF NOT EXISTS citext;
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            username CITEXT NOT NULL UNIQUE CHECK (char_length(username) <= 30),
            password_hash TEXT NOT NULL,
            email CITEXT NOT NULL UNIQUE CHECK (char_length(email) <= 50),
            phone VARCHAR(15) NOT NULL DEFAULT '',
            country_id UUID REFERENCES countries(id) ON DELETE SET NULL,
            currency_id UUID REFERENCES currencies(id) ON DELETE SET NULL,
            vat_number VARCHAR(15) NOT NULL DEFAULT '',
            billing_full_name VARCHAR(48) NOT NULL DEFAULT '',
            billing_company_name VARCHAR(48) NOT NULL DEFAULT '',
            billing_address_1 VARCHAR(48) NOT NULL DEFAULT '',
            billing_address_2 VARCHAR(48) NOT NULL DEFAULT '',
            billing_postalcode VARCHAR(10) NOT NULL DEFAULT '',
            billing_city VARCHAR(37) NOT NULL DEFAULT '',
            shipping_full_name VARCHAR(48) NOT NULL DEFAULT '',
            shipping_company_name VARCHAR(48) NOT NULL DEFAULT '',
            shipping_address_1 VARCHAR(48) NOT NULL DEFAULT '',
            shipping_address_2 VARCHAR(48) NOT NULL DEFAULT '',
            shipping_postalcode VARCHAR(10) NOT NULL DEFAULT '',
            shipping_city VARCHAR(37) NOT NULL DEFAULT '',
            is_verified BOOLEAN NOT NULL DEFAULT FALSE,
            is_business BOOLEAN NOT NULL DEFAULT FALSE,
            is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
            is_subscribed_info BOOLEAN NOT NULL DEFAULT FALSE,
            is_subscribed_promos BOOLEAN NOT NULL DEFAULT FALSE,
            is_admin BOOLEAN NOT NULL DEFAULT FALSE,
            last_login_at TIMESTAMPTZ,
            last_user_change_at TIMESTAMPTZ,
            last_pass_change_at TIMESTAMPTZ,
            last_email_change_at TIMESTAMPTZ,
            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
            updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );`

    if _, err := db.Exec(query); err != nil {
        log.Fatalln("Error creating users table: %v", err)
    }
}

func generateHash(password string) (string, error) {
    // Generate the salted password hash
    passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

    if err != nil {
        return "", err
    }

    return string(passwordHash), nil
}

func validateHash(password, passwordHash string) error {
    if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
        return err
    }

    return nil
}

func UserCreate(db *sqlx.DB, email, username, password string) (uuid.UUID, error) {
    passwordHash, err := generateHash(password)

    if err != nil {
        return uuid.Nil, err
    }

    var userID uuid.UUID

    query := `INSERT INTO users (email, username, password_hash)
              VALUES ($1, $2, $3)
              RETURNING id`

    err = db.Get(&userID, query, email, username, passwordHash)

    if err != nil {
        return uuid.Nil, err
    }

    return userID, nil
}

func UserExistsByID(db *sqlx.DB, userID uuid.UUID) (bool, error) {
    var exists bool

    query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);`

    err := db.Get(&exists, query, userID)

    if err != nil {
        return false, err
    }
    
    return exists, nil
}

func UserExistsByUsername(db *sqlx.DB, username string) (bool, error) {
    var exists bool

    query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);`

    err := db.Get(&exists, query, username)

    if err != nil {
        return false, err
    }
    
    return exists, nil
}

func UserExistsByEmail(db *sqlx.DB, email string) (bool, error) {
    var exists bool

    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);`

    err := db.Get(&exists, query, email)

    if err != nil {
        return false, err
    }
    
    return exists, nil
}

func UserGetIDByEmail(db *sqlx.DB, email string) (uuid.UUID, error) {
    var id uuid.UUID

    query := `SELECT id FROM users WHERE email = $1;`

    err := db.Get(&id, query, email)
    
    if err != nil {
        return uuid.Nil, err
    }

    return id, nil
}

func UserGetIDByUsername(db *sqlx.DB, username string) (uuid.UUID, error) {
    var id uuid.UUID

    query := `SELECT id FROM users WHERE username = $1;`

    err := db.Get(&id, query, username)
    
    if err != nil {
        return uuid.Nil, err
    }

    return id, nil
}

func UserGetEmailByID(db *sqlx.DB, userID uuid.UUID) (string, error) {
    var email string

    query := `SELECT email FROM users WHERE id = $1;`

    err := db.Get(&email, query, userID)
    
    if err != nil {
        return "", err
    }

    return email, nil
}

func UserGetLastPassChangeAtByID(db *sqlx.DB, userID uuid.UUID) (*time.Time, error) {
    var lastPassChangeAt *time.Time

    query := `SELECT last_pass_change_at FROM users WHERE id = $1;`

    err := db.Get(&lastPassChangeAt, query, userID)
    
    if err != nil {
        return nil, err
    }

    return lastPassChangeAt, nil
}

func UserGetIsVerifiedByID(db *sqlx.DB, userID uuid.UUID) (bool, error) {
    var existsAndIsVerified bool

    query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND is_verified = TRUE);`

    err := db.Get(&existsAndIsVerified, query, userID)

    if err != nil {
        return false, err
    }
    
    return existsAndIsVerified, nil
}

func UserGetIsVerifiedByEmail(db *sqlx.DB, email string) (bool, error) {
    var existsAndIsVerified bool

    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND is_verified = TRUE);`

    err := db.Get(&existsAndIsVerified, query, email)

    if err != nil {
        return false, err
    }
    
    return existsAndIsVerified, nil
}

func UserVerifyByID(db *sqlx.DB, userID uuid.UUID) error {
    query := `UPDATE users
              SET is_verified = TRUE, updated_at = NOW()
              WHERE id = $1 AND is_verified = FALSE;`

    res, err := db.Exec(query, userID)

    if err != nil {
        return err
    }

    rowsAffected, err := res.RowsAffected()

    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("User not found or already verified")
    }

    return nil
}

func UserUpdatePassword(db *sqlx.DB, userID uuid.UUID, password string) error {
    passwordHash, err := generateHash(password)

    if err != nil {
        return err
    }

    query := `UPDATE users
        SET password_hash = $1, last_pass_change_at = NOW(), updated_at = NOW()
        WHERE id = $2`

    res, err := db.Exec(query, passwordHash, userID)

    if err != nil {
        return err
    }

    if rows, _ := res.RowsAffected(); rows == 0 {
        return errors.New("Invalid User ID")
    }

    return nil
}

func UserUpdateCountry(db *sqlx.DB, userID uuid.UUID, countryID *uuid.UUID) error {
    if userID == uuid.Nil {
        return errors.New("User ID can't be blank")
    }

    if (countryID == nil) || (*countryID == uuid.Nil) {
        return errors.New("Country ID can't be blank")
    }

    query := `UPDATE users
        SET country_id = $1, updated_at = NOW()
        WHERE id = $2`

    res, err := db.Exec(query, *countryID, userID)

    if err != nil {
        return err
    }

    if rows, _ := res.RowsAffected(); rows == 0 {
        return errors.New("Invalid User ID or Country ID")
    }

    return nil
}

func UserUpdateCurrency(db *sqlx.DB, userID uuid.UUID, currencyID *uuid.UUID) error {
    if userID == uuid.Nil {
        return errors.New("User ID can't be blank")
    }

    if (currencyID == nil) || (*currencyID == uuid.Nil) {
        return errors.New("Currency ID can't be blank")
    }

    query := `UPDATE users
        SET currency_id = $1, updated_at = NOW()
        WHERE id = $2`

    res, err := db.Exec(query, *currencyID, userID)

    if err != nil {
        return err
    }

    if rows, _ := res.RowsAffected(); rows == 0 {
        return errors.New("Invalid User ID or Currency ID")
    }

    return nil
}

func UserUpdateLastLoginAt(db *sqlx.DB, userID uuid.UUID) error {
    query := `UPDATE users
        SET last_login_at = NOW()
        WHERE id = $1`

    res, err := db.Exec(query, userID)

    if err != nil {
        return err
    }

    if rows, _ := res.RowsAffected(); rows == 0 {
        return errors.New("Invalid User ID")
    }

    return nil
}

func UserValidatePassword(db *sqlx.DB, userID uuid.UUID, password string) error {
    var passwordHash string

    query := `SELECT password_hash FROM users WHERE id = $1;`

    if err := db.Get(&passwordHash, query, userID); err != nil {
        return err
    }

    if err := validateHash(password, passwordHash); err != nil {
        return err
    }

    return nil
}
