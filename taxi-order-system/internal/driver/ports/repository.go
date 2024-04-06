package ports

import (
	"context"
	"qwe/internal/driver/models"
)

type IRepository interface {
	//Получить список свободных водителей
	GetAvailableDrivers(context.Context) ([]models.Driver, error)

	//Назначить поездку водителю
	AssignTripToDriver(context.Context, models.CaptureTrip) error

	//Открепить поездку от водителя
	SetDriverToFree(context.Context, models.CaptureTrip) error
}
