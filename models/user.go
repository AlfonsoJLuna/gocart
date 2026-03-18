package models

import (
	"time"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
)

const userBucket = "users"

type Address struct {
	FullName    string 	`json:"full_name"`
	CompanyName string 	`json:"company_name"`
	Address1    string 	`json:"address_1"`
	Address2    string 	`json:"address_2"`
	PostalCode  string 	`json:"postal_code"`
	City        string 	`json:"city"`
}

type User struct {
	Base
	Email               string     	`json:"email"`
	Username            string     	`json:"username"`
	PasswordHash        string     	`json:"password_hash"`
	Phone       		string		`json:"phone"`
	CurrencyISOCode     string		`json:"currency_iso_code"`
	CountryISOCode      string 		`json:"country_iso_code"`
	RegionName          string 		`json:"region_name"`
	VATNumber           string     	`json:"vat_number"`
	BillingAddress      Address		`json:"billing_address"`
	ShippingAddress    	Address   	`json:"shipping_address"`
	IsVerified          bool       	`json:"is_verified"`
	IsBusiness          bool       	`json:"is_business"`
	IsEnabled           bool       	`json:"is_enabled"`
	IsSubscribedInfo    bool       	`json:"is_subscribed_info"`
	IsSubscribedPromos  bool       	`json:"is_subscribed_promos"`
	IsAdmin             bool       	`json:"is_admin"`
	LastLoginAt         *time.Time 	`json:"last_login_at"`
	LastUserChangeAt    *time.Time 	`json:"last_user_change_at"`
	LastPassChangeAt    *time.Time 	`json:"last_pass_change_at"`
	LastEmailChangeAt   *time.Time 	`json:"last_email_change_at"`
}

func hashGenerate(password string) (string, error) {
	// Generate a salted hash from the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func hashValidate(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

func UserCreate(db *bbolt.DB, u *User, password string) (uuid.UUID, error) {
	hash, err := hashGenerate(password)
	if err != nil {
		return uuid.Nil, err
	}
	u.PasswordHash = hash

	if err := db.Update(func(tx *bbolt.Tx) error {
		if err := create(tx, userBucket, u); err != nil {
			return err
		}
		if err := createIndex(tx, userBucket, "email", u.Email, u.ID, true); err != nil {
			return err
		}
		if err := createIndex(tx, userBucket, "username", u.Username, u.ID, true); err != nil {
			return err
		}
		if err := createIndex(tx, userBucket, "currency_iso_code", u.CurrencyISOCode, u.ID, false); err != nil {
			return err
		}
		if err := createIndex(tx, userBucket, "country_iso_code", u.CountryISOCode, u.ID, false); err != nil {
			return err
		}
		if err := createIndex(tx, userBucket, "region_name", u.RegionName, u.ID, false); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return uuid.Nil, err
	}
	return u.ID, nil
}

func UserListAll(db *bbolt.DB, offset, limit int, desc bool) ([]*User, error) {
	return list[User](db, userBucket, offset, limit, desc)
}

func UserListByCountryISOCode(db *bbolt.DB, countryISOCode string, offset, limit int, desc bool) ([]*User, error) {
	return listByIndex[User](db, userBucket, "country_iso_code", countryISOCode, offset, limit, desc)
}

func UserListByCurrencyISOCode(db *bbolt.DB, currencyISOCode string, offset, limit int, desc bool) ([]*User, error) {
	return listByIndex[User](db, userBucket, "currency_iso_code", currencyISOCode, offset, limit, desc)
}

func UserListByRegionName(db *bbolt.DB, regionName string, offset, limit int, desc bool) ([]*User, error) {
	return listByIndex[User](db, userBucket, "region_name", regionName, offset, limit, desc)
}

func UserReadByID(db *bbolt.DB, id uuid.UUID) (*User, error) {
	return read[User](db, userBucket, id)
}

func UserReadByEmail(db *bbolt.DB, email string) (*User, error) {
	return readByIndex[User](db, userBucket, "email", email)
}

func UserReadByUsername(db *bbolt.DB, username string) (*User, error) {
	return readByIndex[User](db, userBucket, "username", username)
}

func UserUpdate(db *bbolt.DB, u *User) error {
	return db.Update(func(tx *bbolt.Tx) error {
		old, err := readTx[User](tx, userBucket, u.ID)
		if err != nil {
			return err
		}
		if err := updateIndex(tx, userBucket, "email", old.Email, u.Email, u.ID, true); err != nil {
			return err
		}
		if err := updateIndex(tx, userBucket, "username", old.Username, u.Username, u.ID, true); err != nil {
			return err
		}
		if err := updateIndex(tx, userBucket, "currency_iso_code", old.CurrencyISOCode, u.CurrencyISOCode, u.ID, false); err != nil {
			return err
		}
		if err := updateIndex(tx, userBucket, "country_iso_code", old.CountryISOCode, u.CountryISOCode, u.ID, false); err != nil {
			return err
		}
		if err := updateIndex(tx, userBucket, "region_name", old.RegionName, u.RegionName, u.ID, false); err != nil {
			return err
		}
		return update(tx, userBucket, u)
	})
}

func UserDelete(db *bbolt.DB, id uuid.UUID) error {
	return db.Update(func(tx *bbolt.Tx) error {
		u, err := readTx[User](tx, userBucket, id)
		if err != nil {
			return err
		}
		if err := deleteIndex(tx, userBucket, "email", u.Email, id); err != nil {
			return err
		}
		if err := deleteIndex(tx, userBucket, "username", u.Username, id); err != nil {
			return err
		}
		if err := deleteIndex(tx, userBucket, "currency_iso_code", u.CurrencyISOCode, id); err != nil {
			return err
		}
		if err := deleteIndex(tx, userBucket, "country_iso_code", u.CountryISOCode, id); err != nil {
			return err
		}
		if err := deleteIndex(tx, userBucket, "region_name", u.RegionName, id); err != nil {
			return err
		}
		return delete(tx, userBucket, id)
	})
}

func UserValidatePassword(db *bbolt.DB, id uuid.UUID, password string) error {
	u, err := UserReadByID(db, id)
	if err != nil {
		return err
	}
	return hashValidate(password, u.PasswordHash)
}
