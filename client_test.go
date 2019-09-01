package wappa

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

type testTable struct {
	name	string
	call	func(ctx context.Context, req requester) (resp interface{}, err error)
	ctx	context.Context
	method	string
	path	endpoint
	body	interface{}
	wantRes	interface{}
}

type testTableError struct {
	name	string
	call	func(req requester) error
	err	error
}

type testRequester struct {
	body	interface{}
	ctx	context.Context
	err	error
	method	string
	output	reflect.Value
	path	endpoint

}

func (t *testRequester) Request(ctx context.Context, method string, path endpoint, body, output interface{}) error {
	t.ctx = ctx
	t.method = method
	t.path = path
	t.body = body

	if t.output.IsValid() {
		out := reflect.ValueOf(output)
		if !out.IsNil() && out.Elem().CanSet() {
			out.Elem().Set(t.output)
		}
	}

	return t.err
}

func test(tc testTable) func(t *testing.T) {
	return func(t *testing.T) {
		req := &testRequester{output: reflect.ValueOf(tc.wantRes).Elem()}

		res, err := tc.call(tc.ctx, req)
		if err != nil {
			t.Fatalf("got error while calling User %s: %s, want nil", tc.name, err.Error())
		}

		if !reflect.DeepEqual(req.ctx, tc.ctx) {
			t.Errorf("got Requester Context %+v; want %+v.", req.ctx, tc.ctx)
		}

		if req.method != tc.method {
			t.Errorf("got request method: %s; want %s.", req.method, tc.method)
		}

		if req.path != tc.path {
			t.Errorf("got request path: %s; want %s.", req.path, tc.path)
		}

		if !reflect.DeepEqual(req.body, tc.body) {
			t.Errorf("got request body: %+v; want %+v.", req.body, tc.body)
		}

		if !reflect.DeepEqual(res, tc.wantRes) {
			t.Errorf("got response: %+v; want %+v.", res, tc.wantRes)
		}
	}
}

func testError(tc testTableError) func(t *testing.T) {
	return func(t *testing.T) {
		req := &testRequester{err: tc.err}

		err := tc.call(req)
		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("got error: %s; want %s.", err, tc.err)
		}
	}
}

var testingMap = map[string]string{"id": "idTest"}

func TestFilter(t *testing.T) {
	testCases := []struct{
		filter Filter
		want url.Values
	}{
		{
			Filter{"id": []string{"123"}},
			url.Values{
				"idTest": []string{"123"},
			},
		},
	}

	for _, tc := range testCases {
		if got := tc.filter.Values(testingMap); !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got from Filter.Values(%+v): %+v; want %+v.", tc.filter, got, tc.want)
		}
	}
}

func TestNew(t *testing.T) {
	testCases := []struct{
		host	*url.URL
		client	*http.Client
		wantClient *http.Client
	}{
		{&url.URL{}, nil, http.DefaultClient},
		{&url.URL{}, &http.Client{}, &http.Client{}},
	}

	for _, tc := range testCases {
		c := New(tc.host, tc.client)

		if c.host != tc.host {
			t.Errorf("got c.host : %s; want %s.", c.host, tc.host)
		}

		if !reflect.DeepEqual(c.client, tc.wantClient) {
			t.Errorf("got client %+v; want %+v.", c.client, tc.wantClient)
		}

		if c.Driver == nil {
			t.Errorf("got Driver == nil; want not nil.")
		}

		if c.Employee == nil {
			t.Errorf("got Employee == nil; want not nil.")
		}

		if c.Quote == nil {
			t.Errorf("got Quote == nil; want not nil.")
		}

		if c.Ride == nil {
			t.Errorf("got Ride == nil; want not nil.")
		}

		if c.Webhook == nil {
			t.Errorf("got WebhookService == nil; want not nil.")
		}

	}
}

func newMockServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

type dummy struct { Name string `json:"name"` }

func TestClientRequest(t *testing.T) {
	emptyObj := []byte(`{}`)

	testCases := []struct{
		endpoint endpoint
		method	 string
		body	 interface{}
		server	 *httptest.Server
		output   interface{}
		wantOut	 interface{}
	}{
		{
			"",
			http.MethodGet,
			&Result{},
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("go Request.Method %s; want %s.", r.Method, http.MethodGet)
				}
				w.Write(emptyObj)
			}),
			&Result{},
			&Result{},
		},
		{
			"foo",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if want := "/api/foo"; r.URL.Path != want {
					t.Errorf("got Request.URL: %s; want %s.", r.URL.Path, want)
				}
				w.Write(emptyObj)
			}),
			&Result{},
			&Result{},
		},
		{
			"foo",
			http.MethodPost,
			&Result{},
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if c := r.Header.Get("Content-Type"); c != "application/json" {
					t.Errorf("got 'Content-Type' Header: '%s'; want 'application/json'.", c)
				}
				w.Write(emptyObj)
			}),
			&Result{},
			&Result{},
		},
		{
			"",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Body != http.NoBody {
					t.Errorf("got Request.Body: %+v, want empty.", r.Body)
				}
				w.Write(emptyObj)
			}),
			&Result{},
			&Result{},
		},
		{
			"",
			http.MethodPost,
			&Result{Success: true},
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Body == http.NoBody {
					t.Error("got Request.Body empty, want not empty.")
				}

				got, _ := ioutil.ReadAll(r.Body)

				if want := []byte(`"Success":true`); !bytes.Contains(got, want) {
					t.Errorf("got body: %s, want %s.", got, want)
				}
				w.Write([]byte(`{"Success":true}`))
			}),
			&Result{},
			&Result{Success: true},
		},
		{
			"",
			http.MethodPost,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"name":"Testing"}`))
			}),
			&dummy{},
			&dummy{Name: "Testing"},
		},
	}

	for _, tc := range testCases {
		output := tc.output

		u, _ := url.Parse(tc.server.URL)
		c := New(u, nil)

		err := c.Request(context.Background(), tc.method, tc.endpoint, tc.body, output)
		if err != nil {
			t.Fatalf("got error calling Client.Request(context.Background(), %s, %s, %+v, %+v): %s; want nil.",
				tc.method, tc.endpoint, tc.body, output, err.Error())
		}

		if !reflect.DeepEqual(output, tc.wantOut) {
			t.Errorf("got output from Client.Request(): %+v; want %+v.", output, tc.wantOut)
		}

		tc.server.Close()
	}
}

func TestClientRequestError(t *testing.T) {
	testCases := []struct{
		path		endpoint
		method		string
		body		interface{}
		server		*httptest.Server
		assertError	func(e error)
	}{
		{
			":",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {}),
			nil,
		},
		{
			"",
			http.MethodGet,
			make(chan int),
			newMockServer(func(w http.ResponseWriter, r *http.Request) {}),
			nil,
		},
		{
			"",
			",",
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {}),
			nil,
		},
		{
			"",
			http.MethodPost,
			nil,
			httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})),
			nil,
		},
		{
			"",
			http.MethodPost,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}),
			func(e error) {
				err, ok := e.(*ApiError)
				if !ok {
					t.Fatal("got error different from *ApiError")
				}

				if err != nil {
					if wantStatus := http.StatusInternalServerError; err.statusCode != wantStatus {
						t.Errorf("got Error.statusCode: %d; want %d.", err.statusCode, wantStatus)
					}

					if wantSubStr := "invalid character 'I'"; !strings.Contains(err.msg, wantSubStr) {
						t.Errorf("go Error.msg: %s; want it to contain `%s` substring.", err.msg, wantSubStr)
					}
				}
			},
		},
	}

	for _, tc := range testCases {
		u, _ := url.Parse(tc.server.URL)
		c := New(u, nil)

		err := c.Request(context.Background(), tc.method, tc.path, tc.body, &dummy{})
		if err == nil {
			t.Errorf("got error nil; want not nil.")
		}

		if tc.assertError != nil {
			tc.assertError(err)
		}

		tc.server.Close()
	}
}

func TestClientRequestWithContext(t *testing.T) {
	s := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		select {
		case <-time.After(1 * time.Second):
			t.Errorf("Expected request to be canceled by context")
		case <-ctx.Done():
			return
		}
	})
	defer s.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	u, _ := url.Parse(s.URL)
	c := New(u, nil)

	if err := c.Request(ctx, http.MethodGet, "/", nil, nil); err == nil {
		t.Errorf("got error nil; want not nil")
	}
}

func TestClientStatus(t *testing.T) {

	testCases := []struct{
		server	 *httptest.Server
		wantOK	 bool
	}{
		{
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if want := "/api/status"; r.URL.Path != want {
					t.Errorf("got Request.URL: %s; want %s.", r.URL.Path, want)
				}
				w.WriteHeader(http.StatusOK)
			}),
			true,
		},
		{
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("go Request.Method %s; want %s.", r.Method, http.MethodGet)
				}
				w.WriteHeader(http.StatusOK)
			}),
			true,
		},
		{
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			false,
		},
	}

	for _, tc := range testCases {
		u, _ := url.Parse(tc.server.URL)
		c := New(u, nil)

		ok, err := c.Status(context.Background())
		if err != nil {
			t.Fatalf("got error calling Client.Status(context.Background()): %s; want nil.", err.Error())
		}

		if ok != tc.wantOK{
			t.Errorf("got output from Client.Status(): %+v; want %+v.", ok, tc.wantOK)
		}

		tc.server.Close()
	}
}

func TestClientStatusError(t *testing.T) {
	testCases := []struct{
		server		*httptest.Server
	}{
		{
			httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})),
		},
	}

	for _, tc := range testCases {
		u, _ := url.Parse(tc.server.URL)
		c := New(u, nil)

		_, err := c.Status(context.Background())
		if err == nil {
			t.Errorf("got error nil; want not nil.")
		}

		tc.server.Close()
	}
}

func TestClientStatusWithContext(t *testing.T) {
	s := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		select {
		case <-time.After(1 * time.Second):
			t.Errorf("Expected request to be canceled by context")
		case <-ctx.Done():
			return
		}
	})
	defer s.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	u, _ := url.Parse(s.URL)
	c := New(u, nil)

	if _, err := c.Status(ctx); err == nil {
		t.Errorf("got error nil; want not nil")
	}
}
