package integration

import (
	"context"
	"strconv"
	"testing"

	wappa "github.com/mobilitee-smartmob/wappa/v2"
)

func TestEmployee(t *testing.T) {
	if !checkAuth("TestEmployee") {
		return
	}

	r, err := wpp.Employee.Read(context.Background(), employeeFilter)
	if err != nil {
		t.Fatalf("got error while calling Employee.Read(%+v): '%s'; want nil.", employeeFilter, err.Error())
	}

	if len(r.Employees) == 0 {
		t.Skip("No Employee found.")
	}

	employee := r.Employees[0]

	_, err = wpp.Employee.Status(context.Background(), employee.ID)
	if err != nil {
		t.Errorf("got error while calling Employee.Status(%d): '%s'; want nil.", employee.ID, err.Error())
	}

	res, err := wpp.Employee.LastRides(context.Background(), wappa.Filter{"employee": []string{strconv.Itoa(employee.ID)}})
	if err != nil {
		t.Errorf("got error while calling Employee.LastRides(): '%s'; want nil.", err.Error())
	}

	if !res.Success {
		t.Errorf("got failed response while reading employee's last rides: '%s'; want it to be successful.", res.Message)
	}
}
