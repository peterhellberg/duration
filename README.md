# duration

[![Build Status](https://travis-ci.org/peterhellberg/duration.svg?branch=master)](https://travis-ci.org/peterhellberg/duration)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/duration)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/duration#license-mit)

Parse a [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) duration string into `time.Duration`

There are probably a few unsupported edge cases still to be fixed, please help me find them :)

The following constants are used to do the calculations for longer durations:

```
HoursPerDay = 24.0
HoursPerWeek = 168.0
HoursPerMonth = 730.4841667
HoursPerYear = 8765.81
```

Look in the test for examples of both valid and invalid duration strings.

## Installation

    go get -u github.com/peterhellberg/duration

Feel free to copy this package into your own codebase.

## Usage

```go
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
```

## RFC3339 grammar for durations

```
   dur-second        = 1*DIGIT "S"
   dur-minute        = 1*DIGIT "M" [dur-second]
   dur-hour          = 1*DIGIT "H" [dur-minute]
   dur-time          = "T" (dur-hour / dur-minute / dur-second)
   dur-day           = 1*DIGIT "D"
   dur-week          = 1*DIGIT "W"
   dur-month         = 1*DIGIT "M" [dur-day]
   dur-year          = 1*DIGIT "Y" [dur-month]
   dur-date          = (dur-day / dur-month / dur-year) [dur-time]

   duration          = "P" (dur-date / dur-time / dur-week)
```

## License (MIT)

Copyright (c) 2015 [Peter Hellberg](http://c7.se/)

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:

> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
