package wappa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestCostCenter(t *testing.T) {
	testCases := []testTable{
		{
			"Read()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&CostCenterService{req}).Read(ctx, Filter{"code": "45"})
				return
			},
			context.Background(),
			http.MethodGet,
			costCenterEndpoint.Action(read).Query(url.Values{"codDescricao": []string{"45"}}),
			nil,
			&CostCenterResponse{
				DefaultResponse: DefaultResponse{Success: true},
				Response: []*CostCenter{
					&CostCenter{ID:123},
				},
			},
		},
		{
			"Create()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&CostCenterService{req}).Create(ctx, &CostCenter{Name: "test"})
				return
			},
			context.Background(),
			http.MethodPost,
			costCenterEndpoint.Action(create),
			&CostCenter{Name: "test"},
			&DefaultResponse{
				Success: true,
			},
		},
		{
			"Update()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&CostCenterService{req}).Update(ctx, &CostCenter{ID: 123})
				return
			},
			context.Background(),
			http.MethodPost,
			costCenterEndpoint.Action(update),
			&CostCenter{ID: 123},
			&DefaultResponse{
				Success: true,
			},
		},
		{
			"Inactivate()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&CostCenterService{req}).Inactivate(ctx, 123)
				return
			},
			context.Background(),
			http.MethodPost,
			costCenterEndpoint.Action(inactivate),
			&CostCenter{ID: 123},
			&DefaultResponse{
				Success: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestCostCenterError(t *testing.T) {
	testCases := []testTableError{
		{
			"Read()",
			func(req requester) error {
				_, err := (&CostCenterService{req}).Read(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Create()",
			func(req requester) error {
				_, err := (&CostCenterService{req}).Create(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Update()",
			func(req requester) error {
				_, err := (&CostCenterService{req}).Update(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Inactivate()",
			func(req requester) error {
				_, err := (&CostCenterService{req}).Inactivate(context.Background(), 1)
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
