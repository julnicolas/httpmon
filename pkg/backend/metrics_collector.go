package backend

import (
	"fmt"

	"github.com/julnicolas/httpmon/pkg/metrics"
	"github.com/julnicolas/httpmon/pkg/trace"
)

// MetricsCollector is an object managing metrics collection
type MetricsCollector struct {
	probers map[string]metrics.Prober
}

// NewMetricsCollector creates a new object to collect metrics on every loop cycle
func NewMetricsCollector(probers []metrics.Prober) *MetricsCollector {
	c := &MetricsCollector{
		probers: make(map[string]metrics.Prober, len(probers)),
	}

	for _, p := range probers {
		c.probers[p.Metric().Name()] = p
	}

	return c
}

// Collect runs all probers on the trace to update all metrics
// Is meant to be repetively called on every loop cycle
func (o *MetricsCollector) Collect(t trace.Trace) error {
	for _, p := range o.probers {
		p.Update(t)
	}

	return nil
}

// DeepCopy returns a deep copy of the metric struct.
// It is thread-safe, more expensive as locking non atomic structures
// on top of deep-copying them
func (o *MetricsCollector) DeepCopy(name string) (metrics.Metric, error) {
	m, err := o.lookupMetric(name)
	if err != nil {
		var null metrics.Metric
		return null, err
	}

	return m.DeepCopy(), err
}

// Metric returns the collected metric of name 'name' if existent, an error otherwise
func (o *MetricsCollector) Metric(name string) (metrics.Metric, error) {
	panic("deprecated")

	/*
		m, err := o.lookupMetric(name)
		if err != nil {
			var null metrics.Metric
			return null, err
		}

		return m.Metric(), nil
	*/
}

func (o *MetricsCollector) lookupMetric(name string) (metrics.Prober, error) {
	m, ok := o.probers[name]
	if !ok {
		return nil, fmt.Errorf("metric %s not found", name)
	}
	return m, nil
}
