package integration

import(
	"context"
	"testing"

	"bitbucket.org/mobilitee/wappa"
)

func TestRole(t *testing.T) {
	desc := randString(5, numberBytes)
	res, err := wpp.Role.Create(context.Background(), &wappa.Role{Description: desc})
	if err != nil {
		t.Fatalf("got error while calling Role.Read(): %s; want nil.", err.Error())
	}

	if !res.Success {
		t.Fatalf("got error while creating an role: %s; want nil.", res.Message)
	}

	f := wappa.Filter{"desc": desc}
	r, err := wpp.Role.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Role.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !r.Success {
		t.Fatalf("got error while reading an role: %s; want nil.", res.Message)
	}
}
