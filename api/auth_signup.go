package api

import (
    "encoding/json"
    "net/http"

    "github.com/google/uuid"
    "github.com/jmoiron/sqlx"

    "gocart/models"
    "gocart/utils"
    "gocart/emails"
)

type SignUpReq struct {
    Email               string      `json:"email"`
    Username            string      `json:"username"`
    Password            string      `json:"password"`
    PasswordConfirm     string      `json:"password_confirm"`
    IsPolicyAccepted    bool        `json:"is_policy_accepted"`
    CountryID           *uuid.UUID  `json:"country_id"`
    CurrencyID          *uuid.UUID  `json:"currency_id"`
}

func (u *SignUpReq) Validate() *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    utils.ValidateEmail(&errRes, u.Email)
    utils.ValidateUsername(&errRes, u.Username)
    utils.ValidatePassword(&errRes, "password", u.Password)
    utils.ValidatePassword(&errRes, "password_confirm", u.PasswordConfirm)

    if u.Password == u.Email {
        errRes.Add("password", "Please don't use your email as password")
    }

    if u.Password == u.Username {
        errRes.Add("password", "Please don't use your username as password")
    }

    if u.PasswordConfirm != u.Password {
        errRes.Add("password_confirm", "Passwords don't match")
    } 

    if !u.IsPolicyAccepted {
        errRes.Add("is_policy_accepted", "Please accept our terms and policies to continue")
    }

    // CountryID, CurrencyID are optional and allowed to be nil
    
    return errRes.HasErrors()
}

func AuthSignUp(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var req SignUpReq

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
        
        // Check that email is not already taken
        emailExists, err := models.UserExistsByEmail(db, req.Email)

        if err != nil {
            answerInternalError(w)
            return
        }

        if emailExists {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Form validation failed")            
            errRes.Add("email", "Please select a different email")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // Check that username is not already taken
        usernameExists, err := models.UserExistsByUsername(db, req.Username)

        if err != nil {
            answerInternalError(w)
            return
        }

        if usernameExists {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Form validation failed")            
            errRes.Add("username", "Please select a different username")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // Create the user
        userID, err := models.UserCreate(db, req.Email, req.Username, req.Password)

        if err != nil {
            answerInternalError(w)
            return
        }

        // User was created ok

        // Update user country and currency settings, if provided
        models.UserUpdateCountry(db, userID, req.CountryID)
        models.UserUpdateCurrency(db, userID, req.CurrencyID)

        // Send verification email
        emails.SendVerificationLink(userID, req.Email)

        // User created ok
        w.WriteHeader(http.StatusCreated)
    }
}
