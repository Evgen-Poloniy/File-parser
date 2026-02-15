package service

import (
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

func (r *ReaderService) GetData(unitGUID string, limit int, page int) ([]dto.Response, error) {
	return r.reader.GetData(unitGUID, limit, page)
}
