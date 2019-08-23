package integration

import(
	"context"
	"testing"

	"bitbucket.org/mobilitee/wappa"
)

func TestRole(t *testing.T) {
	desc := randString(5, letterBytes)
	res, err := wpp.Role.Create(context.Background(), &wappa.Role{ID: 123, Description: desc})
	if err != nil {
		t.Fatalf("got error while calling Role.Read(): %s; want nil.", err.Error())
	}

	if !res.Success {
		t.Errorf("got error while creating a role: %s; want nil.", res.Message)
	}

	f := wappa.Filter{"desc": desc}
	role, err := wpp.Role.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Role.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !role.Success {
		t.Errorf("got error while reading a role: %s; want nil.", res.Message)
	}
}
