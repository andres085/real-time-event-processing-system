package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Clients ClientModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Clients: ClientModel{DB: db},
	}
}
