package duration_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/peterhellberg/duration"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		dur string
		err error
		out float64
	}{
		{"PT1.5M", nil, 90},
		{"PT0.5H", nil, 1800},
		{"PT0.5H29M60S", nil, 3600}, // Probably shouldn’t be valid since only the last value can have fractions
		{"PT15S", nil, 15},
		{"PT1M", nil, 60},
		{"PT3M", nil, 180},
		{"PT130S", nil, 130},
		{"PT2M10S", nil, 130},
		{"P1DT2S", nil, 86402},
		{"PT5M10S", nil, 310},
		{"PT1H30M5S", nil, 5405},
		{"P2DT1H10S", nil, 176410},
		{"PT1004199059S", nil, 1004199059},
		{"P3DT5H20M30.123S", nil, 278430.123},
		{"P1W", nil, 604800},
		{"P0.123W", nil, 74390.4},
		{"P1WT5S", nil, 604805},
		{"P1WT1H", nil, 608400},
		{"P2YT1H30M5S", nil, 63119237},
		{"P1Y2M3DT5H20M30.123S", nil, 37094832.1218},

		// Not supported since negative period
		{"-P1Y", duration.ErrUnsupportedFormat, 0},

		// Not supported since fields in the wrong order
		{"P1M2Y", duration.ErrUnsupportedFormat, 0},

		// Not supported since negative value
		{"P-1Y", duration.ErrUnsupportedFormat, 0},

		// Not supported since negative value
		{"P1YT-1M", duration.ErrUnsupportedFormat, 0},

		// Not supported since missing T
		{"P1S", duration.ErrUnsupportedFormat, 0},

		// Not supported since missing P
		{"1Y", duration.ErrUnsupportedFormat, 0},

		// Not supported since no value is specified for months
		{"P1YM5D", duration.ErrUnsupportedFormat, 0},

		// Not supported since wrong format of string
		{"FOOBAR", duration.ErrUnsupportedFormat, 0},

		// Invalid since empty string
		{"", duration.ErrInvalidString, 0},

		// Invalid since no time fields present
		{"P", duration.ErrInvalidString, 0},

		// Invalid since no time fields present
		{"PT", duration.ErrInvalidString, 0},

		// Invalid since ending with T
		{"P1Y2M3DT", duration.ErrInvalidString, 0},
	} {
		d, err := duration.Parse(tt.dur)
		if err != tt.err {
			t.Fatalf("unexpected error: %s", err)
		}

		if got := d.Seconds(); got != tt.out {
			t.Errorf("Parse(%q) -> d.Seconds() = %f, want %f", tt.dur, got, tt.out)
		}
	}
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
		td, _ := time.ParseDuration(tt.timeStr)
		dd, _ := duration.Parse(tt.durationStr)

		if td != dd {
			t.Errorf(`not equal: %q->%v != %q->%v`, tt.timeStr, td, tt.durationStr, dd)
		}
	}
}

func ExampleParse() {
	if d, err := duration.Parse("PT1M65S"); err == nil {
		fmt.Println(d)
	}

	// Output: 2m5s
}
