/*

Package duration parses RFC3339 duration strings into time.Duration

Installation

Just go get the package:

    go get -u github.com/peterhellberg/duration

Usage

A small usage example

		package main

		import (
			"fmt"

			"github.com/peterhellberg/duration"
		)

		func main() {
			if d, err := duration.Parse("P1DT30H4S"); err == nil {
				fmt.Println(d) // Output: 54h0m4s
			}
		}

*/
package duration

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

var (
	// ErrUnsupportedFormat is returned when parsing fails
	ErrUnsupportedFormat = errors.New("unsupported string format")

	pattern = regexp.MustCompile(`\AP((?P<days>\d+)D)?(T((?P<hours>\d+)H)?((?P<minutes>\d+)M)?((?P<seconds>[\d\.]+?)S)?)?\z`)
)

// Parse a RFC3339 duration string into time.Duration
//
// Not supported for now: Years, Months and Weeks
func Parse(s string) (time.Duration, error) {
	d := time.Duration(0)

	var match []string

	// Extract matches
	if pattern.MatchString(s) {
		match = pattern.FindStringSubmatch(s)
	} else {
		return d, ErrUnsupportedFormat
	}

	// Assign match values to Duration struct fields
	for i, name := range pattern.SubexpNames() {
		value := match[i]
		if i == 0 || name == "" || value == "" {
			continue
		}

		// Ignore error since we accept the zero value
		f, _ := strconv.ParseFloat(value, 64)

		switch name {
		case "days":
			d += time.Duration(f) * time.Hour * 24
		case "hours":
			d += time.Duration(f) * time.Hour
		case "minutes":
			d += time.Duration(f) * time.Minute
		case "seconds":
			d += time.Duration(f) * time.Second
		}
	}

	return d, nil
}
