// Package data provides the database models and query methods
// for the API service, including API key management.
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/andres085/real-time-event-processing-system/internal/api/validator"
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

func ValidateClient(v *validator.Validator, client *Client) {
	v.Check(client.Name != "", "title", "must be provided")
}

type MovieModel struct {
	DB *sql.DB
}

func (c ClientModel) Insert(client *Client) error {
	query := `
	INSERT INTO clients(id, name)
	VALUES (gen_random_uuid(), $1)
	RETURNING id, name, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, client.Name).
		Scan(&client.ID, &client.Name, &client.CreatedAt, &client.Version)
}

func (c ClientModel) Get(id string) (*Client, error) {
	query := `
	SELECT  id, created_at, name, version
	FROM clients
	WHERE id = $1`

	var client Client

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).
		Scan(&client.ID, &client.CreatedAt, &client.Name, &client.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &client, nil
}

func (c ClientModel) GetAll(name string, filters Filters) ([]*Client, Metadata, error) {

	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, name, created_at, version
	FROM clients
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	ORDER BY %s %s, id ASC
	LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{name, filters.limit(), filters.offset()}

	rows, err := c.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	clients := []*Client{}

	for rows.Next() {
		var client Client

		err := rows.Scan(
			&totalRecords,
			&client.ID,
			&client.Name,
			&client.CreatedAt,
			&client.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		clients = append(clients, &client)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return clients, metadata, nil
}

func (c ClientModel) Update(client *Client) error {
	query := `
	UPDATE clients
	SET name = $1, version = version + 1
	WHERE id = $2 AND version = $3
	RETURNING version`

	args := []any{client.Name, client.ID, client.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(&client.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (c ClientModel) Delete(id string) error {
	if id == "" {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM clients
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := c.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
