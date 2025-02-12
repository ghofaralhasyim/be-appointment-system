package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetFriendlyErrorMessage(e validator.FieldError) string {
	field := strings.ToLower(e.Field())
	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must have at least " + e.Param() + " characters"
	case "max":
		return field + " must have no more than " + e.Param() + " characters"
	case "strongpassword":
		return field + " must be at least 8 characters long and include an uppercase letter, lowercase letter, number, and special character"
	default:
		return e.Error()
	}
}
