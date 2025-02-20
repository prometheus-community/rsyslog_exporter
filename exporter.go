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
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type rsyslogType int

const (
	rsyslogUnknown rsyslogType = iota
	rsyslogAction
	rsyslogInput
	rsyslogQueue
	rsyslogResource
	rsyslogDynStat
	rsyslogDynafileCache
	rsyslogInputIMDUP
	rsyslogForward
	rsyslogKubernetes
	rsyslogOmkafka
)

type rsyslogExporter struct {
	scanner *bufio.Scanner
	pointStore
}

func newRsyslogExporter() *rsyslogExporter {
	e := &rsyslogExporter{
		scanner: bufio.NewScanner(os.Stdin),
		pointStore: pointStore{
			pointMap: make(map[string]*point),
			lock:     &sync.RWMutex{},
		},
	}
	return e
}

func (re *rsyslogExporter) handleStatLine(rawbuf []byte) error {
	s := bytes.SplitN(rawbuf, []byte(" "), 4)
	if len(s) != 4 {
		return fmt.Errorf("failed to split log line, expected 4 columns, got: %v", len(s))
	}
	buf := s[3]

	pstatType := getStatType(buf)

	switch pstatType {
	case rsyslogAction:
		a, err := newActionFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range a.toPoints() {
			// nolint:errcheck
			re.set(p)
		}

	case rsyslogInput:
		i, err := newInputFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range i.toPoints() {
			// nolint:errcheck
			re.set(p)
		}

	case rsyslogInputIMDUP:
		u, err := newInputIMUDPFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range u.toPoints() {
			// nolint:errcheck
			re.set(p)
		}

	case rsyslogQueue:
		q, err := newQueueFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range q.toPoints() {
			// nolint:errcheck
			re.set(p)
		}

	case rsyslogResource:
		r, err := newResourceFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range r.toPoints() {
			// nolint:errcheck
			re.set(p)
		}
	case rsyslogDynStat:
		s, err := newDynStatFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range s.toPoints() {
			// nolint:errcheck
			re.set(p)
		}
	case rsyslogDynafileCache:
		d, err := newDynafileCacheFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range d.toPoints() {
			// nolint:errcheck
			re.set(p)
		}
	case rsyslogForward:
		f, err := newForwardFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range f.toPoints() {
			// nolint:errcheck
			re.set(p)
		}
	case rsyslogKubernetes:
		k, err := newKubernetesFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range k.toPoints() {
			// nolint:errcheck
			re.set(p)
		}
	case rsyslogOmkafka:
		o, err := newOmkafkaFromJSON(buf)
		if err != nil {
			return err
		}
		for _, p := range o.toPoints() {
			// nolint:errcheck
			re.set(p)
		}

	default:
		return fmt.Errorf("unknown pstat type: %v", pstatType)
	}
	return nil
}

// Describe sends the description of currently known metrics collected
// by this Collector to the provided channel. Note that this implementation
// does not necessarily send the "super-set of all possible descriptors" as
// defined by the Collector interface spec, depending on the timing of when
// it is called. The rsyslog exporter does not know all possible metrics
// it will export until the first full batch of rsyslog impstats messages
// are received via stdin. This is ok for now.
func (re *rsyslogExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName("", "rsyslog", "scrapes"),
		"times exporter has been scraped",
		nil, nil,
	)

	keys := re.keys()

	for _, k := range keys {
		p, err := re.get(k)
		if err != nil {
			ch <- p.promDescription()
		}
	}
}

// Collect is called by Prometheus when collecting metrics.
func (re *rsyslogExporter) Collect(ch chan<- prometheus.Metric) {
	keys := re.keys()

	for _, k := range keys {
		p, err := re.get(k)
		if err != nil {
			continue
		}

		labelValues := []string{}
		if p.promLabelValue() != "" {
			labelValues = []string{p.promLabelValue()}
		}
		metric := prometheus.MustNewConstMetric(
			p.promDescription(),
			p.promType(),
			p.promValue(),
			labelValues...,
		)

		ch <- metric
	}
}

func (re *rsyslogExporter) run(silent bool) {
	errorPoint := &point{
		Name:        "stats_line_errors",
		Type:        counter,
		Description: "Counts errors during stats line handling",
	}
	// nolint:errcheck
	re.set(errorPoint)
	for re.scanner.Scan() {
		err := re.handleStatLine(re.scanner.Bytes())
		if err != nil {
			errorPoint.Value += 1
			if !silent {
				log.Printf("error handling stats line: %v, line was: %s", err, re.scanner.Bytes())
			}
		}
	}
	if err := re.scanner.Err(); err != nil {
		log.Printf("error reading input: %v", err)
	}
	log.Print("input ended, exiting normally")
	os.Exit(0)
}
