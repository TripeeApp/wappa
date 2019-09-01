package wappa

import (
	"context"
	"net/http"
)

const webhookEndpoint endpoint = `webhook`

// Webhook is the  struct representing the Webhook
// entity in the API.
type Webhook struct {
	URL string `json:"url"`
	Endpoint string `json:"endpoint"`
	AuthKey string `json:"authKey"`
	Active bool `json:"active,omitempty"`
}

// WebhookResult is the API response payload.
type WebhookResult struct {
	Result

	Listeners []*Webhook `json:"listeners"`
}

// WebhookService is responsible for handling
// the requests to the webhook resource.
type WebhookService service

// Read returns the webhooks created in the API.
func (ws *WebhookService) Read(ctx context.Context) (*WebhookResult, error) {
	wr := &WebhookResult{}

	if err := ws.client.Request(ctx, http.MethodGet, webhookEndpoint, nil, wr); err != nil {
		return nil, err
	}

	return wr, nil
}

// Create creates a webhook resource in the API.
func (ws *WebhookService) Create(ctx context.Context, w *Webhook) (*Result, error) {
	res := &Result{}

	if err := ws.client.Request(ctx, http.MethodPost, webhookEndpoint, w, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Updated edits the webhook information.
func (ws *WebhookService) Update(ctx context.Context, w *Webhook) (*Result, error) {
	res := &Result{}

	if err := ws.client.Request(ctx, http.MethodPost, webhookEndpoint.Action(update), w, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Activate activates the current webhook if it has been deactivated in the API.
func (ws *WebhookService) Activate(ctx context.Context) (*Result, error) {
	res := &Result{}

	if err := ws.client.Request(ctx, http.MethodPost, webhookEndpoint.Action(activate), nil, res); err != nil {
		return nil, err
	}

	return res, nil
}

// Deactivate deactivates the current webhook in the API.
func (ws *WebhookService) Deactivate(ctx context.Context) (*Result, error) {
	res := &Result{}

	if err := ws.client.Request(ctx, http.MethodPost, webhookEndpoint.Action(deactivate), nil, res); err != nil {
		return nil, err
	}

	return res, nil
}
