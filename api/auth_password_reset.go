package api

import (
    "encoding/json"
    "net/http"

    "github.com/jmoiron/sqlx"

    "gocart/models"
    "gocart/utils"
)

type PasswordResetReq struct {
    Token           string `json:"token"`
    Password        string `json:"password"`
    PasswordConfirm string `json:"password_confirm"`
}

func (u *PasswordResetReq) Validate() *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    utils.ValidateStrLen(&errRes, "token", u.Token, 100, 500)
    utils.ValidatePassword(&errRes, "password", u.Password)
    utils.ValidatePassword(&errRes, "password_confirm", u.PasswordConfirm)

    if u.PasswordConfirm != u.Password {
        errRes.Add("password_confirm", "Passwords don't match")
    }

    return errRes.HasErrors()
}

func AuthPasswordReset(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var req PasswordResetReq

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

        // Validate token and get user ID and issuedAt timestamp from the token
        userID, issuedAt, err := utils.ValidateJwtToken(req.Token, "reset_password")

        if err != nil {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Token is expired or invalid")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // Check that the user exists
        userExists, err := models.UserExistsByID(db, userID)

        if err != nil {
            answerInternalError(w)
            return
        }

        if !userExists {
            errRes := utils.ErrorResponse{}
            errRes.SetError("User is not registered. Please create a new account.")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        // Check that the token was issued after last password change
        // This makes impossible to reuse a token once the password has been changed
        lastPassChangeAt, err := models.UserGetLastPassChangeAtByID(db, userID)

        if err != nil {
            answerInternalError(w)
            return
        }

        if (lastPassChangeAt != nil) && issuedAt.Before(*lastPassChangeAt) {
            errRes := utils.ErrorResponse{}
            errRes.SetError("Token was issued before the last password change. Please request a new reset link.")
            errRes.Write(w, http.StatusBadRequest)
            return
        }

        err = models.UserUpdatePassword(db, userID, req.Password)

        w.WriteHeader(http.StatusOK)
    }
}
