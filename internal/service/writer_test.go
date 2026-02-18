package service_test

import (
	"context"
	"errors"
	"file-parser/internal/entity"
	"file-parser/internal/service"
	"os"
	"path/filepath"
	"testing"
)

type mockWriter struct {
	CheckRecordAboutFileFunc  func(ctx context.Context, filename string) (int64, error)
	CheckRecordAboutErrorFunc func(ctx context.Context, filename string) (int64, error)
	InsertParsedDataFunc      func(ctx context.Context, filename string, devices []entity.Device, data []entity.DeviceData) error
	InsertErrorsFunc          func(ctx context.Context, filename string, parseErrors *entity.ParseErrors) error
}

func (m *mockWriter) CheckRecordAboutFile(ctx context.Context, filename string) (int64, error) {
	return m.CheckRecordAboutFileFunc(ctx, filename)
}
func (m *mockWriter) CheckRecordAboutError(ctx context.Context, filename string) (int64, error) {
	return m.CheckRecordAboutErrorFunc(ctx, filename)
}
func (m *mockWriter) InsertParsedData(ctx context.Context, filename string, devices []entity.Device, data []entity.DeviceData) error {
	return m.InsertParsedDataFunc(ctx, filename, devices, data)
}
func (m *mockWriter) InsertErrors(ctx context.Context, filename string, parseErrors *entity.ParseErrors) error {
	return m.InsertErrorsFunc(ctx, filename, parseErrors)
}

func TestParseFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.tsv")
	content := "n\tmqtt\tinv_id\tunit_guid\tmsg_id\ttext\tcontext\tclass\tlevel\tarea\taddr\tblock\ttype\tbit\tinvert_bit\n" +
		"1\ttopic\tinv1\tguid1\tmsg1\ttext1\tctx1\tclass1\t2\tarea1\taddr1\tblock1\ttype1\t1\t0\n"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	mock := &mockWriter{
		CheckRecordAboutFileFunc: func(ctx context.Context, filename string) (int64, error) {
			return 0, nil
		},
		CheckRecordAboutErrorFunc: func(ctx context.Context, filename string) (int64, error) {
			return 0, nil
		},
	}

	svc := service.NewWriterService(mock)

	devices, data, err := svc.ParseFile(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	if len(data) != 1 {
		t.Fatalf("expected 1 data row, got %d", len(data))
	}
	if devices[0].UnitGUID != "guid1" {
		t.Fatalf("unexpected UnitGUID: %s", devices[0].UnitGUID)
	}
}

func TestLoadParsedData(t *testing.T) {
	mock := &mockWriter{
		InsertParsedDataFunc: func(ctx context.Context, filename string, devices []entity.Device, data []entity.DeviceData) error {
			if filename == "fail" {
				return errors.New("insert error")
			}
			return nil
		},
	}

	svc := service.NewWriterService(mock)

	devices := []entity.Device{{UnitGUID: "guid1", InvID: "inv1"}}
	data := []entity.DeviceData{{UnitGUID: "guid1"}}

	if err := svc.LoadParsedData(context.Background(), "file1", devices, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := svc.LoadParsedData(context.Background(), "fail", devices, data)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadErrors(t *testing.T) {
	mock := &mockWriter{
		InsertErrorsFunc: func(ctx context.Context, filename string, parseErrors *entity.ParseErrors) error {
			if filename == "fail" {
				return errors.New("insert error")
			}
			return nil
		},
	}

	svc := service.NewWriterService(mock)

	parseErrs := &entity.ParseErrors{
		Errors: []entity.ParseError{
			{Err: errors.New("parse fail")},
		},
	}

	if err := svc.LoadErrors(context.Background(), "file1", parseErrs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := svc.LoadErrors(context.Background(), "fail", parseErrs)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWriteDataToRTF(t *testing.T) {
	svc := service.NewWriterService(nil)

	tmpFile := filepath.Join(t.TempDir(), "data.rtf")
	data := []entity.DeviceData{
		{UnitGUID: "guid1"},
	}

	if err := svc.WriteDataToRTF(tmpFile, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(tmpFile); err != nil {
		t.Fatalf("expected file to exist, got error: %v", err)
	}
}

func TestWriteErrorsToRTF(t *testing.T) {
	svc := service.NewWriterService(nil)

	tmpFile := filepath.Join(t.TempDir(), "errors.rtf")
	parseErrs := &entity.ParseErrors{
		Errors: []entity.ParseError{
			{Err: errors.New("parse fail")},
		},
	}

	if err := svc.WriteErrorsToRTF(tmpFile, parseErrs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(tmpFile); err != nil {
		t.Fatalf("expected file to exist, got error: %v", err)
	}
}
