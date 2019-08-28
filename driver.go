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

// Driver is the  struct representing the driver
// entity in the API.
type Driver struct {
	Long float64 `json:"longitude"`
	Lat float64 `json:"latitude"`
	Bearing int `json:"bearing"`
	Type int `json:"typeId"`
}

// DriverResult is the API response payload.
type DriverResult struct {
	Result

	Drivers []*Driver `json:"listenners"`
}

// DriverService is responsible for handling
// the requests to the driver resource.
type DriverService struct {
	client requester
}

// NearBy returns the driver of a given type that are closest to the given coordinates.
func (ds *DriverService) NearBy(ctx context.Context, f Filter) (*DriverResult, error) {
	d := &DriverResult{}

	if err := ds.client.Request(ctx, http.MethodGet, driverEndpoint.Action(nearBy).Query(f.Values(driverFields)), nil, d); err != nil {
		return nil, err
	}

	return d, nil
}
