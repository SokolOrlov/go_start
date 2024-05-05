package service

import (
	"context"
	"log/slog"
	"qwe/internal/client/models"
	"qwe/internal/client/ports"
	common "qwe/pkg/models"
)

var _ ports.IService = &Service{}
var _ ports.BackgroundWorker = &Service{}

type Service struct {
	repo      ports.IRepository
	produceCh chan models.Trip
	consumeCh chan common.UpdateTripStatus
	log       *slog.Logger
}

func (s *Service) Run() {
	go func(s *Service) {
		for {
			m := <-s.consumeCh
			s.HandleConsumerMessage(context.Background(), m)
		}
	}(s)
}

func (s *Service) Stop() {
	close(s.consumeCh)
}

// Обработать входящее сообщение
func (s *Service) HandleConsumerMessage(ctx context.Context, m common.UpdateTripStatus) {

	s.log.Info("Service.HandleConsumerMessage", slog.Any("trip", m))

	err := s.notificateClient(ctx, m)

	if err != nil {
		s.log.Error("Error handle message", err)
	}

}

// Оповестить клиента об изменение статуса поездки
func (s *Service) notificateClient(ctx context.Context, m common.UpdateTripStatus) error {

	s.log.InfoContext(ctx, "------> notificate client", slog.String("clientid", m.ClientId), slog.String("status", m.Status.String()))

	return nil
}

// Создать поездку
func (s *Service) CreateTrip(ctx context.Context, trip models.Trip) {
	//TODO сохранять информацию о поездке клиента или запрашивать у сервиса Trip
	s.produceCh <- trip
}

func New(r ports.IRepository, log *slog.Logger, produceCh chan models.Trip, consumeCh chan common.UpdateTripStatus) *Service {
	s := &Service{
		repo:      r,
		produceCh: produceCh,
		consumeCh: consumeCh,
		log:       log,
	}

	return s
}
