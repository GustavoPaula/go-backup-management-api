package validator

import (
	"errors"
	"regexp"
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

	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username) {
		return errors.New("o username não pode conter caracteres especiais")
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

func EmailValidate(email string) error {
	if email == "" {
		return errors.New("o e-mail deve ser preenchido")
	}

	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if !regex.MatchString(email) {
		return errors.New("formato de e-mail inválido")
	}

	return nil
}
