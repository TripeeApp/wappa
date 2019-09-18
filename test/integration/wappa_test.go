package integration

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
	"unsafe"

	"github.com/mobilitee-smartmob/wappa"
	"golang.org/x/oauth2"
)

const (
	envKeyWappaHost          = "WAPPA_HOST"
	envKeyWappaToken         = "WAPPA_AUTH_TOKEN"
	envKeyWappaEmployeeEmail = "WAPPA_EMPLOYEE_EMAIL"
	envKeyWappaWebhookHost   = "WAPPA_WEBHOOK_HOST"
	envKeyWappaWebhookPort   = "WAPPA_WEBHOOK_PORT"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberBytes = "0123456789"
)

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

var auth bool

var employeeFilter wappa.Filter

// Wappa Client
var wpp *wappa.Client

// TODO: Generate random lat and lng starting at the center of Sao Paulo
// with destination within a certain radius and distance.
// ref: http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates
var (
	// Lat and Lng of Estacao Paraiso
	latOrigin = -23.5719548
	lngOrigin = -46.647377

	// Lat and Lng of Inferno Club
	latDest = -23.5515681
	lngDest = -46.6529553
)

var src = rand.NewSource(time.Now().UnixNano())

var logging = flag.Bool("log", false, "Define if tests should log the requests")

// Transport used to log the requests.
type transportLogger struct {
	base http.RoundTripper
}

func (t *transportLogger) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		reqBody []byte
		err     error
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

func TestMain(m *testing.M) {
	flag.Parse()

	host, err := url.Parse(os.Getenv(envKeyWappaHost))
	if err != nil {
		panic(err.Error())
	}

	token := os.Getenv(envKeyWappaToken)
	if token == "" {
		fmt.Println("No auth token. Some tests may not run!")
		wpp = wappa.New(host, nil)
	} else {
		host, _ := url.Parse(os.Getenv(envKeyWappaHost))

		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, loggingHTTPClient())
		hc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
		wpp = wappa.New(host, hc)

		auth = true
	}

	if email := os.Getenv(envKeyWappaEmployeeEmail); email != "" {
		employeeFilter = wappa.Filter{"email": []string{email}}
	}

	// Endpoint used to indicate if the server is up. Used when testing the Ride flow.
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		w.WriteHeader(http.StatusOK)
	})

	os.Exit(m.Run())
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

func findEmployeeID(f wappa.Filter) (id int, ok bool) {
	r, err := wpp.Employee.Read(context.Background(), f)
	if err != nil {
		fmt.Printf("got error while calling Employee.Read(%+v): '%s'; want nil.", f, err.Error())
		return
	}

	if len(r.Employees) == 0 {
		fmt.Printf("No Employee found. Skipping test...")
	} else {
		id = r.Employees[0].ID
		ok = true
	}

	return
}
