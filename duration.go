package wappa

import (
	"encoding/json"
	"time"
)

// DurationSec is a custom seconds duration type for
// unmashaling data from the API.
type DurationSec struct {
	time.Duration
}

func (d *DurationSec) UnmarshalJSON(b []byte) (err error) {
	d.Duration, err = parseDuration(b, time.Second)
	return
}

// DurationMin is a custom minutes duration type for
// unmashaling data from the API.
type DurationMin struct {
	time.Duration
}

func (d *DurationMin) UnmarshalJSON(b []byte) (err error) {
	d.Duration, err = parseDuration(b, time.Minute)
	return
}

func parseDuration(b []byte, unit time.Duration) (time.Duration, error) {
	var d time.Duration

	s := string(b)
	// By convention, to approximate the behaviour of Unmarshal itself,
	// Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op.
	if s == "null" {
		return d, nil
	}
	n, err := json.Number(s).Float64()

	d = time.Duration(n) * unit

	return d, err
}
