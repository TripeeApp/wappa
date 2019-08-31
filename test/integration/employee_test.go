package integration

import (
	"context"
	"fmt"
	"testing"
)

func TestEmployee(t *testing.T) {
	r, err := wpp.Employee.Read(context.Background(), nil)
	if err != nil {
		t.Fatalf("got error while calling Employee.Read(): %s; want nil.", err.Error())
	}

	if len(r.Employees) == 0 {
		t.Skip("No Employee found.")
	}

	employee := r.Employees[0]

	status, err := wpp.Employee.Status(context.Background(), employee.ID)
	if err != nil {
		t.Errorf("got error while calling Employee.Status(%d): %s; want nil.", employee.ID, err.Error())
	}

	fmt.Printf("%+v\n", status)

	res, err := wpp.Employee.LastRides(context.Background(), nil)
	if err != nil {
		t.Errorf("got error while calling Employee.LastRides(): %s; want nil.", err.Error())
	}

	if !res.Success {
		t.Errorf("got failed response: '%s'; want it to be successful.", res.Message)
	}
}
