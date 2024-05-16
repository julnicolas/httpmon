package ui

import "github.com/mum4k/termdash/container"

type RequestsPerHost struct {
	repartition *Repartition
}

func NewRequestsPerHost(txt string, topk []float64) (*RequestsPerHost, error) {
	repartition, err := NewRepartition(txt, topk)
	if err != nil {
		return nil, err
	}

	return &RequestsPerHost{
		repartition: repartition,
	}, err
}

func (o *RequestsPerHost) Text(txt string) {
	o.repartition.Text(txt)
}

func (o *RequestsPerHost) Topk(percent []float64) error {
	return o.repartition.Percent(percent)
}

func (o *RequestsPerHost) Layout() container.Option {
	return o.repartition.Layout()
}
