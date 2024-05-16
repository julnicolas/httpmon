package app

import (
	"fmt"
	"time"

	"github.com/julnicolas/httpmon/pkg/backend"
	"github.com/julnicolas/httpmon/pkg/config"
	"github.com/julnicolas/httpmon/pkg/metrics"
	"github.com/julnicolas/httpmon/pkg/ui"
)

type App struct {
	backend  *backend.Backend
	frontend *ui.Renderer
}

func NewApp(c config.Config) *App {
	back := backend.NewBackend(c)
	frontend := ui.NewRenderer()

	return &App{
		backend:  back,
		frontend: frontend,
	}
}

func (o *App) Init() error {
	if err := o.backend.Init(); err != nil {
		return err
	}

	if err := o.frontend.Init(); err != nil {
		return err
	}

	return nil
}

func (o *App) Run() error {
	go o.backend.Run()

	for o.frontend.Running() && o.backend.RunErr() == nil {
		o.updateDashboards() // should be moved in view with view/renderer dependecy reversed
		if err := o.frontend.Render(); err != nil {
			return err
		}

		time.Sleep(250 * time.Millisecond)
	}
	return o.backend.RunErr()
}

// updateDashboards reads current metrics' state then
// feed them to their appropriate view
func (o *App) updateDashboards() error {
	c, err := o.counterMetric(metrics.ReqsPerHost)
	if err != nil {
		return err
	}
	o.frontend.View().ReqsPerHost(c)

	reqPers, err := o.counterVectorMetric(metrics.ReqsPerS)
	if err != nil {
		return err
	}
	o.frontend.View().ReqsPerSec(reqPers)
	o.frontend.View().Alerts(o.backend.Alerts())

	rc, err := o.routePerStatusCounterMetric()
	if err != nil {
		return err
	}
	o.frontend.View().RoutesPerStatus(rc)

	return err
}

func (o *App) counterMetric(name string) (metrics.Counter, error) {
	m, err := o.backend.Metric(name)
	if err != nil {
		return metrics.Counter{}, err
	}
	c, ok := m.(metrics.Counter)
	if !ok {
		return metrics.Counter{}, fmt.Errorf("interface cast error - expected Counter")
	}

	return c, err
}

func (o *App) counterVectorMetric(name string) (metrics.CounterVector, error) {
	m, err := o.backend.Metric(name)
	if err != nil {
		return metrics.CounterVector{}, err
	}
	c, ok := m.(metrics.CounterVector)
	if !ok {
		return metrics.CounterVector{}, fmt.Errorf("interface cast error - expected Counter")
	}

	return c, err
}

func (o *App) routePerStatusCounterMetric() (metrics.RoutePerStatusCounter, error) {
	m, err := o.backend.Metric(metrics.RoutesPerStatusN)
	if err != nil {
		return metrics.RoutePerStatusCounter{}, err
	}
	rc, ok := m.(metrics.RoutePerStatusCounter)
	if !ok {
		return metrics.RoutePerStatusCounter{}, fmt.Errorf("expected Counter for metric")
	}

	return rc, err
}

func (o *App) Close() {
	o.backend.Close()
	o.frontend.Close()
}
