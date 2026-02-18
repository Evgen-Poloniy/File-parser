package repository_test

import (
	"context"
	"errors"
	"file-parser/internal/entity"
	"file-parser/internal/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetDataByUnitGUID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresRepository(db)

	query := regexp.QuoteMeta(`
	SELECT
		dd.n,
		d.unit_guid,
		d.inv_id,
		dd.mqtt,
		dd.msg_id,
		dd.text,
		dd.context,
		dd.class,
		dd.level,
		dd.area,
		dd.addr,
		dd.block,
		dd.type,
		dd.bit,
		dd.invert_bit
	FROM device_data dd
	JOIN devices d ON dd.device_id = d.id
	WHERE d.unit_guid = $1
	ORDER BY dd.id ASC
	LIMIT $2 OFFSET $3
	`)

	rows := sqlmock.NewRows([]string{
		"n", "unit_guid", "inv_id", "mqtt", "msg_id", "text",
		"context", "class", "level", "area", "addr", "block",
		"type", "bit", "invert_bit",
	}).AddRow(
		1, "guid1", "inv1", "mqtt", "msg1", "text",
		"context", "class", 100, "area", "addr", "block",
		"type", 1, 0,
	)

	mock.ExpectQuery(query).
		WithArgs("guid1", 10, 0).
		WillReturnRows(rows)

	result, err := repo.GetDataByUnitGUID(context.Background(), "guid1", 10, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].UnitGUID != "guid1" {
		t.Fatalf("expected guid1, got %s", result[0].UnitGUID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetFileErrorsByFilename(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresRepository(db)

	query := regexp.QuoteMeta(`
        SELECT
            f.filename,
            e.line,
            e.line_data,
            e.error_msg
        FROM file_errors e
        JOIN not_parsed_files f ON e.file_id = f.id
        WHERE f.filename = $1
        ORDER BY f.id ASC
		LIMIT $2 OFFSET $3
    `)

	rows := sqlmock.NewRows([]string{
		"filename", "line", "line_data", "error_msg",
	}).AddRow("file1.tsv", 5, "data", "parse error")

	mock.ExpectQuery(query).
		WithArgs("file1.tsv", 10, 0).
		WillReturnRows(rows)

	result, err := repo.GetFileErrorsByFilename(context.Background(), "file1.tsv", 10, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].Filename != "file1.tsv" {
		t.Fatalf("expected file1.tsv, got %s", result[0].Filename)
	}
}

func TestCheckRecordAboutFile(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresRepository(db)

	query := regexp.QuoteMeta("SELECT id FROM files WHERE filename = $1")

	mock.ExpectQuery(query).
		WithArgs("file1.tsv").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	id, err := repo.CheckRecordAboutFile(context.Background(), "file1.tsv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != 42 {
		t.Fatalf("expected 42, got %d", id)
	}
}

func TestCheckRecordAboutError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresRepository(db)

	query := regexp.QuoteMeta("SELECT id FROM not_parsed_files WHERE filename = $1")

	mock.ExpectQuery(query).
		WithArgs("file1.tsv").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(100))

	id, err := repo.CheckRecordAboutError(context.Background(), "file1.tsv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if id != 100 {
		t.Fatalf("expected 100, got %d", id)
	}
}

func TestInsertParsedData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := repository.NewPostgresRepository(db)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO files (filename) VALUES ($1) RETURNING id`,
	)).
		WithArgs("file1.tsv").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectPrepare("INSERT INTO devices").
		ExpectQuery().
		WithArgs("guid1", "inv1").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	mock.ExpectPrepare("INSERT INTO device_data").
		ExpectExec().
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			int64(10),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	n := 1
	mqtt := "mqtt"
	msgID := "msg"
	text := "text"
	contextVal := "ctx"
	class := "class"
	level := 2
	area := "area"
	addr := "addr"
	block := "block"
	typeVal := "type"
	bit := 1
	invert := 0

	devices := []entity.Device{
		{
			UnitGUID: "guid1",
			InvID:    "inv1",
		},
	}

	data := []entity.DeviceData{
		{
			UnitGUID:  "guid1",
			N:         &n,
			MQTT:      &mqtt,
			MsgID:     &msgID,
			Text:      &text,
			Context:   &contextVal,
			Class:     &class,
			Level:     &level,
			Area:      &area,
			Addr:      &addr,
			Block:     &block,
			Type:      &typeVal,
			Bit:       &bit,
			InvertBit: &invert,
		},
	}

	err = repo.InsertParsedData(context.Background(), "file1.tsv", devices, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestInsertErrors(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := repository.NewPostgresRepository(db)

	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO not_parsed_files(filename) VALUES($1) RETURNING id`,
	)).
		WithArgs("file1.tsv").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectPrepare("INSERT INTO file_errors").
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	lineNum := 1
	parseErrors := &entity.ParseErrors{
		Errors: []entity.ParseError{
			{
				LineNum:  &lineNum,
				LineData: []string{"data1"},
				Err:      errors.New("parse error"),
			},
		},
	}

	err := repo.InsertErrors(context.Background(), "file1.tsv", parseErrors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
