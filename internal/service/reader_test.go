package service_test

import (
	"context"
	"errors"
	"file-parser/internal/dto"
	"file-parser/internal/service"
	"testing"
)

type mockReader struct {
	GetDataByUnitGUIDFunc       func(ctx context.Context, unitGUID string, limit int, page int) ([]dto.Response, error)
	GetFileErrorsByFilenameFunc func(ctx context.Context, filename string, limit int, page int) ([]dto.FileErrorInfo, error)
}

func (m *mockReader) GetDataByUnitGUID(ctx context.Context, unitGUID string, limit int, page int) ([]dto.Response, error) {
	return m.GetDataByUnitGUIDFunc(ctx, unitGUID, limit, page)
}

func (m *mockReader) GetFileErrorsByFilename(ctx context.Context, filename string, limit int, page int) ([]dto.FileErrorInfo, error) {
	return m.GetFileErrorsByFilenameFunc(ctx, filename, limit, page)
}

func TestGetDataByUnitGUID(t *testing.T) {
	mock := &mockReader{
		GetDataByUnitGUIDFunc: func(ctx context.Context, unitGUID string, limit int, page int) ([]dto.Response, error) {
			if unitGUID == "fail" {
				return nil, errors.New("some error")
			}
			return []dto.Response{
				{UnitGUID: unitGUID, InvID: "inv1"},
			}, nil
		},
	}

	svc := service.NewReaderService(mock)

	data, err := svc.GetDataByUnitGUID(context.Background(), "guid1", 10, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 1 || data[0].UnitGUID != "guid1" {
		t.Fatalf("unexpected data returned: %+v", data)
	}

	_, err = svc.GetDataByUnitGUID(context.Background(), "fail", 10, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetFileErrorsByFilename(t *testing.T) {
	mock := &mockReader{
		GetFileErrorsByFilenameFunc: func(ctx context.Context, filename string, limit int, page int) ([]dto.FileErrorInfo, error) {
			if filename == "fail" {
				return nil, errors.New("some error")
			}
			return []dto.FileErrorInfo{
				{Filename: filename, Line: nil, ErrorMsg: "parse error"},
			}, nil
		},
	}

	svc := service.NewReaderService(mock)

	errorsData, err := svc.GetFileErrorsByFilename(context.Background(), "file1.tsv", 10, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(errorsData) != 1 || errorsData[0].Filename != "file1.tsv" {
		t.Fatalf("unexpected data returned: %+v", errorsData)
	}

	_, err = svc.GetFileErrorsByFilename(context.Background(), "fail", 10, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
