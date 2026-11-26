/*
Package duration parses ISO 8601 duration strings into time.Duration

The format is the one described in Appendix A of RFC 3339, for example
"P1Y2M3DT4H5M6.7S". In addition to the strict grammar, the following are
accepted:

  - fractions in any component, using "." or "," as the decimal separator
  - weeks combined with other components, such as "P1WT2H"
  - a leading "-" for negative durations

Years, months, weeks and days are converted using the HoursPer* constants.
Results are rounded to the nearest nanosecond, with halves rounded away
from zero. Fractional digits after the 18th are ignored. Strings longer than 256 bytes
are rejected.

# Installation

Just go get the package:

	go get github.com/peterhellberg/duration@latest

# Usage

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
	"math/bits"
	"regexp"
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

// maxLength is the longest input Parse accepts. Real durations are far
// shorter, and the limit bounds the time spent on hostile input.
const maxLength = 256

// Components are summed as 128-bit fixed-point numbers of nanoseconds,
// with fracDigits decimal digits after the point.
const (
	fracDigits = 18
	fracScale  = 1e18
	halfScale  = fracScale / 2
)

var (
	// ErrInvalidString is returned when passed an invalid string
	ErrInvalidString = errors.New("invalid duration string")

	// ErrUnsupportedFormat is returned when parsing fails
	ErrUnsupportedFormat = errors.New("unsupported duration string format")

	pattern = regexp.MustCompile(`\AP((?P<years>[\d.,]+)Y)?((?P<months>[\d.,]+)M)?((?P<weeks>[\d.,]+)W)?((?P<days>[\d.,]+)D)?(T((?P<hours>[\d.,]+)H)?((?P<minutes>[\d.,]+)M)?((?P<seconds>[\d.,]+)S)?)?\z`)

	// units holds the nanoseconds per unit for each subexpression of pattern,
	// or 0 for unnamed ones. The constant expressions are evaluated exactly,
	// and each result is a whole number of nanoseconds.
	units = unitsBySubexp(map[string]uint64{
		"years":   HoursPerYear * 3600e9,
		"months":  HoursPerMonth * 3600e9,
		"weeks":   HoursPerWeek * 3600e9,
		"days":    HoursPerDay * 3600e9,
		"hours":   uint64(time.Hour),
		"minutes": uint64(time.Minute),
		"seconds": uint64(time.Second),
	})
)

func unitsBySubexp(byName map[string]uint64) []uint64 {
	names := pattern.SubexpNames()
	units := make([]uint64, len(names))

	for i, name := range names {
		units[i] = byName[name]
	}

	return units
}

// Parse parses an ISO 8601 duration string into a time.Duration.
//
// It returns ErrInvalidString if s contains no components or ends with "T",
// and ErrUnsupportedFormat if s is longer than 256 bytes, does not match the
// format, or the result does not fit in a time.Duration.
func Parse(s string) (time.Duration, error) {
	if len(s) > maxLength {
		return 0, ErrUnsupportedFormat
	}

	unsigned := strings.TrimPrefix(s, "-")
	negative := unsigned != s

	if unsigned == "" || unsigned == "P" || unsigned == "PT" || strings.HasSuffix(unsigned, "T") {
		return 0, ErrInvalidString
	}

	match := pattern.FindStringSubmatch(unsigned)
	if match == nil {
		return 0, ErrUnsupportedFormat
	}

	// Sum of all components in units of 1/fracScale nanoseconds.
	var hi, lo uint64

	for i, unit := range units {
		if unit == 0 || match[i] == "" {
			continue
		}

		whole, frac, ok := parseDecimal(match[i])
		if !ok {
			return 0, ErrUnsupportedFormat
		}

		// whole * unit nanoseconds, which must not exceed 1<<63 on its own.
		wh, wl := bits.Mul64(whole, unit)
		if wh != 0 || wl > 1<<63 {
			return 0, ErrUnsupportedFormat
		}

		// Scaled: whole*unit*fracScale + frac*unit, below 1.1e37 in total.
		ah, al := bits.Mul64(wl, fracScale)
		fh, fl := bits.Mul64(frac, unit)

		var c uint64

		al, c = bits.Add64(al, fl, 0)
		ah += fh + c

		lo, c = bits.Add64(lo, al, 0)
		hi += ah + c

		// halfScale<<64 is exactly 1<<63 nanoseconds, so anything with a larger
		// hi overflows. This also keeps the division and rounding below safe.
		if hi > halfScale {
			return 0, ErrUnsupportedFormat
		}
	}

	ns, rem := bits.Div64(hi, lo, fracScale)
	if rem >= halfScale {
		ns++
	}

	if ns > 1<<63 || (!negative && ns == 1<<63) {
		return 0, ErrUnsupportedFormat
	}

	if negative {
		// Converting 1<<63 wraps to math.MinInt64, which negates to itself.
		return time.Duration(-int64(ns)), nil
	}

	return time.Duration(ns), nil
}

// parseDecimal parses a non-negative decimal number with "." or "," as the
// decimal separator. It returns the integer part, and the fractional part as
// fracDigits digits, truncating any further digits. The pattern ensures that
// s contains only digits and separators.
func parseDecimal(s string) (whole, frac uint64, ok bool) {
	var (
		digits    int
		separator bool
		fracCount int
	)

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch {
		case c == '.' || c == ',':
			if separator {
				return 0, 0, false
			}

			separator = true
		default:
			digits++

			d := uint64(c - '0')

			if !separator {
				// Values above 1<<63 overflow even in seconds.
				if whole > (1<<63)/10 {
					return 0, 0, false
				}

				whole = whole*10 + d
			} else if fracCount < fracDigits {
				frac = frac*10 + d
				fracCount++
			}
		}
	}

	for ; fracCount < fracDigits; fracCount++ {
		frac *= 10
	}

	return whole, frac, digits > 0
}
