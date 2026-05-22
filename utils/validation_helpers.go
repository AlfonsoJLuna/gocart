package utils

import (
    "fmt"
    "regexp"
    "unicode/utf8"

    "github.com/google/uuid"
)

func ValidateStrLen(errRes *ErrorResponse, field, value string, min, max int) {
    length := utf8.RuneCountInString(value) // Count UTF-8 characters (runes), not bytes

    if (min > 0) && (length == 0) {
        errRes.Add(field, "Please enter this field")
        return
    }

    if (min == 0 || min == 1) && (length > max) {
        errRes.Add(field, fmt.Sprintf("Please enter %d characters or less", max))
        return
    }

    if ((length < min) || (length > max)) {
        errRes.Add(field, fmt.Sprintf("Please enter between %d and %d characters", min, max))
    }
}

func ValidateStrLenBytes(errRes *ErrorResponse, field, value string, min, max int) {
    length := len(value) // Count ASCII characters or bytes, not UTF-8 runes

    if (min > 0) && (length == 0) {
        errRes.Add(field, "Please enter this field")
        return
    }

    if (min == 0 || min == 1) && (length > max) {
        errRes.Add(field, fmt.Sprintf("Please enter %d characters or less", max))
        return
    }

    if ((length < min) || (length > max)) {
        errRes.Add(field, fmt.Sprintf("Please enter between %d and %d characters", min, max))
    }
}

func ValidateIntRange(errRes *ErrorResponse, field string, value, min, max int) {
    if (value < min) || (value > max) {
        errRes.Add(field, fmt.Sprintf("Please enter a number between %d and %d", min, max))
    }
}

func ValidateFloatRange(errRes *ErrorResponse, field string, value, min, max float64) {
    if (value < min) || (value > max) {
        errRes.Add(field, fmt.Sprintf("Please enter a number between %.1f and %.1f", min, max))
    }
}

func ValidateRegex(errRes *ErrorResponse, field string, value, regex, message string) {
    if !regexp.MustCompile(regex).MatchString(value) {
        errRes.Add(field, message)
    }
}

func ValidateUUID(errRes *ErrorResponse, field string, id uuid.UUID) {
    if id == uuid.Nil {
        errRes.Add(field, "Please select a valid value")
    }
}
