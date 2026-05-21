package models

import (
	"github.com/google/uuid"
)

type Region struct {
	Base
	CountryID	uuid.UUID	`json:"country_id"`	// Parent country UUID
	Name      	string    	`json:"name"`       // Region name in English
	NameAlt   	string    	`json:"name_alt"`   // Region name in local language
	IsEU      	bool      	`json:"is_eu"`      // This region is part of the EU
	VATRate   	float64   	`json:"vat_rate"`   // Region-specific VAT rate
}
