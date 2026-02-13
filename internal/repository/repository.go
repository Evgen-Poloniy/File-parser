package repository

import "database/sql"

type RepositoryAdapter interface {
}

type Repository struct {
	RepositoryAdapter
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		RepositoryAdapter: newPostgreSQLRepository(db),
	}
}
