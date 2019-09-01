package wappa

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

// TokenSource implements the oauth2.TokenSource interface,
// in order to reuse the Wappa token.
type TokenSource struct {
	ctx context.Context
	conf *oauth2.Config
	username string
	password string
}

// NewTokenSource returns a oauth2.TokenSource for issuing a Wappa token.
//
// The provided context optionally controls which HTTP client is used. See the oauth2.HTTPClient variable.
// host should be sufixed by '/'.
func NewTokenSource(ctx context.Context, host, clientID, clientSecret, username, password string) *TokenSource {
	conf := &oauth2.Config{
		ClientID: clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf("%stoken", host),
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}

	return &TokenSource{
		ctx: ctx,
		conf: conf,
		username: username,
		password: password,
	}
}

// Token returns a new oauth2.Token or error.
func (ts *TokenSource) Token() (*oauth2.Token, error) {
	return ts.conf.PasswordCredentialsToken(ts.ctx, ts.username, ts.password)
}
