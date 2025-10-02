package validator

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	ErrValidationRequiredField   = errors.New("VALIDATION_REQUIRED_FIELD")
	ErrValidationTooShort        = errors.New("VALIDATION_TOO_SHORT")
	ErrValidationContainsSpaces  = errors.New("VALIDATION_CONTAINS_SPACES")
	ErrValidationInvalidFormat   = errors.New("VALIDATION_INVALID_FORMAT")
	ErrValidationInvalidUserRole = errors.New("VALIDATION_INVALID_USER_ROLE")
)

func UsernameValidate(username string) error {
	if username == "" {
		return ErrValidationRequiredField
	}

	if utf8.RuneCountInString(username) < 3 {
		return ErrValidationTooShort
	}

	if regexp.MustCompile(`\s`).MatchString(username) {
		return ErrValidationContainsSpaces
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return ErrValidationInvalidFormat
	}

	return nil
}

func PasswordValidate(password string) error {
	if password == "" {
		return ErrValidationRequiredField
	}

	if utf8.RuneCountInString(password) < 6 {
		return ErrValidationTooShort
	}

	if regexp.MustCompile(`\s`).MatchString(password) {
		return ErrValidationContainsSpaces
	}

	return nil
}

func EmailValidate(email string) (string, error) {
	normalized := strings.ToLower(email)

	if email == "" {
		return "", ErrValidationRequiredField
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if !regex.MatchString(email) {
		return "", ErrValidationInvalidFormat
	}

	return normalized, nil
}

func UserRoleValidate(role string) (string, error) {
	normalized := strings.ToLower(role)

	if normalized != "admin" && normalized != "member" {
		return "", ErrValidationInvalidUserRole
	}

	return normalized, nil
}
