package ui

import (
	"github.com/mum4k/termdash/container"
)

type Alerts struct {
	txt  *ListLayout // displays active alerts
	logs *ListLayout // a log of alert past alert states
}

func NewAlerts(t string) (*Alerts, error) {
	txt, err := NewListLayout(t)
	if err != nil {
		return nil, err
	}

	log, err := NewListLayout(t)
	if err != nil {
		return nil, err
	}

	return &Alerts{
		txt:  txt,
		logs: log,
	}, err
}

func (o *Alerts) Alerts(alerts string, logs string) {

	o.txt.Text(alerts)
	o.logs.Text(logs)
}

func (o *Alerts) Layout() container.Option {
	return container.Option(
		container.SplitVertical(
			container.Left(
				o.txt.Layout(),
			),
			container.Right(
				o.logs.Layout(),
			),
			container.SplitPercent(60),
		),
	)
}
