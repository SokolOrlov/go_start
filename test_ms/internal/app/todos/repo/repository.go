package repo

import (
	"context"
	"test_ms/internal/models"
)

type IRepository interface {
	GetAll(context.Context) ([]models.Todo, error)
	Get(context.Context, int) (*models.Todo, error)
	Add(context.Context, *models.Todo) error
	Update(context.Context, *models.Todo) error
}
