package paseto

import (
	"encoding/hex"
	"os"
	"time"

	"aidanwoods.dev/go-paseto"

	"github.com/GustavoPaula/go-backup-management-api/internal/adapter/config"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/domain"
	"github.com/GustavoPaula/go-backup-management-api/internal/core/port"
	"github.com/google/uuid"
)

type PasetoToken struct {
	token    *paseto.Token
	key      paseto.V4SymmetricKey
	parser   *paseto.Parser
	duration time.Duration
}

func New(config *config.Token) (port.TokenService, error) {
	durationStr := config.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, domain.ErrTokenDuration
	}

	key, err := getPersistentKey(config)
	if err != nil {
		return nil, err
	}

	token := paseto.NewToken()
	parser := paseto.NewParser()

	return &PasetoToken{
		token:    &token,
		key:      key,
		parser:   &parser,
		duration: duration,
	}, nil
}

func getPersistentKey(config *config.Token) (paseto.V4SymmetricKey, error) {
	if keyHex := config.KeyHex; keyHex != "" {
		return keyFromHex(keyHex)
	}

	if key, err := keyFromFile("paseto.key"); err == nil {
		return key, nil
	}

	return generateAndSaveKey("paseto.key")
}

func keyFromHex(hexStr string) (paseto.V4SymmetricKey, error) {
	keyBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return paseto.V4SymmetricKey{}, err
	}

	return paseto.V4SymmetricKeyFromBytes(keyBytes)
}

func keyFromFile(filename string) (paseto.V4SymmetricKey, error) {
	keyBytes, err := os.ReadFile(filename)
	if err != nil {
		return paseto.V4SymmetricKey{}, err
	}

	return keyFromHex(string(keyBytes))
}

func generateAndSaveKey(filename string) (paseto.V4SymmetricKey, error) {
	key := paseto.NewV4SymmetricKey()
	keyBytes := key.ExportBytes()
	keyHex := hex.EncodeToString(keyBytes)

	if err := os.WriteFile(filename, []byte(keyHex), 0600); err != nil {
		return paseto.V4SymmetricKey{}, err
	}

	return key, nil
}

func (pt *PasetoToken) CreateToken(user *domain.User) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	payload := &domain.TokenPayload{
		ID:     id,
		UserID: user.ID,
		Role:   user.Role,
	}

	err = pt.token.Set("payload", payload)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(pt.duration)

	pt.token.SetIssuedAt(issuedAt)
	pt.token.SetNotBefore(issuedAt)
	pt.token.SetExpiration(expiredAt)
	token := pt.token.V4Encrypt(pt.key, nil)

	return token, nil
}

func (pt *PasetoToken) VerifyToken(token string) (*domain.TokenPayload, error) {
	var payload *domain.TokenPayload
	parsedToken, err := pt.parser.ParseV4Local(pt.key, token, nil)
	if err != nil {
		if err.Error() == "this token has expired" {
			return nil, domain.ErrExpiredToken
		}
		return nil, domain.ErrInvalidToken
	}

	err = parsedToken.Get("payload", &payload)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	return payload, nil
}
