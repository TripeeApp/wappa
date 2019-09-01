package wappa

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const employeeEndpoint endpoint = `employee`
// indexEndpoint is the index resource used to get general information
// about the API, as list of employees and list of cancellation reasons.
const indexEndpoint endpoint = `index`

var employeeLastRidesFields = map[string]string{
	"ride": "RideId",
	"employee" : "EmployeeId",
	"startedAt": "InitialDate",
	"endedAt": "FinalDate",
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

// Base represents a generic type for describing
// a resource with an ID an Description attributes.
type Base struct {
	ID int `json:"id"`
	Description string `json:"description"`
}

type HistoricalCategory struct {
	Base

	Type Base `json:"type"`
	SubCategory Base `json:"subcategory"`
}

type HistoricalDriver struct {
	Driver

	Category HistoricalCategory `json:"category"`
	Photo string `json:"photo"`
}

// RideHistory represents the historical information of a ride.
type RideHistory struct {
	ID int `json:"rideId"`
	CompanyID int `json:"companyId"`
	Passenger Passenger `json:"passenger"`
	Origin Address `json:"origin"`
	Destiny Address `json:"destiny"`
	Driver HistoricalDriver `json:"driver"`
	Info HistoricalRideInfo  `json:"rideInfo"`
}

// HistoricalRideInfo represents a historical information
// about a ride made by the employee
type HistoricalRideInfo struct {
	Status string `json:"status"`
	StartedAt *time.Time `json:"rideDate"`
	EndedAt *time.Time `json:"finishDate"`
	PaidAt *time.Time `json:"paymentDate"`
	MapURL string `json:"rideMapURL"`
	CancelledBy string `json:"cancelledBy"`
	CancelledReason string `json:"cancelledReason"`
	Value float64 `json:"rideValue"`
	OriginalValue float64 `json:"rideOriginalValue"`
	DIscount float64 `json:"rideDiscount"`
	ExternalID int `json:"externalId"`
	DurationInSeconds int `json:"durationInSeconds"`
	Distance int `json:"distance"`
}

// EmployeeLastRidesResult represents the up to 100 last employee last rides.
type EmployeeLastRidesResult struct {
	Result

	History []*RideHistory `json:"history"`
}

// EmployeeService is responsible for handling
// the requests to the webhook resource.
type EmployeeService service

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
