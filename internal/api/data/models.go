package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	ApiKeys ApiKeyModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		ApiKeys: ApiKeyModel{DB: db},
	}
}
