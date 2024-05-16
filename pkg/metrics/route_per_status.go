package metrics

import (
	"strings"
	"sync"
	"time"

	"github.com/julnicolas/httpmon/pkg/trace"
)

const (
	RoutesPerStatusN string = "RoutesPerStatus"
)

// StatusCodeT is type alias representing http status codes
type StatusCodeT = uint

// SectionT is type alias representing http path section
type SectionT = string

type StatusMap = map[StatusCodeT]map[SectionT]float64

// RoutePerStatus is a status code per route counter
type RoutePerStatus struct {
	lastScrape time.Time
	total      float64
	mutex      sync.Mutex
	perStatus  StatusMap
}

func NewRoutePerStatus() *RoutePerStatus {
	return &RoutePerStatus{
		perStatus: make(map[StatusCodeT]map[SectionT]float64),
	}
}

func (o *RoutePerStatus) Update(t trace.Trace) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.lastScrape = t.Date

	m := o.perStatus[t.Status]
	if m == nil {
		m = make(map[SectionT]float64)
		o.perStatus[t.Status] = m
	}
	o.total += 1

	m[t.Section] += 1
}

func (o *RoutePerStatus) DeepCopy() Metric {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	cp := make(StatusMap, len(o.perStatus))
	for status, sectionMap := range o.perStatus {
		smap := make(map[SectionT]float64, len(sectionMap))
		for section, counter := range sectionMap {
			sec := strings.Clone(section)
			smap[sec] = counter
		}
		cp[status] = smap
	}

	return RoutePerStatusCounter{
		time:   o.lastScrape.Unix(),
		name:   RoutesPerStatusN, // constant so pointer doesn't change
		total:  o.total,          // atomic copy
		labels: cp,               // copied using the above lock
	}
}

func (o *RoutePerStatus) Metric() Metric {
	return RoutePerStatusCounter{
		time:   o.lastScrape.Unix(),
		name:   RoutesPerStatusN,
		total:  o.total,
		labels: o.perStatus,
	}
}

type RoutePerStatusCounter struct {
	time   int64 // scrape time - unix
	name   string
	total  float64
	labels map[StatusCodeT]map[SectionT]float64
}

func (o RoutePerStatusCounter) ScrapeTime() int64 {
	return o.time
}

func (o RoutePerStatusCounter) Name() string {
	return o.name
}

func (o RoutePerStatusCounter) Total() float64 {
	return o.total
}

func (o RoutePerStatusCounter) Labels() interface{} {
	return o.labels
}

func (o RoutePerStatusCounter) TypedLabels() StatusMap {
	return o.labels
}
