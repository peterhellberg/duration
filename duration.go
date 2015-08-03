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
	"strings"
	"time"
)

const (
	// HoursPerDay is the number of hours per day according to Google
	HoursPerDay = 24.0

	// HoursPerWeek is the number of hours per week according to Google
	HoursPerWeek = 168.0

	// HoursPerMonth is the number of hours per month according to Google
	HoursPerMonth = 730.4841667

	// HoursPerYear is the number of hours per year according to Google
	HoursPerYear = 8765.81
)

var (
	// ErrInvalidString is returned when passed an invalid string
	ErrInvalidString = errors.New("invalid duration string")

	// ErrUnsupportedFormat is returned when parsing fails
	ErrUnsupportedFormat = errors.New("unsupported duration string format")

	pattern = regexp.MustCompile(`\A(-)?P((?P<years>[\d\.]+)Y)?((?P<months>[\d\.]+)M)?((?P<weeks>[\d\.]+)W)?((?P<days>[\d\.]+)D)?(T((?P<hours>[\d\.]+)H)?((?P<minutes>[\d\.]+)M)?((?P<seconds>[\d\.]+?)S)?)?\z`)

	invalidStrings = []string{"", "P", "PT"}
)

// Parse a RFC3339 duration string into time.Duration
func Parse(s string) (time.Duration, error) {
	if contains(invalidStrings, s) || strings.HasSuffix(s, "T") {
		return 0, ErrInvalidString
	}

	var (
		match  []string
		prefix string
	)

	if pattern.MatchString(s) {
		match = pattern.FindStringSubmatch(s)
	} else {
		return 0, ErrUnsupportedFormat
	}

	if strings.HasPrefix(s, "-") {
		prefix = "-"
	}

	return durationFromMatchAndPrefix(match, prefix)
}

func durationFromMatchAndPrefix(match []string, prefix string) (time.Duration, error) {
	d := time.Duration(0)

	duration := func(format string, f float64) (time.Duration, error) {
		return time.ParseDuration(fmt.Sprintf(prefix+format, f))
	}

	for i, name := range pattern.SubexpNames() {
		value := match[i]
		if i == 0 || name == "" || value == "" {
			continue
		}

		if f, err := strconv.ParseFloat(value, 64); err == nil {
			switch name {
			case "years":
				if years, err := duration("%fh", f*HoursPerYear); err == nil {
					d += years
				}
			case "months":
				if months, err := duration("%fh", f*HoursPerMonth); err == nil {
					d += months
				}
			case "weeks":
				if weeks, err := duration("%fh", f*HoursPerWeek); err == nil {
					d += weeks
				}
			case "days":
				if days, err := duration("%fh", f*HoursPerDay); err == nil {
					d += days
				}
			case "hours":
				if hours, err := duration("%fh", f); err == nil {
					d += hours
				}
			case "minutes":
				if minutes, err := duration("%fm", f); err == nil {
					d += minutes
				}
			case "seconds":
				if seconds, err := duration("%fs", f); err == nil {
					d += seconds
				}
			}
		}
	}

	return d, nil
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
