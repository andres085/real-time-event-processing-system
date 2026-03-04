package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	RawEvents RawEventModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		RawEvents: RawEventModel{DB: db},
	}
}
