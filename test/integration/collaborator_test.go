package integration

import(
	"context"
	"fmt"
	"testing"

	"bitbucket.org/mobilitee/wappa"
)

func TestCollaborator(t *testing.T) {
	enrollment := randString(5, numberBytes)
	res, err := wpp.Collaborator.Create(context.Background(), &wappa.Collaborator{
		Name: randString(10, letterBytes),
		Enrollment: enrollment,
		Email: fmt.Sprintf("%s@fakemail.com", randString(8, letterBytes)),
		CostCenterID: createCostCenter(t),
		RoleID: createRole(t),
	})
	if err != nil {
		t.Fatalf("got error while calling Collaborator.Create(): %s; want nil.", err.Error())
	}

	if !res.Success {
		t.Fatalf("got error while creating a collaborator: %s; want nil.", res.Message)
	}

	f := wappa.Filter{"enrollment": enrollment}
	c, err := wpp.Collaborator.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Collaborator.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !c.Success {
		t.Fatalf("got error while reading a collaborator: %s; want nil.", res.Message)
	}

	collaborator := c.Response[0]

	newCollaborator := &wappa.Collaborator{ID: collaborator.ID, Enrollment: randString(5, numberBytes)}

	opres, err := wpp.Collaborator.Update(context.Background(), newCollaborator)
	if err != nil {
		t.Fatalf("got error while calling Collaborator.Update(%+v): %s; want nil.", newCollaborator, err.Error())
	}

	if !opres.Success {
		t.Fatalf("got error while updating a collaborator: %s; want nil.", opres.Message)
	}

	f = wappa.Filter{"enrollment": newCollaborator.Enrollment}
	c, err = wpp.Collaborator.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Collaborator.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !c.Success {
		t.Fatalf("got error while reading a collaborator: %s; want nil.", res.Message)
	}

	opres, err = wpp.Collaborator.Activate(context.Background(), collaborator.ID)
	if err != nil {
		t.Fatalf("got error while calling Collaborator.Activate(%d): %s; want nil.", collaborator.ID, err.Error())
	}

	if !opres.Success {
		t.Fatalf("got error while activating a collaborator: %s; want nil.", opres.Message)
	}

	opres, err = wpp.Collaborator.Inactivate(context.Background(), collaborator.ID)
	if err != nil {
		t.Fatalf("got error while calling Collaborator.Inactivate(%d): %s; want nil.", collaborator.ID, err.Error())
	}

	if !opres.Success {
		t.Fatalf("got error while inactivating a collaborator: %s; want nil.", opres.Message)
	}
}

func createCostCenter(t *testing.T) int {
	code := randString(5, letterBytes)
	res, err := wpp.CostCenter.Create(context.Background(), &wappa.CostCenter{Name: "test", Code: code, CNPJ: "27906045000186"})
	if err != nil {
		t.Fatalf("got error while calling CostCenter.Create(): %s; want nil.", err.Error())
	}

	if !res.Success {
		t.Fatalf("got error while creating a cost center: %s; want nil.", res.Message)
	}

	f := wappa.Filter{"code": code}
	c, err := wpp.CostCenter.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling CostCenter.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !c.Success {
		t.Fatalf("got error while reading a cost center: %s; want nil.", res.Message)
	}

	return c.Response[0].ID
}

func createRole(t *testing.T) int {
	desc := randString(5, numberBytes)
	res, err := wpp.Role.Create(context.Background(), &wappa.Role{Description: desc})
	if err != nil {
		t.Fatalf("got error while calling Role.Create(): %s; want nil.", err.Error())
	}

	if !res.Success {
		t.Fatalf("got error while creating a role: %s; want nil.", res.Message)
	}

	f := wappa.Filter{"desc": desc}
	r, err := wpp.Role.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Role.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !r.Success {
		t.Fatalf("got error while reading a role: %s; want nil.", res.Message)
	}

	return r.Response[0].ID
}
