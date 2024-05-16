package metrics

import (
	"fmt"
	"time"
)

type CounterVector struct {
	time   int64 // scrape time - unix seconds
	name   string
	total  []float64
	labels map[string][]float64
}

func (o CounterVector) String() string {
	var last float64
	if len(o.total) > 0 {
		last = o.total[len(o.total)-1]
	}
	msg := fmt.Sprintf("Metric:\ntime: %s\nname: %s\nlast: %f\n",
		time.Unix(o.time, 0), o.name, last)

	return msg
}

func (o CounterVector) ScrapeTime() int64 {
	return o.time
}

func (o CounterVector) Name() string {
	return o.name
}

func (o CounterVector) Total() []float64 {
	return o.total
}

func (o CounterVector) Labels() interface{} {
	return o.labels
}

func (o CounterVector) TypedLabels() map[string][]float64 {
	return o.labels
}
