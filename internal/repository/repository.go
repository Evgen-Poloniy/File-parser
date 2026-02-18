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
	GetDataByUnitGUID(ctx context.Context, unitGUID string, limit int, page int) ([]dto.Response, error)

	// Get information about parse errors
	GetFileErrorsByFilename(ctx context.Context, filename string, limit int, page int) ([]dto.FileErrorInfo, error)
}

// Write data at database from parsed files
type Writer interface {
	// Check record in the database. If returning id equals 0, then file must go to parsing
	CheckRecordAboutFile(ctx context.Context, filename string) (int64, error)

	// Check record in the database. If returning id equals 0, then file must go to parsing
	CheckRecordAboutError(ctx context.Context, filename string) (int64, error)

	// Insert information about file and parsed data into the database by one transaction
	InsertParsedData(ctx context.Context, filename string, devices []entity.Device, deviceData []entity.DeviceData) error

	// Insert information about parse errors into the database
	InsertErrors(ctx context.Context, filename string, parseErrors *entity.ParseErrors) error
}

type Repository struct {
	Reader
	Writer
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Reader: NewPostgresRepository(db),
		Writer: NewPostgresRepository(db),
	}
}
