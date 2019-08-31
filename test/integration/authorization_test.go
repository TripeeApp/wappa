package integration

//import (
//	"context"
//	"fmt"
//	"os"
//	"testing"
//	"net/url"
//
//	"bitbucket.org/mobilitee/wappa"
//	"golang.org/x/oauth2"
//)
//
//const msgEnvMissing = "Skipping test because the required environment variable (%v) is not present."
//
//const(
//	envKeyWappaUsername = "WAPPA_USERNAME"
//	envKeyWappaPassword = "WAPPA_PASSWORD"
//)
//
//func TestAuthorization(t *testing.T) {
//	client := getOauthAppClient(t)
//	desc := randString(5, letterBytes)
//
//	res, err := client.Role.Create(context.Background(), &wappa.Role{ID: 1234, Description: desc})
//	if err != nil {
//		t.Fatalf("go error while calling Role.Read(nil): %s; want nil.", err.Error())
//	}
//
//	if !res.Success {
//		t.Errorf("go error while creating a role: %s; want nil.", err.Error())
//	}
//
//	f := wappa.Filter{"desc": desc}
//	role, err := client.Role.Read(context.Background(), f)
//	if err != nil {
//		t.Fatalf("go error while calling Role.Read(%+v): %s; want nil.",f, err.Error())
//	}
//
//	if !role.Success {
//		t.Errorf("got error while reading a role: %s; want nil.", res.Message)
//	}
//}
//
//func getOauthAppClient(t *testing.T) *wappa.Client {
//	username, ok := os.LookupEnv(envKeyWappaUsername)
//	if !ok {
//		t.Skipf(msgEnvMissing, envKeyWappaUsername)
//	}
//
//	password, ok := os.LookupEnv(envKeyWappaPassword)
//	if !ok {
//		t.Skipf(msgEnvMissing, envKeyWappaPassword)
//	}
//
//	host := os.Getenv(envKeyWappaHost)
//
//	conf := &oauth2.Config{
//		Endpoint: oauth2.Endpoint{
//			TokenURL: fmt.Sprintf("%stoken", host),
//			AuthStyle: oauth2.AuthStyleInParams,
//		},
//	}
//
//	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, loggingHTTPClient())
//	tkn, err := conf.PasswordCredentialsToken(ctx, username, password)
//	if err != nil {
//		t.Fatalf("got error while calling oauth2.PasswordCredentialsToken(%s, %s): %s; want nil.", username, password, err.Error())
//	}
//
//	u, err := url.Parse(host)
//	if err != nil {
//		t.Fatalf("got error while calling url.Parse(%s): %s; want nil.", host, err.Error())
//	}
//
//	return wappa.New(u, conf.Client(ctx, tkn))
//}
