package model

import "time"

type Courier struct {
	ID        		int64
	Name      		string
	Phone     		string
	Status    		CourierStatus
	TransportType 	CourierTransportType
	CreatedAt 		time.Time
	UpdatedAt 		time.Time
}

type CourierStatus string

type CourierTransportType string

const (
	CourierStatusAvailable 	CourierStatus = "available"
	CourierStatusBusy 		CourierStatus = "busy"
)

const (
	TransportTypeOnFoot  CourierTransportType = "on_foot"
	TransportTypeScooter CourierTransportType = "scooter"
	TransportTypeCar     CourierTransportType = "car"
)

func (c *Courier) ChangeStatus(status CourierStatus) {
	c.Status = status
}
