package models

type Currency struct {
	Base
	ISOCode  	string	`json:"iso_code"`	// 3-letter ISO 4217 code
	Name     	string  `json:"name"`     	// Currency name in English
	NameAlt  	string  `json:"name_alt"` 	// Currency name in local language
	Decimals	int     `json:"decimals"` 	// Decimal places to show in the UI
	FXRate   	float64	`json:"fx_rate"`  	// Exchange rate against EUR
}

func (c *Currency) Validate() *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    utils.ValidateRegex(&errRes, "iso_code", c.ISOCode, `^[A-Z]{3}$`, "Please enter 3 uppercase letters")

    utils.ValidateRegex(&errRes, "flag_name", c.FlagName, `^[A-Z]{2}$`, "Please enter 2 uppercase letters")

    utils.ValidateStrLen(&errRes, "name_en", c.NameEN, 1, 100)

    utils.ValidateStrLen(&errRes, "name_es", c.NameES, 1, 100)

    utils.ValidateIntRange(&errRes, "decimals", c.Decimals, 0, 6)

    return errRes.HasErrors()
}
