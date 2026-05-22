package emails

import (
    "log"
)

func sendEmail(to, subject, body string) error {
    log.Printf("Sending email to %s: %s\n%s\n", to, subject, body)

	// To do: Store email in the database (pending_emails) table
	// To do: Add an async worker to send pending emails from the database

    return nil
}
