package utils

import (
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate *validator.Validate

// InitializeValidator sets up the validator with all required validations
func InitializeValidator() {
	validate = validator.New()

	// Register custom validations
	validate.RegisterValidation("datetime", validateDateTime)
	validate.RegisterValidation("uuid4", validateUUID4)
}

// GetValidator returns the initialized validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		InitializeValidator()
	}
	return validate
}

// validateDateTime checks if the string is a valid datetime in RFC3339 format
func validateDateTime(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse(time.RFC3339, dateStr)
	return err == nil
}

// validateUUID4 checks if the UUID is version 4
func validateUUID4(fl validator.FieldLevel) bool {
	// Handle pointer types
	if fl.Field().Kind() == reflect.Ptr {
		if fl.Field().IsNil() {
			return true // nil pointers are valid for omitempty
		}
		// Get the value the pointer points to
		field := fl.Field().Elem()
		if field.Kind() == reflect.String {
			_, err := uuid.Parse(field.String())
			return err == nil
		}
	}

	// Handle direct string values
	if fl.Field().Kind() == reflect.String {
		_, err := uuid.Parse(fl.Field().String())
		return err == nil
	}

	return false
}
