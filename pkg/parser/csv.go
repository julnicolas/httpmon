package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/julnicolas/httpmon/pkg/trace"
)

// CSV parses csv-formatted strings representing
// http calls. It returns a well formed Trace.
type CSV struct {
	delimiter string
	// true if header has been parsed and validated.
	// It means content is ready to be parsed.
	validHeader bool
}

// NewCSV creates a new CSV parser
func NewCSV() *CSV {
	return &CSV{delimiter: ","}
}

// value returns a csv field value with all formatting removed
func (o *CSV) value(field string) string {
	return strings.Trim(field, "\"")
}

// Header parses the header, returns an error if it doesn't follow the
// expected format
//
// If the header is valid, it returns the parser.HeaderData error value
func (o *CSV) validateHeader(header string) error {
	// Number of expected CSV fields
	const nbFields uint = 7

	fields := strings.Split(strings.ToLower(header), o.delimiter)
	if len(fields) != 7 {
		return fmt.Errorf("csv parse error - expected %d fields, add %d", nbFields, len(fields))
	}

	fmtErr := "missing required field %s"
	if o.value(fields[0]) != "remotehost" {
		return fmt.Errorf(fmtErr, fields[0])
	}
	if o.value(fields[1]) != "rfc931" {
		return fmt.Errorf(fmtErr, fields[1])
	}
	if o.value(fields[2]) != "authuser" {
		return fmt.Errorf(fmtErr, fields[2])
	}
	if o.value(fields[3]) != "date" {
		return fmt.Errorf(fmtErr, fields[3])
	}
	if o.value(fields[4]) != "request" {
		return fmt.Errorf(fmtErr, fields[4])
	}
	if o.value(fields[5]) != "status" {
		return fmt.Errorf(fmtErr, fields[5])
	}
	if o.value(fields[6]) != "bytes" {
		return fmt.Errorf(fmtErr, fields[6])
	}

	// All checks have passed so the stream can be parsed
	o.validHeader = true

	return ErrHeaderData
}

// Parse creates a Trace representing an HTTP call
//
// If raw is a valid header line, the parser.HeaderData
// error value is returned.
func (o *CSV) Parse(raw string) (trace.Trace, error) {
	// Intent header validation until it is validated
	// then proceed with data parsing on next calls
	if !o.validHeader {
		return trace.Trace{}, o.validateHeader(raw)
	}

	const nbFields uint = 7

	fields := strings.Split(raw, o.delimiter)
	if len(fields) != 7 {
		return trace.Trace{}, fmt.Errorf("csv parse error - expected %d fields, add %d", nbFields, len(fields))
	}

	t := trace.Trace{}
	if err := parseDate(&t, o.value(fields[3])); err != nil {
		return trace.Trace{}, err
	}
	if err := parseRemoteHost(&t, o.value(fields[0])); err != nil {
		return trace.Trace{}, err
	}
	if err := parseAuthUser(&t, o.value(fields[2])); err != nil {
		return trace.Trace{}, err
	}
	if err := parseRFC931(&t, o.value(fields[1])); err != nil {
		return trace.Trace{}, err
	}
	if err := parseRequest(&t, o.value(fields[4])); err != nil {
		return trace.Trace{}, err
	}
	if err := parseStatus(&t, o.value(fields[5])); err != nil {
		return trace.Trace{}, err
	}
	if err := parseBytes(&t, o.value(fields[6])); err != nil {
		return trace.Trace{}, err
	}

	return t, nil
}

func parseDate(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}

	unix, err := strconv.Atoi(field)
	if err != nil {
		return err
	}
	t.Date = time.Unix(int64(unix), 0)
	return nil
}

func parseRemoteHost(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}
	t.RemoteHost = field
	return nil
}

func parseAuthUser(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}
	t.AuthUser = field
	return nil
}

func parseRFC931(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}
	if field == "-" {
		return nil
	}
	t.RFC931 = field
	return nil
}

func parseRequest(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}

	fields := strings.Fields(field)
	if len(fields) != 3 {
		// TODO: We could make it smarter here
		// Missing a few fields can still be relevant
		return fmt.Errorf("missing request log data")
	}

	if err := parseMethod(t, fields[0]); err != nil {
		return err
	}
	if err := parseSection(t, fields[1]); err != nil {
		return err
	}
	if err := parseVersion(t, fields[2]); err != nil {
		return err
	}

	return nil
}

func parseMethod(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}

	method := strings.ToUpper(field)
	switch method {
	case "GET":
	case "POST":
	case "PUT":
	case "PATCH":
	case "DELETE":
	case "HEAD":
	case "TRACE":
	case "CONNECT":
	case "OPTIONS":
	default:
		return fmt.Errorf("http method verb is invalid")
	}

	t.Method = method
	return nil
}

// parseSection parses the section (first part of the path part of an URL)
// of a URL.
func parseSection(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}
	if field == "" {
		return fmt.Errorf("empty url path")
	}

	t.Section = "/" + strings.Split(field, "/")[1]
	return nil
}

func parseVersion(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}

	version := strings.Split(field, "/")
	if len(version) != 2 {
		return fmt.Errorf("version is ill-formatted")
	}
	t.Version = version[1]
	return nil
}

func parseStatus(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}

	status, err := strconv.Atoi(field)
	if err != nil {
		return err
	}
	if status < 100 || status > 599 {
		return fmt.Errorf("status code is out of bond")
	}

	t.Status = uint(status)
	return nil
}

func parseBytes(t *trace.Trace, field string) error {
	if t == nil {
		return fmt.Errorf("nil receiver")
	}
	bytes, err := strconv.Atoi(field)
	if err != nil {
		return err
	}

	t.Bytes = uint(bytes)
	return nil
}
