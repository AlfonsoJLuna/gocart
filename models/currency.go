package models

type Currency struct {
	Base
	ISOCode  	string	`json:"iso_code"`	// 3-letter ISO 4217 code
	Name     	string  `json:"name"`     	// Currency name in English
	NameAlt  	string  `json:"name_alt"` 	// Currency name in local language
	Decimals	int     `json:"decimals"` 	// Decimal places to show in the UI
	FXRate   	float64	`json:"fx_rate"`  	// Exchange rate against EUR
}
