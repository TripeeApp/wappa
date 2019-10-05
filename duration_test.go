package wappa

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestDurationSecUnmarshal(t *testing.T) {
	testCases := []struct {
		payload []byte
		want    DurationSec
	}{
		{[]byte("250.3"), DurationSec{time.Duration(250) * time.Second}},
		{[]byte("null"), DurationSec{}},
	}

	for _, tc := range testCases {
		got := DurationSec{}
		if err := json.Unmarshal(tc.payload, &got); err != nil {
			t.Fatalf("got error calling json.Unmarshal(%s): '%s'; want nil.", tc.payload, err.Error())
		}

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got DurationSec %+v; want %+v.", got, tc.want)
		}
	}
}

func TestDurationSecMarshalError(t *testing.T) {
	payload := []byte(`"abc"`)

	if err := json.Unmarshal(payload, &DurationSec{}); err == nil {
		t.Fatal("got error nil; want it not nil.")
	}
}

func TestDurationMinUnmarshal(t *testing.T) {
	testCases := []struct {
		payload []byte
		want    DurationMin
	}{
		{[]byte("14.5"), DurationMin{time.Duration(14) * time.Minute}},
		{[]byte("null"), DurationMin{}},
	}

	for _, tc := range testCases {
		got := DurationMin{}
		if err := json.Unmarshal(tc.payload, &got); err != nil {
			t.Fatalf("got error calling json.Unmarshal(%s): '%s'; want nil.", tc.payload, err.Error())
		}

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got DurationMin %+v; want %+v.", got, tc.want)
		}
	}
}

func TestDurationMinMarshalError(t *testing.T) {
	payload := []byte(`"abc"`)

	if err := json.Unmarshal(payload, &DurationMin{}); err == nil {
		t.Fatal("got error nil; want it not nil.")
	}
}
