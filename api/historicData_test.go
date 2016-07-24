package api

import (
	"appengine"
	"appengine/aetest"
	"log"
	"testing"

	"github.com/eqinox76/ownFolio/data"
)

func TestGetInstrument(t *testing.T) {
	called := 0
	getData = func(_ appengine.Context, id string) (data.DataPoints, error) {
		log.Printf("overwritten method called with id %s", id)
		var i data.DataPoints
		called++
		return i, nil
	}

	ctx, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}

	defer ctx.Close()

	GetInstrument(ctx, "DAX")

	if called != 1 {
		t.Fatal("Calling once produced %i called", called)
	}

	GetInstrument(ctx, "DAX")
	GetInstrument(ctx, "DAX")
	GetInstrument(ctx, "DAX")

	if called != 1 {
		t.Fatal("Calling once produced %i called", called)
	}
}
