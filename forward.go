// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
