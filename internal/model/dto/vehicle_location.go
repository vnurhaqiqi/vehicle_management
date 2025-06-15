package dto

import (
	"strconv"

	"github.com/vnurhaqiqi/vehicle_management/internal/model"
)

type VehicleLocation struct {
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"` // Unix timestamp in seconds
}

type ResolveVehicleLocationRequest struct {
	VehicleID string
	Start     string
	End       string
}

func (r ResolveVehicleLocationRequest) ToFilter() model.VehicleLocationFilter {
	filter := model.VehicleLocationFilter{}

	filter.VechileID = r.VehicleID

	if r.Start != "" {
		start, _ := strconv.Atoi(r.Start)
		filter.Start = int64(start)
	}

	if r.End != "" {
		end, _ := strconv.Atoi(r.End)
		filter.End = int64(end)
	}

	return filter
}

func NewVehicleLocationResponse(vechileLocation model.VehicleLocation) VehicleLocation {
	return VehicleLocation{
		VehicleID: vechileLocation.VehicleID,
		Latitude:  vechileLocation.Latitude,
		Longitude: vechileLocation.Longitude,
		Timestamp: vechileLocation.Timestamp,
	}
}

func NewVehicleLocationFromRequest(request VehicleLocation) model.VehicleLocation {
	return model.VehicleLocation{
		VehicleID: request.VehicleID,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
		Timestamp: request.Timestamp,
	}
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type VehicleLocationMessage struct {
	VehicleID string   `json:"vehicle_id"`
	Event     string   `json:"event"`
	Location  Location `json:"location"`
	Timestamp int64    `json:"timestamp"` // Unix timestamp in seconds
}

func NewMessageFromVehicleLocation(location model.VehicleLocation) VehicleLocationMessage {
	message := VehicleLocationMessage{
		VehicleID: location.VehicleID,
		Event:     "geofence_entry",
		Timestamp: location.Timestamp,
		Location: Location{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		},
	}

	return message
}
