package data

import (
	"time"
)

type Holding struct {
	ISIN    string
	Volume  float64
	Price   float64
	Fee     float64
	BuyDate time.Time
}
