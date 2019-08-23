package wappa

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestTimeUnmarshal(t *testing.T) {
	testCases := []struct{
		payload []byte
	}{
		{[]byte(`"2019-08-23T19:00:13"`)},
		{[]byte("null")},
	}

	for _, tc := range testCases {
		got := &Time{time.Time{}}

		if err := json.Unmarshal(tc.payload, got); err != nil {
			t.Fatalf("error while calling Time.Unmarshal(%s): '%s'; want nil.", tc.payload, err.Error())
		}

		tm, _ := time.Parse(timeLayout, string(tc.payload))
		if want := (&Time{tm}); !reflect.DeepEqual(got, want) {
			t.Errorf("got Time %+v; want %+v.", got, want)
		}
	}
}

func TestTimeMarshal(t *testing.T) {
	now := time.Now()
	testCases := []struct{
		tm Time
		want string
	}{
		{Time{now}, string(now.Format(timeLayout))},
		{Time{time.Time{}}, "null"},
	}

	for _, tc := range testCases {
		b, err := json.Marshal(tc.tm)
		if err != nil {
			t.Fatalf("error while calling Time.Marshal(%+v): '%s'; want nil.", tc.tm, err.Error())
		}

		if got := string(b); got != tc.want {
			t.Errorf("got Time %s; want %s.", got, tc.want)
		}
	}
}
