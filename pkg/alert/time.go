package alert

import (
	"fmt"
	"time"

	"github.com/julnicolas/httpmon/pkg/metrics"
)

// UnixSeconds is the number of seconds since epoch
type UnixSeconds = int64

// Timer is an interface to implement a timer
type Timer interface {
	// Start starts the timer now
	Start() UnixSeconds
	// Over returns true if now is later than start + duration
	// Said differently the timer reached 0
	Over() bool
	// Now returns the now moment relevant to the timer
	Now() UnixSeconds
}

// MetricTimer takes a metric bases its time calculation
// based on a reference metrics called 'start metrics'
//
// Do not use time.Time so that the timer can be copied
// time.Time uses internal pointers
type MetricsTimer struct {
	start    UnixSeconds
	duration time.Duration
	current  metrics.Metric
}

func NewMetricsTimer(duration time.Duration) *MetricsTimer {
	return &MetricsTimer{
		duration: duration,
	}
}

func (o *MetricsTimer) String() string {
	str := fmt.Sprintf("start:%s\ncurrent:%s\nduration:%s\n",
		time.Unix(o.start, 0),
		time.Unix(o.current.ScrapeTime(), 0),
		o.duration)
	return str
}

// Now returns current metrics scrape time
func (o *MetricsTimer) Now() UnixSeconds {
	return o.current.ScrapeTime()
}

// Metric sets the current reference metric for time calculations
// It must be called BEFORE Start and Over.
// It starts the time origin.
func (o *MetricsTimer) Metric(m metrics.Metric) {
	o.current = m
}

// Start sets the start time to current metric's ScrapeTime()
func (o *MetricsTimer) Start() UnixSeconds {
	o.start = o.current.ScrapeTime()
	return o.start
}

// Over returns true when current metric's scrape time is after
// Start + period
func (o *MetricsTimer) Over() bool {
	t := time.Unix(o.current.ScrapeTime(), 0)
	s := time.Unix(o.start, 0)

	return t.After(s.Add(o.duration))
}

// SystemTimer is a Timer struct relying on system time.
// Concretely it uses time.Now()
type SystemTimer struct {
	start    UnixSeconds
	duration time.Duration
}

func NewSystemTimer(duration time.Duration) *SystemTimer {
	return &SystemTimer{
		duration: duration,
	}
}

func (o *SystemTimer) Now() UnixSeconds {
	return time.Now().Unix()
}

// Start sets the timer on from time.Now()
func (o *SystemTimer) Start() UnixSeconds {
	o.start = time.Now().Unix()
	return o.start
}

// Over is true when time.Now() is after start + duration
func (o *SystemTimer) Over() bool {
	return time.Now().After(time.Unix(o.start, 0).Add(o.duration))
}
