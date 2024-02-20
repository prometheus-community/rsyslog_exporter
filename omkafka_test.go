package main

import (
	"fmt"
	"testing"
)

var (
	omkafkaLog = []byte(`{ "name": "omkafka", "origin": "omkafka", "submitted": 59, "maxoutqsize": 9, "failures": 0, "topicdynacache.skipped": 57, "topicdynacache.miss": 2, "topicdynacache.evicted": 0, "acked": 55, "failures_msg_too_large": 0, "failures_unknown_topic": 0, "failures_queue_full": 0, "failures_unknown_partition": 0, "failures_other": 0, "errors_timed_out": 0, "errors_transport": 0, "errors_broker_down": 0, "errors_auth": 0, "errors_ssl": 0, "errors_other": 0, "rtt_avg_usec": 0, "throttle_avg_msec": 0, "int_latency_avg_usec": 0 }`)
)

func TestNewOmkafkaFromJSON(t *testing.T) {
	logType := getStatType(omkafkaLog)
	if logType != rsyslogOmkafka {
		t.Errorf("detected pstat type should be %d but is %d", rsyslogOmkafka, logType)
	}

	_, err := newOmkafkaFromJSON([]byte(omkafkaLog))
	if err != nil {
		t.Fatalf("expected parsing action not to fail, got: %v", err)
	}
}

func TestOmkafkaToPoints(t *testing.T) {
	pstat, err := newOmkafkaFromJSON([]byte(omkafkaLog))
	if err != nil {
		t.Fatalf("expected parsing action not to fail, got: %v", err)
	}
	points := pstat.toPoints()

	testCases := []*point{
		{
			Name:       "input_submitted",
			Type:       counter,
			Value:      59,
			LabelValue: "omkafka",
		},
		{
			Name:       "omkafka_messages",
			Type:       counter,
			Value:      59,
			LabelValue: "submitted",
		},
		{
			Name:  "omkafka_maxoutqsize",
			Type:  counter,
			Value: 9,
		},
		{
			Name:       "omkafka_messages",
			Type:       counter,
			Value:      0,
			LabelValue: "failures",
		},
		{
			Name:       "omkafka_topicdynacache",
			Type:       counter,
			Value:      57,
			LabelValue: "skipped",
		},
		{
			Name:       "omkafka_topicdynacache",
			Type:       counter,
			Value:      2,
			LabelValue: "miss",
		},
		{
			Name:       "omkafka_topicdynacache",
			Type:       counter,
			Value:      0,
			LabelValue: "evicted",
		},
		{
			Name:       "omkafka_messages",
			Type:       counter,
			Value:      55,
			LabelValue: "acked",
		},
		{
			Name:       "omkafka_failures",
			Type:       counter,
			Value:      0,
			LabelValue: "msg_too_large",
		},

		{
			Name:       "omkafka_failures",
			Type:       counter,
			Value:      0,
			LabelValue: "unknown_topic",
		},
		{
			Name:       "omkafka_failures",
			Type:       counter,
			Value:      0,
			LabelValue: "queue_full",
		},
		{
			Name:       "omkafka_failures",
			Type:       counter,
			Value:      0,
			LabelValue: "unknown_partition",
		},
		{
			Name:       "omkafka_failures",
			Type:       counter,
			Value:      0,
			LabelValue: "other",
		},
		{
			Name:       "omkafka_errors",
			Type:       counter,
			Value:      0,
			LabelValue: "timed_out",
		},
		{
			Name:       "omkafka_errors",
			Type:       counter,
			Value:      0,
			LabelValue: "transport",
		},
		{
			Name:       "omkafka_errors",
			Type:       counter,
			Value:      0,
			LabelValue: "broker_down",
		},
		{
			Name:       "omkafka_errors",
			Type:       counter,
			Value:      0,
			LabelValue: "auth",
		},
		{
			Name:       "omkafka_errors",
			Type:       counter,
			Value:      0,
			LabelValue: "ssl",
		},
		{
			Name:       "omkafka_errors",
			Type:       counter,
			Value:      0,
			LabelValue: "other",
		},
		{
			Name:  "omkafka_rtt_avg_usec_acg",
			Type:  gauge,
			Value: 0,
		},
		{
			Name:  "omkafka_throttle_avg_msec_avg",
			Type:  gauge,
			Value: 0,
		},
		{
			Name:  "omkafka_int_latency_avg_usec_avg",
			Type:  gauge,
			Value: 0,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("point idx %d", idx), func(t *testing.T) {
			p := points[idx]
			if p.Name != tc.Name {
				t.Errorf("got name %s; want %s", p.Name, tc.Name)
			}
			if p.Type != tc.Type {
				t.Errorf("got type %d;  %d", p.Type, tc.Type)
			}
			if p.Value != tc.Value {
				t.Errorf("got value %d;  %d", p.Value, tc.Value)
			}
			if p.LabelValue != tc.LabelValue {
				t.Errorf("got label value %s;  %s", p.LabelValue, tc.LabelValue)
			}
		})
	}

}
