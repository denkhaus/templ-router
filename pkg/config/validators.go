package config

import (
	"fmt"
	"net/mail"

	"github.com/denkhaus/templ-router/pkg/shared"
)

// PortValidator validates port numbers (1-65535)
type PortValidator struct {
	fieldName string
}

// NewPortValidator creates a new port validator
func NewPortValidator(fieldName string) *PortValidator {
	return &PortValidator{fieldName: fieldName}
}

// ValidatePort checks if a port number is within valid range
func (pv *PortValidator) ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		// Use field-specific message to match existing tests
		message := "Invalid server port"
		if pv.fieldName != "server.port" {
			message = "Invalid " + pv.fieldName + " port"
		}

		return shared.NewValidationError(message).
			WithDetails(fmt.Sprintf("Port %d is outside valid range 1-65535", port)).
			WithContext("field", pv.fieldName).
			WithContext("value", port).
			WithContext("valid_range", "1-65535")
	}
	return nil
}

// EmailValidator validates email addresses
type EmailValidator struct {
	fieldName string
}

// NewEmailValidator creates a new email validator
func NewEmailValidator(fieldName string) *EmailValidator {
	return &EmailValidator{fieldName: fieldName}
}

// ValidateEmail checks if an email address is valid
func (ev *EmailValidator) ValidateEmail(email string) error {
	if email == "" {
		return shared.NewValidationError("Email cannot be empty").
			WithDetails("Email field is required").
			WithContext("field", ev.fieldName)
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return shared.NewValidationError("Invalid email format").
			WithDetails(fmt.Sprintf("Email '%s' is not a valid email address", email)).
			WithCause(err).
			WithContext("field", ev.fieldName).
			WithContext("value", email)
	}

	return nil
}

// PositiveNumberValidator validates positive numbers
type PositiveNumberValidator struct {
	fieldName string
	minimum   int
}

// NewPositiveNumberValidator creates a new positive number validator
func NewPositiveNumberValidator(fieldName string, minimum int) *PositiveNumberValidator {
	return &PositiveNumberValidator{
		fieldName: fieldName,
		minimum:   minimum,
	}
}

// ValidatePositiveNumber checks if a number is positive and above minimum
func (pnv *PositiveNumberValidator) ValidatePositiveNumber(value int) error {
	if value < pnv.minimum {
		return shared.NewValidationError("Invalid number").
			WithDetails(fmt.Sprintf("Value %d is below minimum of %d", value, pnv.minimum)).
			WithContext("field", pnv.fieldName).
			WithContext("value", value).
			WithContext("minimum", pnv.minimum)
	}
	return nil
}

// RequiredStringValidator validates required string fields
type RequiredStringValidator struct {
	fieldName string
}

// NewRequiredStringValidator creates a new required string validator
func NewRequiredStringValidator(fieldName string) *RequiredStringValidator {
	return &RequiredStringValidator{fieldName: fieldName}
}

// ValidateRequiredString checks if a string is not empty
func (rsv *RequiredStringValidator) ValidateRequiredString(value string) error {
	if value == "" {
		return shared.NewValidationError("Required field cannot be empty").
			WithDetails(fmt.Sprintf("Field '%s' is required", rsv.fieldName)).
			WithContext("field", rsv.fieldName)
	}
	return nil
}

// ConditionalValidator validates a field conditionally
type ConditionalValidator struct {
	condition   bool
	fieldName   string
	validators []Validator
	errorMsg    string
}

// NewConditionalValidator creates a new conditional validator
func NewConditionalValidator(condition bool, fieldName string, errorMsg string, validators ...Validator) *ConditionalValidator {
	return &ConditionalValidator{
		condition:   condition,
		fieldName:   fieldName,
		validators:  validators,
		errorMsg:    errorMsg,
	}
}

// Validate runs all validators if condition is true
func (cv *ConditionalValidator) Validate(value interface{}) error {
	if !cv.condition {
		return nil
	}

	for _, validator := range cv.validators {
		if err := validator.Validate(value); err != nil {
			if appErr, ok := err.(*shared.AppError); ok {
				return appErr.WithContext("conditional_validation", cv.errorMsg)
			}
			return err
		}
	}

	return nil
}

// Validator interface for validation functions
type Validator interface {
	Validate(value interface{}) error
}

// PasswordValidator validates password requirements
type PasswordValidator struct {
	fieldName        string
	minLength        int
	strongRequired   bool
}

// NewPasswordValidator creates a new password validator
func NewPasswordValidator(fieldName string, minLength int, strongRequired bool) *PasswordValidator {
	return &PasswordValidator{
		fieldName:      fieldName,
		minLength:      minLength,
		strongRequired: strongRequired,
	}
}

// ValidatePassword checks password requirements
func (pv *PasswordValidator) ValidatePassword(password string) error {
	if len(password) < pv.minLength {
		return shared.NewValidationError("Password too short").
			WithDetails(fmt.Sprintf("Password must be at least %d characters", pv.minLength)).
			WithContext("field", pv.fieldName).
			WithContext("current_length", len(password)).
			WithContext("required_length", pv.minLength)
	}

	// TODO: Add strong password validation if required
	if pv.strongRequired {
		// Add complexity requirements here
	}

	return nil
}