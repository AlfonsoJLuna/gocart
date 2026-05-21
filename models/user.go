package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Base
	Email        		string		`json:"email"`
	Username     		string		`json:"username"`
	PasswordHash 		string		`json:"-"`

	CurrencyID 			*uuid.UUID	`json:"currency_id"`
	CountryID  			*uuid.UUID	`json:"country_id"`
	RegionID   			*uuid.UUID	`json:"region_id"`

	VATNumber  			string     	`json:"vat_number"`
	Phone        		string 		`json:"phone"`

	BillingAddress  	Address		`json:"billing_address"`
	ShippingAddress 	Address		`json:"shipping_address"`

	IsVerified         	bool 		`json:"is_verified"`
	IsBusiness         	bool 		`json:"is_business"`
	IsSubscribedInfo   	bool 		`json:"is_subscribed_info"`
	IsSubscribedPromos	bool		`json:"is_subscribed_promos"`
	IsAdmin            	bool 		`json:"is_admin"`

	LastLoginAt       	*time.Time 	`json:"last_login_at"`
	LastUserChangeAt  	*time.Time 	`json:"last_user_change_at"`
	LastPassChangeAt  	*time.Time 	`json:"last_pass_change_at"`
	LastEmailChangeAt 	*time.Time	`json:"last_email_change_at"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)

	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte(password),
	)
}
