package wappa

import (
	"context"
	"net/http"
)

const driverEndpoint endpoint = `driver`

var driverFields = map[string]string{
	"lat": "Latitude",
	"long": "Longitude",
	"type": "TypeIds",
	"employee": "EmployeeId",
}

// DriverLocation is the  struct representing the driver location
// entity in the API.
type DriverLocation struct {
	Location

	Bearing int `json:"bearing,omitempty"`
	Type int `json:"typeId,omitempty"`
}

// DriverResult is the API response payload.
type DriverResult struct {
	Result

	Drivers []*DriverLocation `json:"listenners"`
}

// DriverService is responsible for handling
// the requests to the driver resource.
type DriverService struct {
	client requester
}

// Nearby returns the driver of a given type that are closest to the given coordinates.
func (ds *DriverService) Nearby(ctx context.Context, f Filter) (*DriverResult, error) {
	d := &DriverResult{}

	if err := ds.client.Request(ctx, http.MethodGet, driverEndpoint.Action(nearby).Query(f.Values(driverFields)), nil, d); err != nil {
		return nil, err
	}

	return d, nil
}
