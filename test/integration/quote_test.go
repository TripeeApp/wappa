package integration

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/rdleal/wappa"
)

func TestQuote(t *testing.T) {
	if !checkAuth("TestQuote") {
		return
	}

	empID, ok := findEmployeeID(employeeFilter)
	if !ok {
		return
	}

	f := wappa.Filter{
		"latOrigin": []string{fmt.Sprintf("%.7f", latOrigin)},
		"lngOrigin": []string{fmt.Sprintf("%.7f", lngOrigin)},
		"latDest": []string{fmt.Sprintf("%.7f", latDest)},
		"lngDest": []string{fmt.Sprintf("%.7f", lngDest)},
		"employee": []string{strconv.Itoa(empID)},
	}

	e, err := wpp.Quote.Estimate(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Quote.Estimate(%+v): %s; want nil.", f, err.Error())
	}

	if len(e.Categories) == 0 {
		t.Errorf("found no categories for the estimate: %+v.", f)
	}
}
