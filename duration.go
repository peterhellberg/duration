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
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	// HoursPerDay is the number of hours per day according to Google
	HoursPerDay = 24.0

	// HoursPerWeek is the number of hours per week according to Google
	HoursPerWeek = 168.0

	// HoursPerMonth is the number of hours per month according to Google
	HoursPerMonth = 730.4841667

	// HoursPerYear is the number of hours per year according to Google
	HoursPerYear = 8765.81

	// ErrUnsupportedFormat is returned when parsing fails
	ErrUnsupportedFormat = errors.New("unsupported string format")

	pattern = regexp.MustCompile(`\AP((?P<years>[\d\.]+)Y)?((?P<months>[\d\.]+)M)?((?P<weeks>[\d\.]+)W)?((?P<days>[\d\.]+)D)?(T((?P<hours>[\d\.]+)H)?((?P<minutes>[\d\.]+)M)?((?P<seconds>[\d\.]+?)S)?)?\z`)
)

// Parse a RFC3339 duration string into time.Duration
func Parse(s string) (time.Duration, error) {
	d := time.Duration(0)

	var match []string

	if pattern.MatchString(s) {
		match = pattern.FindStringSubmatch(s)
	} else {
		return d, ErrUnsupportedFormat
	}

	for i, name := range pattern.SubexpNames() {
		value := match[i]
		if i == 0 || name == "" || value == "" {
			continue
		}

		if f, err := strconv.ParseFloat(value, 64); err == nil {
			switch name {
			case "years":
				if years, err := time.ParseDuration(fmt.Sprintf("%fh", f*HoursPerYear)); err == nil {
					d += years
				}
			case "months":
				if months, err := time.ParseDuration(fmt.Sprintf("%fh", f*HoursPerMonth)); err == nil {
					d += months
				}
			case "weeks":
				if weeks, err := time.ParseDuration(fmt.Sprintf("%fh", f*HoursPerWeek)); err == nil {
					d += weeks
				}
			case "days":
				if days, err := time.ParseDuration(fmt.Sprintf("%fh", f*HoursPerDay)); err == nil {
					d += days
				}
			case "hours":
				if hours, err := time.ParseDuration(fmt.Sprintf("%fh", f)); err == nil {
					d += hours
				}
			case "minutes":
				if minutes, err := time.ParseDuration(fmt.Sprintf("%fm", f)); err == nil {
					d += minutes
				}
			case "seconds":
				if seconds, err := time.ParseDuration(fmt.Sprintf("%fs", f)); err == nil {
					d += seconds
				}
			}
		}
	}

	return d, nil
}
