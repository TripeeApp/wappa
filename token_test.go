package wappa

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestNewTokenSource(t *testing.T) {
	host := "http://testing-host/"
	username := "test"
	password := "testing-passwod"
	clientID := ""
	clientSecret := ""

	ts := NewTokenSource(context.Background(), host, username, password)

	if ts.username != username {
		t.Errorf("got username: '%s'; want '%s'", ts.username, username)
	}

	if ts.password != password {
		t.Errorf("got password: '%s'; want '%s'", ts.password, password)
	}

	if ts.conf == nil {
		t.Fatalf("got nil conf; want not nil")
	}

	if ts.conf.ClientID != clientID {
		t.Errorf("got conf.ClientID: '%s'; want '%s'", ts.conf.ClientID, clientID)
	}

	if ts.conf.ClientSecret != clientSecret {
		t.Errorf("got conf.ClientSecret: '%s'; want '%s'", ts.conf.ClientSecret, clientSecret)
	}

	if want := fmt.Sprintf("%stoken", host); ts.conf.Endpoint.TokenURL != want {
		t.Errorf("got conf.Endpoint.TokenURL: '%s'; want '%s'", ts.conf.Endpoint.TokenURL, want)
	}
}

func TestTokenSourceToken(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/token"
		if r.URL.String() != expected {
			t.Errorf("URL = %q; want %q", r.URL, expected)
		}
		mySigningKey := []byte("AllYourBase")

		// Create the Claims
		claims := &jwt.StandardClaims{
		    ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
		    Issuer:    "test",
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, _ := token.SignedString(mySigningKey)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"access_token": "%s"}`, ss)))
	}))
	defer s.Close()

	ts := NewTokenSource(context.Background(), s.URL + "/", "", "")

	if _, err := ts.Token(); err != nil {
		t.Errorf("got error calling Token(): '%s'; want nil.", err.Error())
	}
}
