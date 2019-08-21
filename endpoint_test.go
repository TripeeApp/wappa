package wappa

import(
	"testing"
	"net/url"
)

func TestEndpointAction(t *testing.T) {
	testCases := []struct{
		path endpoint
		action string
		want endpoint
	}{
		{endpoint("test"), list, "api/test/listar"},
		{endpoint("test"), create, "api/test/cadastrar"},
	}

	for _, tc := range testCases {
		if got := tc.path.Action(tc.action); tc.want != got {
			t.Errorf("go path from endpoint.Action('%s'): '%s'; want '%s'.", tc.action, got, tc.want)
		}
	}
}

func TestEndpointQuery(t *testing.T) {
	testCases := []struct{
		path endpoint
		vals url.Values
		want endpoint
	}{
		{
			"test",
			url.Values{
				"filter": []string{"123"},
			},
			"api/test?filter=123",
		},
		{
			"test",
			url.Values{
				"filter": []string{"123"},
				"name": []string{"Joao"},
			},
			"api/test?filter=123&name=Joao",
		},
	}

	for _, tc := range testCases {
		if got := tc.path.Query(tc.vals); tc.want != got {
			t.Errorf("got path from endpoint.Query(%+v): '%s'; want '%s'.", tc.vals, got, tc.want)
		}
	}
}

func TestEndpointString(t *testing.T) {
	testCases := []struct{
		path endpoint
		want string
	}{
		{
			endpoint("test").Action(list).Query(url.Values{"filter": []string{"123"}}),
			"api/test/listar?filter=123",
		},
	}

	for _, tc := range testCases {
		if got := tc.path.String(); got != tc.want {
			t.Errorf("got path from endpoint.String(): '%s'; want '%s'.", got, tc.want)
		}
	}
}
