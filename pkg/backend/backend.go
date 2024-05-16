package backend

import (
	"fmt"
	"strings"

	"github.com/julnicolas/httpmon/pkg/alert"
	"github.com/julnicolas/httpmon/pkg/config"
	"github.com/julnicolas/httpmon/pkg/metrics"
	"github.com/julnicolas/httpmon/pkg/parser"
	"github.com/julnicolas/httpmon/pkg/reader"
)

type Backend struct {
	ingestor  *Ingestor
	collector *MetricsCollector
	alertor   *AlertManager
	ingestErr error // not nil if an error occurs in ingestor.Ingest()
	pollErr   error // not nil if an error occurs in ingestor.Poll()
}

// Creates a new backend object
// after configuration has been properly implemented
// this would select appropriate ingestion parameters
// func NewBackend(file string, readBufferLen uint, alertor *AlertManager) *Backend {
func NewBackend(conf config.Config) *Backend {
	// TODO: could be moved in MetricsCollector based on a config object?
	probers := make([]metrics.Prober, 0, 3)
	probers = append(probers, metrics.NewRequestsPerHost())
	probers = append(probers, metrics.NewRoutePerStatus())
	probers = append(probers, metrics.NewRequestsPerSecond(conf.Period)) // Atta

	var r reader.Reader
	if strings.ToLower(conf.File) == "stdin" {
		r = reader.NewStdin(conf.ReadBufferSize)
	} else {
		r = reader.NewTailer(conf.ReadBufferSize)
	}

	return &Backend{
		ingestor: NewIngestor(
			conf.File,
			r,
			parser.NewCSV(),
			conf.ReadBufferSize,
		),
		collector: NewMetricsCollector(probers),
		alertor: NewAlertManager(conf.Alert.Period,
			alert.RequestsPerSecondInput{
				Period:    conf.Alert.RequestsPerSecond.Period,
				Threshold: conf.Alert.RequestsPerSecond.Threshold,
			}),
	}
}

func (o *Backend) Init() error {
	return o.ingestor.Init()
}

// runErr concatenates all errors that occured in run in a single message
type runErr struct {
	ingestErr error
	pollErr   error
}

func (o runErr) Error() string {
	err := ""
	if o.ingestErr != nil {
		err += fmt.Sprintf("ingest: %s; ", o.ingestErr)
	}
	if o.pollErr != nil {
		err += fmt.Sprintf("poll: %s;", o.pollErr)
	}
	return err
}

// RunErr returns non nil if an error occured in Run
// Useful if backend is run in another goroutine
func (o *Backend) RunErr() error {
	if o.ingestErr != nil || o.pollErr != nil {
		return runErr{ingestErr: o.ingestErr, pollErr: o.pollErr}
	}

	return nil
}

func (o *Backend) Run() error {
	go func() {
		for {
			if err := o.ingestor.Ingest(); err != nil {
				o.ingestErr = err
				return
			}
		}
	}()

	for {
		t := o.ingestor.Poll()
		if err := o.collector.Collect(t); err != nil {
			o.pollErr = err
			return err
		}

		m, err := o.Metric(metrics.ReqsPerS)
		if err != nil {
			panic(err)
		}

		// Evaluate all alerts from updated metric values
		o.alertor.Eval(m)
	}
}

// Metric returns a copy of a metric so that it can be threadsafe
func (o *Backend) Metric(name string) (metrics.Metric, error) {
	// TODO: Cache metric per 1s as UI and alerts are directed to humans
	// and human reaction time is 1s on avaerage.
	return o.collector.DeepCopy(name)
}

// Alerts only exposes alerts which state's have changed
// Every alert is at least sent once when it is initialises as inactive
func (o *Backend) Alerts() <-chan AlertStateTransition {
	return o.alertor.Alerts()
}

func (o *Backend) Close() error {
	return o.ingestor.Close()
}
