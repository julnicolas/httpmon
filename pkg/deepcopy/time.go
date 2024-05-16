package deepcopy

import "time"

// package deepcopy provides different utilities to
// copy complex datastructures

// Time deep copies a time.Time struct
func Time(t time.Time) time.Time {
	n, err := time.Parse(time.RFC1123, t.Format(time.RFC1123))
	if err != err {
		// this a format mismatch above
		panic(err)
	}

	return n
}
