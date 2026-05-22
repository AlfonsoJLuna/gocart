package emails

import (
    "le_backend/utils"

    "github.com/google/uuid"
)

func SendVerificationLink(userID uuid.UUID, email string) error {
    jwt, err := utils.GenerateJwtToken(userID, "validate_email", 24)

    if err != nil {
        return err
    }

    sendEmail(email, "Verify your email", jwt)

    return nil
}
