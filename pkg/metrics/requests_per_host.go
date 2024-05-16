package metrics

import (
	"strings"
	"sync"
	"time"

	"github.com/julnicolas/httpmon/pkg/trace"
)

const (
	ReqsPerHost string = "ReqsPerHost"
)

type RequestsPerHost struct {
	mutex      sync.Mutex
	lastScrape time.Time
	total      float64
	perHost    map[string]float64
}

func NewRequestsPerHost() *RequestsPerHost {
	return &RequestsPerHost{
		perHost: make(map[string]float64),
	}
}

// Count counts http calls per host and globally
func (o *RequestsPerHost) Update(t trace.Trace) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.lastScrape = t.Date
	o.total += 1
	o.perHost[t.RemoteHost] += 1
}

// DeepCopy Returns a metric out of a deep copy of internal structures.
// It is thread-safe though more expensive as locking Update on top of a copy
func (o *RequestsPerHost) DeepCopy() Metric {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	new_ := make(map[string]float64, len(o.perHost))
	for k, v := range o.perHost {
		section := strings.Clone(k)
		new_[section] = v
	}

	return Counter{
		time:   o.lastScrape.Unix(),
		name:   ReqsPerHost,
		total:  o.total,
		labels: new_,
	}
}

// Metric exposes a metric representing the number of requests
// emitted by a host and the total of requests.
func (o *RequestsPerHost) Metric() Metric {
	return Counter{
		time:   o.lastScrape.Unix(),
		name:   ReqsPerHost,
		total:  o.total,
		labels: o.perHost,
	}
}
