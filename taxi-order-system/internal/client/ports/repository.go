package ports

import (
	"context"
	"qwe/internal/client/models"
)

type IRepository interface {
	Get(context.Context, string) (*models.Client, error)
	Create(context.Context, string) (*models.Client, error)
}
