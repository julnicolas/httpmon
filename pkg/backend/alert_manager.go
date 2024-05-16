package backend

import (
	"fmt"
	"time"

	"github.com/julnicolas/httpmon/pkg/alert"
	"github.com/julnicolas/httpmon/pkg/metrics"
)

type AlertManager struct {
	// period is the evaluation period for alerting rules
	period time.Duration
	// enabled list enabled alerts, keys are alert names
	enabled map[alert.NameT]bool
	alerts  map[alert.NameT]AlertStateTransition
	reqPerS *alert.RequestsPerSecond
	states  chan AlertStateTransition
}

type AlertStateTransition struct {
	// Prev is the previous alert state, if different
	// than current, send it.
	// Initialised to alert.Inactive
	Prev          alert.State
	Alert         alert.Alert
	Time          int64 // Evaluation time in Unix seconds
	publishedOnce bool
}

func NewAlertManager(eval time.Duration, in alert.RequestsPerSecondInput) *AlertManager {
	enabled := make(map[alert.NameT]bool, 1)
	enabled[alert.ReqPerS] = true

	return &AlertManager{
		period:  eval,
		enabled: enabled,
		alerts:  make(map[alert.NameT]AlertStateTransition),
		reqPerS: alert.NewRequestsPerSecond(in),
		states:  make(chan AlertStateTransition, 100),
	}
}

// Eval evaluates alerts, making them available in Alerts()
// if alert state has changed
func (o *AlertManager) Eval(m metrics.Metric) {
	for name, enabled := range o.enabled {
		if enabled {
			// Try to publish every evaluated alerts
			// Keep in memory their previous state
			a := o.evalAlert(m)

			t, ok := o.alerts[name]
			t.Alert = a
			t.Time = a.EvalTime().Unix()
			o.publish(t)

			if ok {
				t.Prev = t.Alert.State()
				t.publishedOnce = true
			}
			o.alerts[name] = t
		}
	}
}

// publish publishes an alert with its previous state if the current state is different
// publish all alerts once in their inactive state so that clients can now what alerts
// are going to be published.
func (o *AlertManager) publish(t AlertStateTransition) {
	if !t.publishedOnce || t.Prev != t.Alert.State() {
		// Alerts are pointer receivers for most parts so they can be changed
		// even after being transmitted
		n := AlertStateTransition{
			Prev:  t.Prev,
			Alert: t.Alert.DeepCopy(),
			Time:  t.Time,
		}
		o.states <- n
	}
}

func (o *AlertManager) evalAlert(m metrics.Metric) alert.Alert {
	switch m.Name() {
	case metrics.ReqsPerS:
		o.reqPerS.Eval(m)
		return o.reqPerS
	default:
		// Should never happen so let's get notified early on
		err := fmt.Errorf("configuration error - alert %s is not supported", m.Name())
		panic(err)
	}
}

// Alerts only exposes alerts which state's have changed
// Every alert is at least sent once when it is initialised as inactive
func (o *AlertManager) Alerts() <-chan AlertStateTransition {
	return o.states
}
