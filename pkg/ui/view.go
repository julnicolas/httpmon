package ui

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/julnicolas/httpmon/pkg/alert"
	"github.com/julnicolas/httpmon/pkg/backend"
	"github.com/julnicolas/httpmon/pkg/metrics"
	"github.com/mum4k/termdash/container"
)

type View struct {
	main   *MainWindow
	layout container.Option
}

type kv struct {
	K string
	V float64
}

// kvslice is a key/value slice sorted by value in descending order
type kvslice []kv

func (o kvslice) Len() int           { return len(o) }
func (o kvslice) Less(i, j int) bool { return o[i].V > o[j].V }
func (o kvslice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

// kvslice is a kvslice sorted in lexical order
type kvv struct {
	K  string
	VV []float64
}
type kvvlexslice []kvv

func (o kvvlexslice) Len() int           { return len(o) }
func (o kvvlexslice) Less(i, j int) bool { return o[i].K < o[j].K }
func (o kvvlexslice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

func (o *View) ReqsPerHost(m metrics.Counter) {
	// Total number of requests
	txt := fmt.Sprintf("total: %d\n\nHosts:\n", int(m.Total()))

	// Sort by top contributor O(nlog(n))
	sorted := make(kvslice, 0, len(m.TypedLabels()))
	for k, v := range m.TypedLabels() {
		sorted = append(sorted, kv{K: k, V: v})
	}
	sort.Sort(sorted)

	for _, kv_ := range sorted {
		txt += fmt.Sprintf("  %s: %d\n", kv_.K, int(kv_.V))
	}

	// Compute top 5 repartition
	repartition := [5]float64{}
	for i := 0; i < min(len(sorted), 5); i++ {
		repartition[i] = (sorted[i].V / m.Total()) * 100
	}

	o.main.ReqsPerHost(txt, repartition[:])
}

func (o *View) ReqsPerSec(m metrics.CounterVector) {
	// Sort requests in lexicographic order
	perSection := m.TypedLabels()
	reqsPerSec := make(kvvlexslice, 0, len(perSection))
	for k, v := range perSection {
		reqsPerSec = append(reqsPerSec, kvv{K: k, VV: v})
	}
	sort.Sort(reqsPerSec)

	// TODO: Display average's time window
	// TODO: Display time window
	txt := "Last:\n  Per Section:\n"
	for _, req := range reqsPerSec {
		avg := 0.0
		if len(req.VV) > 0 {
			avg = req.VV[len(req.VV)-1]
		}
		txt += fmt.Sprintf("    %s: %.2f\n", req.K, avg)
	}

	o.main.ReqsPerSec(txt, m.Total())
}

func (o *View) RoutesPerStatus(m metrics.RoutePerStatusCounter) {
	// Data structures to be sorted to display the request listing
	sortedStatuses := make([]metrics.StatusCodeT, 0, len(m.TypedLabels()))
	sortedSections := make(map[metrics.StatusCodeT]kvslice, len(m.TypedLabels()))

	// compute repartition - aggregate similar status codes, summing their counters
	// similar status codes -> 2xx, 3xx, 4xx and 5xx
	repartition := [4]float64{} // 5xx, 4xx, 3xx, 2xx
	sectionPerStatus := m.TypedLabels()
	i := 0
	for status, sectionMap := range sectionPerStatus {
		// Use aggregation loop to layout data for request listing
		sortedStatuses = append(sortedStatuses, status)
		sortedSections[status] = make(kvslice, 0, len(sectionMap))

		// aggregate status codes
		switch {
		case status >= 500:
			i = 0
		case status >= 400:
			i = 1
		case status >= 300:
			i = 2
		case status >= 200:
			i = 3
		}

		for section, count := range sectionMap {
			// use aggregation loop to collect sections in a counter-sortable slice
			// used to display request listing
			sortedSections[status] = append(sortedSections[status], kv{K: section, V: count})

			// aggregate status codes
			repartition[i] += count
		}
	}

	// compute status code repartition in percent
	for i := range repartition {
		if m.Total() > 0 {
			repartition[i] = (repartition[i] / m.Total()) * 100
		}
	}

	// layout repartition for later rendering (last call)
	status := StatusPercent{
		Http5xx: repartition[0],
		Http4xx: repartition[1],
		Http3xx: repartition[2],
		Http2xx: repartition[3],
	}

	// DISPLAY REQUESTS BY STATUS CODE
	// requests (sections) are grouped by status (decreasing order)
	// and within status groups, by counter value (decreasing order)

	// Sort by status
	slices.Sort(sortedStatuses)

	// Sort by section counter
	for _, secSlice := range sortedSections {
		sort.Sort(secSlice)
	}

	// Generate string for to display
	txt := ""
	if len(sortedStatuses) > 0 {
		for i := len(sortedStatuses) - 1; i >= 0; i-- {
			txt += fmt.Sprintf("Code %d:\n", sortedStatuses[i])

			for _, secCount := range sortedSections[sortedStatuses[i]] {
				txt += fmt.Sprintf("    %s: %d\n", secCount.K, int(secCount.V))
			}

			txt += "\n"
		}
	}

	o.main.RoutesPerStatus(txt, status)
}

type alertLog struct {
	Date  time.Time
	Name  alert.NameT
	State alert.State
}

// TODO: rm global variable
type alertLogSlice []alertLog

func (o alertLogSlice) Len() int           { return len(o) }
func (o alertLogSlice) Less(i, j int) bool { return o[i].Date.Before(o[j].Date) }
func (o alertLogSlice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

var gAlertLogs alertLogSlice
var gActivityMonitor ActivityMonitor = *NewActivityMonitor()

func (o *View) Alerts(alerts <-chan backend.AlertStateTransition) {
	select {
	case a := <-alerts:
		switch a.Alert.State() {
		case alert.Inactive:
			gActivityMonitor.Monitor(a.Alert)

			// logs
			if a.Prev == alert.Active {
				gAlertLogs = append(
					gAlertLogs,
					alertLog{
						Date:  time.Unix(a.Time, 0),
						Name:  a.Alert.Name(),
						State: alert.Inactive,
					})
			}
		case alert.Pending:
			gActivityMonitor.Monitor(a.Alert)
		case alert.Active:
			gActivityMonitor.Monitor(a.Alert)

			// logs
			gAlertLogs = append(
				gAlertLogs,
				alertLog{
					Date:  time.Unix(a.Time, 0),
					Name:  a.Alert.Name(),
					State: alert.Active,
				})
		default:
			// critical logical error because enums do not exist
			err := fmt.Errorf("unsupported alert state %v", a)
			panic(err)
		}
	default:
		// Nothing to read let's try some other time
		return
	}

	logs := "Log:\n"
	for _, l := range gAlertLogs {
		logs += "  " + alertLogActivityTxt(l)
	}

	o.main.Alerts(gActivityMonitor.ActivityTxt(), logs)
}

func alertLogActivityTxt(log alertLog) string {
	switch log.State {
	case alert.Active:
		return fmt.Sprintf("%s\n    State: Active\n    Time: %s\n", log.Name, log.Date)
	case alert.Inactive:
		return fmt.Sprintf("%s\n    State: Inactive\n    Time: %s\n", log.Name, log.Date)
	default:
		// critical logical bug if that happens
		panic("pending alert or unsupported type has been placed in alert log slice")
	}
}

func NewView() (*View, error) {
	main, err := NewMainWindow()
	if err != nil {
		return nil, err
	}

	return &View{
		main: main,
	}, err
}

func (o *View) Layout() container.Option {
	o.layout = o.main.Layout()

	return o.layout
}
