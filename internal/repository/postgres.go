package repository

import (
	"context"
	"database/sql"
	"file-parser/internal/dto"
	"file-parser/internal/entity"
	"strings"
)

type PostgresRepository struct {
	db *sql.DB
}

func newPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Get data about device by unit_guid
func (p *PostgresRepository) GetDataByUnitGUID(ctx context.Context, unitGUID string, limit int, page int) ([]dto.Response, error) {
	offset := (page - 1) * limit

	query := `
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
	`

	rows, err := p.db.QueryContext(ctx, query, unitGUID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dto.Response

	for rows.Next() {
		var r dto.Response

		err := rows.Scan(
			&r.N,
			&r.UnitGUID,
			&r.InvID,
			&r.MQTT,
			&r.MsgID,
			&r.Text,
			&r.Context,
			&r.Class,
			&r.Level,
			&r.Area,
			&r.Addr,
			&r.Block,
			&r.Type,
			&r.Bit,
			&r.InvertBit,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Get information about parse errors
func (p *PostgresRepository) GetFileErrorsByFilename(ctx context.Context, filename string, limit int, page int) ([]dto.FileErrorInfo, error) {
	offset := (page - 1) * limit

	rows, err := p.db.QueryContext(ctx, `
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
    `, filename, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var errors []dto.FileErrorInfo

	for rows.Next() {
		var fe dto.FileErrorInfo
		err := rows.Scan(&fe.Filename, &fe.Line, &fe.LineData, &fe.ErrorMsg)
		if err != nil {
			return nil, err
		}
		errors = append(errors, fe)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return errors, nil
}

// Check record in the database. If returning id equals 0, then file must go to parsing
func (p *PostgresRepository) CheckRecordAboutFile(ctx context.Context, filename string) (int64, error) {
	var id int64
	query := "SELECT id FROM files WHERE filename = $1"
	if err := p.db.QueryRowContext(ctx, query, filename).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

// Check record in the database. If returning id equals 0, then file must go to parsing
func (p *PostgresRepository) CheckRecordAboutError(ctx context.Context, filename string) (int64, error) {
	var id int64
	query := "SELECT id FROM not_parsed_files WHERE filename = $1"
	if err := p.db.QueryRowContext(ctx, query, filename).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return id, nil
}

// Insert information about file and parsed data into the database by one transaction
func (p *PostgresRepository) InsertParsedData(
	ctx context.Context,
	filename string,
	devices []entity.Device,
	deviceData []entity.DeviceData,
) error {

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// files
	var fileID int64
	if err := tx.QueryRowContext(
		ctx,
		`INSERT INTO files (filename) VALUES ($1) RETURNING id`,
		filename,
	).Scan(&fileID); err != nil {
		return err
	}

	deviceMap := make(map[string]int64)

	// devices
	deviceStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO devices (unit_guid, inv_id)
		VALUES ($1, $2)
		ON CONFLICT (unit_guid)
		DO UPDATE SET unit_guid = EXCLUDED.unit_guid
		RETURNING id
	`)
	if err != nil {
		return err
	}
	defer deviceStmt.Close()

	for _, device := range devices {
		var id int64

		if err := deviceStmt.QueryRowContext(
			ctx,
			device.UnitGUID,
			device.InvID,
		).Scan(&id); err != nil {
			return err
		}

		deviceMap[device.UnitGUID] = id
	}

	// device_data
	dataStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO device_data (
			n, mqtt, msg_id, text, context, class, level,
			area, addr, block, type, bit, invert_bit, device_id
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`)
	if err != nil {
		return err
	}
	defer dataStmt.Close()

	for _, dd := range deviceData {
		deviceID, ok := deviceMap[dd.UnitGUID]
		if !ok {
			return sql.ErrNoRows
		}

		if _, err := dataStmt.ExecContext(
			ctx,
			dd.N,
			dd.MQTT,
			dd.MsgID,
			dd.Text,
			dd.Context,
			dd.Class,
			dd.Level,
			dd.Area,
			dd.Addr,
			dd.Block,
			dd.Type,
			dd.Bit,
			dd.InvertBit,
			deviceID,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Insert information about parse errors into the database
func (p *PostgresRepository) InsertErrors(ctx context.Context, filename string, parseErrors *entity.ParseErrors) error {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fileID int64
	if err := tx.QueryRowContext(
		ctx,
		`INSERT INTO not_parsed_files(filename) VALUES($1) RETURNING id`,
		filename,
	).Scan(&fileID); err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO file_errors(file_id, line, line_data, error_msg)
		VALUES($1, $2, $3, $4)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range parseErrors.Errors {
		var lineData *string
		if len(e.LineData) > 0 {
			s := strings.Join(e.LineData, "\t")
			lineData = &s
		}

		var lineNum interface{}
		if e.LineNum != nil {
			lineNum = *e.LineNum
		} else {
			lineNum = nil
		}

		_, err = stmt.ExecContext(ctx, fileID, lineNum, lineData, e.Err.Error())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
