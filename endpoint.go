package wappa

import(
	"fmt"
	"net/url"
	"strings"
)

// Actions
const (
	read		= `listar`
	create		= `cadastrar`
	update		= `alterar`
	activate	= `ativar`
	inactivate	= `inativar`
)

// endpoint represents the path to the resource.
type endpoint string

// Action returns a new endpoint suffixed with the action.
func (e endpoint) Action(a string) endpoint {
	if strings.Index(e.String(), "api") != 0 {
		e = e.prefix()
	}
	return endpoint(fmt.Sprintf("%s/%s", e, a))
}

// Query returns a new endpoint with the query parameters appended.
func (e endpoint) Query(v url.Values) endpoint {
	if strings.Index(e.String(), "api") != 0 {
		e = e.prefix()
	}
	return endpoint(fmt.Sprintf("%s?%s", e, v.Encode()))
}

// String returns a string of endpoint type
func (e endpoint) String() string {
	return string(e)
}

func (e endpoint) prefix() endpoint {
	return endpoint(fmt.Sprintf("api/%s", e))
}



