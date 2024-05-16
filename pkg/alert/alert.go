package alert

import (
	"fmt"
	"time"

	"github.com/julnicolas/httpmon/pkg/metrics"
)

type Alert interface {
	// Name is the alert's name
	Name() NameT
	// State returns the current alert state
	State() State
	// Eval evaluates the current alert state.
	// State can be:
	// - inactive -> the alert condition is false
	// - pending -> alert condition has been verified but
	// 		not on the whole alert period
	// - active -> the alert is on
	Eval(m metrics.Metric) State
	// EvalTime is the last evaluation time
	EvalTime() time.Time
	// Description is a human readable text describing
	// the alert.
	Description() string
	// Creates an entirely new copy, see it as a copy constructor
	//
	// TODO: Isolate this into a new Specialised interface
	// DeepCopier with specialised alert copier objects
	// since this feature is only used by AlertManager for
	// internal reasons
	DeepCopy() Alert
}

// NameT describes an alert name
// Alerts and metrics can be called the same so that's good
// to let the compiler disambiguate
type NameT string

// State describes an alert state
type State uint

const (
	// Inactive means the alert is not active
	Inactive State = 0
	// Pending means that the alert rule has been returning true but not yet on
	// whole evaluation interval. This is a potential alert at this stage.
	Pending State = 1
	// Active means that the alert is considered active - the alert rule has checked
	// for the whole evaluation interval.
	Active State = 2
)

// String returns a human readable value for states
func (o State) String() string {
	switch o {
	case Inactive:
		return "Inactive"
	case Pending:
		return "Pending"
	case Active:
		return "Active"
	default:
		return "Not Supported"
	}
}

// Others returns the two other possible states
func (o State) Others() (State, State) {
	switch o {
	case Inactive:
		return Pending, Active
	case Pending:
		return Inactive, Active
	case Active:
		return Inactive, Pending
	default:
		panic(fmt.Errorf("invalid alert state %d", o))
	}
}
