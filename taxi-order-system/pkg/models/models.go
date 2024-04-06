package models

type TripStatus int

const (
	NEW TripStatus = iota
	CREATED
	DRIVER_FOUND
	ON_POSITION
	STARTED
	ENDED
)

func (e TripStatus) String() string {
	switch e {
	case NEW:
		return "NEW"
	case CREATED:
		return "CREATED"
	case DRIVER_FOUND:
		return "DRIVER_FOUND"
	case ON_POSITION:
		return "ON_POSITION"
	case STARTED:
		return "STARTED"
	case ENDED:
		return "ENDED"
	default:
		return ""
	}
}

const (
	TRIP_MESSAGE               = "Trip"
	UPDATE_STATUS_TRIP_MESSAGE = "UpdateTripStatus"
)

type KafkaMessage struct {
	Type string
	Body []byte
}

type Trip struct {
	Id       int64
	ClientId int64
	From     string
	To       string
}

type UpdateTripStatus struct {
	TripId   int64      `sql:"id"`
	DriverId string     `sql:"driverid"`
	ClientId string     `sql:"clientid"`
	Status   TripStatus `sql:"trip_status"`
}
