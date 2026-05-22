package api

import (  
    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"

    "gocart/utils"
    "gocart/models"
)

type UserUpdateAddressesReq struct {
    CountryID           uuid.UUID   `json:"country_id"`
    IsBusiness          bool        `json:"is_business"`
    VATNumber           string      `json:"vat_number"`
    Phone               string      `json:"phone"`
    ShippingFullName    string      `json:"shipping_full_name"`
    ShippingCompanyName string      `json:"shipping_company_name"`
    ShippingAddress1    string      `json:"shipping_address_1"`
    ShippingAddress2    string      `json:"shipping_address_2"`
    ShippingPostalCode  string      `json:"shipping_postalcode"`
    ShippingCity        string      `json:"shipping_city"`
    BillingFullName     string      `json:"billing_full_name"`
    BillingCompanyName  string      `json:"billing_company_name"`
    BillingAddress1     string      `json:"billing_address_1"`
    BillingAddress2     string      `json:"billing_address_2"`
    BillingPostalCode   string      `json:"billing_postalcode"`
    BillingCity         string      `json:"billing_city"`
}

func (u *UserUpdateAddressesReq) Validate(db *sqlx.DB) *utils.ErrorResponse {
    var errRes utils.ErrorResponse
    var country models.Country

    utils.ValidateUUID(&errRes, "country_id", u.CountryID)
 
    if err := country.GetByID(db, u.CountryID); err != nil {
        errRes.Add("country_id", "Please select a valid country")
        return &errRes
    }

    utils.ValidateVATNumber(&errRes, u.VATNumber, u.IsBusiness, country.ISOCode)

    utils.ValidatePhone(&errRes, u.Phone)

    utils.ValidatePostalCode(&errRes, "shipping_postalcode", u.ShippingPostalCode, country.ISOCode)
    utils.ValidatePostalCode(&errRes, "billing_postalcode",  u.BillingPostalCode,  country.ISOCode)

    utils.ValidateStrLen(&errRes, "shipping_full_name",    u.ShippingFullName,    1, 48)
    utils.ValidateLatin(&errRes,  "shipping_full_name",    u.ShippingFullName)
    utils.ValidateStrLen(&errRes, "shipping_company_name", u.ShippingCompanyName, 0, 48)
    utils.ValidateLatin(&errRes,  "shipping_company_name", u.ShippingCompanyName)
    utils.ValidateStrLen(&errRes, "shipping_address_1",    u.ShippingAddress1,    1, 48)
    utils.ValidateLatin(&errRes,  "shipping_address_1",    u.ShippingAddress1)
    utils.ValidateStrLen(&errRes, "shipping_address_2",    u.ShippingAddress2,    0, 48)
    utils.ValidateLatin(&errRes,  "shipping_address_2",    u.ShippingAddress2)
    utils.ValidateStrLen(&errRes, "shipping_city",         u.ShippingCity,        1, 37)
    utils.ValidateLatin(&errRes,  "shipping_city",         u.ShippingCity)

    utils.ValidateStrLen(&errRes, "billing_full_name",     u.BillingFullName,     1, 48)
    utils.ValidateLatin(&errRes,  "billing_full_name",     u.BillingFullName) 
    utils.ValidateStrLen(&errRes, "billing_company_name",  u.BillingCompanyName,  0, 48)
    utils.ValidateLatin(&errRes,  "billing_company_name",  u.BillingCompanyName) 
    utils.ValidateStrLen(&errRes, "billing_address_1",     u.BillingAddress1,     1, 48)
    utils.ValidateLatin(&errRes,  "billing_address_1",     u.BillingAddress1) 
    utils.ValidateStrLen(&errRes, "billing_address_2",     u.BillingAddress2,     0, 48)
    utils.ValidateLatin(&errRes,  "billing_address_2",     u.BillingAddress2) 
    utils.ValidateStrLen(&errRes, "billing_city",          u.BillingCity,         1, 37)
    utils.ValidateLatin(&errRes,  "billing_city",          u.BillingCity)

    return errRes.HasErrors()
}
