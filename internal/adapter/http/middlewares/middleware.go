package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/http/response"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
)

type contextKey string

const (
	authorizationHeaderKey  = "authorization"
	authorizationType       = "bearer"
	authorizationPayloadKey = contextKey("authorization_payload")
)

func AuthMiddleware(token port.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorizationHeader := r.Header.Get(authorizationHeaderKey)

			isEmpty := len(authorizationHeader) == 0
			if isEmpty {
				response.JSON(w, http.StatusUnauthorized, "Falha na autenticação!", nil, domain.ErrEmptyAuthorizationHeader.Error())
				return
			}

			fields := strings.Fields(authorizationHeader)
			isValid := len(fields) == 2
			if !isValid {
				response.JSON(w, http.StatusUnauthorized, "Falha na autenticação!", nil, domain.ErrInvalidAuthorizationHeader.Error())
				return
			}

			currentAuthorizationType := strings.ToLower(fields[0])
			if currentAuthorizationType != authorizationType {
				response.JSON(w, http.StatusUnauthorized, "Falha na autenticação!", nil, domain.ErrInvalidAuthorizationHeader.Error())
				return
			}

			accessToken := fields[1]
			payload, err := token.VerifyToken(accessToken)
			if err != nil {
				response.JSON(w, http.StatusUnauthorized, "Falha na autenticação!", nil, domain.ErrUnauthorized.Error())
				return
			}

			ctx := context.WithValue(r.Context(), authorizationPayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload, ok := r.Context().Value(authorizationPayloadKey).(*domain.TokenPayload)
			if !ok {
				response.JSON(w, http.StatusUnauthorized, "Falha na autenticação!", nil, domain.ErrInvalidAuthorizationPayload.Error())
				return
			}

			isAdmin := payload.Role == domain.Admin
			if !isAdmin {
				response.JSON(w, http.StatusForbidden, "Falha na autenticação!", nil, domain.ErrForbidden.Error())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
