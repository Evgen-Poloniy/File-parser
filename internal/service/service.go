package service

import (
	"context"
	"file-parser/internal/dto"
	"file-parser/internal/entity"
	"file-parser/internal/repository"
)

type Reader interface {
	// Get data about device by unit_guid
	GetData(unitGUID string, limit int, page int) ([]dto.Response, error)
}

type Writer interface {
	// Parse tsv file and return slices of data
	ParseFile(ctx context.Context, path string) ([]entity.Device, []entity.DeviceData, error)

	// Load data into database
	LoadParsedData(ctx context.Context, file *entity.File, devices []entity.Device, deviceData []entity.DeviceData) error
}

type QueryService struct {
	Reader
}

type ParserService struct {
	Writer
}

func NewQueryService(reader repository.Reader) *QueryService {
	return &QueryService{
		Reader: newReaderService(reader),
	}
}

func NewParserService(writer repository.Writer) *ParserService {
	return &ParserService{
		Writer: newWriterService(writer),
	}
}
