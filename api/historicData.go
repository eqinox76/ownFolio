package api

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"strconv"

	"appengine"
	"appengine/memcache"

	"github.com/eqinox76/ownFolio/data"
	"github.com/eqinox76/ownFolio/importers"
)

func getDataF(ctx appengine.Context, id string) (data.DataPoints, error) {
	url := importers.GenerateURL("YAHOO", id)
	instr, err := importers.GetHistory(ctx, url)
	return instr.Data, err
}

var getData = getDataF

func GetInstrumentLimited(ctx appengine.Context, id string, limit int) ([]byte, error) {

	data, err := getInstrument(ctx, id)
	if err != nil {
		return nil, err
	}
	//TODO why no min for int ?
	pos := math.Max(0, float64(len(data)-limit))

	return formatData(data[int(pos):], id), nil
}

func GetInstrument(ctx appengine.Context, id string) ([]byte, error) {
	data, err := getInstrument(ctx, id)
	if err != nil {
		return nil, err
	}
	return formatData(data, id), nil
}

func getInstrument(ctx appengine.Context, id string) (data.DataPoints, error) {
	// get item from memcache
	if item, err := memcache.Get(ctx, id); err == memcache.ErrCacheMiss {
		// item not in cache
		instr, err := getData(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("could not get data for '%s' because: %s", id, err)
		}
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err = enc.Encode(instr)
		if err != nil {
			return nil, fmt.Errorf("could not encode %s", instr)
		}
		cacheEntry := &memcache.Item{
			Key:   id,
			Value: buffer.Bytes(),
		}
		err = memcache.Add(ctx, cacheEntry)
		if err != nil {
			return nil, fmt.Errorf("something was wrong when adding a new item %s err:%s", ctx, err)
		}

		ctx.Infof("Got '%s' from source", id)
		return instr, nil
	} else if err != nil {
		return nil, fmt.Errorf("error getting id:%s err:%s", id, err)
		instr, err := getData(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("could not get data for %s because %s", id, err)
		}
		return instr, nil
	} else {
		buffer := bytes.NewBuffer(item.Value)
		var instr data.DataPoints
		dec := gob.NewDecoder(buffer)
		err := dec.Decode(&instr)
		if err != nil {
			return nil, fmt.Errorf("Could not decode data for %s err: %s", id, err)
			instr, err = getData(ctx, id)
			if err != nil {
				return nil, fmt.Errorf("could not get data for %s because %s", id, err)
			}
		}

		ctx.Infof("Got '%s' from cache", id)
		return instr, nil
	}
}

/*
target example:
{
	columns: [
		['data3', 400, 500, 450, 700, 600, 500],
		['x', '2013-01-01', '2013-01-02', '2013-02-02', '2013-01-04', '2013-01-09', '2013-01-06'],
	]
}
*/
func formatData(instruments data.DataPoints, id string) []byte {
	// TODO the desired format is a 'bit' strange therefore we need to generate the json by hand
	var data bytes.Buffer
	var time bytes.Buffer
	data.WriteString("[\"")
	data.WriteString(id)
	data.WriteString("\"")

	time.WriteString("[\"")
	time.WriteString("date")
	time.WriteString("\"")
	for _, instr := range instruments {
		data.WriteString(", ")
		// TODO write at least directly into the bytes.Buffer strconv.AppendFloat seems to have the wrong interface
		// TODO lets convert everything to float64 and thereafter tell the function which float we used
		data.WriteString(strconv.FormatFloat(float64(instr.Close), 'f', -1, 32))

		time.WriteString(", \"")
		time.WriteString(instr.Time.Format("2006-01-02T15:04:05Z"))
		time.WriteString("\"")
	}
	data.WriteString("],")
	time.WriteString("]")

	var result bytes.Buffer
	result.WriteString("[")
	result.WriteString(data.String())
	result.WriteString(time.String())
	result.WriteString("]")

	return result.Bytes()
}
