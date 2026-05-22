package utils

import (
    "fmt"
    "strings"
    "encoding/json"
    "net/http"
)

type ErrorResponse struct {
    ErrorMsg    string          `json:"error,omitempty"`
    Fields      []FieldError    `json:"fields,omitempty"`
}

type FieldError struct {
    Field       string          `json:"field"`
    Message     string          `json:"message"`
}

// Set main/general error
func (err *ErrorResponse) SetError(msg string) {
    err.ErrorMsg = msg
}

// Add a field-specific error
func (err *ErrorResponse) Add(field string, msg string) {
    // We check if the field already has an error
    for _, fe := range err.Fields {
        if fe.Field == field {
            // Field already has an error -> ignore new one
            return
        }
    }

    // Add new field with error, if not already present
    err.Fields = append(err.Fields, FieldError{Field: field, Message: msg})
}

// Check if there are any errors, and return nil if there is none
func (err *ErrorResponse) HasErrors() *ErrorResponse {
    if (err.ErrorMsg == "") && (len(err.Fields) == 0) {
        return nil
    }

    return err
}

// Send error as a JSON HTTP response
func (err ErrorResponse) Write(w http.ResponseWriter, status int) {
    // Auto-fill general error if not explicitly set but fields exist
    if (err.ErrorMsg == "") && (len(err.Fields) > 0) {
        err.ErrorMsg = "Form validation failed"
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(err)
}

// Convert error to a raw string, satisfying Go's interface
func (err ErrorResponse) Error() string {
    var parts []string

    // Add main/general error if present
    if err.ErrorMsg != "" {
        parts = append(parts, err.ErrorMsg)
    }

    // Add field-specific errors
    for _, fe := range err.Fields {
        parts = append(parts, fmt.Sprintf("%s: %s", fe.Field, fe.Message))
    }

    // If nothing at all, return fallback
    if len(parts) == 0 {
        return "Unknown error."
    }

    // Return a string including all info in the ErrorResponse
    return strings.Join(parts, "; ")
}
