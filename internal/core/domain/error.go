package domain

import "errors"

var (
	ErrBadRequest                  = errors.New("ERR_BAD_REQUEST")
	ErrServiceUnavailable          = errors.New("ERR_SERVICE_UNAVAILABLE")
	ErrForbidden                   = errors.New("ERR_FORBIDDEN")
	ErrInternal                    = errors.New("ERR_INTERNAL_ERROR")
	ErrDataNotFound                = errors.New("ERR_DATA_NOT_FOUND")
	ErrConflictingData             = errors.New("ERR_CONFLICTING_DATA")
	ErrInvalidCredentials          = errors.New("ERR_INVALID_CREDENTIALS")
	ErrUnauthorized                = errors.New("ERR_UNAUTHORIZED")
	ErrTokenCreation               = errors.New("ERR_TOKEN_CREATION_ERROR")
	ErrTokenDuration               = errors.New("ERR_TOKEN_DURATION_ERROR")
	ErrExpiredToken                = errors.New("ERR_EXPIRED_TOKEN")
	ErrInvalidToken                = errors.New("ERR_INVALID_TOKEN")
	ErrEmptyAuthorizationHeader    = errors.New("ERR_EMPTY_AUTH_HEADER")
	ErrInvalidAuthorizationHeader  = errors.New("ERR_INVALID_AUTH_HEADER")
	ErrInvalidAuthorizationType    = errors.New("ERR_INVALID_AUTH_TYPE")
	ErrInvalidAuthorizationPayload = errors.New("ERR_INVALID_AUTH_PAYLOAD")
)
