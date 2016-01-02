package api

import (
	"log"
	"encoding/gob"
	"bytes"

	"appengine/memcache"
	"appengine"

	"github.com/eqinox76/ownFolio/importers"
	"github.com/eqinox76/ownFolio/data"
)

func getDataF(ctx appengine.Context, id string) (data.C3Data, error){
	url := importers.GenerateURL("YAHOO", id)
	instr, err := importers.GetHistory(ctx, url)
	return instr.C3Data(), err
}

var getData = getDataF


func GetInstrument(ctx appengine.Context, id string) data.C3Data {
	// get item from memcache
	if item, err := memcache.Get(ctx, id); err == memcache.ErrCacheMiss {
		// item not in cache
		instr, err := getData(ctx, id)
		if err != nil{
			log.Panicf("could not get data for %s because %s", id, err)
		}
		var buffer bytes.Buffer
		enc := gob.NewEncoder(&buffer)
		err = enc.Encode(instr)
		if err != nil{
			log.Panicf("could not encode %s", instr)
		}
		cacheEntry := &memcache.Item{
			Key: id,
			Value: buffer.Bytes(),
		}
		err = memcache.Add(ctx, cacheEntry)
		if err != nil{
			log.Printf("something was wrong when adding a new item %s err:%s", ctx, err)
		}
		return instr
	} else if err != nil {
		log.Panicf("error getting id:%s err:%s", id, err)
		instr, err := getData(ctx, id)
		if err != nil{
			log.Panicf("could not get data for %s because %s", id, err)
		}
		return instr
	} else {
		buffer := bytes.NewBuffer(item.Value)
		var instr data.C3Data
		dec := gob.NewDecoder(buffer)
		err := dec.Decode(&instr)
		if err != nil{
			log.Panicf("Could not decode data for %s err: %s", id, err)
			instr, err = getData(ctx, id)
			if err != nil{
				log.Panicf("could not get data for %s because %s", id, err)
			}
		}
		return instr
	}
}
