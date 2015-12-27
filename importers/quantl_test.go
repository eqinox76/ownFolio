package importers

import (
	"bufio"
	"os"
	"testing"
	"time"
)

func TestParseHistoric(t *testing.T) {
	f, err := os.Open("dax.data.json")
	if err != nil {
		t.Error("Could not open data file. %s", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	instr, err := ParseQuantlJson(reader)

	if err != nil {
		t.Error(err)
	}

	if instr.Name != "DAX Index (Germany)" {
		t.Errorf("Name is wrong: %s", instr.Name)
	}

	if len(instr.Data) != 6352 {
		t.Errorf("Wrong number of Datapoints parsed: %d", len(instr.Data))
	}
	dp := &instr.Data[len(instr.Data)/2]
	if dp.High != 3378.14 {
		t.Errorf("High parsed wrong '%f'", dp.High)
	}

	if instr.Data[0].Time != time.Date(1990, time.November, 26, 0, 0, 0, 0, time.UTC) {
		t.Errorf("First date seems wrong %s", instr.Data[0].Time)
	}
	if instr.Data[len(instr.Data)-1].Time != time.Date(2015, time.December, 23, 0, 0, 0, 0, time.UTC) {
		t.Errorf("Last date seems wrong %s", instr.Data[len(instr.Data)-1].Time)
	}
}

// https://www.quandl.com/api/v3/datasets/YAHOO/INDEX_GDAXI.json
