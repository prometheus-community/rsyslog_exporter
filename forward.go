package main

import (
	"encoding/json"
	"fmt"
)

type forward struct {
	Name      string `json:"name"`
	BytesSent int64  `json:"bytes.sent"`
}

func newForwardFromJSON(b []byte) (*forward, error) {
	var pstat forward
	err := json.Unmarshal(b, &pstat)
	if err != nil {
		return nil, fmt.Errorf("failed to decode forward stat `%v`: %v", string(b), err)
	}
	return &pstat, nil
}

func (f *forward) toPoints() []*point {
	points := make([]*point, 1)

	points[0] = &point{
		Name:        "forward_bytes_total",
		Type:        counter,
		Value:       f.BytesSent,
		Description: "bytes forwarded to destination",
		LabelName:   "destination",
		LabelValue:  f.Name,
	}

	return points
}
