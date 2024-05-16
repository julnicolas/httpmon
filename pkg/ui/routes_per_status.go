package ui

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/widgets/gauge"
)

type RoutesPerStatus struct {
	repartition *Repartition
}

func NewRoutesPerStatus(txt string, _2xx, _3xx, _4xx, _5xx float64) (*RoutesPerStatus, error) {
	status := []float64{
		_2xx, _3xx, _4xx, _5xx,
	}
	opts := [][]gauge.Option{
		{gauge.TextLabel("5xx"), gauge.Color(cell.ColorRed)},
		{gauge.TextLabel("4xx"), gauge.Color(cell.ColorRed)},
		{gauge.TextLabel("3xx")},
		{gauge.TextLabel("2xx")},
	}
	repartition, err := NewRepartition(txt, status, opts...)
	if err != nil {
		return nil, err
	}

	return &RoutesPerStatus{
		repartition: repartition,
	}, err
}

// Text defines request listing on left pannel
func (o *RoutesPerStatus) Text(txt string) {
	o.repartition.Text(txt)
}

// StatusPercent describes repartion of http status codes in percent
type StatusPercent struct {
	Http2xx float64
	Http3xx float64
	Http4xx float64
	Http5xx float64
}

// StatusPercent draws status repartition in percent
func (o *RoutesPerStatus) StatusPercent(s StatusPercent) error {
	status := []float64{s.Http5xx, s.Http4xx, s.Http3xx, s.Http2xx}
	return o.repartition.Percent(status)
}

// Layout returns the page's layout
func (o *RoutesPerStatus) Layout() container.Option {
	return o.repartition.Layout()
}
