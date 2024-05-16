package ui

import (
	"fmt"

	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/widgets/gauge"
)

type HBars struct {
	bars []*gauge.Gauge
}

func NewHBars(values []float64, opts ...[]gauge.Option) (*HBars, error) {
	bars, err := newHBarSlice(values, opts...)
	if err != nil {
		return nil, err
	}

	return &HBars{
		bars: bars,
	}, err
}

func (o *HBars) Percent(values []float64) error {
	if len(values) != len(o.bars) {
		return fmt.Errorf("cannot update bars - %d bars and %d values", len(o.bars), len(values))
	}

	for i := 0; i < len(o.bars); i++ {
		if err := o.bars[i].Percent(int(values[i])); err != nil {
			return err
		}
	}

	return nil
}

func (o *HBars) Layout() container.Option {
	return barSliceLayout(o.bars, 1)
}

func newHbar(v int, opts ...gauge.Option) (*gauge.Gauge, error) {
	g, err := gauge.New(opts...)
	if err != nil {
		return nil, err
	}

	if err := g.Percent(v); err != nil {
		return nil, err
	}

	return g, err
}

func newHBarSlice(values []float64, opts ...[]gauge.Option) ([]*gauge.Gauge, error) {
	if len(opts) != 0 && len(opts) != len(values) {
		return nil, fmt.Errorf("options must be provided for every gauges or none (empty lists accepted)")
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("NewHBars - empty list of hbars")
	}

	bars := make([]*gauge.Gauge, 0, len(values))
	var barOpts []gauge.Option
	for i, v := range values {
		if len(opts) > 0 {
			barOpts = opts[i]
		}
		bar, err := newHbar(int(v), barOpts...)
		if err != nil {
			return nil, err
		}

		bars = append(bars, bar)
	}

	return bars, nil
}

func barSliceLayout(bars []*gauge.Gauge, cellLen int) container.Option {
	if len(bars) == 0 {
		var zero container.Option
		return zero
	}

	if len(bars) == 1 {
		return barLayout(container.PlaceWidget(bars[0]), container.AlignHorizontal(align.Horizontal(align.VerticalTop)), cellLen)
	}

	return barLayout(container.PlaceWidget(bars[0]), barSliceLayout(bars[1:], cellLen), cellLen)
}

func barLayout(top, bottom container.Option, cellLen int) container.Option {
	return container.SplitHorizontal(
		container.Top(
			top,
			container.AlignVertical(align.VerticalTop),
		),
		container.Bottom(
			bottom,
			container.AlignVertical(align.VerticalTop),
			container.MarginTop(1),
		),
		container.SplitFixed(cellLen),
	)
}
