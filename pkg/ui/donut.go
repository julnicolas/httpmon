package ui

import (
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/widgets/donut"
)

// DonutLayout is a circular progress gauge chart
type DonutLayout struct {
	Donut *donut.Donut
}

func (o *DonutLayout) Layout() container.Option {
	return container.PlaceWidget(o.Donut)
}
