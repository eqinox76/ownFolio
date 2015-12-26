package importers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
	"sort"
)

func ParseQuantlJson(data []byte) (Instrument, error) {
	var instr Instrument

	var parsed interface{}
	err := json.Unmarshal(data, &parsed)
	if err != nil {
		return instr, fmt.Errorf("Could not parse json: %s", err)
	}

	m, ok := parsed.(map[string]interface{})
	if !ok {
		return instr, fmt.Errorf("Wrong Format. First element is %s", reflect.TypeOf(data).String())
	}

	dataset, ok := m["dataset"]
	if !ok {
		return instr, errors.New("Wrong Format.")
	}

	datasetm, ok := dataset.(map[string]interface{})
	if !ok {
		return instr, errors.New("Wrong Format.")
	}

	name, ok := datasetm["name"]
	if !ok {
		return instr, errors.New("Wrong Format.")
	}
	instr.Name = name.(string)

	columnes, ok := datasetm["column_names"]
	if !ok {
		return instr, errors.New("Could not find 'column_names'.")
	}

	columnesAy, ok := columnes.([]interface{})
	if !ok {
		return instr, errors.New("Could not convert 'column_names' to array")
	}

	datePos := -1
	highPos := -1
	lowPos := -1
	closePos := -1
	volPos := -1

	for index, value := range columnesAy{
		v, ok := value.(string)
		if !ok {
			return instr, fmt.Errorf("Unknown type in column_names array: %s", reflect.TypeOf(value))
		}
		switch (v){
			case "Date":
				datePos = index
			case "High":
				highPos = index
			case "Low":
				lowPos = index
			case "Adjusted Close":
				closePos = index
			case "Volume":
				volPos = index
		}
	}

	if datePos == -1 || highPos == -1 || lowPos == -1 || closePos == -1 || volPos == -1 {
		return instr, fmt.Errorf("Could not identify position for a datapoint %s.", columnesAy)
	}

	entries, ok := datasetm["data"]
	if ! ok {
		return instr, fmt.Errorf("Could not find 'data' entry")
	}

	entriesAy, ok := entries.([]interface{})
	if ! ok {
		return instr, fmt.Errorf("'data' entry is no array")
	}

	for _, entry := range entriesAy {
		entryAy, ok := entry.([]interface{})
		if !ok {
			return instr, fmt.Errorf("Could not interpret type '%s'", reflect.TypeOf(entry))
		}

		var dp DataPoint

		pos := highPos
		val, ok := entryAy[pos].(float64)
		if ! ok {
			return instr, fmt.Errorf("highPos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.High = float32(val)

		pos = lowPos
		val, ok = entryAy[pos].(float64)
		if ! ok {
			return instr, fmt.Errorf("lowPos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.Low = float32(val)

		pos = volPos
		val, ok = entryAy[pos].(float64)
		if ! ok {
			return instr, fmt.Errorf("volPos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.Volume = uint64(val)

		pos = closePos
		val, ok = entryAy[pos].(float64)
		if ! ok {
			return instr, fmt.Errorf("closePos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.Close = float32(val)

		pos = datePos
		date, ok := entryAy[pos].(string)
		if ! ok {
			return instr, fmt.Errorf("datePos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.Time, err = time.Parse("2006-01-02", date)
		if err != nil {
			return instr, fmt.Errorf("Could not parse datetime '%s': %e", date, err)
		}

		instr.add(&dp)
	}

	sort.Sort(instr.Data)

	return instr, nil

}
