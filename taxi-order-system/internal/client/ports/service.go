package ports

import (
	"context"
	"qwe/internal/client/models"
	common "qwe/pkg/models"
)

type IService interface {
	//Создать поездку
	CreateTrip(context.Context, models.Trip)
	//Обработаать входящее сообщение
	HandleConsumerMessage(context.Context, common.UpdateTripStatus)
}

type BackgroundWorker interface {
	Run()
	Stop()
}
