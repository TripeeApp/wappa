package integration

import (
	"context"
	"os"
	"testing"
	"net/url"

	"github.com/rdleal/wappa"
	"golang.org/x/oauth2"
)

const msgEnvMissing = "Skipping test because the required environment variable (%v) is not present."

const(
	envKeyWappaClientID = "WAPPA_CLIENT_ID"
	envKeyWappaClientSecret ="WAPPA_CLIENT_SECRET"
	envKeyWappaUsername = "WAPPA_USERNAME"
	envKeyWappaPassword = "WAPPA_PASSWORD"
)

func TestAuthorization(t *testing.T) {
	client := getOauthAppClient(t)

	role, err := client.Role.Read(context.Background(), nil)
	if err != nil {
		t.Fatalf("go error while calling Role.Read(): %s; want nil.", err.Error())
	}

	if !role.Success {
		t.Errorf("got error while reading a role: %s; want nil.", role.Message)
	}
}

func getOauthAppClient(t *testing.T) *wappa.Client {
	username, ok := os.LookupEnv(envKeyWappaUsername)
	if !ok {
		t.Skipf(msgEnvMissing, envKeyWappaClientID)
	}

	password, ok := os.LookupEnv(envKeyWappaPassword)
	if !ok {
		t.Skipf(msgEnvMissing, envKeyWappaPassword)
	}

	clientID, ok := os.LookupEnv(envKeyWappaClientID)
	if !ok {
		t.Skipf(msgEnvMissing, envKeyWappaClientID)
	}

	clientSecret, ok := os.LookupEnv(envKeyWappaClientSecret)
	if !ok {
		t.Skipf(msgEnvMissing, envKeyWappaClientSecret)
	}

	host := os.Getenv(envKeyWappaHost)

	u, err := url.Parse(host)
	if err != nil {
		t.Fatalf("got error while calling url.Parse(%s): %s; want nil.", host, err.Error())
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, loggingHTTPClient())

	return wappa.New(u, oauth2.NewClient(ctx, wappa.NewTokenSource(ctx, host, clientID, clientSecret, username, password)))
}
