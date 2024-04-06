package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"qwe/internal/driver/models"
	"qwe/internal/driver/ports"
	common "qwe/pkg/models"
)

var _ ports.IService = &Service{}

type Service struct {
	repo      ports.IRepository
	produceCh chan common.UpdateTripStatus
	consumeCh chan common.Trip
	log       *slog.Logger
}

// Обработать входящее сообщение
func (s *Service) HandleConsumerMessage(ctx context.Context, trip common.Trip) error {

	s.log.Info("Service.HandleConsumerMessage", slog.Any("trip", trip))

	drivers, err := s.getAvailableDrivers(ctx)

	if err != nil {
		s.log.Error("Error handle message", err)
		return err
	}

	//TODO если все водители заняты, trip отослать обратно?

	err = s.sendOfferToTrip(ctx, drivers, trip)

	if err != nil {
		return err
	}

	return nil
}

// Получить список доступных водителей
func (s *Service) getAvailableDrivers(ctx context.Context) ([]models.Driver, error) {

	s.log.Info("Service.getAvailableDrivers")

	drivers, err := s.repo.GetAvailableDrivers(ctx)

	return drivers, err
}

// Отослать предложение поезди доступным водителям
func (s *Service) sendOfferToTrip(ctx context.Context, drivers []models.Driver, trip common.Trip) error {
	//TODO send to websocket
	s.log.InfoContext(ctx, "Service.sendOfferToTrip", slog.Any("drivers", drivers), slog.Any("trip", trip))

	return nil
}

// Захватить заявку
func (s *Service) CaptureTrip(ctx context.Context, model models.CaptureTrip) error {
	err := s.repo.AssignTripToDriver(ctx, model)

	if err != nil {
		s.log.Error("Error capture trip", err)
		return err
	}

	s.UpdateTripStatus(ctx, model, common.DRIVER_FOUND)

	return nil
}

// обновить статус заявки
func (s *Service) UpdateTripStatus(ctx context.Context, model models.CaptureTrip, st common.TripStatus) error {
	modelstr, _ := json.Marshal(model)
	s.log.Info("Service.UpdateTripStatus", slog.String("CaptureTrip", string(modelstr)), slog.String("status", st.String()))

	if st == common.ENDED {
		err := s.repo.SetDriverToFree(ctx, model)
		if err != nil {
			s.log.Error("Error free driver", err)
			return err
		}
	}

	updateTripStatus := common.UpdateTripStatus{TripId: model.TripId, Status: st, DriverId: model.DriverId}
	s.produceCh <- updateTripStatus
	return nil
}

func New(r ports.IRepository, log *slog.Logger, produceCh chan common.UpdateTripStatus, consumeCh chan common.Trip) *Service {

	s := &Service{
		repo:      r,
		produceCh: produceCh,
		consumeCh: consumeCh,
		log:       log,
	}

	go func(s *Service) {
		for {
			m := <-consumeCh
			err := s.HandleConsumerMessage(context.Background(), m)

			if err != nil {
				log.Error("Error service handle message", err)
			}
		}
	}(s)

	return s
}
