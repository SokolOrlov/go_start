package ports

import (
	"context"
	"qwe/internal/driver/models"

	common "qwe/pkg/models"
)

type IService interface {
	//Обработать входящее сообщение
	HandleConsumerMessage(context.Context, common.Trip) error
	//обновить статус заявки
	UpdateTripStatus(context.Context, models.CaptureTrip, common.TripStatus) error
	//Захватить заявку
	CaptureTrip(context.Context, models.CaptureTrip) error
}

type BackgroundWorker interface {
	Run()
	Stop()
}
