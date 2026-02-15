package repository

import (
	"context"
	"database/sql"
	"file-parser/internal/dto"
	"file-parser/internal/entity"
)

type PostgreSQLRepository struct {
	db *sql.DB
}

func newPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		db: db,
	}
}

// Check record in the database. If returning id equals 0, then file must go to parsing
func (p *PostgreSQLRepository) CheckRecordAboutFile(ctx context.Context, filename string) (int64, error) {
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

// Insert information about file and parsed data into the database by one transaction
func (p *PostgreSQLRepository) InsertParsedData(ctx context.Context, file *entity.File, devices []entity.Device, deviceData []entity.DeviceData) error {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert information about parsed file
	var fileID int64
	if err := tx.QueryRowContext(
		ctx,
		`
		INSERT INTO files (filename)
		VALUES ($1)
		RETURNING id
		`,
		file.Filename,
	).Scan(&fileID); err != nil {
		return err
	}

	// Insert metadata about devices
	deviceMap := make(map[string]int64)
	for _, device := range devices {
		var id int64
		if err := tx.QueryRowContext(
			ctx,
			`
			INSERT INTO devices (unit_guid, inv_id, file_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id
			`,
			device.UnitGUID, device.InvID, fileID,
		).Scan(&id); err != nil {
			return err
		}
		deviceMap[device.UnitGUID] = id
	}

	// Insert data about devices
	for _, dd := range deviceData {
		_, err := tx.ExecContext(
			ctx,
			`
			INSERT INTO device_data (
				n, mtqq, msg_id, text, context, class, level, area, addr, block, type, bit, invert_bit, device_id
			) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
			`,
			dd.N, dd.MQTT, dd.MsgID, dd.Text, dd.Context, dd.Class, dd.Level, dd.Area,
			dd.Addr, dd.Block, dd.Type, dd.Bit, dd.InvertBit, deviceMap[dd.UnitGUID],
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Get data about device by unit_guid
func (p *PostgreSQLRepository) GetData(unitGUID string, limit int, page int) ([]dto.Response, error) {
	return nil, nil
}
