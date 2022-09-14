package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type pointType int

const (
	counter pointType = iota
	gauge
)

type point struct {
	Name        string
	Description string
	Type        pointType
	Value       int64
	LabelName   string
	LabelValue  string
}

func (p *point) promDescription() *prometheus.Desc {
	variableLabels := []string{}
	if p.promLabelName() != "" {
		variableLabels = []string{p.promLabelName()}
	}
	return prometheus.NewDesc(
		prometheus.BuildFQName("", "rsyslog", p.Name),
		p.Description,
		variableLabels,
		nil,
	)
}

func (p *point) promType() prometheus.ValueType {
	if p.Type == counter {
		return prometheus.CounterValue
	}
	return prometheus.GaugeValue
}

func (p *point) promValue() float64 {
	return float64(p.Value)
}

func (p *point) promLabelValue() string {
	return p.LabelValue
}

func (p *point) promLabelName() string {
	return p.LabelName
}

func (p *point) key() string {
	return fmt.Sprintf("%s.%s", p.Name, p.LabelValue)
}
