package data

import (
	"fmt"
	"time"
)

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

type DataPoints []DataPoint

func (slice DataPoints) Len() int {
	return len(slice)
}

func (slice DataPoints) Less(i, j int) bool {
	return slice[i].Time.Before(slice[j].Time)
}

func (slice DataPoints) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Instrument struct {
	Name string
	ISIN string
	Data DataPoints
}

func (instr *Instrument) Add(p *DataPoint) {
	instr.Data = append(instr.Data, *p)
}

type C3Data struct{
	Date []string
	High []float32
	Low []float32
	Close []float32
}

func (instr *Instrument) C3Data() C3Data{
	var result C3Data
	for _, value := range instr.Data {
		result.Date = append(result.Date, value.Time.Format("2006-01-02"))
		result.High = append(result.High, value.High)
		result.Low = append(result.Low, value.Low)
		result.Close = append(result.Close, value.Close)
	}
	return result
}
