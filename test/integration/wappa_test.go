package integration

import (
	"context"
	"bytes"
	"flag"
	"log"
	"math/rand"
	"net/url"
	"net/http"
	"io/ioutil"
	"os"
	"time"
	"unsafe"

	"bitbucket.org/mobilitee/wappa"
	"golang.org/x/oauth2"
)

// Transport used to log the requests.
type transportLogger struct {
	base http.RoundTripper
}

func (t *transportLogger) RoundTrip(r *http.Request) (*http.Response, error) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	res, err := t.base.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	resBody, _ := ioutil.ReadAll(res.Body)
	log.Printf("Request /%s %s %s --> Response %s %s",
		r.Method, r.URL.String(), string(reqBody), res.Status, string(resBody))

	res.Body = ioutil.NopCloser(bytes.NewBuffer(resBody))

	return res, nil
}

const (
	envKeyWappaHost = "WAPPA_HOST"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes = "0123456789"
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax = 63 / letterIdxBits
)

// Ligue Taxi Client
var wappa *wappa.Client

var src = rand.NewSource(time.Now().UnixNano())

var logging = flag.Bool("log", false, "Define if tests should log the requests")

func init() {
	flag.Parse()

	token := os.Getenv("WAPPA_AUTH_TOKEN")
	if token == "" {
		panic("No auth token defined!!")
	}

	host, _ := url.Parse(os.Getenv(envKeyWappaHost))

	tc := oauth2.NewClient(loggingContext, &oauth2.StaticTokenSource(
		&oauth2.Token{accessToken: token},
	))

	wappa = wappa.New(host, tc)
}

func randString(max int, rangeBytes string) string {
	b := make([]byte, max)

	for i, cache, remain := max-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMax); idx < len(rangeBytes) {
			b[i] = rangeBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func checkOperation(delay time.Duration, retries int, check func() error ) (success bool, err error) {
	for i := 0; i < retries; i++ {
		<-time.Tick(delay)

		err = check()
		if success := err == nil; success {
			return success, err
		}
	}

	return success, err
}

func loggingContext() context.Context {
	ctx := context.Background()
	if *logging {
		hc := &http.Client{
			Transport: &transportLogger{http.DefaultTransport}
		}
		ctx = ctx.WithValue(context.Background(), oauth2.HTTPClient, hc)
	}

	return ctx
}
