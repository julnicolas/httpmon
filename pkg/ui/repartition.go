package ui

import (
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/widgets/gauge"
)

type Repartition struct {
	text *ListLayout
	bars *HBars
}

func NewRepartition(txt string, topk []float64, opts ...[]gauge.Option) (*Repartition, error) {
	txtP, err := NewListLayout(txt)
	if err != nil {
		return nil, err
	}
	barsP, err := NewHBars(topk, opts...)
	if err != nil {
		return nil, err
	}

	return &Repartition{
		text: txtP,
		bars: barsP,
	}, err
}

func (o *Repartition) Text(txt string) {
	o.text.Text(txt)
}

// Percent draws k bar charts on top of each other.
// Values must be expressed in percent.
func (o *Repartition) Percent(values []float64) error {
	return o.bars.Percent(values)
}

func (o *Repartition) Layout() container.Option {
	if o.text == nil || o.bars == nil {
		var zero container.Option
		return zero
	}

	return container.Option(
		container.SplitVertical(
			container.Left(
				o.text.Layout(),
			),
			container.Right(
				o.bars.Layout(),
			),
			container.SplitPercent(80),
		),
	)
}
