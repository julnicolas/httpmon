package parser

import (
	"errors"

	"github.com/julnicolas/httpmon/pkg/trace"
)

// Parser is an interface used to parse a string into a usable
// Trace structure
type Parser interface {
	// Parse takes an input string, converts it to a Trace.
	//
	// Some parsers read a header text to initialise.
	// If the format you parse uses a header or if you don't know,
	// check if err != ErrHeaderData.
	Parse(string) (trace.Trace, error)
}

// HeaderData is an error returned when header content is being processed
// and still considered valid.
var ErrHeaderData error = errors.New("parser.Parser: header data")
