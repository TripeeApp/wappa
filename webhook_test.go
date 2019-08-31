package wappa

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

func TestWebhook(t *testing.T) {
	testCases := []testTable{
		{
			"Read()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&WebhookService{req}).Read(ctx)
				return
			},
			context.Background(),
			http.MethodGet,
			webhookEndpoint,
			nil,
			&WebhookResult{
				Result: Result{Success: true},
				Listeners: []*Webhook{
					&Webhook{URL:"testing.url"},
				},
			},
		},
		{
			"Create()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&WebhookService{req}).Create(ctx, &Webhook{URL: "testing.url"})
				return
			},
			context.Background(),
			http.MethodPost,
			webhookEndpoint,
			&Webhook{URL: "testing.url"},
			&Result{Success: true},
		},
		{
			"Update()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&WebhookService{req}).Update(ctx, &Webhook{URL: "testing.new.url"})
				return
			},
			context.Background(),
			http.MethodPost,
			webhookEndpoint.Action(update),
			&Webhook{URL: "testing.new.url"},
			&Result{Success: true},
		},
		{
			"Activate()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&WebhookService{req}).Activate(ctx)
				return
			},
			context.Background(),
			http.MethodPost,
			webhookEndpoint.Action(activate),
			nil,
			&Result{Success: true},
		},
		{
			"Deactivate()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&WebhookService{req}).Deactivate(ctx)
				return
			},
			context.Background(),
			http.MethodPost,
			webhookEndpoint.Action(deactivate),
			nil,
			&Result{Success: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, test(tc))
	}
}

func TestWebhookError(t *testing.T) {
	testCases := []testTableError{
		{
			"Read()",
			func(req requester) error {
				_, err := (&WebhookService{req}).Read(context.Background())
				return err
			},
			errors.New("Error"),
		},
		{
			"Create()",
			func(req requester) error {
				_, err := (&WebhookService{req}).Create(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Update()",
			func(req requester) error {
				_, err := (&WebhookService{req}).Update(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Activate()",
			func(req requester) error {
				_, err := (&WebhookService{req}).Activate(context.Background())
				return err
			},
			errors.New("Error"),
		},
		{
			"Deactivate()",
			func(req requester) error {
				_, err := (&WebhookService{req}).Deactivate(context.Background())
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, testError(tc))
	}
}
