package repository

import (
	"context"
	"qwe/internal/driver/config"
	"qwe/internal/driver/models"
	"qwe/internal/driver/ports"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var _ ports.IRepository = (*Repo)(nil)

type Repo struct {
	db *sqlx.DB
}

// Получить список свободных водителей
func (r *Repo) GetAvailableDrivers(ctx context.Context) ([]models.Driver, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, driverid FROM taxi.drivers WHERE tripid IS null")

	if err != nil {
		return nil, err
	}

	drivers := make([]models.Driver, 0)

	for rows.Next() {
		var driver models.Driver
		rows.Scan(&driver.Id, &driver.DriverId)
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

// Назначить поездку водителю
func (r *Repo) SetDriverToFree(ctx context.Context, model models.CaptureTrip) error {
	_, err := r.db.QueryContext(ctx, "UPDATE taxi.drivers SET tripid=null WHERE driverid=$1", model.DriverId)

	return err
}

// Открепить поездку от водителя
func (r *Repo) AssignTripToDriver(ctx context.Context, model models.CaptureTrip) error {
	_, err := r.db.QueryContext(ctx, "UPDATE taxi.drivers SET tripid=$1 WHERE driverid=$2", model.TripId, model.DriverId)

	return err
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
