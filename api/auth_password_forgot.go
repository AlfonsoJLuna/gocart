package api

import (
    "encoding/json"
    "net/http"

    "github.com/jmoiron/sqlx"

    "gocart/models"
    "gocart/emails"
    "gocart/utils"
)

type PasswordForgotReq struct {
    EmailOrUsername string `json:"email_or_username"`
}

func (u *PasswordForgotReq) Validate() *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    // We validate loosely to avoid leaking internal limits
    utils.ValidateStrLen(&errRes, "email_or_username", u.EmailOrUsername, 1, 250)

    return errRes.HasErrors()
}

func AuthPasswordForgot(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var req PasswordForgotReq

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
                errRes.SetError("Invalid email or username")
                errRes.Write(w, http.StatusBadRequest)
                return
            }
        }

        // Get the email from the user ID
        email, err := models.UserGetEmailByID(db, userID)

        if err != nil {
            answerInternalError(w)
            return
        }

        emails.SendPasswordResetLink(userID, email)

        w.WriteHeader(http.StatusOK)
    }
}
