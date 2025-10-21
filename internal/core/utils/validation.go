package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func TranslateValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("O campo '%s' é obrigatório", strings.ToLower(e.Field()))
	case "email":
		return fmt.Sprintf("O campo '%s' deve ser um email válido", strings.ToLower(e.Field()))
	case "oneof":
		return fmt.Sprintf("O campo '%s' deve ser um dos valores permitidos: %s", strings.ToLower(e.Field()), e.Param())
	case "uuid4":
		return fmt.Sprintf("O campo '%s' deve conter um UUID v4 válido", strings.ToLower(e.Field()))
	case "min":
		return fmt.Sprintf("O campo '%s' deve ter no mínimo '%s' caracteres", strings.ToLower(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("O campo '%s' deve ter no máximo '%s' caracteres", strings.ToLower(e.Field()), e.Param())
	case "alphanum":
		return fmt.Sprintf("O campo '%s' deve contém apenas caracteres alfanuméricos", strings.ToLower(e.Field()))
	default:
		return fmt.Sprintf("O campo '%s' é inválido (%s)", strings.ToLower(e.Field()), e.Tag())
	}
}

func ValidationErrorsToMap(err error) map[string]string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		errorsMap := make(map[string]string)
		for _, e := range errs {
			errorsMap[strings.ToLower(e.Field())] = TranslateValidationError(e)
		}
		return errorsMap
	}
	return nil
}
