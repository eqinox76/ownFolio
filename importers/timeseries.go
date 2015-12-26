package importers

import "time"

type DataPoint struct {
	Time time.Time
	High      float32
	Low       float32
	Volume    uint64
	Close     float32
}

type DataPoints []DataPoint

func (slice DataPoints) Len() int {
    return len(slice)
}

func (slice DataPoints) Less(i, j int) bool {
    return slice[i].Time.Before(slice[j].Time);
}

func (slice DataPoints) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}



type Instrument struct {
	Name string
	ISIN string
	Data DataPoints
}

func (instr *Instrument) add(p *DataPoint){
	instr.Data = append(instr.Data, *p)
}
