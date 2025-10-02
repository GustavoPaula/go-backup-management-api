package domain

import "errors"

var (
	ErrInternal                    = errors.New("INTERNAL_ERROR")
	ErrDataNotFound                = errors.New("DATA_NOT_FOUND")
	ErrConflictingData             = errors.New("CONFLICTING_DATA")
	ErrInvalidCredentials          = errors.New("INVALID_CREDENTIALS")
	ErrUnauthorized                = errors.New("UNAUTHORIZED")
	ErrTokenCreation               = errors.New("TOKEN_CREATION_ERROR")
	ErrTokenDuration               = errors.New("TOKEN_DURATION_ERROR")
	ErrExpiredToken                = errors.New("EXPIRED_TOKEN")
	ErrInvalidToken                = errors.New("INVALID_TOKEN")
	ErrEmptyAuthorizationHeader    = errors.New("EMPTY_AUTH_HEADER")
	ErrInvalidAuthorizationHeader  = errors.New("INVALID_AUTH_HEADER")
	ErrInvalidAuthorizationType    = errors.New("INVALID_AUTH_TYPE")
	ErrInvalidAuthorizationPayload = errors.New("INVALID_AUTH_PAYLOAD")
	ErrForbidden                   = errors.New("FORBIDDEN")
)
