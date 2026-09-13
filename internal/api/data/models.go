package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	APIKeys APIKeyModel
	Clients ClientModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		APIKeys: APIKeyModel{DB: db},
		Clients: ClientModel{DB: db},
	}
}
