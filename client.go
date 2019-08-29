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

// perser is the interface that performs the parsing of the body request.
type parser interface {
	Parse(io.Reader) error
}
// requester is the interface that performs a request
// to the server and delegates the persing to the parser interface. 
type requester interface{
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
type ResponseError struct{
	Code int
	Message int
}

// Result contains the base porperties
// all requests to the Wappa API returns.
type Result struct{
	Success bool
	Error ResponseError
	// Returned by the API when error occurs
	Message string
}

// ApiError implements the error interface
// and returns infos from the request.
type ApiError struct{
	statusCode int
	msg        string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("Error Status Code: %d; Message: %s.", e.statusCode, e.msg)
}


// Client is responsible for handling requests to the Wappa API.
type Client struct{
	// client to comunicate with the API.
	client *http.Client

	// Host used for API requets.
	// host should always be specified with a trailing slash.
	host *url.URL

	// Services implemented
	Driver *DriverService
	Employee *EmployeeService
	Ride *RideService
	Webhook *WebhookService
}

// Client returns a new Wappa API client with provided host URL and HTTP client.
func New(host *url.URL, client *http.Client) *Client {
	if client == nil {
		client = &http.Client{}
	}
	c := &Client{client: client, host: host}

	// Sets services.
	c.Webhook = &WebhookService{c}
	c.Driver = &DriverService{c}
	c.Ride = &RideService{c}
	c.Employee = &EmployeeService{c}

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

	// If the output is a parser, pass it the control over the body parsing.
	if p, ok := output.(parser); ok {
		if err := p.Parse(res.Body); err != nil {
			return &ApiError{
				statusCode: res.StatusCode,
				msg:        err.Error(),
			}
		}
	// Otherwise parses the body using the default API return type (JSON).
	} else {
		b, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(b, output); err != nil {
			return &ApiError{
				statusCode: res.StatusCode,
				msg: fmt.Sprintf("Couldn't unmarshal body: '%s'. Message: '%s'.",  string(b), err.Error()),
			}
		}
	}

	return nil
}
