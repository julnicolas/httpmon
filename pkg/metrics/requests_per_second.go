package metrics

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/julnicolas/httpmon/pkg/trace"
)

const (
	ReqsPerS string = "ReqsPerSecond"
)

// RequestsPerSecond -> rename avg requests per second
type RequestsPerSecond struct {
	mutex sync.Mutex
	// time last capture started
	// capture are spaced of scrape period time
	lastCapture time.Time
	start       time.Time            // start of ongoing capture period
	period      time.Duration        // Period is the collection period to compute
	total       []float64            // data points, series of previous req/s values
	perSection  map[string][]float64 // per-section series of previous req/s values
}

func NewRequestsPerSecond(duration time.Duration) *RequestsPerSecond {
	// This should have been validated before, should never happen
	if duration < time.Second {
		err := fmt.Errorf("critical, duration is below 1s, input : %s", duration)
		panic(err)
	}

	return &RequestsPerSecond{
		period:     duration,
		perSection: make(map[string][]float64),
	}
}

// Update computes the request rate, it is assumed entries are time-sorted
// in increasing order (increasingly recent)
func (o *RequestsPerSecond) Update(t trace.Trace) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	zero := time.Time{}
	if o.start == zero {
		o.start = t.Date
		o.total = make([]float64, 1)
	}

	// If true, data arrive from more recent time window
	// so we need a new one to compute recent values
	if t.Date.After(o.start.Add(o.period)) {
		o.lastCapture = o.start
		o.start = t.Date
		o.total[len(o.total)-1] /= o.period.Seconds()
		o.total = append(o.total, 0.0)

		// The new slice reference would expire if using the value
		// so let's make sure to store it in the object's map
		for section := range o.perSection {
			o.perSection[section][len(o.perSection[section])-1] /= o.period.Seconds()
			o.perSection[section] = append(o.perSection[section], 0.0)
		}
	}

	// Count requests globally and per section on active time window
	o.total[len(o.total)-1] += 1
	if len(o.perSection[t.Section]) == 0 {
		o.perSection[t.Section] = make([]float64, 1)
	}
	section := o.perSection[t.Section]
	section[len(section)-1] += 1
}

// DeepCopy Returns a metric out of a deep copy of internal structures.
// It is thread-safe though more expensive as locking Update on top of a copy
func (o *RequestsPerSecond) DeepCopy() Metric {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	var newTotal []float64
	if len(o.total) > 0 {
		newTotal = make([]float64, len(o.total)-1)
		copy(newTotal, o.total)
	}

	newPerSection := make(map[string][]float64, len(o.perSection))
	for k, vector := range o.perSection {
		section := strings.Clone(k)

		if len(o.perSection[k]) > 0 {
			newVector := make([]float64, len(vector)-1)
			copy(newVector, vector)
			newPerSection[section] = newVector
		}
	}

	return CounterVector{
		time:   o.lastCapture.Unix(),
		name:   ReqsPerS,
		total:  newTotal,
		labels: newPerSection,
	}
}

// Metric outputs a metric measuring the number or requests per second
func (o *RequestsPerSecond) Metric() Metric {
	t := o.total
	if len(o.total) > 0 {
		t = t[:len(t)-1]
	}
	return CounterVector{
		time:   o.lastCapture.Unix(),
		name:   ReqsPerS,
		total:  t,            // last window is ongoing so result is not an average yet
		labels: o.perSection, // false here because last cell must not be considered
	}
}
