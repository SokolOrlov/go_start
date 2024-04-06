package ports

import (
	"context"
	"qwe/internal/trip/models"
	common "qwe/pkg/models"
)

type IRepository interface {
	CreateTrip(context.Context, *models.Trip) (*models.Trip, error)
	UpdateTrip(context.Context, *common.UpdateTripStatus) (*common.UpdateTripStatus, error)
}
