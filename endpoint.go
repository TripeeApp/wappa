package wappa

import (
	"fmt"
	"net/url"
	"strings"
)

// Actions
const (
	nearby             = `nearby`
	status             = `status`
	lastRides          = `last-rides`
	cancellationReason = `cancellation-reason`
	employee           = `employee`
	cancel             = `cancel`
	rate               = `rate`
	update             = `update`
	deactivate         = `deactivate`
	activate           = `activate`
	qrcode             = `qrcode`
)

// endpoint represents the path to the resource.
type endpoint string

// Action returns a new endpoint suffixed with the action.
func (e endpoint) Action(a string) endpoint {
	return endpoint(fmt.Sprintf("%s/%s", e, a))
}

// Query returns a new endpoint with the query parameters appended.
func (e endpoint) Query(v url.Values) endpoint {
	return endpoint(fmt.Sprintf("%s?%s", e, v.Encode()))
}

// String returns a string of endpoint type
func (e endpoint) String() string {
	s := string(e)
	if strings.Index(s, "api") != 0 {
		s = fmt.Sprintf("api/%s", s)
	}
	return s
}
