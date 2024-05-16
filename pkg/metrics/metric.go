package metrics

import (
	"github.com/julnicolas/httpmon/pkg/trace"
)

// Prober is an interface used to measure, then expose metrics
type Prober interface {
	// Update updates the metrics measurement
	// Prober are meant to be called on every cycles so that
	// metrics can be updated iteratively
	Update(trace.Trace)
	// Metric returns a Metric struct representing a measured
	// phenomenom
	Metric() Metric
	// DeepCopy returns a deep copy of the metric so that it can
	// be used concurrently by other consummers.
	// Keep in mind that calling this function may lock metric
	// update during the copy
	DeepCopy() Metric
}

// Metric is a general Metrics interface, cast it to a concrete type
// to have a clear view on available values.
type Metric interface {
	// Unix time representation so that metrics can be copied
	ScrapeTime() int64
	Name() string
	Labels() interface{}
}
