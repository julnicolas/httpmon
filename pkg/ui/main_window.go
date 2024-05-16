package ui

import (
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/button"
)

// MainWindow is the main control window
type MainWindow struct {
	activeTab        uint
	reqsPerHost      *RequestsPerHost
	reqsPerSec       *RequestsPerSecond
	routesPerStatus  *RoutesPerStatus
	alerts           *Alerts
	reqsPerHostB     *button.Button
	reqsPerSecB      *button.Button
	routesPerStatusB *button.Button
	alertsB          *button.Button
}

func (o *MainWindow) ReqsPerHost(reqList string, topk []float64) {
	o.reqsPerHost.Text(reqList)
	o.reqsPerHost.Topk(topk)
}

func (o *MainWindow) ReqsPerSec(perSectionTxt string, values []float64) {
	o.reqsPerSec.Text(perSectionTxt)
	o.reqsPerSec.Values(values)
}

func (o *MainWindow) RoutesPerStatus(routeList string, repartition StatusPercent) {
	o.routesPerStatus.Text(routeList)
	o.routesPerStatus.StatusPercent(repartition)
}

func (o *MainWindow) Alerts(txt, logs string) {
	o.alerts.Alerts(txt, logs)
}

func NewMainWindow() (*MainWindow, error) {
	rPerHost, err := NewRequestsPerHost(
		"no incomming requests",
		// Top 5 hosts
		[]float64{0.0, 0.0, 0.0, 0.0, 0.0},
	)
	if err != nil {
		return nil, err
	}

	rPerS, err := NewRequestsPerSecond()
	if err != nil {
		return nil, err
	}

	status, err := NewRoutesPerStatus(
		"no incomming requests",
		75.0,
		10,
		15,
		2,
	)
	if err != nil {
		return nil, err
	}

	alerts, err := NewAlerts("")
	if err != nil {
		return nil, err
	}
	b := &MainWindow{
		reqsPerHost:     rPerHost,
		reqsPerSec:      rPerS,
		routesPerStatus: status,
		alerts:          alerts,
	}

	if err := b.newButtons(); err != nil {
		return nil, err
	}

	return b, nil
}

type childPage struct {
	Name string
	P    Page
}

func (o *MainWindow) Layout() container.Option {
	page := o.activePage()
	return container.SplitHorizontal(
		container.Top(
			container.Border(linestyle.Light),
			container.BorderTitle(page.Name),
			page.P.Layout(),
		),
		container.Bottom(
			o.tabLayout(),
		),
		container.SplitPercent(97),
	)
}

func (o *MainWindow) activePage() childPage {
	return o.getChildPage(o.activeTab)
}

func (o *MainWindow) getChildPage(id uint) childPage {
	switch id {
	case 0:
		return childPage{
			Name: "Requests/Host",
			P:    o.reqsPerHost,
		}
	case 1:
		return childPage{
			Name: "Requests/s",
			P:    o.reqsPerSec,
		}
	case 2:
		return childPage{
			Name: "Status",
			P:    o.routesPerStatus,
		}
	case 3:
		return childPage{
			Name: "Alerts",
			P:    o.alerts,
		}
	default:
		return childPage{
			Name: "Requests/Host",
			P:    o.reqsPerHost,
		}
	}
}

func (o *MainWindow) tabLayout() container.Option {
	return buttonLayout(
		container.PlaceWidget(o.reqsPerHostB),
		buttonLayout(
			container.PlaceWidget(o.reqsPerSecB),
			buttonLayout(
				container.PlaceWidget(o.routesPerStatusB),
				container.PlaceWidget(o.alertsB),
				10,
			),
			14,
		),
		17,
	)
}

func buttonLayout(left, right container.Option, cellLen int) container.Option {
	return container.SplitVertical(
		container.Left(
			left,
			container.AlignHorizontal(align.HorizontalLeft),
		),
		container.Right(
			right,
			container.AlignHorizontal(align.HorizontalLeft),
		),
		container.SplitFixed(cellLen),
	)
}

func (o *MainWindow) newButtons() error {
	opts := []button.Option{
		button.FillColor(cell.ColorNumber(220)),
		button.Height(1),
		button.DisableShadow(),
	}
	r1, err := button.New(o.getChildPage(0).Name, func() error {
		o.activeTab = 0
		return nil
	}, opts...)
	if err != nil {
		return err
	}
	r2, err := button.New(o.getChildPage(1).Name, func() error {
		o.activeTab = 1
		return nil
	}, opts...)
	if err != nil {
		return err
	}
	r3, err := button.New(o.getChildPage(2).Name, func() error {
		o.activeTab = 2
		return nil
	}, opts...)
	if err != nil {
		return err
	}
	r4, err := button.New(o.getChildPage(3).Name, func() error {
		o.activeTab = 3
		return nil
	}, opts...)
	if err != nil {
		return err
	}

	o.reqsPerHostB = r1
	o.reqsPerSecB = r2
	o.routesPerStatusB = r3
	o.alertsB = r4
	return err
}
