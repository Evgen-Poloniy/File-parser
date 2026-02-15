package repository

import (
	"context"
	"database/sql"
	"file-parser/internal/dto"
	"file-parser/internal/entity"
)

// Read data from database
type Reader interface {
	// Get data about device by unit_guid
	GetData(unitGUID string, limit int, page int) ([]dto.Response, error)
}

// Write data at database from parsed files
type Writer interface {
	// Check record in the database. If returning id equals 0, then file must go to parsing
	CheckRecordAboutFile(ctx context.Context, filename string) (int64, error)

	//Insert information about file and parsed data into the database by one transaction
	InsertParsedData(ctx context.Context, file *entity.File, devices []entity.Device, deviceData []entity.DeviceData) error
}

type Repository struct {
	Reader
	Writer
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Reader: newPostgreSQLRepository(db),
		Writer: newPostgreSQLRepository(db),
	}
}
