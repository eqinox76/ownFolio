package api

import (
	"log"
	"testing"
	"appengine/aetest"

	"eqinox76/ownFolio/data"
)

func TestGetInstrument(t *testing.T){
	called := 0
	getData = func (id string) (data.Instrument, error){
		log.Printf("overwritten method called with id %s", id)
		var i data.Instrument
		called++
		return i, nil
	}

	ctx, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}

	defer ctx.Close()

	GetInstrument("DAX", ctx)

	if called != 1{
		t.Fatal("Calling once produced %i called", called)
	}

	GetInstrument("DAX", ctx)
	GetInstrument("DAX", ctx)
	GetInstrument("DAX", ctx)

	if called != 1{
		t.Fatal("Calling once produced %i called", called)
	}
}
