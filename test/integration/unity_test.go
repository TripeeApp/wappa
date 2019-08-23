package integration

import(
	"context"
	"testing"

	"bitbucket.org/mobilitee/wappa"
)

func TestUnity(t *testing.T) {
	code := randString(5, numberBytes)
	res, err := wpp.Unity.Create(context.Background(), &wappa.Unity{CustomerCompany: 123, Code: code})
	if err != nil {
		t.Fatalf("got error while calling Unity.Create(): %s; want nil.", err.Error())
	}

	if !res.Success {
		t.Fatalf("got error while creating an unity: %s; want nil.", res.Message)
	}

	f := wappa.Filter{"code": code}
	u, err := wpp.Unity.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Unity.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !u.Success {
		t.Fatalf("got error while reading an unity: %s; want nil.", res.Message)
	}

	unity := u.Response[0]

	newUnity := &wappa.Unity{ID: unity.ID, Code: randString(5, numberBytes)}

	opres, err := wpp.Unity.Update(context.Background(), newUnity)
	if err != nil {
		t.Fatalf("got error while calling Unity.Update(%+v): %s; want nil.", newUnity, err.Error())
	}

	if !opres.Success {
		t.Fatalf("got error while updating an unity: %s; want nil.", opres.Message)
	}

	f = wappa.Filter{"code": newUnity.Code}
	u, err = wpp.Unity.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Unity.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !u.Success {
		t.Fatalf("got error while reading an unity: %s; want nil.", res.Message)
	}

	opres, err = wpp.Unity.Inactivate(context.Background(), unity.ID)
	if err != nil {
		t.Fatalf("got error while calling Unity.Inactivate(%+v): %s; want nil.", f, err.Error())
	}

	if !opres.Success {
		t.Fatalf("got error while inactivating an unity: %s; want nil.", opres.Message)
	}
}
