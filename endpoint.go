package wappa

import(
	"fmt"
	"net/url"
	"strings"
)

const (
	list		= `listar`
	create		= `cadastrar`
	edit		= `alterar`
	activate	= `ativar`
	inactivate	= `inativar`
)

// endpoint represents the path to the resource.
type endpoint string

func (e endpoint) Action(a string) endpoint {
	if strings.Index(e.String(), "api") != 0 {
		e = e.prefix()
	}
	return endpoint(fmt.Sprintf("%s/%s", e, a))
}

func (e endpoint) Query(v url.Values) endpoint {
	if strings.Index(e.String(), "api") != 0 {
		e = e.prefix()
	}
	return endpoint(fmt.Sprintf("%s?%s", e, v.Encode()))
}

func (e endpoint) prefix() endpoint {
	return endpoint(fmt.Sprintf("api/%s", e))
}

func (e endpoint) String() string {
	return string(e)
}


