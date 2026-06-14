package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	RequestVolumeAggregationData RequestVolumeAggregationDataModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		RequestVolumeAggregationData: RequestVolumeAggregationDataModel{DB: db},
	}
}
