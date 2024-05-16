package reader

import (
	"bufio"
	"os"
	"time"
)

// Reads Stdin line-by-line
type Stdin struct {
	scanner *bufio.Scanner
	once    bool        // once is used to start bufferize() on Read only once
	lines   chan string // buffered channel of lines
}

func NewStdin(nbLines uint) *Stdin {
	return &Stdin{
		lines: make(chan string, nbLines),
	}
}

// Open opens stdin for reading, filename is ignored here
// Though it is present in the reader interface
func (o *Stdin) Open(filename string) error {
	o.scanner = bufio.NewScanner(os.Stdin)
	return nil
}

func (o *Stdin) bufferize() {
	for {
		if err := o.buffering(); err != nil {
			// FIXME: add proper error management here
			// But I've spent enough time on the project :)
			panic(err)
		}
	}
}

// buffering is the action method of bufferize
// it makes sure new strings are allocated on every loop run
// so that memory is not corrupted
func (o *Stdin) buffering() error {
	line := ""
	if o.scanner.Scan() {
		line = o.scanner.Text()
	}

	if err := o.scanner.Err(); err != nil {
		return err
	}

	if line == "" {
		// Nothing's going on on stdin, check out later
		time.Sleep(time.Second)
		return nil
	}

	o.lines <- line
	return nil
}

// Read reads stdin, blocks when nothing is comming through
// Note: an empty string is considered as no entry by the Reader
// (default behaviour of bufio.Scan when reading stdin)
func (o *Stdin) Read() (line string, err error) {
	if !o.once {
		go o.bufferize()
		o.once = true
	}

	return <-o.lines, nil
}

func (o *Stdin) Close() error {
	return nil
}
