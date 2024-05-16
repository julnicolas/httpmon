package ui

import (
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/widgets/text"
)

type ListLayout struct {
	txt *text.Text
}

// NewListLayout creates a new layout struct which displays a list
// (aka a rolling text). Content is reset on every call to Text.
func NewListLayout(txt string) (*ListLayout, error) {
	w, err := newTextWidget(txt)
	if err != nil {
		return nil, err
	}

	return &ListLayout{
		txt: w,
	}, err
}

// Text sets the text in the list, resetting content
func (o *ListLayout) Text(txt string) {
	if o.txt == nil {
		return
	}

	o.txt.Write(txt, text.WriteReplace())
}

// Append appends new content to the text already displayed
func (o *ListLayout) Append(txt string) {
	if o.txt == nil {
		return
	}

	o.txt.Write(txt)
}

// Layout exposes the page's layout
func (o *ListLayout) Layout() container.Option {
	if o.txt == nil {
		var zero container.Option
		return zero
	}

	return container.PlaceWidget(o.txt)
}

// newTextWidget creates a new text widget
func newTextWidget(txt string) (*text.Text, error) {
	t, err := text.New(text.RollContent())
	if err != nil {
		return nil, err
	}
	t.Write(txt)

	return t, nil
}
