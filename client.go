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

const (
	// Error message format.
	errFmt = `Error while request the LigueTaxi API: %s; Status Code: %d; Body: %s.`
)

// Filter is used to filter the requests to the API.
type Filter map[string]string

// Values returns a url.Values mapped between Filter and values.
func (f Filter) Values(fields map[string]string) url.Values {
	vals := url.Values{}
	for k, v := range f {
		if q, ok := fields[k]; ok {
			vals.Set(q, v)
		}
	}
	return vals
}

// ResponseError contains the default error porperties
// returnes by the API.
type ResponseError struct {
	Code int
	Message int
}

// DefaultResponse contains the basics porperties
// all requests to the Wappa API returns.
type DefaultResponse struct {
	Success bool
	Quantity int `json:"quantidade"`
	Error ResponseError
}

// requester is the interface that performs a request
// to the server and parses the payload.
type requester interface{
	Request(ctx context.Context, method string, path endpoint, body, output interface{}) error
}

// ApiError implements the error interface
// and returns infos from the request.
type ApiError struct {
	statusCode int
	body       []byte
	msg        string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf(errFmt, e.msg, e.statusCode, e.body)
}


// Client is responsible for handling requests to the Wappa API.
type Client struct {
	// client to comunicate with the API.
	client *http.Client

	// Host used for API requets.
	// host should always be specified with a trailing slash.
	host *url.URL
}

// Client returns a new Wappa API client with provided host URL and HTTP client.
func New(host *url.URL, client *http.Client) *Client {
	if client == nil {
		client = &http.Client{}
	}
	return &Client{client, host}
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

	// TODO: add tests for error on reading body
	r, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(r, output); err != nil {
		return &ApiError{
			statusCode: res.StatusCode,
			body:       r,
			msg:        err.Error(),
		}
	}

	return nil
}
