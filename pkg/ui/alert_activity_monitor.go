package ui

import (
	"fmt"
	"sort"

	"github.com/julnicolas/httpmon/pkg/alert"
)

// ActivityMonitor monitors alert states
type ActivityMonitor struct {
	states map[alert.State]map[alert.NameT]alert.Alert
}

// alert_slice is a sortable alert slice. Items are sorted by Name().
type alert_slice []alert.Alert

func (o alert_slice) Len() int           { return len(o) }
func (o alert_slice) Less(i, j int) bool { return o[i].Name() < o[j].Name() }
func (o alert_slice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

func NewActivityMonitor() *ActivityMonitor {
	o := ActivityMonitor{
		states: make(map[alert.State]map[alert.NameT]alert.Alert),
	}

	o.states[alert.Inactive] = make(map[alert.NameT]alert.Alert)
	o.states[alert.Pending] = make(map[alert.NameT]alert.Alert)
	o.states[alert.Active] = make(map[alert.NameT]alert.Alert)

	return &o
}

func (o *ActivityMonitor) Monitor(a alert.Alert) {
	name := a.Name()
	current := a.State()
	u, v := current.Others()

	delete(o.states[u], name)
	delete(o.states[v], name)
	o.states[current][name] = a
}

// ActivityTxt returns a string describing current alert activity
// usable as is for a view display
func (o *ActivityMonitor) ActivityTxt() string {
	inactive, pending, active := gActivityMonitor.activitySlices()

	txt := fmt.Sprintf("Enabled Alerts: %d\n\n", o.Enabled())
	txt += o.activitySliceTxt("Active", active)
	txt += o.activitySliceTxt("Pending", pending)
	txt += o.activitySliceTxt("Inactive", inactive)
	return txt
}

// Enabled returns the number of enabled metrics
func (o *ActivityMonitor) Enabled() int {
	return len(o.states[alert.Inactive]) + len(o.states[alert.Pending]) + len(o.states[alert.Active])
}

func (o *ActivityMonitor) activitySliceTxt(header string, alerts alert_slice) string {
	txt := header + ":\n"
	for _, a := range alerts {
		txt += fmt.Sprintf("  Name: %s\n  Description: %s\n\n", a.Name(), a.Description())
	}
	return txt
}

func (o *ActivityMonitor) activitySlices() (inactive alert_slice, pending alert_slice, active alert_slice) {
	return o.activitySlice(alert.Inactive),
		o.activitySlice(alert.Pending),
		o.activitySlice(alert.Active)
}

// fill populates the input alert slice, then sorts its content
func (o *ActivityMonitor) activitySlice(state alert.State) alert_slice {
	s := make(alert_slice, 0, len(o.states[state]))

	for _, alert := range o.states[state] {
		s = append(s, alert)
	}
	sort.Sort(s)

	return s
}
