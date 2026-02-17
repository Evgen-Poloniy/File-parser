package service

import (
	"context"
	"file-parser/internal/dto"
	"file-parser/internal/repository"
)

type ReaderService struct {
	reader repository.Reader
}

func newReaderService(reader repository.Reader) *ReaderService {
	return &ReaderService{
		reader: reader,
	}
}

// Get data about device by unit_guid
func (r *ReaderService) GetDataByUnitGUID(ctx context.Context, unitGUID string, limit int, page int) ([]dto.Response, error) {
	return r.reader.GetDataByUnitGUID(ctx, unitGUID, limit, page)
}

// Get information about parse errors
func (r *ReaderService) GetFileErrorsByFilename(ctx context.Context, filename string, limit int, page int) ([]dto.FileErrorInfo, error) {
	return r.reader.GetFileErrorsByFilename(ctx, filename, limit, page)
}
