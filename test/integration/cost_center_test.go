package integration

import(
	"context"
	"testing"

	"github.com/rdleal/wappa"
)

func TestCostCenter(t *testing.T) {
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
		t.Fatalf("got error while reading a cost center: %s; want nil.", c.Message)
	}

	costCenter := c.Response[0]

	newCostCenter := &wappa.CostCenter{ID: costCenter.ID, Code: randString(5, letterBytes)}

	opres, err := wpp.CostCenter.Update(context.Background(), newCostCenter)
	if err != nil {
		t.Fatalf("got error while calling CostCenter.Update(%+v): %s; want nil.", newCostCenter, err.Error())
	}

	if !opres.Success {
		t.Fatalf("got error while updating a cost center: %s; want nil.", opres.Message)
	}

	f = wappa.Filter{"code": newCostCenter.Code}
	c, err = wpp.CostCenter.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling CostCenter.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !c.Success {
		t.Fatalf("got error while reading a cost center: %s; want nil.", res.Message)
	}

	opres, err = wpp.CostCenter.Inactivate(context.Background(), costCenter.ID)
	if err != nil {
		t.Fatalf("got error while calling CostCenter.Inactivate(%+v): %s; want nil.", f, err.Error())
	}

	if !opres.Success {
		t.Fatalf("got error while inactivating a cost center: %s; want nil.", opres.Message)
	}
}
