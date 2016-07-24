package data

import (
	"testing"
	"time"
	"encoding/json"
	"reflect"
)

func TestC3DataJson(t *testing.T){
	var timeseries TimeSeries
	timeseries.Add(&DataPoint{
		time.Date(2010, 11, 21, 0, 0, 0, 0, time.UTC),
		22,
		33,
		44,
		5500,
	})
	timeseries.Add(&DataPoint{
		time.Date(2010, 1, 21, 0, 0, 0, 0, time.UTC),
		2.2,
		3.3,
		4.4,
		100,
	})
	if len(timeseries.Data) != 2{
		t.Errorf("%s dates", len(timeseries.Data))
	}
	jsonData, _ := json.Marshal(&timeseries)
	expected := []uint8("{\"Name\":\"\",\"ISIN\":\"\",\"Data\":[{\"Time\":\"2010-11-21T00:00:00Z\",\"High\":22,\"Low\":33,\"Close\":44,\"Volume\":5500},{\"Time\":\"2010-01-21T00:00:00Z\",\"High\":2.2,\"Low\":3.3,\"Close\":4.4,\"Volume\":100}]}")
	if ! reflect.DeepEqual(expected, jsonData){
		t.Errorf("json did not match expected:\n%s\nbut got\n%s\n", expected, jsonData)
	}
}
