package backend

import (
	"github.com/julnicolas/httpmon/pkg/parser"
	"github.com/julnicolas/httpmon/pkg/reader"
	"github.com/julnicolas/httpmon/pkg/trace"
)

// Ingestor is an object able to ingest Traces from
// various configurable raw formats
type Ingestor struct {
	reader reader.Reader // reads the input stream
	parser parser.Parser // parses incomming data
	traces chan trace.Trace
	source string // Ingestion source
	// implement filters
}

// Creates a new ingestor
func NewIngestor(source string, r reader.Reader, p parser.Parser, bufferLen uint) *Ingestor {
	return &Ingestor{
		reader: r,
		parser: p,
		traces: make(chan trace.Trace, bufferLen),
		source: source,
	}
}

func (o *Ingestor) Init() error {
	if err := o.reader.Open(o.source); err != nil {
		return err
	}

	return nil
}

func (o *Ingestor) Ingest() error {
	raw, err := o.reader.Read()
	if err != nil {
		return err
	}

	trace, err := o.parser.Parse(raw)
	if err == parser.ErrHeaderData {
		return nil
	}
	if err != nil {
		return err
	}

	o.traces <- trace
	return err
}

func (o *Ingestor) Poll() trace.Trace {
	return <-o.traces
}

func (o *Ingestor) Close() error {
	return o.reader.Close()
}
