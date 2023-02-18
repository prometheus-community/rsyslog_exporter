// { "name": "global", "origin": "percentile", "values": { "host_statistics.new_metric_add": 1, "host_statistics.ops_overflow": 0 } }

// { "name": "host_statistics", "origin": "percentile.bucket", "values": { "msg_per_host|p95": 1950, "msg_per_host|p50": 1500, "msg_per_host|p99": 1990, "msg_per_host|window_min": 1001, "msg_per_host|window_max": 2000, "msg_per_host|window_sum": 1500500, "msg_per_host|window_count": 1000 } }

package main

import (
	"reflect"
	"testing"
)

func TestGetPercentile(t *testing.T) {
	log := []byte(`{ "name": "global", "origin": "percentile", "values": { "host_statistics.new_metric_add": 1, "host_statistics.ops_overflow": 0 } }`)
	values := map[string]int64{
		"host_statistics.ops_overflow":   0,
		"host_statistics.new_metric_add": 1,
	}

	if want, got := rsyslogPercentile, getStatType(log); want != got {
		t.Errorf("detected pstat type should be %d but is %d", want, got)
	}

	pstat, err := newPercentileStatFromJSON(log)
	if err != nil {
		t.Fatalf("expected parsing dynamic stat not to fail, got: %v", err)
	}

	if want, got := "global", pstat.Name; want != got {
		t.Errorf("invalid name, want '%s', got '%s'", want, got)
	}

	if want, got := values, pstat.Values; !reflect.DeepEqual(want, got) {
		t.Errorf("unexpected values, want: %+v got: %+v", want, got)
	}
}

func TestPercentileToPoints(t *testing.T) {
	log := []byte(`{ "name": "global", "origin": "percentile", "values": { "host_statistics.new_metric_add": 1, "host_statistics.ops_overflow": 0 } }`)
	wants := map[string]point{
		"host_statistics.ops_overflow": point{
			Name:        "percentile_global",
			Type:        counter,
			Value:       0,
			Description: "percentile statistics global",
			LabelName:   "counter",
			LabelValue:  "host_statistics.ops_overflow",
		},
		"host_statistics.new_metric_add": point{
			Name:        "percentile_global",
			Type:        counter,
			Value:       1,
			Description: "percentile statistics global",
			LabelName:   "counter",
			LabelValue:  "host_statistics.new_metric_add",
		},
	}

	seen := map[string]bool{}
	for name := range wants {
		seen[name] = false
	}

	pstat, err := newPercentileStatFromJSON(log)
	if err != nil {
		t.Fatalf("expected parsing percentile stat not to fail, got: %v", err)
	}

	points := pstat.toPoints()
	for _, got := range points {
		key := got.LabelValue
		want, ok := wants[key]
		if !ok {
			t.Errorf("unexpected point, got: %+v", got)
			continue
		}

		if !reflect.DeepEqual(want, *got) {
			t.Errorf("expected point to be %+v, got %+v", want, got)
		}

		if seen[key] {
			t.Errorf("point seen multiple times: %+v", got)
		}
		seen[key] = true
	}

	for name, ok := range seen {
		if !ok {
			t.Errorf("expected to see point with key %s, but did not", name)
		}
	}
}

func TestGetPercentileBucket(t *testing.T) {
	log := []byte(`{ "name": "host_statistics", "origin": "percentile.bucket", "values": { "msg_per_host|p95": 1950, "msg_per_host|p50": 1500, "msg_per_host|p99": 1990, "msg_per_host|window_min": 1001, "msg_per_host|window_max": 2000, "msg_per_host|window_sum": 1500500, "msg_per_host|window_count": 1000 } }`)
	values := map[string]int64{
		"msg_per_host|p95":          1950,
		"msg_per_host|p50":          1500,
		"msg_per_host|p99":          1990,
		"msg_per_host|window_min":   1001,
		"msg_per_host|window_max":   2000,
		"msg_per_host|window_sum":   1500500,
		"msg_per_host|window_count": 1000,
	}

	if want, got := rsyslogPercentileBucket, getStatType(log); want != got {
		t.Errorf("detected pstat type should be %d but is %d", want, got)
	}

	pstat, err := newPercentileStatFromJSON(log)
	if err != nil {
		t.Fatalf("expected parsing dynamic stat not to fail, got: %v", err)
	}

	if want, got := "host_statistics", pstat.Name; want != got {
		t.Errorf("invalid name, want '%s', got '%s'", want, got)
	}

	if want, got := values, pstat.Values; !reflect.DeepEqual(want, got) {
		t.Errorf("unexpected values, want: %+v got: %+v", want, got)
	}
}

func TestPercentileBucketToPoints(t *testing.T) {
	log := []byte(`{ "name": "host_statistics", "origin": "percentile.bucket", "values": { "msg_per_host|p95": 1950, "msg_per_host|p50": 1500, "msg_per_host|p99": 1990, "msg_per_host|window_min": 1001, "msg_per_host|window_max": 2000, "msg_per_host|window_sum": 1500500, "msg_per_host|window_count": 1000 } }`)
	wants := map[string]point{
		"msg_per_host|p95": point{
			Name:        "host_statistics_percentile_bucket",
			Type:        gauge,
			Value:       1950,
			Description: "percentile bucket statistics host_statistics",
			LabelName:   "bucket",
			LabelValue:  "msg_per_host|p95",
		},
		"msg_per_host|p50": point{
			Name:        "host_statistics_percentile_bucket",
			Type:        gauge,
			Value:       1500,
			Description: "percentile bucket statistics host_statistics",
			LabelName:   "bucket",
			LabelValue:  "msg_per_host|p50",
		},
		"msg_per_host|p99": point{
			Name:        "host_statistics_percentile_bucket",
			Type:        gauge,
			Value:       1990,
			Description: "percentile bucket statistics host_statistics",
			LabelName:   "bucket",
			LabelValue:  "msg_per_host|p99",
		},
		"msg_per_host|window_min": point{
			Name:        "host_statistics_percentile_bucket",
			Type:        gauge,
			Value:       1001,
			Description: "percentile bucket statistics host_statistics",
			LabelName:   "bucket",
			LabelValue:  "msg_per_host|window_min",
		},
		"msg_per_host|window_max": point{
			Name:        "host_statistics_percentile_bucket",
			Type:        gauge,
			Value:       2000,
			Description: "percentile bucket statistics host_statistics",
			LabelName:   "bucket",
			LabelValue:  "msg_per_host|window_max",
		},
		"msg_per_host|window_sum": point{
			Name:        "host_statistics_percentile_bucket",
			Type:        gauge,
			Value:       1500500,
			Description: "percentile bucket statistics host_statistics",
			LabelName:   "bucket",
			LabelValue:  "msg_per_host|window_sum",
		},
		"msg_per_host|window_count": point{
			Name:        "host_statistics_percentile_bucket",
			Type:        gauge,
			Value:       1000,
			Description: "percentile bucket statistics host_statistics",
			LabelName:   "bucket",
			LabelValue:  "msg_per_host|window_count",
		},
	}

	seen := map[string]bool{}
	for name := range wants {
		seen[name] = false
	}

	pstat, err := newPercentileStatFromJSON(log)
	if err != nil {
		t.Fatalf("expected parsing percentile stat not to fail, got: %v", err)
	}

	points := pstat.toPoints()
	for _, got := range points {
		key := got.LabelValue
		want, ok := wants[key]
		if !ok {
			t.Errorf("unexpected point, got: %+v", got)
			continue
		}

		if !reflect.DeepEqual(want, *got) {
			t.Errorf("expected point to be %+v, got %+v", want, got)
		}

		if seen[key] {
			t.Errorf("point seen multiple times: %+v", got)
		}
		seen[key] = true
	}

	for name, ok := range seen {
		if !ok {
			t.Errorf("expected to see point with key %s, but did not", name)
		}
	}
}
