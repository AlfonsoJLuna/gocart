package api

import (    
    "github.com/google/uuid"

    "gocart/utils"
)

type UserUpdateLocationReq struct {
    CountryID   uuid.UUID   `json:"country_id"`
    CurrencyID  uuid.UUID   `json:"currency_id"`
}

func (u *UserUpdateLocationReq) Validate() *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    utils.ValidateUUID(&errRes, "country_id", u.CountryID)
    utils.ValidateUUID(&errRes, "currency_id", u.CurrencyID)

    return errRes.HasErrors()
}
