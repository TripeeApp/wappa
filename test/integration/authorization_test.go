package integration

import (
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/rdleal/wappa"
	"golang.org/x/oauth2"
)

const msgEnvMissing = "Skipping test because the required environment variable (%v) is not present."

const(
	envKeyWappaUsername = "WAPPA_USERNAME"
	envKeyWappaPassword = "WAPPA_PASSWORD"
)

func TestAuthorization(t *testing.T) {
	client := getOauthAppClient(t)

	cr, err := client.Ride.CancellationReason(context.Background())
	if err != nil {
		t.Fatalf("got error while calling Ride.CancellationReason(): %s; want nil.", err.Error())
	}

	if len(cr.Reasons) == 0 {
		t.Error("got empy cancellation reasons; want not empty.")
	}
}

func getOauthAppClient(t *testing.T) *wappa.Client {
	username, ok := os.LookupEnv(envKeyWappaUsername)
	if !ok {
		t.Skipf(msgEnvMissing, envKeyWappaUsername)
	}

	password, ok := os.LookupEnv(envKeyWappaPassword)
	if !ok {
		t.Skipf(msgEnvMissing, envKeyWappaPassword)
	}

	host := os.Getenv(envKeyWappaHost)

	u, err := url.Parse(host)
	if err != nil {
		t.Fatalf("got error while calling url.Parse(%s): %s; want nil.", host, err.Error())
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, loggingHTTPClient())

	return wappa.New(u, oauth2.NewClient(ctx, wappa.NewTokenSource(ctx, host, username, password)))
}
