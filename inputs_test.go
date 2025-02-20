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

import "testing"

var (
	inputLog = []byte(`{"name":"test_input", "origin":"imuxsock", "submitted":1000}`)
)

func TestGetInput(t *testing.T) {
	logType := getStatType(inputLog)
	if logType != rsyslogInput {
		t.Errorf("detected pstat type should be %d but is %d", rsyslogInput, logType)
	}

	pstat, err := newInputFromJSON([]byte(inputLog))
	if err != nil {
		t.Fatalf("expected parsing input stat not to fail, got: %v", err)
	}

	if want, got := "test_input", pstat.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := int64(1000), pstat.Submitted; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}
}

func TestInputtoPoints(t *testing.T) {
	pstat, err := newInputFromJSON([]byte(inputLog))
	if err != nil {
		t.Fatalf("expected parsing input stat not to fail, got: %v", err)
	}

	points := pstat.toPoints()

	point := points[0]
	if want, got := "input_submitted", point.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := int64(1000), point.Value; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "test_input", point.LabelValue; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}
}
