// Package data provides the database models and query methods
// for the API Key shared data, including API key management.
package data

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/andres085/real-time-event-processing-system/internal/apikey/validator"
)

type APIKey struct {
	ID          int       `json:"id"`
	KeyHash     []byte    `json:"-"`
	ClientID    string    `json:"client_id"`
	Description string    `json:"description"`
	Environment string    `json:"environment"`
	IsActive    bool      `json:"is_active"`
	RateLimit   int32     `json:"rate_limit_per_minute"`
	CreatedAt   time.Time `json:"-"`
	Version     int32     `json:"version"`
}

func (a *APIKey) Generate(env string) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	plainText := "pk_" + env + "_" + base64.URLEncoding.EncodeToString(bytes)

	a.KeyHash = HashKey([]byte("pepper"), []byte(plainText))

	return plainText, nil
}

func HashKey(pepper, key []byte) []byte {
	mac := hmac.New(sha256.New, pepper)
	mac.Write([]byte(key))

	return mac.Sum(nil)
}

type APIKeyModel struct {
	DB *sql.DB
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 44, "token", "must be 45 bytes long")
}

func (a APIKeyModel) Insert(apiKey *APIKey) error {
	query := `
	INSERT INTO api_keys(client_id, key_hash, description, environment, is_active, rate_limit_per_minute, version)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, version`

	args := []any{
		apiKey.ClientID,
		apiKey.KeyHash,
		apiKey.Description,
		apiKey.Environment,
		apiKey.IsActive,
		apiKey.RateLimit,
		apiKey.Version,
	}

	return a.DB.QueryRow(query, args...).Scan(&apiKey.ID, &apiKey.CreatedAt, &apiKey.Version)
}

func (a APIKeyModel) GetByHash(key []byte) bool {
	query := `
	SELECT EXISTS (
		SELECT 1 FROM api_keys WHERE key_hash = $1
	);`

	var exists bool

	err := a.DB.QueryRow(query, key).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}

func (a APIKeyModel) Get(id int) (*APIKey, error) {
	return nil, nil
}

func (a APIKeyModel) Update(apiKey *APIKey) error {
	return nil
}

func (a APIKeyModel) Delete(id int) error {
	return nil
}
