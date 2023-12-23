package todoService

import (
	"context"
	"test_ms/internal/app/todos/repo"
	"test_ms/internal/models"
)

type Service struct {
	repo repo.IRepository
}

func NewService(r repo.IRepository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]models.Todo, error) {
	all, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return all, nil
}

func (s *Service) Get(ctx context.Context, id int) (*models.Todo, error) {
	t, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) Add(ctx context.Context, t *models.Todo) error {
	err := s.repo.Add(ctx, t)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(ctx context.Context, t *models.Todo) error {
	err := s.repo.Update(ctx, t)

	if err != nil {
		return err
	}

	return nil
}
