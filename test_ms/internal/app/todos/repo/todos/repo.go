package todos

import (
	"context"
	"errors"
	"test_ms/internal/pkg/models"
)

type TodoRepository struct {
	db models.Todos
}

func NewRepo(db models.Todos) *TodoRepository {
	return &TodoRepository{
		db: db,
	}
}

func (r *TodoRepository) GetAll(ctx context.Context) ([]models.Todo, error) {

	return r.db, nil
}

func (r *TodoRepository) Get(ctx context.Context, id int) (*models.Todo, error) {

	for _, v := range r.db {
		if v.Id == id {
			return &v, nil
		}
	}

	return nil, errors.New("todo not found")
}

func (r *TodoRepository) Add(ctx context.Context, t *models.Todo) error {
	r.db = append(r.db, *t)
	return nil
}

func (r *TodoRepository) Update(ctx context.Context, t *models.Todo) error {

	for _, v := range r.db {
		if v.Id == t.Id {
			v.Task = t.Task
			v.Complete = t.Complete
			return nil
		}
	}
	return errors.New("todo not found")
}
