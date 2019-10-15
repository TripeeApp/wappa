package wappa

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Duration is a time duration represented as hh:mm:ss.
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var dur string
	if err = json.Unmarshal(b, &dur); err != nil {
		return
	}

	parts := strings.Split(dur, ":")
	if len(parts) < 3 {
		return nil
	}

	d.Duration, err = time.ParseDuration(fmt.Sprintf("%sh%sm%ss", parts[0], parts[1], parts[2]))
	return nil
}

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
	var f float64
	var d time.Duration

	err := json.Unmarshal(b, &f)

	d = time.Duration(f) * unit

	return d, err
}
