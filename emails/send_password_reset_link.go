package emails

import (
    "le_backend/utils"

    "github.com/google/uuid"
)

func SendPasswordResetLink(userID uuid.UUID, email string) error {
    jwt, err := utils.GenerateJwtToken(userID, "reset_password", 1)

    if err != nil {
        return err
    }

    sendEmail(email, "Reset your password", jwt)

    return nil
}
