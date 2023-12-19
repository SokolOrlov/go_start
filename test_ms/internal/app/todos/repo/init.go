package repo

import "test_ms/internal/pkg/models"

func NewDB() models.Todos {
	res := make(models.Todos, 0, 10)
	return res
}
