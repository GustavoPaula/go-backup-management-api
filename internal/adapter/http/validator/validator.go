package validator

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"
)

func UsernameValidate(username string) error {
	if username == "" {
		return errors.New("username deve ser preenchido")
	}

	if utf8.RuneCountInString(username) < 3 {
		return errors.New("o username deve ter no mínimo 3 caracteres")
	}

	if regexp.MustCompile(`\s`).MatchString(username) {
		return errors.New("o username não pode conter espaços")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return errors.New("o username não pode conter caracteres especiais sem ser o underscore")
	}

	return nil
}

func PasswordValidate(password string) error {
	if password == "" {
		return errors.New("password deve ser preenchido")
	}

	if utf8.RuneCountInString(password) < 6 {
		return errors.New("o password deve ter no mínimo 6 caracteres")
	}

	if regexp.MustCompile(`\s`).MatchString(password) {
		return errors.New("o password não pode conter espaços")
	}

	return nil
}

func EmailValidate(email string) (string, error) {
	normalized := strings.ToLower(email)

	if email == "" {
		return "", errors.New("o e-mail deve ser preenchido")
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if !regex.MatchString(email) {
		return "", errors.New("formato de e-mail inválido")
	}

	return normalized, nil
}

func UserRoleValidate(role string) (string, error) {
	normalized := strings.ToLower(role)

	if normalized != "admin" && normalized != "member" {
		return "", errors.New("valor de role precisa ser 'admin' ou 'member'")
	}

	return normalized, nil
}
