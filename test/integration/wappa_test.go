package integration

import (
	"context"
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"net/http"
	"io/ioutil"
	"os"
	"time"
	"unsafe"

	"github.com/rdleal/wappa"
	"golang.org/x/oauth2"
)


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

var auth bool

// Ligue Taxi Client
var wpp *wappa.Client

var src = rand.NewSource(time.Now().UnixNano())

var logging = flag.Bool("log", false, "Define if tests should log the requests")

// Transport used to log the requests.
type transportLogger struct {
	base http.RoundTripper
}

func (t *transportLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		reqBody []byte
		err error
	)

	if req.Body != nil {
		reqBody, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body.Close()

		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
	}

	res, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	resBody, _ := ioutil.ReadAll(res.Body)

	log.Println("Request")
	log.Println("Headers:")
	for k, v := range req.Header {
		log.Printf("%s: %s\n", k, v)
	}
	log.Printf("/%s %s %s --> Response %s %s",
		req.Method, req.URL.String(), string(reqBody), res.Status, string(resBody))

	res.Body = ioutil.NopCloser(bytes.NewBuffer(resBody))

	return res, nil
}

func init() {
	flag.Parse()

	host, err := url.Parse(os.Getenv(envKeyWappaHost))
	if err != nil {
		panic(err.Error())
	}

	token := os.Getenv("WAPPA_AUTH_TOKEN")
	if token == "" {
		fmt.Println("No auth token. Some tests may not run!")
		wpp = wappa.New(host, nil)
	} else {
		host, _ := url.Parse(os.Getenv(envKeyWappaHost))

		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, loggingHTTPClient())
		tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
		wpp = wappa.New(host, tc)

		auth = true
	}
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

func loggingHTTPClient() *http.Client {
	if !*logging {
		return http.DefaultClient
	}

	hc := &http.Client{
		Transport: &transportLogger{http.DefaultTransport},
	}

	return hc
}

func checkAuth(name string) bool {
	if !auth {
		fmt.Printf("Skipping test %s for no authorization token was set.", name)
	}
	return auth
}
