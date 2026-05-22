package models

import (
	"github.com/google/uuid"
)

type Country struct {
	Base
	ISOCode    	string    	`json:"iso_code"`           // 2-letter ISO 3166-1 alpha-2 code
	Name       	string    	`json:"name"`               // Country name in English
	NameAlt    	string    	`json:"name_alt"`           // Country name in local language
	CurrencyID	*uuid.UUID	`json:"currency_id"`        // Default currency for this country
	IsEU       	bool      	`json:"is_eu"`              // This country is part of the EU
	VATRate    	float64   	`json:"vat_rate"`           // Default VAT rate for this country
}

func (c *Country) Validate() *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    utils.ValidateRegex(&errRes, "iso_code", c.ISOCode, `^[A-Z]{2}$`, "Please enter 2 uppercase letters")

    utils.ValidateStrLen(&errRes, "name_en", c.NameEN, 1, 100)

    utils.ValidateStrLen(&errRes, "name_es", c.NameES, 1, 100)

    utils.ValidateFloatRange(&errRes, "vat_rate", c.VATRate, 0.0, 50.0)

    return errRes.HasErrors()
}
