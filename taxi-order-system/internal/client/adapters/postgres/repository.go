package repository

import (
	"context"
	"qwe/internal/client/config"
	"qwe/internal/client/models"
	"qwe/internal/client/ports"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var _ ports.IRepository = (*Repo)(nil)

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) Create(ctx context.Context, iemai string) (*models.Client, error) {
	data, err := r.db.ExecContext(ctx, "INSERT INTO Clients (iemai) VALUES (?)", iemai)

	if err != nil {
		return nil, err
	}

	id, err := data.LastInsertId()

	if err != nil {
		return nil, err
	}

	var client = &models.Client{
		Id:   id,
		Imei: iemai,
	}

	return client, nil
}

func (r *Repo) Get(ctx context.Context, id string) (*models.Client, error) {
	data := r.db.QueryRowContext(ctx, "SELECT id, iemai FROM Clients WHERE id=?", id)

	if data == nil {
		return nil, ErrNotFound
	}

	var client models.Client

	data.Scan(&client.Id, &client.Imei)

	return &client, nil
}

func New(cfg *config.DB) (ports.IRepository, error) {

	db, err := openConnection(cfg)
	if err != nil {
		return nil, err
	}

	migrations, err := readMigrations(cfg)
	if err != nil {
		return nil, err
	}

	if err = migrate(db, migrations); err != nil {
		return nil, err
	}

	return &Repo{
		db: db,
	}, nil
}
