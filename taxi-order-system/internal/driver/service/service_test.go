package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"qwe/internal/driver/models"
	common "qwe/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	produceCh = make(chan common.UpdateTripStatus)
	consumeCh = make(chan common.Trip)
	repo = NewTestRepo()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	st = Service{repo: &repo, log: log, produceCh: produceCh, consumeCh: consumeCh}
}

func TestCaptureTrip(t *testing.T) {

	go func() {
		for {
			m := <-produceCh
			fmt.Println(m)
		}
	}()

	captureTrip := models.CaptureTrip{DriverId: "driver1", TripId: 1}
	st.CaptureTrip(context.Background(), captureTrip)

	var expected models.Driver
	for _, d := range repo.Drivers {
		if d.DriverId == captureTrip.DriverId {
			expected = d
		}
	}

	assert.NotEmpty(t, expected)
	assert.Equal(t, "driver1", expected.DriverId)
	assert.Equal(t, int64(1), *expected.TripId)
}

func TestEndTrip(t *testing.T) {
	go func() {
		for {
			m := <-produceCh
			fmt.Println(m)
		}
	}()
	repo = NewTestRepo()

	captureTrip := models.CaptureTrip{DriverId: "driver2", TripId: 2}

	st.UpdateTripStatus(context.Background(), captureTrip, common.ENDED)

	var expected models.Driver
	for _, d := range repo.Drivers {
		if d.DriverId == captureTrip.DriverId {
			expected = d
		}
	}

	assert.NotEmpty(t, expected)
	assert.Empty(t, expected.TripId)
}
