package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"qwe/internal/trip/config"
	"qwe/internal/trip/models"
	"qwe/internal/trip/ports"
	common "qwe/pkg/models"
	"strconv"
)

var _ ports.IService = &Service{}

type Service struct {
	repo      ports.IRepository
	cfg       *config.Config
	log       *slog.Logger
	produceCh chan models.Message
	consumeCh chan common.KafkaMessage
}

// обработать входящее сообщение
func (s *Service) HandleConsumerMessage(ctx context.Context, m common.KafkaMessage) error {

	var newtrip models.Trip
	var updatetrip common.UpdateTripStatus
	var statustrip common.TripStatus

	switch m.Type {
	case common.TRIP_MESSAGE:
		json.Unmarshal(m.Body, &newtrip)
		statustrip = newtrip.Status
	case common.UPDATE_STATUS_TRIP_MESSAGE:
		json.Unmarshal(m.Body, &updatetrip)
		statustrip = updatetrip.Status
	default:
		fmt.Printf("I don't know about type %T!\n", m.Type)
	}

	s.log.Info("Service.HandleConsumerMessage", slog.String("status", statustrip.String()))

	switch statustrip {
	case common.NEW:
		trip, err := s.repo.CreateTrip(ctx, &newtrip)

		if err != nil {
			s.log.Error("Create trip error.", err)
			return err
		}

		s.log.Info("Trip created", slog.Any("trip", trip))

		s.produceCh <- models.Message{Model: newtrip, Topic: s.cfg.KAFKA.PRODUCER.DRIVERTOPIC}

		updateTripStatus := common.UpdateTripStatus{TripId: trip.Id, ClientId: newtrip.ClientId, DriverId: trip.DriverId, Status: common.CREATED}
		s.produceCh <- models.Message{Model: updateTripStatus, Topic: s.cfg.KAFKA.PRODUCER.CLIENTTOPIC}

		return nil
	case common.DRIVER_FOUND, common.ON_POSITION, common.STARTED, common.ENDED:

		trip, err := s.repo.UpdateTrip(ctx, &updatetrip)

		if err != nil {
			s.log.Error("Update trip error.", err)
			return err
		}

		s.log.Info("Trip updated", slog.Any("trip", trip))

		s.produceCh <- models.Message{Model: &trip, Topic: s.cfg.KAFKA.PRODUCER.CLIENTTOPIC}
		return nil
	default:
		s.log.Error("Trip handle error", errors.New("Invalid status "+strconv.Itoa(int(statustrip))))
		return nil
	}

}
func New(cfg *config.Config, log *slog.Logger, r ports.IRepository, produceCh chan models.Message, consumeCh chan common.KafkaMessage) *Service {
	s := &Service{
		repo:      r,
		cfg:       cfg,
		log:       log,
		produceCh: produceCh,
		consumeCh: consumeCh,
	}

	go func(s *Service) {
		for {
			m := <-consumeCh
			s.HandleConsumerMessage(context.Background(), m)
		}
	}(s)

	return s
}
