package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"qwe/internal/trip/config"
	"qwe/internal/trip/models"
	common "qwe/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	produceCh = make(chan models.Message)
	consumeCh = make(chan common.KafkaMessage)
	repo = NewTestRepo()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	cfg := &config.Config{KAFKA: config.KAFKA{PRODUCER: config.PRODUCER{DRIVERTOPIC: "driver"}}}

	st = Service{repo: &repo, cfg: cfg, log: log, produceCh: produceCh, consumeCh: consumeCh}
}

func TestCreateTrip(t *testing.T) {

	repo.Trips = []models.Trip{}

	go func() {
		for {
			m := <-produceCh
			fmt.Println(m)
		}
	}()

	trip, _ := json.Marshal(models.Trip{ClientId: "client1", From: "from", To: "to", Status: common.NEW})
	kafkaMessage := common.KafkaMessage{Type: common.TRIP_MESSAGE, Body: trip}

	st.HandleConsumerMessage(context.Background(), kafkaMessage)

	assert.Equal(t, 1, len(repo.Trips))
	assert.Equal(t, repo.Trips[0].Status, common.CREATED)
}

func TestUpdateStatusTrip(t *testing.T) {

	repo.Trips = []models.Trip{}
	repo.Trips = append(repo.Trips, models.Trip{Id: 1, ClientId: "client1", From: "from", To: "to", Status: common.CREATED})

	go func() {
		for {
			m := <-produceCh
			fmt.Println(m)
		}
	}()

	trip, _ := json.Marshal(common.UpdateTripStatus{TripId: 1, DriverId: "driver1", ClientId: "client1", Status: common.DRIVER_FOUND})
	kafkaMessage := common.KafkaMessage{Type: common.UPDATE_STATUS_TRIP_MESSAGE, Body: trip}

	st.HandleConsumerMessage(context.Background(), kafkaMessage)

	assert.Equal(t, 1, len(repo.Trips))
	assert.Equal(t, repo.Trips[0].Status, common.DRIVER_FOUND)
	assert.Equal(t, repo.Trips[0].DriverId, "driver1")
}
