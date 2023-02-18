package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type percentileStat struct {
	Name   string           `json:"name"`
	Origin string           `json:"origin"`
	Values map[string]int64 `json:"values"`
}

func newPercentileStatFromJSON(b []byte) (*percentileStat, error) {
	var pstat percentileStat
	err := json.Unmarshal(b, &pstat)
	if err != nil {
		return nil, fmt.Errorf("error decoding values stat `%v`: %v", string(b), err)
	}
	return &pstat, nil
}

func (i *percentileStat) toPoints() []*point {
	points := make([]*point, 0, len(i.Values))

	for name, value := range i.Values {
		if i.Origin == "percentile.bucket" {
			bucketMetricType := gauge
			if strings.Contains(name, "count") {
				bucketMetricType = counter
			}
			points = append(points, &point{
				Name:        fmt.Sprintf("%s_percentile_bucket", i.Name),
				Type:        bucketMetricType,
				Value:       value,
				Description: fmt.Sprintf("percentile bucket statistics %s", i.Name),
				LabelName:   "bucket",
				LabelValue:  name,
			})
		} else {
			points = append(points, &point{
				Name:        fmt.Sprintf("percentile_%s", i.Name),
				Type:        counter,
				Value:       value,
				Description: fmt.Sprintf("percentile statistics %s", i.Name),
				LabelName:   "counter",
				LabelValue:  name,
			})
		}
	}

	return points
}
