package parser_test

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"file-parser/internal/config"
	"file-parser/internal/entity"
	"file-parser/internal/parser"

	"github.com/sirupsen/logrus"
)

type mockWriter struct {
	ParseFileFunc        func(ctx context.Context, path string) ([]entity.Device, []entity.DeviceData, error)
	LoadParsedDataFunc   func(ctx context.Context, filename string, devices []entity.Device, data []entity.DeviceData) error
	LoadErrorsFunc       func(ctx context.Context, filename string, parseErrs *entity.ParseErrors) error
	WriteDataToRTFFunc   func(path string, data []entity.DeviceData) error
	WriteErrorsToRTFFunc func(path string, parseErrs *entity.ParseErrors) error
}

func (m *mockWriter) ParseFile(ctx context.Context, path string) ([]entity.Device, []entity.DeviceData, error) {
	return m.ParseFileFunc(ctx, path)
}

func (m *mockWriter) LoadParsedData(ctx context.Context, filename string, devices []entity.Device, data []entity.DeviceData) error {
	return m.LoadParsedDataFunc(ctx, filename, devices, data)
}

func (m *mockWriter) LoadErrors(ctx context.Context, filename string, parseErrs *entity.ParseErrors) error {
	return m.LoadErrorsFunc(ctx, filename, parseErrs)
}

func (m *mockWriter) WriteDataToRTF(path string, data []entity.DeviceData) error {
	return m.WriteDataToRTFFunc(path, data)
}

func (m *mockWriter) WriteErrorsToRTF(path string, parseErrs *entity.ParseErrors) error {
	return m.WriteErrorsToRTFFunc(path, parseErrs)
}

func TestParseFile(t *testing.T) {
	ctx := context.Background()
	jobs := make(chan string, 1)
	jobs <- "testfile.tsv"

	mockSvc := &mockWriter{
		ParseFileFunc: func(ctx context.Context, path string) ([]entity.Device, []entity.DeviceData, error) {
			return []entity.Device{{UnitGUID: "guid1", InvID: "inv1"}}, []entity.DeviceData{{UnitGUID: "guid1"}}, nil
		},
		LoadParsedDataFunc: func(ctx context.Context, filename string, devices []entity.Device, data []entity.DeviceData) error {
			if filename != "testfile.tsv" {
				t.Fatalf("unexpected filename: %s", filename)
			}
			return nil
		},
		WriteDataToRTFFunc: func(path string, data []entity.DeviceData) error {
			if !strings.Contains(path, "testfile") {
				t.Fatalf("unexpected RTF path: %s", path)
			}
			return nil
		},
		LoadErrorsFunc:       func(ctx context.Context, filename string, parseErrs *entity.ParseErrors) error { return nil },
		WriteErrorsToRTFFunc: func(path string, parseErrs *entity.ParseErrors) error { return nil },
	}

	logger := logrus.New()
	logger.Out = io.Discard
	parser := parser.NewFileParser(mockSvc, &config.ParserConfig{
		InputDir:  "/input",
		OutputDir: "/output",
		ErrorDir:  "/error",
	}, logger)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go parser.ParseFile(ctx, jobs)

	time.Sleep(50 * time.Millisecond)
}
