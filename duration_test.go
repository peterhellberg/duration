package duration_test

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/peterhellberg/duration"
)

type parseTest struct {
	dur string
	err error
	out time.Duration
}

func runParseTests(t *testing.T, tests []parseTest) {
	t.Helper()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.dur, func(t *testing.T) {
			d, err := duration.Parse(tt.dur)
			if d != tt.out || !errors.Is(err, tt.err) {
				t.Errorf("Parse(%q) = (%v, %v), want (%v, %v)", tt.dur, d, err, tt.out, tt.err)
			}
		})
	}
}

// d parses a time.ParseDuration string, for readable table entries.
func d(s string) time.Duration {
	v, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}

	return v
}

func TestParse(t *testing.T) {
	runParseTests(t, []parseTest{
		{"PT1.5M", nil, d("90s")},
		{"PT0.5H", nil, d("1800s")},
		{"PT0.5H29M60S", nil, d("3600s")}, // Probably shouldn’t be valid since only the last value can have fractions
		{"PT15S", nil, d("15s")},
		{"PT1M", nil, d("60s")},
		{"PT3M", nil, d("180s")},
		{"PT130S", nil, d("130s")},
		{"PT2M10S", nil, d("130s")},
		{"P1DT2S", nil, d("86402s")},
		{"PT5M10S", nil, d("310s")},
		{"PT1H30M5S", nil, d("5405s")},
		{"P2DT1H10S", nil, d("176410s")},
		{"PT1004199059S", nil, d("1004199059s")},
		{"P3DT5H20M30.123S", nil, d("278430.123s")},
		{"P1W", nil, d("604800s")},
		{"P0.123W", nil, d("74390.4s")},
		{"P1WT5S", nil, d("604805s")},
		{"P1WT1H", nil, d("608400s")},
		{"P2YT1H30M5S", nil, d("63119237s")},
		{"P1Y2M3DT5H20M30.123S", nil, d("37094832.12324s")},
		{"-P1Y2M3DT5H20M30.123S", nil, d("-37094832.12324s")},
		{"-P1WT1H", nil, d("-608400s")},
		{"-P1DT2S", nil, d("-86402s")},
		{"-PT1M5S", nil, d("-65s")},
		{"-P0.123W", nil, d("-74390.4s")},

		// M means months before T and minutes after it
		{"P1M", nil, d("2629743.00012s")},
		{"PT1M", nil, d("60s")},
		{"P1MT1M", nil, d("2629803.00012s")},

		// Not supported since fields in the wrong order
		{"P1M2Y", duration.ErrUnsupportedFormat, 0},

		// Not supported since negative value
		{"P-1Y", duration.ErrUnsupportedFormat, 0},

		// Not supported since negative value
		{"P1YT-1M", duration.ErrUnsupportedFormat, 0},

		// Not supported since more than one sign
		{"--P1D", duration.ErrUnsupportedFormat, 0},

		// Not supported since missing T
		{"P1S", duration.ErrUnsupportedFormat, 0},

		// Not supported since missing P
		{"1Y", duration.ErrUnsupportedFormat, 0},

		// Not supported since no value is specified for months
		{"P1YM5D", duration.ErrUnsupportedFormat, 0},

		// Not supported since wrong format of string
		{"FOOBAR", duration.ErrUnsupportedFormat, 0},

		// Not supported since lowercase
		{"p1d", duration.ErrUnsupportedFormat, 0},

		// Invalid since empty string
		{"", duration.ErrInvalidString, 0},
		{"-", duration.ErrInvalidString, 0},

		// Invalid since no time fields present
		{"P", duration.ErrInvalidString, 0},
		{"-P", duration.ErrInvalidString, 0},

		// Invalid since no time fields present
		{"PT", duration.ErrInvalidString, 0},
		{"-PT", duration.ErrInvalidString, 0},

		// Invalid since ending with T
		{"P1Y2M3DT", duration.ErrInvalidString, 0},
		{"-P1Y2M3DT", duration.ErrInvalidString, 0},
	})
}

func TestParseCommaDecimalSeparator(t *testing.T) {
	runParseTests(t, []parseTest{
		{"PT1,5S", nil, d("1.5s")},
		{"P0,5D", nil, d("12h")},
		{"-PT1,5M", nil, d("-90s")},
		{"PT,5S", nil, d("0.5s")},
		{"PT1,,5S", duration.ErrUnsupportedFormat, 0},
		{"PT1,5.5S", duration.ErrUnsupportedFormat, 0},
		{"PT,S", duration.ErrUnsupportedFormat, 0},
	})
}

func TestParsePrecision(t *testing.T) {
	runParseTests(t, []parseTest{
		{"PT0.000000001S", nil, time.Nanosecond},
		{"PT0.0000001S", nil, 100 * time.Nanosecond},
		{"PT0.0000015S", nil, 1500 * time.Nanosecond},
		{"PT0.0000000004S", nil, 0},
		{"PT0.0000000005S", nil, time.Nanosecond},
		{"-PT0.0000000005S", nil, -time.Nanosecond},

		// Fractions are summed before rounding
		{"PT0.0000000000001H0.0000000004S", nil, time.Nanosecond},

		// Fractional digits after the 18th are ignored
		{"PT1.5000000000000000000000000000000000001S", nil, d("1.5s")},
		{"PT0.000000000499999999999999999999S", nil, 0},
		{"PT0.00000000000000000000000000000000001S", nil, 0},
		{"PT0" + strings.Repeat("0", 240) + ".5S", nil, d("0.5s")},

		{"PT9223372036.854775807S", nil, math.MaxInt64},
		{"PT2562047H47M16.854775807S", nil, math.MaxInt64},
		{"PT2562047H47M16.854775808S", duration.ErrUnsupportedFormat, 0},
		{"-PT2562047H47M16.854775808S", nil, math.MinInt64},
		{"PT9223372036.854775808S", duration.ErrUnsupportedFormat, 0},
		{"-PT9223372036.854775808S", nil, math.MinInt64},
		{"-PT9223372036.854775809S", duration.ErrUnsupportedFormat, 0},
		{"PT9223372036.8547758074S", nil, math.MaxInt64},
		{"PT9223372036.8547758075S", duration.ErrUnsupportedFormat, 0},
	})
}

func TestParseLengthLimit(t *testing.T) {
	runParseTests(t, []parseTest{
		{"PT" + strings.Repeat("0", 252) + "1S", nil, time.Second},
		{"-PT" + strings.Repeat("0", 251) + "1S", nil, -time.Second},
		{"PT" + strings.Repeat("0", 253) + "1S", duration.ErrUnsupportedFormat, 0},
		{"-PT" + strings.Repeat("0", 252) + "1S", duration.ErrUnsupportedFormat, 0},
		{strings.Repeat("x", 1024), duration.ErrUnsupportedFormat, 0},
	})
}

func TestCompareWithTimeParseDuration(t *testing.T) {
	for _, tt := range []struct {
		timeStr     string
		durationStr string
	}{
		{"1h", "PT1H"},
		{"9m60s", "PT10.0M"},
		{"1h2m", "PT1H2M"},
		{"2h15s", "PT1H60M15S"},
		{"169h", "P1WT1H"},
	} {
		td, err := time.ParseDuration(tt.timeStr)
		if err != nil {
			t.Fatalf("time.ParseDuration(%q): %v", tt.timeStr, err)
		}

		dd, err := duration.Parse(tt.durationStr)
		if err != nil {
			t.Fatalf("duration.Parse(%q): %v", tt.durationStr, err)
		}

		if td != dd {
			t.Errorf("not equal: %q->%v != %q->%v", tt.timeStr, td, tt.durationStr, dd)
		}
	}
}

func ExampleParse() {
	if d, err := duration.Parse("PT1M65S"); err == nil {
		fmt.Println(d)
	}

	// Output: 2m5s
}

func TestParseNumericComponents(t *testing.T) {
	runParseTests(t, []parseTest{
		{"P.Y", duration.ErrUnsupportedFormat, 0},
		{"P1..5M", duration.ErrUnsupportedFormat, 0},
		{"P..W", duration.ErrUnsupportedFormat, 0},
		{"P.5.5D", duration.ErrUnsupportedFormat, 0},
		{"PT1..5H", duration.ErrUnsupportedFormat, 0},
		{"PT..M", duration.ErrUnsupportedFormat, 0},
		{"PT1..5S", duration.ErrUnsupportedFormat, 0},
		{"P1Y..M", duration.ErrUnsupportedFormat, 0},
		{"P..YT1S", duration.ErrUnsupportedFormat, 0},
		{"-P1DT1..5S", duration.ErrUnsupportedFormat, 0},
		{"PT.5S", nil, 500 * time.Millisecond},
		{"PT1.S", nil, time.Second},
		{"PT0.5H29M60S", nil, time.Hour},
		{"P1WT0.5H", nil, 168*time.Hour + 30*time.Minute},
		{"-P1WT0.5H", nil, -(168*time.Hour + 30*time.Minute)},
	})
}

func TestParseOverflow(t *testing.T) {
	runParseTests(t, []parseTest{
		{"P300Y", duration.ErrUnsupportedFormat, 0},
		{"-P300Y", duration.ErrUnsupportedFormat, 0},
		{"PT3000000H", duration.ErrUnsupportedFormat, 0},
		{"PT10000000000000000000000000000000000000000S", duration.ErrUnsupportedFormat, 0},
		{"PT1" + strings.Repeat("0", 250) + "S", duration.ErrUnsupportedFormat, 0},
		{"PT2000000H2000000H", duration.ErrUnsupportedFormat, 0},
		{"P100YT2000000H", duration.ErrUnsupportedFormat, 0},
		{"-P100YT2000000H", duration.ErrUnsupportedFormat, 0},
		{"PT2562047H", nil, 2562047 * time.Hour},
		{"-PT2562047H", nil, -2562047 * time.Hour},
		{"P100Y", nil, 876581 * time.Hour},
	})
}
