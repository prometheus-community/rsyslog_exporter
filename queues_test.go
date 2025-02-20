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
	queueStat = []byte(`{"name":"main Q","size":10,"enqueued":20,"full":30,"discarded.full":40,"discarded.nf":50,"maxqsize":60}`)
)

func TestNewQueueFromJSON(t *testing.T) {
	logType := getStatType(queueStat)
	if logType != rsyslogQueue {
		t.Errorf("detected pstat type should be %d but is %d", rsyslogQueue, logType)
	}

	pstat, err := newQueueFromJSON([]byte(queueStat))
	if err != nil {
		t.Fatalf("expected parsing queue stat not to fail, got: %v", err)
	}

	if want, got := "main Q", pstat.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := int64(10), pstat.Size; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(20), pstat.Enqueued; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(30), pstat.Full; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(40), pstat.DiscardedFull; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(50), pstat.DiscardedNf; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(60), pstat.MaxQsize; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}
}

func TestQueueToPoints(t *testing.T) {
	pstat, err := newQueueFromJSON([]byte(queueStat))
	if err != nil {
		t.Fatalf("expected parsing queue stat not to fail, got: %v", err)
	}
	points := pstat.toPoints()

	point := points[0]
	if want, got := "queue_size", point.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := "main Q", point.LabelValue; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[1]
	if want, got := "queue_enqueued", point.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := int64(20), point.Value; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := counter, point.Type; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "main Q", point.LabelValue; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[2]
	if want, got := "queue_full", point.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := int64(30), point.Value; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := counter, point.Type; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "main Q", point.LabelValue; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[3]
	if want, got := "queue_discarded_full", point.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := int64(40), point.Value; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := counter, point.Type; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "main Q", point.LabelValue; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[4]
	if want, got := "queue_discarded_not_full", point.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := int64(50), point.Value; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := counter, point.Type; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "main Q", point.LabelValue; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[5]
	if want, got := "queue_max_size", point.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := int64(60), point.Value; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := gauge, point.Type; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := "main Q", point.LabelValue; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}
}
