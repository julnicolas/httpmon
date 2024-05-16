package ui

import (
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/widgets/linechart"
)

type Line struct {
	name string
	line *linechart.LineChart
}

// NewLine creates a new line chart
// name corresponds to the graph's name
func NewLine(name string) (*Line, error) {
	line, err := newLineWidget(name)
	if err != nil {
		return nil, err
	}

	return &Line{
		name: name,
		line: line,
	}, err
}

// Values provides the values to plot as a line chart
func (o *Line) Values(values []float64) {
	o.line.Series(o.name, values)
}

// Layout returns this pages layout so that it can be rendered
func (o *Line) Layout() container.Option {
	return container.PlaceWidget(o.line)
}

// newLineWidget creates a new line chart widget
func newLineWidget(name string) (*linechart.LineChart, error) {
	line, err := linechart.New(
		linechart.AxesCellOpts(cell.FgColor(cell.ColorRed)),
		linechart.YLabelCellOpts(cell.FgColor(cell.ColorGreen)),
		linechart.XLabelCellOpts(cell.FgColor(cell.ColorGreen)),
	)
	if err != nil {
		return nil, err
	}

	line.Series(
		name,
		[]float64{0},
		linechart.SeriesCellOpts(cell.FgColor(cell.ColorNumber(33))),
	)

	return line, err
}
