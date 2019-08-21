package wappa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestRole(t *testing.T) {
	testCases := []testTable{
		{
			"Read()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&RoleService{req}).Read(ctx, "123", "test")
				return
			},
			context.Background(),
			http.MethodGet,
			roleEndpoint.Action(read).Query(url.Values{"idCargo": []string{"123"}, "descricao": []string{"test"}}),
			nil,
			&RoleResponse{
				DefaultResponse: DefaultResponse{Success: true},
				Response: []*Role{
					&Role{ID:123},
				},
			},
		},
		{
			"Create()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&RoleService{req}).Create(ctx, &Role{ID: 123})
				return
			},
			context.Background(),
			http.MethodPost,
			roleEndpoint.Action(create),
			&Role{ID: 123},
			&DefaultResponse{
				Success: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestRoleError(t *testing.T) {
	testCases := []testTableError{
		{
			"Read()",
			func(req requester) error {
				_, err := (&RoleService{req}).Read(context.Background(), "123", "test")
				return err
			},
			errors.New("Error"),
		},
		{
			"Create()",
			func(req requester) error {
				_, err := (&RoleService{req}).Create(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
