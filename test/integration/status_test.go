package integration

import (
	"context"
	"testing"
)

func TestStatus(t *testing.T) {
	ok, err := wpp.Status(context.Background())
	if err != nil {
		t.Fatalf("got error calling Wappa.Status(): '%s'; want nil.", err.Error())
	}

	if !ok {
		t.Error("got api status not OK; want it to be OK.")
	}
}
