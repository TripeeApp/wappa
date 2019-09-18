package integration

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	wappa "github.com/mobilitee-smartmob/wappa/v2"
)

func TestDriver(t *testing.T) {
	if !checkAuth("TestDriver") {
		return
	}

	empID, ok := findEmployeeID(employeeFilter)
	if !ok {
		return
	}

	f := wappa.Filter{
		"lat":      []string{fmt.Sprintf("%.7f", latOrigin)},
		"lng":      []string{fmt.Sprintf("%.7f", lngOrigin)},
		"type":     []string{"1"}, // TODO: Get from Env and defaults to 6 (Carro Particular)
		"employee": []string{strconv.Itoa(empID)},
	}

	d, err := wpp.Driver.Nearby(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Driver.Nearby(%+v): %s; want nil.", f, err.Error())
	}

	if !d.Success {
		t.Errorf("got failed response while reading nearby drivers: '%s'; want it to be successful.", d.Message)
	}

	if len(d.Drivers) == 0 {
		fmt.Printf("no drivers for location (%.7f, %.7f).", latOrigin, lngOrigin)
	}
}
