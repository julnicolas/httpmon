package reader

import (
	"fmt"

	"github.com/nxadm/tail"
	"gopkg.in/tomb.v1"
)

// stopSignal used to stop background file read
type stopSignal struct{}

// Tail is a file ingestor which tails a file
// as reading method.
type Tail struct {
	tailer *tail.Tail  // tail reader
	buffer chan string // buffer containing read lines
	stop   chan stopSignal
}

// NewTailer creates a new tailer with a buffer able to contain
// nbLines lines before blocking reads.
func NewTailer(nbLines uint) *Tail {
	return &Tail{
		buffer: make(chan string, nbLines),
		stop:   make(chan stopSignal),
	}
}

// Open opens a file in tail mode
func (o *Tail) Open(path string) error {
	t, err := tail.TailFile(path, tail.Config{
		Follow: true,
		ReOpen: true,
		//Poll: true, // do not use inotify
	})
	if err != nil {
		return err
	}

	o.tailer = t

	go o.bufferize()
	return err
}

// bufferize reads at most len(buffer) lines then stores them
// That enables to control the read rate: stopping reads if our
// allocated memory is full
//
// this function is meant to be called asynchronously.
// It waits if buffer is full
// cycles waiting for logs indefinitely unless stopSignal{} is sent
func (o *Tail) bufferize() {
	for {
		select {
		case line := <-o.tailer.Lines:
			if line != nil {
				o.buffer <- line.Text
			}
			if err := o.tailer.Err(); err != tomb.ErrStillAlive {
				panic(err)
			}
		case <-o.stop:
			o.tailer.Stop()
			return
		}
	}
}

// Poll returns available lines in read-order.
// Returned value can be safely converted to string.
// The function blocks if the buffer is empty.
func (o *Tail) Read() (string, error) {
	return <-o.buffer, nil
}

// Close removes inotify watches added by the tail package.
// The linux kernel may not clean it at process exit.
func (o *Tail) Close() error {
	if o.tailer == nil {
		return fmt.Errorf("tailer is nil")
	}
	o.stop <- stopSignal{}
	o.tailer.Cleanup()

	return nil
}
