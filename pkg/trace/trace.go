package trace

import (
	"time"
)

// Trace is a structure representing a http call
type Trace struct {
	// Date when the trace has been received or written
	// (in case this comes from a file)
	Date time.Time
	// RemoteHost is the hostname or IP
	RemoteHost string
	// AuthUser is a user name from an authenticated user
	AuthUser string
	RFC931   string
	// Method is the HTTP method used for the call
	Method string // would be great to use an enum here
	// A Section is a first part between the two '/' of an URL path
	// for /foo/bar that would be /foo (first '/' is included)
	Section string
	// Version is the HTTP version
	Version string
	// Status is the request's status code
	Status uint
	// Bytes corresponds to the number of bytes sent
	Bytes uint // bytes sent
}
