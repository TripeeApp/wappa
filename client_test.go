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

func TestNew(t *testing.T) {
	testCases := []struct{
		host	*url.URL
		client	*http.Client
		wantClient *http.Client
	}{
		{&url.URL{}, nil, &http.Client{}},
		{&url.URL{}, &http.Client{}, &http.Client{}},
	}

	for _, tc := range testCases {
		c := New(tc.host, tc.client)

		if c.host != tc.host {
			t.Errorf("got c.host : %s; want %s.", c.host, tc.host)
		}

		if !reflect.DeepEqual(c.client, tc.wantClient) {
			t.Errorf("got Client %+v; want %+v.", c.client, tc.wantClient)
		}
	}
}

func newMockServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

type dummy struct {
	Name string `json:"name"`
}

func TestClientRequest(t *testing.T) {
	emptyObj := []byte(`{}`)

	testCases := []struct{
		endpoint endpoint
		method	 string
		body	 interface{}
		server	 *httptest.Server
		wantOut	 interface{}
	}{
		{
			"",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("go Request.Method %s; want %s.", r.Method, http.MethodGet)
				}
				w.Write(emptyObj)
			}),
			dummy{},
		},
		{
			"foo",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if want := "/foo"; r.URL.Path != want {
					t.Errorf("got Request.URL: %s; want %s.", r.URL.Path, want)
				}
				w.Write(emptyObj)
			}),
			dummy{},
		},
		{
			"foo",
			http.MethodPost,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if c := r.Header.Get("Content-Type"); c != "application/json" {
					t.Errorf("got 'Content-Type' Header: '%s'; want 'application/json'.", c)
				}
				w.Write(emptyObj)
			}),
			dummy{},
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
			dummy{},
		},
		{
			"",
			http.MethodPost,
			dummy{"Testing"},
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Body == http.NoBody {
					t.Error("got Request.Body empty, want not empty.")
				}

				got, _ := ioutil.ReadAll(r.Body)

				if want := []byte(`{"name":"Testing"}`); !bytes.Contains(got, want) {
					t.Errorf("got body: %s, want %s.", got, want)
				}
				w.Write([]byte(`{"name":"Testing"}`))
			}),
			dummy{"Testing"},
		},
	}

	for _, tc := range testCases {
		var output dummy

		u, _ := url.Parse(tc.server.URL)
		c := New(u, nil)

		err := c.Request(context.Background(), tc.method, tc.endpoint, tc.body, &output)
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

					if wantBody := []byte(http.StatusText(http.StatusInternalServerError)); !bytes.Equal(wantBody, err.body) {
						t.Errorf("got Error.body: %s; want %s.", err.body, wantBody)
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
