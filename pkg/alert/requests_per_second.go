package alert

import (
	"fmt"
	"time"

	"github.com/julnicolas/httpmon/pkg/metrics"
)

// ReqPerS is the alert's name
const ReqPerS NameT = "Requests Per Second Threshold"

type RequestsPerSecond struct {
	MetricsTimeAlert
	threshold float64 // req/s thresohold, only used for Description
}

func (o *RequestsPerSecond) String() string {
	m := "RequestsPerSecondAlert:\n"
	m += fmt.Sprintf("time: %s", o.timer)
	m += fmt.Sprintf("\nstate: %s", o.State())
	m += fmt.Sprintf("\nname: %s\n", o.Name())
	m += fmt.Sprintf("\threshold: %f\n", o.threshold)

	return m
}

type RequestsPerSecondInput struct {
	Period time.Duration
	// This data would be made an interface to generalise alerting
	Threshold float64 // req/s threshold, if greater trigger alert
}

func NewRequestsPerSecond(in RequestsPerSecondInput) *RequestsPerSecond {
	return &RequestsPerSecond{
		MetricsTimeAlert: *NewMetricsTimeAlert(in.Period,
			func(m metrics.Metric) bool { return checkRate(m, in.Threshold) }),
		threshold: in.Threshold,
	}
}

// Name returns the alert's name
func (o *RequestsPerSecond) Name() NameT {
	return ReqPerS
}

// Description returns a human readable description of the alert
func (o *RequestsPerSecond) Description() string {
	return fmt.Sprintf(
		"active if reqs/s >= %.2f for %s",
		o.threshold,
		o.period)
}

func (o *RequestsPerSecond) DeepCopy() Alert {
	n := new(RequestsPerSecond)
	base := o.MetricsTimeAlert.DeepCopy()
	n.MetricsTimeAlert = *base.(*MetricsTimeAlert)
	n.threshold = o.threshold

	return n
}

func checkRate(m metrics.Metric, reqRateS float64) bool {
	// Stop immediatly, configuration problem
	// Must be checked at runtime now because typing is not strong enough
	// TODO: Improve typing so that it can be detected at build time
	if m.Name() != metrics.ReqsPerS {
		err := fmt.Errorf("configuration error, expecting %s as alert name but had %s instead", metrics.ReqsPerS, m.Name())
		panic(err)
	}
	c := m.(metrics.CounterVector) // same remark as above
	total := c.Total()
	if len(total) == 0 {
		return false
	}

	// Compare last measured value
	return total[len(total)-1] >= reqRateS
}
