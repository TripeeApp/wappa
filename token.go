package wappa

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

// TokenSource implements the oauth2.TokenSource interface,
// in order to reuse the Wappa token.
type TokenSource struct {
	ctx      context.Context
	conf     *oauth2.Config
	username string
	password string
}

// NewTokenSource returns a oauth2.TokenSource for issuing a Wappa token.
//
// The provided context optionally controls which HTTP client is used. See the oauth2.HTTPClient variable.
// host should be sufixed by '/'.
func NewTokenSource(ctx context.Context, host, username, password string) *TokenSource {
	conf := &oauth2.Config{
		// ClientID and ClientSecret not used in this version
		ClientID:     "",
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf("%stoken", host),
		},
	}

	return &TokenSource{
		ctx:      ctx,
		conf:     conf,
		username: username,
		password: password,
	}
}

// Token returns a new oauth2.Token from a JWT or error.
func (ts *TokenSource) Token() (*oauth2.Token, error) {
	tk, err := ts.conf.PasswordCredentialsToken(ts.ctx, ts.username, ts.password)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tk.AccessToken, nil)
	if token == nil {
		return nil, fmt.Errorf("invalid JWT token: '%s'.", tk.AccessToken)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("token '%s' with invalid claims: '%s'.", tk.AccessToken, err)
	}

	var expTime time.Time
	if exp, ok := claims["exp"].(float64); ok {
		expTime = time.Unix(int64(exp), 0)
	}

	tk.Expiry = expTime

	return tk, nil
}
