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
	"reflect"
	"testing"
)

var (
	dynafileCacheLog = []byte(`{ "name": "dynafile cache cluster", "origin": "omfile", "requests": 1783254, "level0": 1470906, "missed": 2625, "evicted": 2525, "maxused": 100, "closetimeouts": 10 }`)
)

func TestNewDynafileCacheFromJSON(t *testing.T) {
	logType := getStatType(dynafileCacheLog)
	if logType != rsyslogDynafileCache {
		t.Errorf("detected pstat type should be %d but is %d", rsyslogDynafileCache, logType)
	}

	pstat, err := newDynafileCacheFromJSON([]byte(dynafileCacheLog))
	if err != nil {
		t.Fatalf("expected parsing dynafile cache stat not to fail, got: %v", err)
	}

	if want, got := "cluster", pstat.Name; want != got {
		t.Errorf("want '%s', got '%s'", want, got)
	}

	if want, got := int64(1783254), pstat.Requests; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(1470906), pstat.Level0; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(2625), pstat.Missed; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(2525), pstat.Evicted; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(100), pstat.MaxUsed; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}

	if want, got := int64(10), pstat.CloseTimeouts; want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}
}

func TestDynafileCacheToPoints(t *testing.T) {

	wants := map[string]point{
		"dynafile_cache_requests": point{
			Name:        "dynafile_cache_requests",
			Type:        counter,
			Value:       1783254,
			Description: "number of requests made to obtain a dynafile",
			LabelName:   "cache",
			LabelValue:  "cluster",
		},
		"dynafile_cache_level0": point{
			Name:        "dynafile_cache_level0",
			Type:        counter,
			Value:       1470906,
			Description: "number of requests for the current active file",
			LabelName:   "cache",

			LabelValue: "cluster",
		},
		"dynafile_cache_missed": point{
			Name:        "dynafile_cache_missed",
			Type:        counter,
			Value:       2625,
			Description: "number of cache misses",
			LabelName:   "cache",
			LabelValue:  "cluster",
		},
		"dynafile_cache_evicted": point{
			Name:        "dynafile_cache_evicted",
			Type:        counter,
			Value:       2525,
			Description: "number of times a file needed to be evicted from cache",
			LabelName:   "cache",
			LabelValue:  "cluster",
		},
		"dynafile_cache_maxused": point{
			Name:        "dynafile_cache_maxused",
			Type:        counter,
			Value:       100,
			Description: "maximum number of cache entries ever used",
			LabelName:   "cache",
			LabelValue:  "cluster",
		},
		"dynafile_cache_closetimeouts": point{
			Name:        "dynafile_cache_closetimeouts",
			Type:        counter,
			Value:       10,
			Description: "number of times a file was closed due to timeout settings",
			LabelName:   "cache",
			LabelValue:  "cluster",
		},
	}

	seen := map[string]bool{}
	for name := range wants {
		seen[name] = false
	}

	pstat, err := newDynafileCacheFromJSON(dynafileCacheLog)
	if err != nil {
		t.Fatalf("expected parsing dynafile cache stat not to fail, got: %v", err)
	}

	points := pstat.toPoints()
	for _, got := range points {
		want, ok := wants[got.Name]
		if !ok {
			t.Errorf("unexpected point, got: %+v", got)
			continue
		}

		if !reflect.DeepEqual(want, *got) {
			t.Errorf("expected point to be %+v, got %+v", want, got)
		}

		if seen[got.Name] {
			t.Errorf("point seen multiple times: %+v", got)
		}
		seen[got.Name] = true
	}

	for name, ok := range seen {
		if !ok {
			t.Errorf("expected to see point with key %s, but did not", name)
		}
	}
}
