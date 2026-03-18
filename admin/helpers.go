package admin

import (
	"strconv"
	"errors"

	"gocart/models"
)

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func friendlyError(err error) string {
	if errors.Is(err, models.ErrAlreadyExists) {
		return "One of the values you entered is already in use by another record."
	}
	if errors.Is(err, models.ErrNotFound) {
		return "Record not found."
	}
	if errors.Is(err, models.ErrBadInput) {
		return "Invalid input."
	}
	return "An unexpected error occurred."
}
