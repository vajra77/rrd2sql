package main

import (
	"fmt"
	"math"
	"time"

	"github.com/ziutek/rrd"
)

type Result struct {
	Timestamp int64
	AvgIn     float64
	AvgOut    float64
	MaxIn     float64
	MaxOut    float64
}

func FetchDailyValues(rrdFile string) (Result, error) {

	var result Result

	end := time.Now()
	start := end.AddDate(0, 0, -1)
	result.Timestamp = end.Unix()

	// Compute average
	fetchAvg, err := rrd.Fetch(rrdFile, "AVERAGE", start, end, 300*time.Second)
	if err != nil {
		return Result{}, fmt.Errorf("unable to fetch data: %v", err)
	}
	values := fetchAvg.Values()
	nVal := len(values)
	for i := 0; i < nVal; i += 2 {
		if math.IsNaN(values[i]) || math.IsNaN(values[i+1]) {
			continue
		}

		if values[i] > result.MaxIn {
			result.MaxIn = values[i]
		}

		if values[i+1] > result.MaxOut {
			result.MaxOut = values[i+1]
		}

		result.AvgIn = result.AvgIn + values[i]
		result.AvgOut = result.AvgOut + values[i+1]
	}

	result.AvgIn = 8 * result.AvgIn / float64(nVal) * 2
	result.AvgOut = 8 * result.AvgOut / float64(nVal) * 2
	result.MaxIn = 8 * result.MaxIn
	result.MaxOut = 8 * result.MaxOut

	return result, nil
}
