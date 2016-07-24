package data

import (
	"fmt"
	"time"
)

// descripes the value for a certain day for a certain instrument
type DataPoint struct {
	Time   time.Time
	High   float32
	Low    float32
	Close  float32
	Volume uint64
}

func (d DataPoint) String() string {
	return fmt.Sprintf("%s high:%f low:%f close:%f vol:%d", d.Time.Format("2006-01-02"), d.High, d.Low, d.Close, d.Volume)
}

// wrapper for simplicity
type DataPoints []DataPoint

// needed for sorting TODO: boilerplate
func (slice DataPoints) Len() int {
	return len(slice)
}

// needed for sorting
func (slice DataPoints) Less(i, j int) bool {
	return slice[i].Time.Before(slice[j].Time)
}

// needed for sorting
func (slice DataPoints) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// describes a instrument with all known values on day basis
type Instrument struct {
	Name string
	ISIN string
	Data DataPoints
}

func (instr *Instrument) Add(p *DataPoint) {
	instr.Data = append(instr.Data, *p)
}
