package service

import (
	"context"
	"errors"
	"qwe/internal/driver/models"
	"qwe/internal/driver/ports"
	common "qwe/pkg/models"
)

var st Service
var repo TestRepo
var consumeCh chan common.Trip
var produceCh chan common.UpdateTripStatus

var _ ports.IRepository = (*TestRepo)(nil)

type TestRepo struct {
	Drivers []models.Driver
}

// AssignTripToDriver implements ports.IRepository.
func (t *TestRepo) AssignTripToDriver(ctx context.Context, m models.CaptureTrip) error {

	for _, d := range t.Drivers {
		if d.TripId == &m.TripId {
			return errors.New("trip already captured")
		}
	}

	for i, d := range t.Drivers {
		if d.DriverId == m.DriverId {
			t.Drivers[i].TripId = &m.TripId
			return nil
		}
	}

	return errors.New("no available drivers")
}

// GetAvailableDrivers implements ports.IRepository.
func (t *TestRepo) GetAvailableDrivers(context.Context) ([]models.Driver, error) {
	drivers := make([]models.Driver, 0, 5)

	for _, d := range t.Drivers {
		if d.TripId == nil {
			drivers = append(drivers, d)
		}
	}

	return drivers, nil
}

// SetDriverToFree implements ports.IRepository.
func (t *TestRepo) SetDriverToFree(ctx context.Context, m models.CaptureTrip) error {
	for i, d := range t.Drivers {
		if d.DriverId == m.DriverId {
			t.Drivers[i].TripId = nil
		}
	}
	return nil
}

func NewTestRepo() TestRepo {
	t := TestRepo{}
	t.Drivers = make([]models.Driver, 0, 10)

	t2 := int64(2)

	t.Drivers = append(t.Drivers, models.Driver{Id: 1, DriverId: "driver1"})
	t.Drivers = append(t.Drivers, models.Driver{Id: 2, DriverId: "driver2", TripId: &t2})
	t.Drivers = append(t.Drivers, models.Driver{Id: 3, DriverId: "driver3"})

	return t
}
