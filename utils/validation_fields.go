package utils

import (

)

func ValidateUsername(errRes *ErrorResponse, username string) {
    ValidateStrLen(errRes, "username", username, 3, 20)
    ValidateRegex(errRes, "username", username, `(?:.*[A-Za-z]){3,}`, "Please enter 3 letters or more")
    ValidateRegex(errRes, "username", username, `^[A-Za-z0-9_]+$`, "Please use only letters, numbers and underscores")
    ValidateRegex(errRes, "username", username, `^[A-Za-z0-9]`, "Please don't use underscore as the first character")
}

func ValidatePassword(errRes *ErrorResponse, field, password string) {
    ValidateStrLen(errRes, field, password, 8, 72)
    ValidateStrLenBytes(errRes, field, password, 8, 72) // Bcrypt max is 72 bytes
    ValidateRegex(errRes, field, password, `[a-z]`, "Please include lowercase, uppercase and numbers")
    ValidateRegex(errRes, field, password, `[A-Z]`, "Please include lowercase, uppercase and numbers")
    ValidateRegex(errRes, field, password, `[0-9]`, "Please include lowercase, uppercase and numbers")
}

func ValidateEmail(errRes *ErrorResponse, email string) {
    ValidateStrLen(errRes, "email", email, 5, 50)
    ValidateRegex(errRes, "email", email, `^[^@\s]+@[^@\s]+\.[^@\s]+$`, "Please enter a valid email")
}

func ValidatePhone(errRes *ErrorResponse, phone string) {
    ValidateStrLen(errRes, "phone", phone, 3, 15)
    ValidateRegex(errRes, "phone", phone, `^[0-9+\-(). ]+$`, "Please enter a valid phone")
}

func ValidateLatin(errRes *ErrorResponse, field, value string) {
    ValidateRegex(errRes, field, value, `^[0-9A-Za-zÀ-ÖØ-öø-ÿ !&'()*+,\-./:=?@\\_~¡¿ªº]+$`, "Please use only latin characters")
}

func ValidatePostalCode(errRes *ErrorResponse, field, postalCode, countryCode string) {
    ValidateStrLen(errRes, field, postalCode, 1, 10)
    ValidateRegex(errRes, field, postalCode, `^[A-Za-z0-9 \-/.()]+$`, "Please enter a valid postal code")

    switch countryCode {
        case "AD": ValidateRegex(errRes, field, postalCode, `^AD[0-9]{3}$`,           "Please enter AD followed by 3 digits")
        case "DE": ValidateRegex(errRes, field, postalCode, `^[0-9]{5}$`,             "Please enter a valid German postal code")
        case "ES": ValidateRegex(errRes, field, postalCode, `^[0-5][0-9]{4}$`,        "Please enter a valid Spanish postal code")
        case "US": ValidateRegex(errRes, field, postalCode, `^[0-9]{5}(-[0-9]{4})?$`, "Please enter a valid US postal code")
    }
}

func ValidateVATNumber(errRes *ErrorResponse, vatNumber string, isBusiness bool, countryCode string) {
    if countryCode == "ES" {
        if isBusiness {
            ValidateStrLen(errRes, "vat_number", vatNumber, 1, 11)
        }
        ValidateRegex(errRes, "vat_number", vatNumber, `^(?:ES)?[A-Z0-9][0-9]{7}[A-Z0-9]$|^$`, "Please enter a valid Spanish NIF")
    } else {
        ValidateStrLen(errRes, "vat_number", vatNumber, 0, 15)
        ValidateRegex(errRes, "vat_number", vatNumber, `^[A-Z0-9 \-.]+$|^$`, "Please enter a valid VAT number or leave blank")
    }
}
