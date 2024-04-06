package models

import status "qwe/pkg/models"

type Trip struct {
	Id       int64             `sql:"id" json:"Id"`
	ClientId string            `sql:"clientid" json:"ClientId"`
	DriverId string            `sql:"driverid" json:"DriverId"`
	From     string            `sql:"begin_route" json:"From"`
	To       string            `sql:"end_route" json:"To"`
	Status   status.TripStatus `sql:"trip_status" json:"Status"`
}

type Message struct {
	Model interface{}
	Topic string
}
