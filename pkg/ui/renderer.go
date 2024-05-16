package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
)

type Renderer struct {
	running   bool
	term      terminalapi.Terminal
	container *container.Container
	//renderer  *termdash.Controller
	ctx    context.Context
	cancel context.CancelFunc
	errRun error
	view   *View // should be the opposite -> renderer a dep
}

func NewRenderer() *Renderer {
	return &Renderer{}
}

// rootID is the ID of the root Renderer container
const rootID string = "root"

// View return's the renderer's view
// Todo: rm method by inversing dependency between view and renderer
func (o *Renderer) View() *View {
	return o.view
}

// Init initialises the Renderer
func (o *Renderer) Init() (err error) {
	o.term, err = tcell.New(tcell.ColorMode(terminalapi.ColorMode256))
	if err != nil {
		return err
	}

	o.container, err = container.New(o.term, container.ID(rootID))
	if err != nil {
		return err
	}

	o.view, err = NewView()
	if err != nil {
		return err
	}

	o.ctx, o.cancel = context.WithCancel(context.Background())
	go func() {
		o.errRun = termdash.Run(o.ctx, o.term, o.container, termdash.KeyboardSubscriber(o.quitter), termdash.RedrawInterval(16*time.Millisecond))
	}()

	o.running = true
	return err
}

func (o *Renderer) quitter(k *terminalapi.Keyboard) {
	if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
		o.running = false
		o.cancel()
	}
}

func (o *Renderer) Running() bool {
	// mutex?
	return o.errRun == nil && o.running
}

func (o *Renderer) Render() error {
	if o.Running() {
		if err := o.container.Update(rootID, o.view.Layout()); err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("renderer is not running")
}

func (o *Renderer) Close() {
	o.running = false
	if o.term != nil {
		o.term.Close()
	}
}
