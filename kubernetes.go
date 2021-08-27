package main

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var (
	apiNameRegexp = regexp.MustCompile(`mmkubernetes\((\S+)\)`)
)

type kubernetes struct {
	Name                  string `json:"name"`
	Url                   string
	RecordSeen            int64 `json:"recordseen"`
	NamespaceMetaSuccess  int64 `json:"namespacemetadatasuccess"`
	NamespaceMetaNotFound int64 `json:"namespacemetadatanotfound"`
	NamespaceMetaBusy     int64 `json:"namespacemetadatabusy"`
	NamespaceMetaError    int64 `json:"namespacemetadataerror"`
	PodMetaSuccess        int64 `json:"podmetadatasuccess"`
	PodMetaNotFound       int64 `json:"podmetadatanotfound"`
	PodMetaBusy           int64 `json:"podmetadatabusy"`
	PodMetaError          int64 `json:"podmetadataerror"`
}

func newKubernetesFromJSON(b []byte) (*kubernetes, error) {
	var pstat kubernetes
	err := json.Unmarshal(b, &pstat)
	if err != nil {
		return nil, fmt.Errorf("failed to decode kubernetes stat `%v`: %v", string(b), err)
	}
	matches := apiNameRegexp.FindSubmatch([]byte(pstat.Name))
	if matches != nil {
		pstat.Url = string(matches[1])
	}
	return &pstat, nil
}

func (k *kubernetes) toPoints() []*point {
	points := make([]*point, 9)

	points[0] = &point{
		Name:        "kubernetes_namespace_metadata_success_total",
		Type:        counter,
		Value:       k.NamespaceMetaSuccess,
		Description: "successful fetches of namespace metadata",
		LabelName:   "url",
		LabelValue:  k.Url,
	}

	points[1] = &point{
		Name:        "kubernetes_namespace_metadata_notfound_total",
		Type:        counter,
		Value:       k.NamespaceMetaNotFound,
		Description: "notfound fetches of namespace metadata",
		LabelName:   "url",
		LabelValue:  k.Url,
	}

	points[2] = &point{
		Name:        "kubernetes_namespace_metadata_busy_total",
		Type:        counter,
		Value:       k.NamespaceMetaBusy,
		Description: "busy fetches of namespace metadata",
		LabelName:   "url",
		LabelValue:  k.Url,
	}

	points[3] = &point{
		Name:        "kubernetes_namespace_metadata_error_total",
		Type:        counter,
		Value:       k.NamespaceMetaError,
		Description: "error fetches of namespace metadata",
		LabelName:   "url",
		LabelValue:  k.Url,
	}

	points[4] = &point{
		Name:        "kubernetes_pod_metadata_success_total",
		Type:        counter,
		Value:       k.PodMetaSuccess,
		Description: "successful fetches of pod metadata",
		LabelName:   "url",
		LabelValue:  k.Url,
	}

	points[5] = &point{
		Name:        "kubernetes_pod_metadata_notfound_total",
		Type:        counter,
		Value:       k.PodMetaNotFound,
		Description: "notfound fetches of pod metadata",
		LabelName:   "url",
		LabelValue:  k.Url,
	}

	points[6] = &point{
		Name:        "kubernetes_pod_metadata_busy_total",
		Type:        counter,
		Value:       k.PodMetaBusy,
		Description: "busy fetches of pod metadata",
		LabelName:   "url",
		LabelValue:  k.Url,
	}

	points[7] = &point{
		Name:        "kubernetes_pod_metadata_error_total",
		Type:        counter,
		Value:       k.PodMetaError,
		Description: "error fetches of pod metadata",
		LabelName:   "url",
		LabelValue:  k.Url,
	}

	points[8] = &point{
		Name:        "kubernetes_record_seen_total",
		Type:        counter,
		Value:       k.RecordSeen,
		Description: "records fetched from the api",
		LabelName:   "url",
		LabelValue:  k.Url,
	}

	return points
}
