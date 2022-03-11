package main

import "testing"

var (
	forwardLog = []byte(`{ "name": "TCP-FQDN-6514", "origin": "omfwd", "bytes.sent": 666 }`)
)

func TestNewForwardFromJSON(t *testing.T) {
	logType := getStatType(forwardLog)
	if logType != rsyslogForward {
		t.Errorf("detected pstat type should be %d but is %d", rsyslogForward, logType)
	}

	pstat, err := newForwardFromJSON([]byte(forwardLog))
	if err != nil {
		t.Fatalf("expected parsing action not to fail, got: %v", err)
	}

	if want, got := "TCP-FQDN-6514", pstat.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	if want, got := int64(666), pstat.BytesSent; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}
}

func TestForwardToPoints(t *testing.T) {
	pstat, err := newForwardFromJSON([]byte(forwardLog))
	if err != nil {
		t.Fatalf("expected parsing action not to fail, got: %v", err)
	}
	points := pstat.toPoints()

	point := points[0]
	if want, got := "forward_bytes_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	if want, got := int64(666), point.Value; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := counter, point.Type; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := "TCP-FQDN-6514", point.LabelValue; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}
}
