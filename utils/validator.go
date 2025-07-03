package utils

import (
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
	str := fl.Field().String()
	_, err := uuid.Parse(str)
	return err == nil
}
