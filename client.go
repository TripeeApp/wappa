package wappa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// requester is the interface that performs a request
// to the server and delegates the persing to the parser interface.
type requester interface {
	Request(ctx context.Context, method string, path endpoint, body, output interface{}) error
}

// Filter is used to filter the requests to the API.
type Filter map[string][]string

// Values returns a url.Values mapped between Filter and values.
func (f Filter) Values(fields map[string]string) url.Values {
	vals := url.Values{}
	for k, filters := range f {
		if q, ok := fields[k]; ok {
			for _, v := range filters {
				vals.Add(q, v)
			}
		}
	}
	return vals
}

// ResponseError contains the default error porperties
// returnes by the API.
type ResponseError struct {
	Code    int
	Message int
}

// Result contains the base porperties
// all requests to the Wappa API returns.
type Result struct {
	Success bool
	Error   ResponseError
	// Returned by the API when error occurs
	Message string
}

// ApiError implements the error interface
// and returns infos from the request.
type ApiError struct {
	statusCode int
	msg        string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("Error Status Code: %d; Message: %s.", e.statusCode, e.msg)
}

type service struct {
	client requester
}

// Client is responsible for handling requests to the Wappa API.
type Client struct {
	// client to comunicate with the API.
	client *http.Client

	// Host used for API requets.
	// host should always be specified with a trailing slash.
	host *url.URL

	// reuse a single struct intead of allocation one for each service on the heap.
	common service

	// Services implemented
	Driver   *DriverService
	Employee *EmployeeService
	Quote    *QuoteService
	Ride     *RideService
	Webhook  *WebhookService
}

// NewClient returns a new Wappa API client with provided host URL and HTTP client.
func NewClient(host *url.URL, client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}
	c := &Client{client: client, host: host}

	c.common.client = c
	// Sets services.
	c.Driver = (*DriverService)(&c.common)
	c.Employee = (*EmployeeService)(&c.common)
	c.Quote = (*QuoteService)(&c.common)
	c.Ride = (*RideService)(&c.common)
	c.Webhook = (*WebhookService)(&c.common)

	return c
}

// Request created an API request. A relative path can be providaded
// in which case it is resolved relative to the host of the Client.
func (c *Client) Request(ctx context.Context, method string, path endpoint, body, output interface{}) error {
	u, err := c.host.Parse(path.String())
	if err != nil {
		return err
	}

	var b io.ReadWriter
	if body != nil {
		b = new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, u.String(), b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bd, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(bd, output); err != nil {
		return &ApiError{
			statusCode: res.StatusCode,
			msg:        fmt.Sprintf("Couldn't unmarshal body: '%s'. Message: '%s'.", string(bd), err.Error()),
		}
	}

	return nil
}

// endpoint for checking the API status. Pulled off for teting.
var statusEndpoint endpoint = `status`

// Status returns the API status.
func (c *Client) Status(ctx context.Context) (ok bool, err error) {
	u, err := c.host.Parse(statusEndpoint.String())
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return
	}

	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		return
	}

	ok = res.StatusCode == http.StatusOK

	return
}
