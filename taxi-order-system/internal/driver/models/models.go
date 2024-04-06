package models

type Driver struct {
	Id       int64  `sql:"id"`
	DriverId string `sql:"driverid"`
	TripId   int64  `sql:"tripid"`
}

type CaptureTrip struct {
	DriverId string
	TripId   int64
}
