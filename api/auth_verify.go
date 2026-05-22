package api

import (
    "encoding/json"
    "net/http"

    "github.com/jmoiron/sqlx"

    "gocart/models"
    "gocart/utils"
)

type VerifyReq struct {
    Token string `json:"token"`
}

func (u *VerifyReq) Validate() *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    utils.ValidateStrLen(&errRes, "token", u.Token, 100, 500)

    return errRes.HasErrors()
}

func AuthVerify(db *sqlx.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var req VerifyReq

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

        // Validate token and get user ID from the token
        userID, _, err := utils.ValidateJwtToken(req.Token, "validate_email")

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

        // Check that the user is still pending verification
        isVerified, err := models.UserGetIsVerifiedByID(db, userID)

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

        // Mark the user as verified
        if err := models.UserVerifyByID(db, userID); err != nil {
            answerInternalError(w)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}
