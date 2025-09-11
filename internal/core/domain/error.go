package domain

import "errors"

var (
	ErrInternal                    = errors.New("erro interno")
	ErrDataNotFound                = errors.New("dados não encontrado")
	ErrInvalidCredentials          = errors.New("credenciais inválidas")
	ErrConflictingData             = errors.New("dados em conflito")
	ErrUnauthorized                = errors.New("usuário não tem permissão para acessar o recurso")
	ErrTokenCreation               = errors.New("erro ao criar o token do usuário")
	ErrTokenDuration               = errors.New("erro a duração do token do usuário")
	ErrExpiredToken                = errors.New("erro de token expirado")
	ErrInvalidToken                = errors.New("erro de token inválido")
	ErrEmptyAuthorizationHeader    = errors.New("erro de header vazio")
	ErrInvalidAuthorizationHeader  = errors.New("erro de autorização do header inválida")
	ErrInvalidAuthorizationType    = errors.New("erro de tipo de autorização inválida")
	ErrInvalidAuthorizationPayload = errors.New("erro de payload inválido")
	ErrForbidden                   = errors.New("não tem permissão")
)
