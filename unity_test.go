package wappa

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

type testRequester struct{
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

func TestUnity(t *testing.T) {
	testCases := []struct{
		name	string
		call	func(ctx context.Context, req requester) (resp interface{}, err error)
		ctx	context.Context
		method	string
		path	endpoint
		body	interface{}
		wantRes	interface{}
	}{
		{
			"Read()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UnityService{req}).Read(ctx, "123", "test")
				return
			},
			context.Background(),
			http.MethodGet,
			unityEndpoint.Action(read).Query(url.Values{"idUnidade": []string{"123"}, "codDescricao": []string{"test"}}),
			nil,
			&UnityResponse{
				DefaultResponse: DefaultResponse{Success: true},
				Response: []*Unity{
					&Unity{ID:123},
				},
			},
		},
		{
			"Create()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UnityService{req}).Create(ctx, &Unity{Code: "123"})
				return
			},
			context.Background(),
			http.MethodPost,
			unityEndpoint.Action(create),
			&Unity{Code: "123"},
			&DefaultResponse{
				Success: true,
			},
		},
		{
			"Update()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UnityService{req}).Update(ctx, &Unity{Code: "123"})
				return
			},
			context.Background(),
			http.MethodPost,
			unityEndpoint.Action(update),
			&Unity{Code: "123"},
			&DefaultResponse{
				Success: true,
			},
		},
		{
			"Inactivate()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UnityService{req}).Inactivate(ctx, "123")
				return
			},
			context.Background(),
			http.MethodPost,
			unityEndpoint.Action(inactivate).Query(url.Values{"idUnidade": []string{"123"}}),
			nil,
			&DefaultResponse{
				Success: true,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // creates scoped test case
		t.Run(tc.name, func(t *testing.T) {
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
		})
	}
}

func TestUnityError(t *testing.T) {
	testCases := []struct{
		name	string
		call	func(req requester) error
		err	error
	}{
		{
			"Read()",
			func(req requester) error {
				_, err := (&UnityService{req}).Read(context.Background(), "123", "test")
				return err
			},
			errors.New("Error"),
		},
		{
			"Create()",
			func(req requester) error {
				_, err := (&UnityService{req}).Create(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Update()",
			func(req requester) error {
				_, err := (&UnityService{req}).Update(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"ReadClassifier()",
			func(req requester) error {
				_, err := (&UnityService{req}).Inactivate(context.Background(), "1")
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		tc := tc // creates scoped test case
		t.Run(tc.name, func(t *testing.T) {
			req := &testRequester{err: tc.err}

			err := tc.call(req)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("got error: %s; want %s.", err, tc.err)
			}
		})
	}
}
