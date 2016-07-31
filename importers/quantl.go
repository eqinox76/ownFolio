package importers

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"time"

	"appengine"
	"appengine/urlfetch"

	"github.com/eqinox76/ownFolio/data"
)

var ErrEmptyResponse = fmt.Errorf("Empty data returned.")

type Root struct {
	Dataset struct {
		Name         string
		Column_names []string
		Data         [][]interface{}
	}
}

func ParseQuantlJson(reader io.Reader) (data.TimeSeries, error) {
	var instr data.TimeSeries

	var r Root

	dec := json.NewDecoder(reader)
	if err := dec.Decode(&r); err != nil {
		return instr, err
	}

	instr.Name = r.Dataset.Name

	datePos := -1
	highPos := -1
	lowPos := -1
	closePos := -1
	volPos := -1

	if len(r.Dataset.Column_names) == 0 {
		// the response is empty. most likely there is no data.
		return instr, nil
	}

	for index, value := range r.Dataset.Column_names {
		switch value {
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
		return instr, fmt.Errorf("Could not identify position for a datapoint %s.", r.Dataset.Column_names)
	}

	for _, entryAy := range r.Dataset.Data {

		var dp data.DataPoint
		var err error

		pos := highPos
		val, ok := entryAy[pos].(float64)
		if !ok {
			return instr, fmt.Errorf("highPos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.High = float32(val)

		pos = lowPos
		val, ok = entryAy[pos].(float64)
		if !ok {
			return instr, fmt.Errorf("lowPos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.Low = float32(val)

		pos = volPos
		val, ok = entryAy[pos].(float64)
		if !ok {
			return instr, fmt.Errorf("volPos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.Volume = uint64(val)

		pos = closePos
		val, ok = entryAy[pos].(float64)
		if !ok {
			return instr, fmt.Errorf("closePos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.Close = float32(val)

		pos = datePos
		date, ok := entryAy[pos].(string)
		if !ok {
			return instr, fmt.Errorf("datePos '%d' not parsable %s", pos, entryAy[pos])
		}
		dp.Time, err = time.Parse("2006-01-02", date)
		if err != nil {
			return instr, fmt.Errorf("Could not parse datetime '%s': %e", date, err)
		}

		instr.Add(&dp)
	}
	sort.Sort(instr.Data)
	return instr, nil
}

func GenerateURL(database string, dataset string) string {
	return fmt.Sprintf("https://www.quandl.com/api/v3/datasets/%s/%s.json", database, dataset)
}

func GetHistory(ctx appengine.Context, url string) (data.TimeSeries, error) {
	// example https://www.quandl.com/api/v3/datasets/YAHOO/INDEX_GDAXI.json

	var instr data.TimeSeries
	client := urlfetch.Client(ctx)
	r, err := client.Get(url)
	if err != nil {
		return instr, err
	}
	defer r.Body.Close()

	return ParseQuantlJson(r.Body)
}
