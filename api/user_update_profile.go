package api

import (   
    //"github.com/google/uuid"
    "github.com/jmoiron/sqlx"

    "gocart/utils"
    //"gocart/models"
)

type UserUpdateProfileReq struct {
    Email               string  `json:"email"`
    Username            string	`json:"username"`
    PasswordCurrent     string  `json:"password_current"`
    PasswordNew         string  `json:"password_new"`
    PasswordNewConfirm  string  `json:"password_new_confirm"`
    IsSubscribedInfo    bool    `json:"is_subscribed_info"`
    IsSubscribedPromos	bool    `json:"is_subscribed_promos"`
    IsDeleteRequested   bool    `json:"is_delete_requested"`
}

func (u *UserUpdateProfileReq) Validate(db *sqlx.DB) *utils.ErrorResponse {
    var errRes utils.ErrorResponse

    utils.ValidateEmail(&errRes, u.Email)
    utils.ValidateUsername(&errRes, u.Username)

    utils.ValidateStrLen(&errRes, "password_current", u.PasswordCurrent, 1, 250)

    utils.ValidatePassword(&errRes, "password_new", u.PasswordNew)
    utils.ValidatePassword(&errRes, "password_new_confirm", u.PasswordNewConfirm)

    return errRes.HasErrors()
}

// Don't forget to update these fields!!
//            last_user_change_at TIMESTAMPTZ,
//            last_pass_change_at TIMESTAMPTZ,
//            last_email_change_at TIMESTAMPTZ,
