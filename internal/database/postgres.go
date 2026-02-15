package database

import (
	"database/sql"
	"file-parser/internal/config"
	"fmt"
)

func NewPostgreSQLDatabase(config *config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	queries := []string{
		`
		CREATE TABLE files (
			id BIGSERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		`,
		`
		CREATE TABLE devices (
			id BIGSERIAL PRIMARY KEY,
			unit_guid VARCHAR(36) NOT NULL UNIQUE,
			inv_id VARCHAR(100),
			file_id BIGINT REFERENCES files(id) ON DELETE CASCADE
		);
		`,
		`
		CREATE TABLE device_data (
			id BIGSERIAL PRIMARY KEY,
			n INT,
			mqtt VARCHAR(255),
			msg_id VARCHAR(128),
			text VARCHAR(255),
			context VARCHAR(255),
			class VARCHAR(20),
			level INT,
			area VARCHAR(50),
			addr TEXT,
			block VARCHAR(50),
			type VARCHAR(50),
			bit INT,
			invert_bit INT,
			device_id BIGINT REFERENCES devices(id) ON DELETE CASCADE
		);
		`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return db, nil
}
