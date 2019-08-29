package wappa

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

const employeeEndpoint endpoint = `employee`
// indexEndpoint is the index resource used to get general information
// about the API, as list of employees and list of cancellation reasons.
const indexEndpoint endpoint = `index`

var employeeLastRidesFields = map[string]string{
	"ride": "RideId",
	"employee" : "EmployeeId",
	"startDate": "InitialDate",
	"endDate": "FinalDate",
	"externalID": "ExternalID",
}

var employeeFields = map[string]string{
	"id": "EmployeeID",
	"name": "Name",
	"email": "Email",
}

// Employe is the  struct representing the employee
// entity in the API.
type Employee struct {
	ID int `json:"employeeId"`
	Email string `json:"email"`
	DDD string `json:"ddd"`
	Phone string `json:"phone"`
	Registration string `json:"registration"`
}

// WebhookResult is the API response payload.
type EmployeeResult struct {
	Employees []*Employee `json:"employees"`
}

// EmployeeStatusResult represents the current status of the employee.
type EmployeeStatusResult struct {
	// The ride id if the current ride.
	RideID int `json:"rideId"`
	// Employee ride status.
	// Possible values: Free, AwaitingPickup, OnRide, RideCompleted or OnAuction.
	Status string `json:"status"`
}

type RideHistory struct {
}

// EmployeeLastRidesResult represents the up to 100 last employee last rides.
type EmployeeLastRidesResult struct {
	Result

	History []*RideHistory
}

// EmployeeService is responsible for handling
// the requests to the webhook resource.
type EmployeeService struct {
	client requester
}

// Status returns the current status of employees in the API.
func (es *EmployeeService) Status(ctx context.Context, id int) (*EmployeeStatusResult, error) {
	e := &EmployeeStatusResult{}
	vals := url.Values{}
	vals.Set("employeeId", strconv.Itoa(id))

	if err := es.client.Request(ctx, http.MethodGet, employeeEndpoint.Action(status).Query(vals), nil, e); err != nil {
		return nil, err
	}

	return e, nil
}

// LastRides returns up to 100 last rides of a given employee.
func (es *EmployeeService) LastRides(ctx context.Context, f Filter) (*EmployeeLastRidesResult, error) {
	res := &EmployeeLastRidesResult{}

	if err := es.client.Request(ctx, http.MethodGet, employeeEndpoint.Action(lastRides).Query(f.Values(employeeLastRidesFields)), nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Read returns the list of employees associated with the credentials.
func (es *EmployeeService) Read(ctx context.Context, f Filter) (*EmployeeResult, error) {
	res := &EmployeeResult{}

	if err := es.client.Request(ctx, http.MethodGet, indexEndpoint.Action(employee).Query(f.Values(employeeFields)), nil, res); err != nil {
		return nil, err
	}

	return res, nil
}
