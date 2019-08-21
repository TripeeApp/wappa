package wappa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestUnity(t *testing.T) {
	testCases := []testTable{
		{
			"Read()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UnityService{req}).Read(ctx, "123", "test")
				return
			},
			context.Background(),
			http.MethodGet,
			unityEndpoint.Action(read).Query(url.Values{"idUnidade": []string{"123"}, "codDescricao": []string{"test"}}),
			nil,
			&UnityResponse{
				DefaultResponse: DefaultResponse{Success: true},
				Response: []*Unity{
					&Unity{ID:123},
				},
			},
		},
		{
			"Create()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UnityService{req}).Create(ctx, &Unity{Code: "123"})
				return
			},
			context.Background(),
			http.MethodPost,
			unityEndpoint.Action(create),
			&Unity{Code: "123"},
			&DefaultResponse{
				Success: true,
			},
		},
		{
			"Update()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UnityService{req}).Update(ctx, &Unity{Code: "123"})
				return
			},
			context.Background(),
			http.MethodPost,
			unityEndpoint.Action(update),
			&Unity{Code: "123"},
			&DefaultResponse{
				Success: true,
			},
		},
		{
			"Inactivate()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UnityService{req}).Inactivate(ctx, "123")
				return
			},
			context.Background(),
			http.MethodPost,
			unityEndpoint.Action(inactivate).Query(url.Values{"idUnidade": []string{"123"}}),
			nil,
			&DefaultResponse{
				Success: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestUnityError(t *testing.T) {
	testCases := []testTableError{
		{
			"Read()",
			func(req requester) error {
				_, err := (&UnityService{req}).Read(context.Background(), "123", "test")
				return err
			},
			errors.New("Error"),
		},
		{
			"Create()",
			func(req requester) error {
				_, err := (&UnityService{req}).Create(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Update()",
			func(req requester) error {
				_, err := (&UnityService{req}).Update(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"ReadClassifier()",
			func(req requester) error {
				_, err := (&UnityService{req}).Inactivate(context.Background(), "1")
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
