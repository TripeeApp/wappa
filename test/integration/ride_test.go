package integration

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/mobilitee-smartmob/wappa"
)

func TestRide(t *testing.T) {
	if !checkAuth("TestRide") {
		return
	}

	empID, ok := findEmployeeID(employeeFilter)
	if !ok {
		return
	}

	qr, err := wpp.Ride.QRCode(context.Background(), wappa.Filter{"employee": []string{strconv.Itoa(empID)}})
	if err != nil {
		t.Fatalf("got error while calling Ride.QRCode(): %s; want nil.", err.Error())
	}

	if !qr.Success {
		t.Errorf("got failed response while creating QR code: '%s'; want it to be successful.", qr.Message)
	}

	fmt.Println(qr.QRCode)

	r, err := wpp.Ride.CancellationReason(context.Background())
	if err != nil {
		t.Fatalf("got error while calling Ride.CancellationReasion(): %s; want nil.", err.Error())
	}

	if len(r.Reasons) == 0 {
		t.Fatal("got empty cancellation reasons; want not empty")
	}

	reason := r.Reasons[0].ID

	ride := &wappa.Ride{
		Employee:       empID,
		TaxiType:       1,
		TaxiCategoryId: 1,
		LatOrigin:      latOrigin,
		LngOrigin:      lngOrigin,
		LatDestiny:     latDest,
		LngDestiny:     lngDest,
		ExternalID:     randString(5, numberBytes),
	}

	newRide, err := wpp.Ride.Create(context.Background(), ride)
	if err != nil {
		t.Fatalf("got error while calling Ride.Create(%+v): %s; want nil.", newRide, err.Error())
	}

	if !newRide.Success {
		t.Errorf("got failed response while creating a ride: '%s'; want it to be successful.", newRide.Message)
	}

	cancelled, err := wpp.Ride.Cancel(context.Background(), newRide.ID, reason)
	if err != nil {
		t.Fatalf("got error while calling Ride.Cancel(%d): %s; want nil.", newRide.ID, err.Error())
	}

	if !cancelled.Success {
		t.Errorf("got failed response while cancelling a ride: '%s'; want it to be successful.", cancelled.Message)
	}

	f := wappa.Filter{"id": []string{strconv.Itoa(newRide.ID)}}
	rd, err := wpp.Ride.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Ride.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !rd.Success {
		t.Errorf("got failed response while reading a ride: '%s'; want it to be successful.", rd.Message)
	}

	if s := rd.Info.Status; s != wappa.RideStatusCancelled {
		t.Errorf("got RideInfo.Status: %s; want 'ride-canceled'.", s)
	}

	if cb := rd.Info.CancelledBy; cb != wappa.RideCancelledByUser {
		t.Errorf("got RideInfo.CancelleBy: %s; want '1'.", cb)
	}
}
