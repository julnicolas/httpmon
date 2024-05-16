package ui

import "github.com/mum4k/termdash/container"

// Page is an interface representing a page displayed by the UI
type Page interface {
	// Layout returns a container option describing page's layout
	//
	// Warning: make sure related widgets are managed by the page (best
	// solution) or available at all times to avoid segfaults.
	Layout() container.Option
}
