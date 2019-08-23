package wappa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestCollaborator(t *testing.T) {
	testCases := []testTable{
		{
			"Read()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&CollaboratorService{req}).Read(ctx, Filter{"name": "test"})
				return
			},
			context.Background(),
			http.MethodGet,
			collaboratorEndpoint.Action(read).Query(url.Values{"nome": []string{"test"}}),
			nil,
			&CollaboratorResponse{
				DefaultResponse: DefaultResponse{Success: true},
				Response: []*Collaborator{
					&Collaborator{ID:123},
				},
			},
		},
		{
			"Create()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&CollaboratorService{req}).Create(ctx, &Collaborator{Name: "123"})
				return
			},
			context.Background(),
			http.MethodPost,
			collaboratorEndpoint.Action(create),
			&Collaborator{Name: "123"},
			&DefaultResponse{
				Success: true,
			},
		},
		{
			"Update()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&CollaboratorService{req}).Update(ctx, &Collaborator{Company: 1})
				return
			},
			context.Background(),
			http.MethodPost,
			collaboratorEndpoint.Action(update),
			&Collaborator{Company: 1},
			&OperationDefaultResponse{
				DefaultResponse{Success: true},
			},
		},
		{
			"Inactivate()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&CollaboratorService{req}).Inactivate(ctx, 1)
				return
			},
			context.Background(),
			http.MethodPost,
			collaboratorEndpoint.Action(inactivate),
			&Collaborator{ID: 1},
			&OperationDefaultResponse{
				DefaultResponse{Success: true},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestCollaboratorError(t *testing.T) {
	testCases := []testTableError{
		{
			"Read()",
			func(req requester) error {
				_, err := (&CollaboratorService{req}).Read(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Create()",
			func(req requester) error {
				_, err := (&CollaboratorService{req}).Create(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Update()",
			func(req requester) error {
				_, err := (&CollaboratorService{req}).Update(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Inactivate()",
			func(req requester) error {
				_, err := (&CollaboratorService{req}).Inactivate(context.Background(), 1)
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
