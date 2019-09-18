package integration

import (
	"context"
	"testing"

	wappa "github.com/mobilitee-smartmob/wappa/v2"
)

func TestWebhook(t *testing.T) {
	if !checkAuth("TestWebhook") {
		return
	}

	r, err := wpp.Webhook.Read(context.Background())
	if err != nil {
		t.Fatalf("got error calling Webhook.Read(): '%s'; want nil.", err.Error())
	}

	if !r.Success {
		t.Errorf("got failed response while reading Webhooks: '%s'; want it to be successful.", r.Message)
	}

	var curWebhook *wappa.Webhook
	// If there's already a webhook in the API, uses its data from now on.
	// Otherwise create a new one.
	if len(r.Listeners) > 0 {
		curWebhook = r.Listeners[0]
	} else {
		curWebhook = &wappa.Webhook{
			URL:      "test.com",
			Endpoint: "path/for/testing",
			AuthKey:  "auth-key",
		}
		res, err := wpp.Webhook.Create(context.Background(), curWebhook)
		if err != nil {
			t.Fatalf("got error calling Webhook.Create(%+v): '%s'; want nil.", curWebhook, err.Error())
		}

		if !res.Success {
			t.Errorf("got failed response while creating Webhooks (%+v): '%s'; want it to be successful.", curWebhook, res.Message)
		}

		curWebhook.Active = true
	}

	// Checks current webhhok status and toggles it, then set it to its original state.
	if curWebhook.Active {
		deactivateWebhook(t)
		activateWebhook(t)
	} else {
		activateWebhook(t)
		deactivateWebhook(t)
	}

	res, err := wpp.Webhook.Update(context.Background(), curWebhook)
	if err != nil {
		t.Fatalf("got error calling Webhook.Update(%+v): %s; want nil.", curWebhook, err.Error())
	}

	if !res.Success {
		t.Errorf("got failed response while updating Webhooks (%+v): '%s'; want it to be successful.", curWebhook, res.Message)
	}
}

func activateWebhook(t *testing.T) {
	res, err := wpp.Webhook.Activate(context.Background())
	if err != nil {
		t.Fatalf("got error calling Webhook.Activate(): '%s'; want nil.", err.Error())
	}

	if !res.Success {
		t.Errorf("got failed response while activating Webhooks: '%s'; want it to be successful.", res.Message)
	}
}

func deactivateWebhook(t *testing.T) {
	res, err := wpp.Webhook.Deactivate(context.Background())
	if err != nil {
		t.Fatalf("got error calling Webhook.Deactivate(): '%s'; want nil.", err.Error())
	}

	if !res.Success {
		t.Errorf("got failed response while deactivating Webhooks: '%s'; want it to be successful.", res.Message)
	}
}
