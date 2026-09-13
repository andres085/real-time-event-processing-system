// Package data provides the database models and query methods
// for the API service, including API key management.
package data

import (
	"database/sql"
	"time"
)

type Client struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}

type ClientModel struct {
	DB *sql.DB
}

func (c ClientModel) Insert(client *Client) error {
	query := `
	INSERT INTO clients(id, name, version)
	VALUES (gen_random_uuid(), $1, $2)
	RETURNING id, name, created_at, version`

	// args := []any{
	// 	client.Name,
	// 	client.Version,
	// }

	return c.DB.QueryRow(query, client.Name, client.Version).
		Scan(&client.ID, &client.Name, &client.CreatedAt, &client.Version)
}

func (c ClientModel) Get(id string) (*Client, error) {
	return nil, nil
}
