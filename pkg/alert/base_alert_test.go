package alert

import (
	"testing"
	"time"

	"github.com/julnicolas/httpmon/pkg/metrics"
	"github.com/stretchr/testify/assert"
)

type FalseTimer struct{}

func (o FalseTimer) Start() UnixSeconds { return 0 }
func (o FalseTimer) Over() bool         { return false }
func (o FalseTimer) Now() UnixSeconds   { return 0 }

type TrueTimer struct{}

func (o TrueTimer) Start() UnixSeconds { return 0 }
func (o TrueTimer) Over() bool         { return true }
func (o TrueTimer) Now() UnixSeconds   { return 0 }

// It would be better to mock time but... that takes a bit of time!
func TestEvalReturnsInactiveWithFalseRule(t *testing.T) {
	placeholder := metrics.Counter{}
	expected := Inactive
	a := NewBaseAlert(time.Second,
		func(metrics.Metric) bool { return false })
	a.timer = FalseTimer{}

	// exercise
	state0 := a.Eval(placeholder)
	state1 := a.Eval(placeholder)
	state2 := a.Eval(placeholder)

	// verification
	assert.Equal(t, expected, state0)
	assert.Equal(t, expected, state1)
	assert.Equal(t, expected, state2)
}

func TestEvalReturnsActiveIfRuleIsTrueAndTimeCOnditionsAreMet(t *testing.T) {
	placeholder := metrics.Counter{}
	a := NewBaseAlert(time.Second,
		func(metrics.Metric) bool { return true }) // the alert is on
	a.timer = FalseTimer{} // timer is not over yet

	// rule checks but timer is not over so state must be pending
	state1 := a.Eval(placeholder)
	// now the timer is over so the alert should be active
	a.timer = TrueTimer{}
	state2 := a.Eval(placeholder)

	// verification
	assert.Equal(t, Pending, state1)
	assert.Equal(t, Active, state2)
}

/*

To the reviewers: WIP on a more complex testing scenario

This is not necessarily in the right file but I left it there
for your information.

// iterationTimer iterates over a slice
// it is over when the slice reaches the end
type iterationTimer struct {
	I     int       // current index
	Slice []float64 // slice of value to iterate over
}

// Crashes intentionaly if slice is nil/len 0
func (o *iterationTimer) Start() int64 { o.I = 1; return 0 }
func (o *iterationTimer) Over() bool   { o.I++; return o.I >= len(o.Slice) }
func (o *iterationTimer) Now() int64   { return 0 }

type testMetric struct {
	V float64
}

func (o testMetric) ScrapeTime() int64   { return 1 }
func (o testMetric) Name() string        { return "test metric" }
func (o testMetric) Labels() interface{} { return "labels are not used" }

func TestEvalReturnsAppropriateValuesRemovingTimeDependency(t *testing.T) {

	threshold := 10.0
	a := NewBaseAlert(time.Second, func(m metrics.Metric) bool { return m.(testMetric).V >= threshold })

	inactive1 := []float64{0, 1, 2, 3}
	pending1 := []float64{10, 11, 12, 13, 14, 15}
	active1 := []float64{10, 11, 12, 13, 14}
	inactive2 := []float64{9, 8, 7}
	//pending2 := [7]float64{10, 11, 12}
	//inactive3 := [3]float64{9, 8, 7}

	a.timer = &iterationTimer{
		Slice: inactive1,
	}

	for range inactive1 {
		a.Eval()
		assert.Equal(t, alert.Inactive, a.State())
	}
}
*/
