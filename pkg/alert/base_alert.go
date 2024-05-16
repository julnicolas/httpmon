package alert

import (
	"time"

	"github.com/julnicolas/httpmon/pkg/metrics"
)

// AlertRule is a predicate run against the alert.
// if true the alert condition if verified, which modifies its state.
type AlertRule func(metrics.Metric) bool

// BaseAlert defines fields and methods common to all alerts
// it should be embedded in any deriving struct
type BaseAlert struct {
	timer    Timer
	period   time.Duration
	rule     AlertRule
	evalTime int64 // evaluation time
	state    State
}

// NewBaseAlert creates a new alert
func NewBaseAlert(period time.Duration, rule AlertRule) *BaseAlert {
	return &BaseAlert{
		timer:  NewSystemTimer(period),
		period: period,
		rule:   rule,
	}
}

func (o *BaseAlert) EvalTime() time.Time {
	return time.Unix(o.evalTime, 0)
}

// State returns the alert rule's last-evaluated-state
func (o *BaseAlert) State() State {
	return o.state
}

// Eval runs the alert rule then returns its current state
func (o *BaseAlert) Eval(m metrics.Metric) State {
	if o.rule(m) {
		switch o.state {
		case Inactive:
			o.timer.Start()
			o.state = Pending
		case Pending:
			if o.timer.Over() {
				o.state = Active
			}
		}
	} else {
		o.state = Inactive
	}

	o.evalTime = o.timer.Now()
	return o.state
}

func (o *BaseAlert) Description() string {
	return "BaseAlert"
}

func (o *BaseAlert) Name() NameT {
	return "BaseAlert"
}

func (o *BaseAlert) DeepCopy() Alert {
	n := new(BaseAlert)
	*n = *o
	return n
}
