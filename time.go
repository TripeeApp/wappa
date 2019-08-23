package wappa

import  "time"

const timeLayout = `"2006-01-02T15:04:05"`

// Time is a custom time type for
// unmashaling data from the API.
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var err error

	s := string(b)
	if s == "null" {
		return err
	}
	t.Time, err = time.Parse(timeLayout, s)
	return err
}

func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(timeLayout))
	return t.AppendFormat(b, timeLayout), nil
}
