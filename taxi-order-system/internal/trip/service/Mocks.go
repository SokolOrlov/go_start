package service

import (
	"context"
	"errors"
	"qwe/internal/trip/models"
	"qwe/internal/trip/ports"
	common "qwe/pkg/models"
)

var st Service
var repo TestRepo
var produceCh chan models.Message
var consumeCh chan common.KafkaMessage

var _ ports.IRepository = (*TestRepo)(nil)

type TestRepo struct {
	Trips []models.Trip
	id    int64
}

// CreateTrip implements ports.IRepository.
func (t *TestRepo) CreateTrip(ctx context.Context, m *models.Trip) (*models.Trip, error) {
	t.id += 1
	newTrip := models.Trip{
		Id:       t.id,
		ClientId: m.ClientId,
		DriverId: m.DriverId,
		From:     m.From,
		To:       m.To,
		Status:   common.CREATED,
	}

	t.Trips = append(t.Trips, newTrip)

	return &newTrip, nil
}

// UpdateTrip implements ports.IRepository.
func (t *TestRepo) UpdateTrip(ctx context.Context, m *common.UpdateTripStatus) (*common.UpdateTripStatus, error) {
	for i, n := range t.Trips {
		if n.Id == m.TripId {
			t.Trips[i].DriverId = m.DriverId
			t.Trips[i].Status = m.Status

			return &common.UpdateTripStatus{TripId: n.Id, DriverId: n.DriverId, ClientId: n.ClientId, Status: m.Status}, nil
		}
	}

	return nil, errors.New("")
}

func NewTestRepo() TestRepo {
	t := TestRepo{}
	t.Trips = make([]models.Trip, 0, 10)
	return t
}
