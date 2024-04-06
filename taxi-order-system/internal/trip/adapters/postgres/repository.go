package repository

import (
	"context"
	"qwe/internal/trip/config"
	"qwe/internal/trip/models"
	"qwe/internal/trip/ports"
	common "qwe/pkg/models"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var _ ports.IRepository = (*Repo)(nil)

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) CreateTrip(ctx context.Context, m *models.Trip) (*models.Trip, error) {
	_, err := r.db.ExecContext(ctx, "INSERT INTO taxi.Trips (clientid, begin_route, end_route, trip_status) VALUES ($1,$2,$3,$4)",
		m.ClientId, m.From, m.To, common.CREATED)

	if err != nil {
		return nil, err
	}

	m.Status = common.CREATED
	return m, nil
}

func (r *Repo) UpdateTrip(ctx context.Context, m *common.UpdateTripStatus) (*common.UpdateTripStatus, error) {

	_, err := r.db.ExecContext(ctx, "UPDATE taxi.Trips SET driverid=$1, trip_status=$2 WHERE id=$3",
		m.DriverId, m.Status, m.TripId)

	if err != nil {
		return nil, err
	}

	data := r.db.QueryRowContext(ctx, "SELECT id, driverid, clientid, trip_status FROM taxi.Trips WHERE id=$1", m.TripId)

	var trip common.UpdateTripStatus

	data.Scan(&trip.TripId, &trip.DriverId, &trip.ClientId, &trip.Status)

	return &trip, nil
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
