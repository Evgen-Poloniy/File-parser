package repository

import "database/sql"

type PostgreSQLRepository struct {
	db *sql.DB
}

func newPostgreSQLRepository(db *sql.DB) *PostgreSQLRepository {
	return &PostgreSQLRepository{
		db: db,
	}
}
