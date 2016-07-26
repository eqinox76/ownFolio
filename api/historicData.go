package api

import (
	"bytes"
	"encoding/gob"
	"log"
	"math"

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

func GetInstrumentLimited(ctx appengine.Context, id string, limit int) data.DataPoints {

	data := GetInstrument(ctx, id)
	//TODO why no min for int ?
	pos := math.Max(0, float64(len(data)-limit))

	return data[int(pos):]
}

func GetInstrument(ctx appengine.Context, id string) data.DataPoints {
	// get item from memcache
	if item, err := memcache.Get(ctx, id); err == memcache.ErrCacheMiss {
		// item not in cache
		instr, err := getData(ctx, id)
		if err != nil {
			log.Panicf("could not get data for '%s' because: %s", id, err)
		}
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err = enc.Encode(instr)
		if err != nil {
			log.Panicf("could not encode %s", instr)
		}
		cacheEntry := &memcache.Item{
			Key:   id,
			Value: buffer.Bytes(),
		}
		err = memcache.Add(ctx, cacheEntry)
		if err != nil {
			ctx.Infof("something was wrong when adding a new item %s err:%s", ctx, err)
		}

		ctx.Infof("Got '%s' from source", id)
		return instr
	} else if err != nil {
		log.Panicf("error getting id:%s err:%s", id, err)
		instr, err := getData(ctx, id)
		if err != nil {
			log.Panicf("could not get data for %s because %s", id, err)
		}
		return instr
	} else {
		buffer := bytes.NewBuffer(item.Value)
		var instr data.DataPoints
		dec := gob.NewDecoder(buffer)
		err := dec.Decode(&instr)
		if err != nil {
			log.Panicf("Could not decode data for %s err: %s", id, err)
			instr, err = getData(ctx, id)
			if err != nil {
				log.Panicf("could not get data for %s because %s", id, err)
			}
		}

		ctx.Infof("Got '%s' from cache", id)
		return instr
	}
}
