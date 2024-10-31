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
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type pointType int

const (
	counter pointType = iota
	gauge
)

type point struct {
	Name        string
	Description string
	Type        pointType
	Value       int64
	LabelName   string
	LabelValue  string
}

func (p *point) promDescription() *prometheus.Desc {
	variableLabels := []string{}
	if p.promLabelName() != "" {
		variableLabels = []string{p.promLabelName()}
	}
	return prometheus.NewDesc(
		prometheus.BuildFQName("", "rsyslog", p.Name),
		p.Description,
		variableLabels,
		nil,
	)
}

func (p *point) promType() prometheus.ValueType {
	if p.Type == counter {
		return prometheus.CounterValue
	}
	return prometheus.GaugeValue
}

func (p *point) promValue() float64 {
	return float64(p.Value)
}

func (p *point) promLabelValue() string {
	return p.LabelValue
}

func (p *point) promLabelName() string {
	return p.LabelName
}

func (p *point) key() string {
	return fmt.Sprintf("%s.%s", p.Name, p.LabelValue)
}
