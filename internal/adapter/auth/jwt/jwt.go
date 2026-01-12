package jwt

import (
	"time"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/config"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtToken struct {
	secretKey []byte
	duration  time.Duration
}

type jwtClaims struct {
	ID     uuid.UUID       `json:"id"`
	UserID uuid.UUID       `json:"user_id"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func New(config *config.Token) (port.TokenService, error) {
	if config.JwtSecretKey == "" {
		return nil, domain.ErrTokenRequired
	}

	durationStr := config.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, domain.ErrTokenDuration
	}

	return &JwtToken{
		secretKey: []byte(config.JwtSecretKey),
		duration:  duration,
	}, nil
}

func (j *JwtToken) CreateToken(user *domain.User) (string, error) {
	if user == nil {
		return "", domain.ErrDataNotFound
	}

	tokenID := uuid.New()

	claims := jwtClaims{
		ID:     tokenID,
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-backup-management-api",
			Subject:   user.ID.String(),
			ID:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	return signedToken, nil
}

func (j *JwtToken) VerifyToken(tokenString string) (*domain.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, domain.ErrExpiredToken
	}

	return &domain.TokenPayload{
		ID:     claims.ID,
		UserID: claims.UserID,
		Role:   claims.Role,
	}, nil
}
