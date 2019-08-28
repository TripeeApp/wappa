package wappa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestDriver(t *testing.T) {
	testCases := []testTable{
		{
			"Nearby()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&DriverService{req}).Nearby(ctx, Filter{"lat": []string{"3.14"}})
				return
			},
			context.Background(),
			http.MethodGet,
			driverEndpoint.Action(nearby).Query(url.Values{"Latitude": []string{"3.14"}}),
			nil,
			&DriverResult{
				Result: Result{Success: true},
				Drivers: []*DriverLocation{
					&DriverLocation{Location: Location{Lat: 3.14}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestDriverError(t *testing.T) {
	testCases := []testTableError{
		{
			"Nearby()",
			func(req requester) error {
				_, err := (&DriverService{req}).Nearby(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
