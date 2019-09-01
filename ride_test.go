package wappa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestRide(t *testing.T) {
	testCases := []testTable{
		{
			"Read()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&RideService{req}).Read(ctx, Filter{"id": []string{"1"}})
				return
			},
			context.Background(),
			http.MethodGet,
			rideEndpoint.Action(status).Query(url.Values{"rideId": []string{"1"}}),
			nil,
			&RideResult{
				Result: Result{Success: true},
				ID: 1,
			},
		},
		{
			"Create()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&RideService{req}).Create(ctx, &Ride{LatOrigin: 3.14})
				return
			},
			context.Background(),
			http.MethodPost,
			rideEndpoint,
			&Ride{LatOrigin: 3.14},
			&RideResult{
				Result: Result{Success: true},
				ID: 2,
			},
		},
		{
			"CancellationReason()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&RideService{req}).CancellationReason(ctx)
				return
			},
			context.Background(),
			http.MethodGet,
			indexEndpoint.Action(cancellationReason),
			nil,
			&CancellationReasonResult{
				Reasons: []Base{
					Base{1, "Test"},
				},
			},
		},
		{
			"Cancel()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&RideService{req}).Cancel(ctx, 1, 2)
				return
			},
			context.Background(),
			http.MethodPost,
			rideEndpoint.Action(cancel),
			&rideCancel{1, 2},
			&Result{
				Success: true,
			},
		},
		{
			"Rate()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&RideService{req}).Rate(ctx, 1, 5)
				return
			},
			context.Background(),
			http.MethodPost,
			rideEndpoint.Action(rate),
			&rideRate{1, 5},
			&Result{
				Success: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestRideError(t *testing.T) {
	testCases := []testTableError{
		{
			"Read()",
			func(req requester) error {
				_, err := (&RideService{req}).Read(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Create()",
			func(req requester) error {
				_, err := (&RideService{req}).Create(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"CancellationReason()",
			func(req requester) error {
				_, err := (&RideService{req}).CancellationReason(context.Background())
				return err
			},
			errors.New("Error"),
		},
		{
			"Cancel()",
			func(req requester) error {
				_, err := (&RideService{req}).Cancel(context.Background(), 1, 1)
				return err
			},
			errors.New("Error"),
		},
		{
			"Rate()",
			func(req requester) error {
				_, err := (&RideService{req}).Rate(context.Background(), 1, 1)
				return err
			},
			errors.New("Error"),
		},

	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
