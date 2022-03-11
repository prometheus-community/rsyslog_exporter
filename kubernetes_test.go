package main

import "testing"

var (
	kubernetesLog = []byte(`{ "name": "mmkubernetes(https://host.domain.tld:6443)", "origin": "mmkubernetes", "recordseen": 477943, "namespacemetadatasuccess": 7, "namespacemetadatanotfound": 0, "namespacemetadatabusy": 0, "namespacemetadataerror": 0, "podmetadatasuccess": 26, "podmetadatanotfound": 0, "podmetadatabusy": 0, "podmetadataerror": 0 }`)
)

func TestNewKubernetesFromJSON(t *testing.T) {
	logType := getStatType(kubernetesLog)
	if logType != rsyslogKubernetes {
		t.Errorf("detected pstat type should be %d but is %d", rsyslogKubernetes, logType)
	}

	pstat, err := newKubernetesFromJSON([]byte(kubernetesLog))
	if err != nil {
		t.Fatalf("expected parsing action not to fail, got: %v", err)
	}

	if want, got := "mmkubernetes(https://host.domain.tld:6443)", pstat.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	if want, got := "https://host.domain.tld:6443", pstat.Url; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	if want, got := int64(477943), pstat.RecordSeen; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := int64(7), pstat.NamespaceMetaSuccess; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := int64(0), pstat.NamespaceMetaNotFound; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := int64(0), pstat.NamespaceMetaBusy; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := int64(0), pstat.NamespaceMetaError; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := int64(26), pstat.PodMetaSuccess; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := int64(0), pstat.PodMetaNotFound; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := int64(0), pstat.PodMetaBusy; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

	if want, got := int64(0), pstat.PodMetaError; want != got {
		t.Errorf("wanted '%d', got '%d'", want, got)
	}

}

func TestKubernetesToPoints(t *testing.T) {
	pstat, err := newKubernetesFromJSON([]byte(kubernetesLog))
	if err != nil {
		t.Fatalf("expected parsing action not to fail, got: %v", err)
	}
	points := pstat.toPoints()

	point := points[0]
	if want, got := "kubernetes_namespace_metadata_success_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	if want, got := "https://host.domain.tld:6443", point.LabelValue; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[1]
	if want, got := "kubernetes_namespace_metadata_notfound_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[2]
	if want, got := "kubernetes_namespace_metadata_busy_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[3]
	if want, got := "kubernetes_namespace_metadata_error_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[4]
	if want, got := "kubernetes_pod_metadata_success_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[5]
	if want, got := "kubernetes_pod_metadata_notfound_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[6]
	if want, got := "kubernetes_pod_metadata_busy_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[7]
	if want, got := "kubernetes_pod_metadata_error_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}

	point = points[8]
	if want, got := "kubernetes_record_seen_total", point.Name; want != got {
		t.Errorf("wanted '%s', got '%s'", want, got)
	}
}
