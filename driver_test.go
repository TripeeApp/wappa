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
			"NearBy()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&DriverService{req}).NearBy(ctx, Filter{"lat": []string{"3.14"}})
				return
			},
			context.Background(),
			http.MethodGet,
			driverEndpoint.Action(nearBy).Query(url.Values{"Latitude": []string{"3.14"}}),
			nil,
			&DriverResult{
				Result: Result{Success: true},
				Drivers: []*Driver{
					&Driver{Lat: 3.14},
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
			"NearBy()",
			func(req requester) error {
				_, err := (&DriverService{req}).NearBy(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
