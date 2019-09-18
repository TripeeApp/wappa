package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/mobilitee-smartmob/wappa"
)

func TestRide(t *testing.T) {
	if !checkAuth("TestRide") {
		return
	}

	// Setting up the webhook.
	wh, svr := webhook(t)

	empID, ok := findEmployeeID(employeeFilter)
	if !ok {
		return
	}

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

	fmt.Println("Waiting for Webhook response.")

	ok, got, want := checkWebhookStatuses(wh, []string{
		wappa.RideStatusSearchingForDriver,
		wappa.RideStatusDriverFound,
		wappa.RideStatusInProgress,
		wappa.RideStatusCompleted,
	})
	svr.Shutdown(context.Background())
	if !ok {
		t.Errorf("got ride status: '%s'; want '%s'.", got, want)

		cancelRide(t, newRide.ID)
	}

	f := wappa.Filter{"id": []string{strconv.Itoa(newRide.ID)}}
	rd, err := wpp.Ride.Read(context.Background(), f)
	if err != nil {
		t.Fatalf("got error while calling Ride.Read(%+v): %s; want nil.", f, err.Error())
	}

	if !rd.Success {
		t.Errorf("got failed response while reading a ride: '%s'; want it to be successful.", rd.Message)
	}

	var wantStatus = wappa.RideStatusCompleted
	if !ok {
		wantStatus = wappa.RideStatusCancelled
	}

	if got := rd.Info.Status; got != wantStatus {
		t.Errorf("got RideInfo.Status: %s; want '%s'.", got, wantStatus)
	}
}

//func TestRideCancelledByUser(t *testing.T) {
//	if !checkAuth("TestRide") {
//		return
//	}
//
//	// Setting up the webhook.
//	wh, svr := webhook(t)
//
//	empID, ok := findEmployeeID(employeeFilter)
//	if !ok {
//		return
//	}
//
//	ride := &wappa.Ride{
//		Employee:       empID,
//		TaxiType:       1,
//		TaxiCategoryId: 1,
//		LatOrigin:      latOrigin,
//		LngOrigin:      lngOrigin,
//		LatDestiny:     latDest,
//		LngDestiny:     lngDest,
//		ExternalID:     randString(5, numberBytes),
//	}
//
//	newRide, err := wpp.Ride.Create(context.Background(), ride)
//	if err != nil {
//		t.Fatalf("got error while calling Ride.Create(%+v): %s; want nil.", newRide, err.Error())
//	}
//
//	if !newRide.Success {
//		t.Errorf("got failed response while creating a ride: '%s'; want it to be successful.", newRide.Message)
//	}
//
//	// User cancels the ride after 500 milliseconds.
//	go func() {
//		time.Sleep(1 * time.Second)
//		cancelRide(t, newRide.ID)
//	}()
//
//	fmt.Println("Waiting for Webhook response.")
//
//	ok, got, want := checkWebhookStatuses(wh, []string{
//		wappa.RideStatusCancelled,
//	})
//	svr.Shutdown(context.Background())
//	// Status flow not ok. Try to cancel the ride.
//	if !ok {
//		t.Fatalf("got ride status: '%s'; want '%s'.", got, want)
//	}
//
//	f := wappa.Filter{"id": []string{strconv.Itoa(newRide.ID)}}
//	rd, err := wpp.Ride.Read(context.Background(), f)
//	if err != nil {
//		t.Fatalf("got error while calling Ride.Read(%+v): %s; want nil.", f, err.Error())
//	}
//
//	if !rd.Success {
//		t.Errorf("got failed response while reading a ride: '%s'; want it to be successful.", rd.Message)
//	}
//
//	if got := rd.Info.Status; got != wappa.RideStatusCancelled {
//		t.Errorf("got RideInfo.Status: %s; want '%s'.", got, wappa.RideStatusCancelled)
//	}
//}
//
//func TestRideCancelledByDriver(t *testing.T) {
//	if !checkAuth("TestRide") {
//		return
//	}
//
//	// Setting up the webhook.
//	wh, svr := webhook(t)
//
//	empID, ok := findEmployeeID(employeeFilter)
//	if !ok {
//		return
//	}
//
//	ride := &wappa.Ride{
//		Employee:       empID,
//		TaxiType:       1,
//		TaxiCategoryId: 1,
//		LatOrigin:      latOrigin,
//		LngOrigin:      lngOrigin,
//		LatDestiny:     latDest,
//		LngDestiny:     lngDest,
//		ExternalID:     randString(5, numberBytes),
//	}
//
//	newRide, err := wpp.Ride.Create(context.Background(), ride)
//	if err != nil {
//		t.Fatalf("got error while calling Ride.Create(%+v): %s; want nil.", newRide, err.Error())
//	}
//
//	if !newRide.Success {
//		t.Errorf("got failed response while creating a ride: '%s'; want it to be successful.", newRide.Message)
//	}
//
//	fmt.Println("Waiting for driver to cancel.")
//
//	ok, got, want := checkWebhookStatuses(wh, []string{
//		wappa.RideStatusCancelled,
//	})
//	svr.Shutdown(context.Background())
//	// Status flow not ok. Try to cancel the ride.
//	if !ok {
//		t.Fatalf("got ride status: '%s'; want '%s'.", got, want)
//	}
//
//	f := wappa.Filter{"id": []string{strconv.Itoa(newRide.ID)}}
//	rd, err := wpp.Ride.Read(context.Background(), f)
//	if err != nil {
//		t.Fatalf("got error while calling Ride.Read(%+v): %s; want nil.", f, err.Error())
//	}
//
//	if !rd.Success {
//		t.Errorf("got failed response while reading a ride: '%s'; want it to be successful.", rd.Message)
//	}
//
//	if got := rd.Info.Status; got != wappa.RideStatusCancelled {
//		t.Errorf("got RideInfo.Status: %s; want '%s'.", got, wappa.RideStatusCancelled)
//	}
//}
//
//expectedStatuses := []string{
//  RideStatusSearchingForDriver = "searching-for-driver"
//  RideStatusDriverNotFound     = "driver-not-found"
//  RideStatusCancelled          = "ride-cancelled"
//  RideStatusDriverFound        = "driver-found"
//  RideStatusWaitingForDriver   = "waiting-for-driver"
//  RideStatusInProgress         = "on-ride"
//  RideStatusPaid               = "ride-paid"
//  RideStatusCompleted          = "ride-completed"
//}

//	if cb := rd.Info.CancelledBy; cb != wappa.RideCancelledByUser {
//		t.Errorf("got RideInfo.CancelleBy: %s; want '1'.", cb)
//	}

func webhook(t *testing.T) (<-chan []byte, *http.Server) {
	host, ok := os.LookupEnv(envKeyWappaWebhookHost)
	if !ok {
		t.Skipf(msgEnvMissing, envKeyWappaWebhookHost)
	}

	port, ok := os.LookupEnv(envKeyWappaWebhookPort)
	if !ok {
		t.Skipf(msgEnvMissing, envKeyWappaWebhookPort)
	}

	r, err := wpp.Webhook.Read(context.Background())
	if err != nil {
		t.Fatalf("got error calling Webhook.Read(): '%s'; want nil.", err.Error())
	}

	if !r.Success {
		t.Fatalf("got failed response while reading Webhooks: '%s'; want it to be successful.", r.Message)
	}

	// New webhook endpoint
	endpoint := fmt.Sprintf("%s", randString(5, letterBytes))
	wh := &wappa.Webhook{
		URL:      host,
		Endpoint: endpoint,
		AuthKey:  "auth-key",
	}

	var curWebhook *wappa.Webhook
	// If Webhook already exists, updates it. Otherwise creates it.
	if len(r.Listeners) > 0 {
		curWebhook = r.Listeners[0]

		updateWebhook(t, wh)

		if !wh.Active {
			activateWebhook(t)
		}
	} else {
		createWebhook(t, wh)
	}

	ch := make(chan []byte)

	svr := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	http.HandleFunc("/"+endpoint, func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("go error while reading body %+v: %s.", r.Body, err.Error())
		}

		fmt.Printf("received payload: '%s'.\n", b)

		ch <- b
	})

	go func() {
		defer func() {
			// Replaces the newly created webhook
			// with the old one, if it exists.
			if curWebhook != nil {
				updateWebhook(t, curWebhook)

				if !curWebhook.Active {
					deactivateWebhook(t)
				}
			}
			close(ch)
		}()

		if err := svr.ListenAndServe(); err != http.ErrServerClosed {
			t.Fatalf("got error while starting server at localhost:%s: %s.", port, err.Error())
		}
	}()

	// Blocks until server is ready.
	for {
		_, err := http.Get(fmt.Sprintf("http://localhost:%s/ping", port))
		if err == nil {
			break
		}
	}

	fmt.Printf("Tunnelling at '%s:%s/%s' ready.\n", host, port, endpoint)

	return ch, svr
}

func updateWebhook(t *testing.T, wh *wappa.Webhook) {
	res, err := wpp.Webhook.Update(context.Background(), wh)
	if err != nil {
		t.Fatalf("got error calling Webhook.Update(%+v): %s; want nil.", wh, err.Error())
	}

	if !res.Success {
		t.Fatalf("got failed response while updating Webhooks (%+v): '%s'; want it to be successful.", wh, res.Message)
	}

	// Wait for the webhook url to be synched before returning.
	for {
		r, err := wpp.Webhook.Read(context.Background())
		if err != nil {
			t.Fatalf("got error calling Webhook.Read(): '%s'; want nil.", err.Error())
		}

		if !r.Success {
			t.Fatalf("got failed response while reading Webhooks: '%s'; want it to be successful.", r.Message)
		}

		if l := r.Listeners; len(l) > 0 && l[0].Endpoint == wh.Endpoint {
			return
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func createWebhook(t *testing.T, wh *wappa.Webhook) {
	res, err := wpp.Webhook.Create(context.Background(), wh)
	if err != nil {
		t.Fatalf("got error calling Webhook.Create(%+v): '%s'; want nil.", wh, err.Error())
	}

	if !res.Success {
		t.Fatalf("got failed response while creating Webhooks (%+v): '%s'; want it to be successful.", wh, res.Message)
	}
}

func cancelRide(t *testing.T, rideID int) {
	r, err := wpp.Ride.CancellationReason(context.Background())
	if err != nil {
		t.Fatalf("got error while calling Ride.CancellationReasion(): %s; want nil.", err.Error())
	}

	if len(r.Reasons) == 0 {
		t.Fatal("got empty cancellation reasons; want not empty")
	}

	reason := r.Reasons[0].ID

	cancelled, err := wpp.Ride.Cancel(context.Background(), rideID, reason)
	if err != nil {
		t.Errorf("got error while calling Ride.Cancel(%d): %s; want nil.", rideID, err.Error())
	}

	if !cancelled.Success {
		t.Errorf("got failed response while cancelling a ride: '%s'; want it to be successful.", cancelled.Message)
	}
}

func checkWebhookStatuses(c <-chan []byte, statuses []string) (ok bool, got, want string) {
	var i int
	for {
		want = statuses[i]

		select {
		case <-time.After(time.Minute):
			return
		case body := <-c:
			wr := &wappa.WebhookRide{}
			if err := json.Unmarshal(body, wr); err != nil {
				got = err.Error()
				return
			}
			fmt.Printf("Status arrived by webhook: '%s'; want '%s'.\n", wr.Status, want)

			if got = wr.Status; got == want {
				if i == len(statuses)-1 {
					ok = true
					return
				}
				i += 1
			}
		}
	}
	return
}
