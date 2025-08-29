package domain

import "errors"

var (
	ErrInternal           = errors.New("erro interno")
	ErrDataNotFound       = errors.New("dados não encontrado")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrConflictingData    = errors.New("dados em conflito")
	ErrUnauthorized       = errors.New("usuário não tem permissão para acessar o recurso")
	ErrTokenCreation      = errors.New("erro ao criar o token do usuário")
)
