package utils

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func GetFriendlyErrorMessage(e validator.FieldError, s interface{}) (string, string) {

	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, _ := t.FieldByName(e.StructField())
	jsonTag := field.Tag.Get("json")

	if jsonTag == "" {
		jsonTag = strings.ToLower(e.Field())
	}

	switch e.Tag() {
	case "required":
		return jsonTag, jsonTag + " is required"
	case "email":
		return jsonTag, jsonTag + " must be a valid email address"
	case "min":
		return jsonTag, jsonTag + " must have at least " + e.Param() + " characters"
	case "max":
		return jsonTag, jsonTag + " must have no more than " + e.Param() + " characters"
	case "ISOdate":
		return jsonTag, jsonTag + "must in ISO 8601 date format"
	default:
		return "", e.Error()
	}
}

func ISODateValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() != reflect.Struct {
		return false
	}

	t, ok := field.Interface().(time.Time)
	if !ok {
		return false
	}

	return !t.IsZero()
}
