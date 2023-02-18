package main

import (
	"fmt"
	"testing"
)

func testHelper(t *testing.T, line []byte, testCase []*testUnit) {
	exporter := newRsyslogExporter()
	exporter.handleStatLine(line)

	for _, k := range exporter.keys() {
		t.Logf("have key: '%s'", k)
	}

	for _, item := range testCase {
		p, err := exporter.get(item.key())
		if err != nil {
			t.Error(err)
		}

		if want, got := item.Val, p.promValue(); want != got {
			t.Errorf("%s: want '%f', got '%f'", item.Name, want, got)
		}
	}

	exporter.handleStatLine(line)

	for _, item := range testCase {
		p, err := exporter.get(item.key())
		if err != nil {
			t.Error(err)
		}

		var wanted float64
		switch p.Type {
		case counter:
			wanted = item.Val
		case gauge:
			wanted = item.Val
		default:
			t.Errorf("%d is not a valid metric type", p.Type)
			continue
		}

		if want, got := wanted, p.promValue(); want != got {
			t.Errorf("%s: want '%f', got '%f'", item.Name, want, got)
		}
	}
}

type testUnit struct {
	Name       string
	Val        float64
	LabelValue string
}

func (t *testUnit) key() string {
	return fmt.Sprintf("%s.%s", t.Name, t.LabelValue)
}

func TestHandleLineWithAction(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name:       "action_processed",
			Val:        100000,
			LabelValue: "test_action",
		},
		&testUnit{
			Name:       "action_failed",
			Val:        2,
			LabelValue: "test_action",
		},
		&testUnit{
			Name:       "action_suspended",
			Val:        1,
			LabelValue: "test_action",
		},
		&testUnit{
			Name:       "action_suspended_duration",
			Val:        1000,
			LabelValue: "test_action",
		},
		&testUnit{
			Name:       "action_resumed",
			Val:        1,
			LabelValue: "test_action",
		},
	}

	actionLog := []byte(`2017-08-30T08:10:04.786350+00:00 some-node.example.org rsyslogd-pstats: {"name":"test_action","processed":100000,"failed":2,"suspended":1,"suspended.duration":1000,"resumed":1}`)
	testHelper(t, actionLog, tests)
}

func TestHandleLineWithResource(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name:       "resource_utime",
			Val:        10,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name:       "resource_stime",
			Val:        20,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name:       "resource_maxrss",
			Val:        30,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name:       "resource_minflt",
			Val:        40,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name:       "resource_majflt",
			Val:        50,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name:       "resource_inblock",
			Val:        60,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name:       "resource_oublock",
			Val:        70,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name:       "resource_nvcsw",
			Val:        80,
			LabelValue: "resource-usage",
		},
		&testUnit{
			Name:       "resource_nivcsw",
			Val:        90,
			LabelValue: "resource-usage",
		},
	}

	resourceLog := []byte(`2017-08-30T08:10:04.786350+00:00 some-node.example.org rsyslogd-pstats: {"name":"resource-usage","utime":10,"stime":20,"maxrss":30,"minflt":40,"majflt":50,"inblock":60,"oublock":70,"nvcsw":80,"nivcsw":90}`)
	testHelper(t, resourceLog, tests)
}

func TestHandleLineWithInput(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name:       "input_submitted",
			Val:        1000,
			LabelValue: "test_input",
		},
	}

	inputLog := []byte(`2017-08-30T08:10:04.786350+00:00 some-node.example.org rsyslogd-pstats: {"name":"test_input", "origin":"imuxsock", "submitted":1000}`)
	testHelper(t, inputLog, tests)
}

func TestHandleLineWithQueue(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name:       "queue_size",
			Val:        10,
			LabelValue: "main Q",
		},
		&testUnit{
			Name:       "queue_enqueued",
			Val:        20,
			LabelValue: "main Q",
		},
		&testUnit{
			Name:       "queue_full",
			Val:        30,
			LabelValue: "main Q",
		},
		&testUnit{
			Name:       "queue_discarded_full",
			Val:        40,
			LabelValue: "main Q",
		},
		&testUnit{
			Name:       "queue_discarded_not_full",
			Val:        50,
			LabelValue: "main Q",
		},
		&testUnit{
			Name:       "queue_max_size",
			Val:        60,
			LabelValue: "main Q",
		},
	}

	queueLog := []byte(`2017-08-30T08:10:04.786350+00:00 some-node.example.org rsyslogd-pstats: {"name":"main Q","size":10,"enqueued":20,"full":30,"discarded.full":40,"discarded.nf":50,"maxqsize":60}`)
	testHelper(t, queueLog, tests)
}

func TestHandleLineWithGlobal(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name:       "dynstat_global",
			Val:        1,
			LabelValue: "msg_per_host.ops_overflow",
		},
		&testUnit{
			Name:       "dynstat_global",
			Val:        3,
			LabelValue: "msg_per_host.new_metric_add",
		},
		&testUnit{
			Name:       "dynstat_global",
			Val:        0,
			LabelValue: "msg_per_host.no_metric",
		},
		&testUnit{
			Name:       "dynstat_global",
			Val:        0,
			LabelValue: "msg_per_host.metrics_purged",
		},
		&testUnit{
			Name:       "dynstat_global",
			Val:        0,
			LabelValue: "msg_per_host.ops_ignored",
		},
	}

	log := []byte(`2018-01-18T09:39:12.763025+00:00 some-node.example.org rsyslogd-pstats: { "name": "global", "origin": "dynstats", "values": { "msg_per_host.ops_overflow": 1, "msg_per_host.new_metric_add": 3, "msg_per_host.no_metric": 0, "msg_per_host.metrics_purged": 0, "msg_per_host.ops_ignored": 0 } }`)

	testHelper(t, log, tests)
}

func TestHandleLineWithDynafileCache(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name:       "dynafile_cache_requests",
			Val:        412044,
			LabelValue: "cluster",
		},
		&testUnit{
			Name:       "dynafile_cache_level0",
			Val:        294002,
			LabelValue: "cluster",
		},
		&testUnit{
			Name:       "dynafile_cache_missed",
			Val:        210,
			LabelValue: "cluster",
		},
		&testUnit{
			Name:       "dynafile_cache_evicted",
			Val:        14,
			LabelValue: "cluster",
		},
	}

	dynafileCacheLog := []byte(`2019-07-03T17:04:01.312432+00:00 some-node.example.org rsyslogd-pstats: { "name": "dynafile cache cluster", "origin": "omfile", "requests": 412044, "level0": 294002, "missed": 210, "evicted": 14, "maxused": 100, "closetimeouts": 0 }`)
	testHelper(t, dynafileCacheLog, tests)
}

func TestHandleLineWithPercentileGlobal(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name:       "percentile_global",
			Val:        0,
			LabelValue: "host_statistics.ops_overflow",
		},
		&testUnit{
			Name:       "percentile_global",
			Val:        1,
			LabelValue: "host_statistics.new_metric_add",
		},
	}

	log := []byte(`2022-02-18T19:11:12.672935+00:00 some-node.example.org rsyslogd-pstats: { "name": "global", "origin": "percentile", "values": { "host_statistics.new_metric_add": 1, "host_statistics.ops_overflow": 0 } }`)

	testHelper(t, log, tests)
}

func TestHandleLineWithPercentileBucket(t *testing.T) {
	tests := []*testUnit{
		&testUnit{
			Name:       "host_statistics_percentile_bucket",
			Val:        1950,
			LabelValue: "msg_per_host|p95",
		},
		&testUnit{
			Name:       "host_statistics_percentile_bucket",
			Val:        1500,
			LabelValue: "msg_per_host|p50",
		},
		&testUnit{
			Name:       "host_statistics_percentile_bucket",
			Val:        1990,
			LabelValue: "msg_per_host|p99",
		},
		&testUnit{
			Name:       "host_statistics_percentile_bucket",
			Val:        1001,
			LabelValue: "msg_per_host|window_min",
		},
		&testUnit{
			Name:       "host_statistics_percentile_bucket",
			Val:        2000,
			LabelValue: "msg_per_host|window_max",
		},
		&testUnit{
			Name:       "host_statistics_percentile_bucket",
			Val:        1500500,
			LabelValue: "msg_per_host|window_sum",
		},
		&testUnit{
			Name:       "host_statistics_percentile_bucket",
			Val:        1000,
			LabelValue: "msg_per_host|window_count",
		},
	}

	log := []byte(`2022-02-18T19:11:12.672935+00:00 some-node.example.org rsyslogd-pstats: { "name": "host_statistics", "origin": "percentile.bucket", "values": { "msg_per_host|p95": 1950, "msg_per_host|p50": 1500, "msg_per_host|p99": 1990, "msg_per_host|window_min": 1001, "msg_per_host|window_max": 2000, "msg_per_host|window_sum": 1500500, "msg_per_host|window_count": 1000 } }`)

	testHelper(t, log, tests)
}

func TestHandleUnknown(t *testing.T) {
	unknownLog := []byte(`2017-08-30T08:10:04.786350+00:00 some-node.example.org rsyslogd-pstats: {"a":"b"}`)

	exporter := newRsyslogExporter()
	exporter.handleStatLine(unknownLog)

	if want, got := 0, len(exporter.keys()); want != got {
		t.Errorf("want '%d', got '%d'", want, got)
	}
}
