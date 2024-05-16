package ui

import "github.com/mum4k/termdash/container"

type RequestsPerSecond struct {
	line *Line
	txt  *ListLayout
}

func NewRequestsPerSecond() (*RequestsPerSecond, error) {
	line, err := NewLine("Requests/s")
	if err != nil {
		return nil, err
	}

	txt, err := NewListLayout("")
	if err != nil {
		return nil, err
	}

	return &RequestsPerSecond{
		line: line,
		txt:  txt,
	}, err
}

// Values provides the values to plot as a line chart
func (o *RequestsPerSecond) Values(values []float64) {
	o.line.Values(values)
}

func (o *RequestsPerSecond) Text(txt string) {
	o.txt.Text(txt)
}

func (o *RequestsPerSecond) Layout() container.Option {
	return container.Option(
		container.SplitVertical(
			container.Left(
				o.line.Layout(),
			),
			container.Right(
				o.txt.Layout(),
			),
			container.SplitPercent(80),
		),
	)
}
