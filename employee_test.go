package wappa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestEmployee(t *testing.T) {
	testCases := []testTable{
		{
			"Status()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&EmployeeService{req}).Status(ctx, 1)
				return
			},
			context.Background(),
			http.MethodGet,
			employeeEndpoint.Action(status).Query(url.Values{"employeeId": []string{"1"}}),
			nil,
			&EmployeeStatusResult{
				RideID: 1,
				Status: "Test",
			},
		},
		{
			"LastRides()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&EmployeeService{req}).LastRides(ctx, Filter{"ride": []string{"1"}})
				return
			},
			context.Background(),
			http.MethodGet,
			employeeEndpoint.Action(lastRides).Query(url.Values{"RideId": []string{"1"}}),
			nil,
			&EmployeeLastRidesResult{Result: Result{Success: true}},
		},
		{
			"Read()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&EmployeeService{req}).Read(ctx, Filter{"id": []string{"1"}})
				return
			},
			context.Background(),
			http.MethodGet,
			indexEndpoint.Action(employee).Query(url.Values{"EmployeeID": []string{"1"}}),
			nil,
			&EmployeeResult{
				Employees: []*Employee{
					&Employee{ID: 1},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestEmployeeError(t *testing.T) {
	testCases := []testTableError{
		{
			"Status()",
			func(req requester) error {
				_, err := (&EmployeeService{req}).Status(context.Background(), 0)
				return err
			},
			errors.New("Error"),
		},
		{
			"LastRide()",
			func(req requester) error {
				_, err := (&EmployeeService{req}).LastRides(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Read()",
			func(req requester) error {
				_, err := (&EmployeeService{req}).Read(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
