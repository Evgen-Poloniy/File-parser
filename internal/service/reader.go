package service

import "file-parser/internal/repository"

type RepositoryService struct {
	repository repository.RepositoryAdapter
}

func newRepositoryService(repository repository.RepositoryAdapter) *RepositoryService {
	return &RepositoryService{
		repository: repository,
	}
}
