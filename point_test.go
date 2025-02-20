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
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestCounter(t *testing.T) {
	p1 := &point{
		Name:  "my counter",
		Type:  counter,
		Value: int64(10),
	}

	if want, got := float64(10), p1.promValue(); want != got {
		t.Errorf("want '%f', got '%f'", want, got)
	}

	if want, got := prometheus.ValueType(1), p1.promType(); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	wanted := `Desc{fqName: "rsyslog_my counter", help: "", constLabels: {}, variableLabels: {}}`
	if want, got := wanted, p1.promDescription().String(); want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}
}

func TestGauge(t *testing.T) {
	p1 := &point{
		Name:  "my gauge",
		Type:  gauge,
		Value: int64(10),
	}

	if want, got := float64(10), p1.promValue(); want != got {
		t.Errorf("want '%f', got '%f'", want, got)
	}

	if want, got := prometheus.ValueType(2), p1.promType(); want != got {
		t.Errorf("want '%v', got '%v'", want, got)
	}

	wanted := `Desc{fqName: "rsyslog_my gauge", help: "", constLabels: {}, variableLabels: {}}`
	if want, got := wanted, p1.promDescription().String(); want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

}
