package data

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ApiKey struct {
	ID          int       `json:"id"`
	KeyHash     key       `json:"-"`
	ClientId    int       `json:"client_id"`
	Description string    `json:"description"`
	Environment string    `json:"environment"`
	IsActive    bool      `json:"is_active"`
	RateLimit   int32     `json:"rate_limit_per_minute"`
	CreatedAt   time.Time `json:"-"`
	Version     int32     `json:"version"`
}

type key []byte

func (k *key) Generate(env string) (string, error) {

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	plainText := "pk_" + env + "_" + base64.URLEncoding.EncodeToString(bytes)

	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	*k = hash

	return plainText, nil
}

type ApiKeyModel struct {
	DB *sql.DB
}

func (a ApiKeyModel) Insert(apiKey *ApiKey) error {
	query := `
	INSERT INTO api_keys(client_id, key_hash, description, environment, is_active, rate_limit_per_minute, version)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, created_at, version`

	args := []any{apiKey.ClientId, apiKey.KeyHash, apiKey.Description, apiKey.Environment, apiKey.IsActive, apiKey.RateLimit, apiKey.Version}

	return a.DB.QueryRow(query, args...).Scan(&apiKey.ID, &apiKey.CreatedAt, &apiKey.Version)
}

func (a ApiKeyModel) Get(id int) (*ApiKey, error) {
	return nil, nil
}
func (a ApiKeyModel) Update(apiKey *ApiKey) error {
	return nil
}
func (a ApiKeyModel) Delete(id int) error {
	return nil
}
