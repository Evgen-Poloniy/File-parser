package service

import (
	"context"
	"file-parser/internal/dto"
	"file-parser/internal/entity"
	"file-parser/internal/repository"
)

type Reader interface {
	// Get data about device by unit_guid
	GetDataByUnitGUID(ctx context.Context, unitGUID string, limit int, page int) ([]dto.Response, error)

	// Get information about parse errors
	GetFileErrorsByFilename(ctx context.Context, filename string, limit int, page int) ([]dto.FileErrorInfo, error)
}

type Writer interface {
	// Parse tsv file and return slices of data
	ParseFile(ctx context.Context, path string) ([]entity.Device, []entity.DeviceData, error)

	// Load data into the database
	LoadParsedData(ctx context.Context, filename string, devices []entity.Device, deviceData []entity.DeviceData) error

	// Load information about parse errors into the database
	LoadErrors(ctx context.Context, filename string, parseErrors *entity.ParseErrors) error

	// Write data from tsv file into rtf file
	WriteDataToRTF(path string, data []entity.DeviceData) error

	//  Write information about errors from tsv file into rtf file
	WriteErrorsToRTF(path string, parseErrs *entity.ParseErrors) error
}

type QueryService struct {
	Reader
}

type ParserService struct {
	Writer
}

func NewQueryService(reader repository.Reader) *QueryService {
	return &QueryService{
		Reader: NewReaderService(reader),
	}
}

func NewParserService(writer repository.Writer) *ParserService {
	return &ParserService{
		Writer: NewWriterService(writer),
	}
}
