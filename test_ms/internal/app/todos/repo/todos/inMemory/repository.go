package inmemory

import (
	"context"
	"errors"
	"test_ms/internal/app/todos/repo"
	"test_ms/internal/models"
)

var _ repo.IRepository = (*InMemoryRepository)(nil)

type InMemoryRepository struct {
	db []models.Todo
}

func NewRepo(db []models.Todo) *InMemoryRepository {
	return &InMemoryRepository{
		db: db,
	}
}

func (r *InMemoryRepository) GetAll(ctx context.Context) ([]models.Todo, error) {

	return r.db, nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id int) (*models.Todo, error) {

	for _, v := range r.db {
		if v.Id == id {
			return &v, nil
		}
	}

	return nil, errors.New("todo not found")
}

func (r *InMemoryRepository) Add(ctx context.Context, t *models.Todo) error {

	id := 1

	if len(r.db) > 0 {
		id = r.db[len(r.db)-1].Id + 1
	}

	t.Id = id

	r.db = append(r.db, *t)
	return nil
}

func (r *InMemoryRepository) Update(ctx context.Context, t *models.Todo) error {

	for _, v := range r.db {
		if v.Id == t.Id {
			v.Task = t.Task
			v.Complete = t.Complete
			return nil
		}
	}
	return errors.New("todo not found")
}
