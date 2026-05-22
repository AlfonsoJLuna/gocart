package api

import (
    "encoding/json"
    "net/http"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"

    "gocart/models"
    "gocart/utils"
)

type SignInReq struct {
    EmailOrUsername string      `json:"email_or_username"`
    Password        string      `json:"password"`
    CountryID       *uuid.UUID  `json:"country_id"`
    CurrencyID      *uuid.UUID  `json:"currency_id"`
}

type SignInRes struct {
    Token           string      `json:"token"`
}

func (u *SignInReq) Validate() *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    // We validate loosely to avoid leaking internal limits
    utils.ValidateStrLen(&errRes, "email_or_username", u.EmailOrUsername, 1, 250)
    utils.ValidateStrLen(&errRes, "password", u.Password, 1, 250)

    // CountryID, CurrencyID are optional and allowed to be nil

    return errRes.HasErrors()
}

func AuthSignIn(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var req SignInReq

        // Validate JSON format and parse if valid
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("JSON format is invalid")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // Validate user input format
        if errRes := req.Validate(); errRes != nil {
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // Check if the user exists and get the user ID from the email or username

        // Try by email first
        userID, err := models.UserGetIDByEmail(db, req.EmailOrUsername)

        if err != nil {
            // Email doesn't exist, let's try by username
            userID, err = models.UserGetIDByUsername(db, req.EmailOrUsername)

            if err != nil {
                errRes := utils.ErrorResponse{}
                errRes.SetError("Invalid email, username or password")
                errRes.Write(w, http.StatusBadRequest)
                return
            }
        }

        // User exists, let's check if the entered password is valid
        if err := models.UserValidatePassword(db, userID, req.Password); err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Invalid email, username or password")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // User was identified ok

        // Update user country and currency settings, if provided
        models.UserUpdateCountry(db, userID, req.CountryID)
        models.UserUpdateCurrency(db, userID, req.CurrencyID)

        // Check that the user is verified
        isVerified, err := models.UserGetIsVerifiedByID(db, userID)

        if err != nil {
            answerInternalError(w)
            return
        }

        if !isVerified {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Please verify the email before signing in")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // Update last login date
        if err := models.UserUpdateLastLoginAt(db, userID); err != nil {
            answerInternalError(w)
            return
        }

        // Generate JWT token and return it
        jwt, err := utils.GenerateJwtToken(userID, "auth", 7 * 24)

        if err != nil {
            answerInternalError(w)
            return
        }

        res := SignInRes{Token: jwt}

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(res)
    }
}
