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

type input struct {
	Name      string `json:"name"`
	Submitted int64  `json:"submitted"`
}

func newInputFromJSON(b []byte) (*input, error) {
	var pstat input
	err := json.Unmarshal(b, &pstat)
	if err != nil {
		return nil, fmt.Errorf("error decoding input stat `%v`: %v", string(b), err)
	}
	return &pstat, nil
}

func (i *input) toPoints() []*point {
	points := make([]*point, 1)

	points[0] = &point{
		Name:        "input_submitted",
		Type:        counter,
		Value:       i.Submitted,
		Description: "messages submitted",
		LabelName:   "input",
		LabelValue:  i.Name,
	}

	return points
}
