package service

import "file-parser/internal/repository"

type RepositoryAdapter interface {
}

type Service struct {
	RepositoryAdapter
}

func NewService(repository repository.RepositoryAdapter) *Service {
	return &Service{
		RepositoryAdapter: newRepositoryService(repository),
	}
}
