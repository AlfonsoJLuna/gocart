package api

import (
    "encoding/json"
    "net/http"

    "github.com/jmoiron/sqlx"

    "gocart/models"
    "gocart/emails"
    "gocart/utils"
)

type VerifyResendReq struct {
    Email string `json:"email"`
}

func (u *VerifyResendReq) Validate() *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    utils.ValidateEmail(&errRes, u.Email)

    return errRes.HasErrors()
}

func AuthVerifyResend(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var req VerifyResendReq

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

        // Check that the email is registered
        emailExists, err := models.UserExistsByEmail(db, req.Email)

        if err != nil {
            answerInternalError(w)
            return
        }

        if !emailExists {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Email is not registered. Please create a new account.")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // Check that the email is still pending verification
        isVerified, err := models.UserGetIsVerifiedByEmail(db, req.Email)

        if err != nil {
            answerInternalError(w)
            return
        }

        if isVerified {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Email is already verified. Please sign in.")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // Get user UUID from the email
        userID, err := models.UserGetIDByEmail(db, req.Email)

        if err != nil {
            answerInternalError(w)
            return
        }

        emails.SendVerificationLink(userID, req.Email)

        w.WriteHeader(http.StatusOK)
    }
}
