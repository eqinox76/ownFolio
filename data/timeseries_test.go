package data

import (
	"testing"
	"time"
	"encoding/json"
	"reflect"
)

func TestC3DataJson(t *testing.T){
	var instrument Instrument
	instrument.Add(&DataPoint{
		time.Date(2010, 11, 21, 0, 0, 0, 0, time.UTC),
		22,
		33,
		44,
		5500,
	})
	instrument.Add(&DataPoint{
		time.Date(2010, 1, 21, 0, 0, 0, 0, time.UTC),
		2.2,
		3.3,
		4.4,
		100,
	})
	data := instrument.C3Data()
	if len(data.Date) != 2{
		t.Errorf("%s dates", len(data.Date))
	}
	jsonData, _ := json.Marshal(&data)
	expected := []uint8("{\"Date\":[\"2010-11-21\",\"2010-01-21\"],\"High\":[22,2.2],\"Low\":[33,3.3],\"Close\":[44,4.4]}")
	if ! reflect.DeepEqual(expected, jsonData){
		t.Errorf("json did not match expected:\n%s\nbut got\n%s\n", expected, jsonData)
	}
}
