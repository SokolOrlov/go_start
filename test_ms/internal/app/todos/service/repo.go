package todoService

import (
	"context"
	"test_ms/internal/pkg/models"
)

type ITodoRepository interface {
	GetAll(context.Context) ([]models.Todo, error)
	Get(context.Context, int) (*models.Todo, error)
	Add(context.Context, *models.Todo) error
	Update(context.Context, *models.Todo) error
}
