package data

import (
	"database/sql"
	"time"
)

type ApiKey struct {
	ID          int       `json:"id"`
	KeyHash     string    `json:"-"`
	ClientId    int       `json:"client_id"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	RateLimit   int32     `json:"rate_limit_per_minute"`
	CreatedAt   time.Time `json:"-"`
	Version     int32     `json:"version"`
}

type ApiKeyModel struct {
	DB *sql.DB
}

func (a ApiKeyModel) Insert(apiKey *ApiKey) error {
	query := `
	INSERT INTO api_keys(client_id, key_hash, description, is_active, rate_limit_per_minute, version)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, version`

	args := []any{apiKey.ClientId, apiKey.KeyHash, apiKey.Description, apiKey.IsActive, apiKey.RateLimit, apiKey.Version}

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
