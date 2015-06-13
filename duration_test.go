package duration_test

import (
	"fmt"
	"testing"

	"github.com/peterhellberg/duration"
)

func TestParse(t *testing.T) {
	for _, tt := range []struct {
		dur string
		err error
		out float64
	}{
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
		{"P3DT5H20M30.123S", nil, 278430},
		{"P2YT1H30M5S", duration.ErrUnsupportedFormat, 0},
		{"-P1Y", duration.ErrUnsupportedFormat, 0},
	} {
		d, err := duration.Parse(tt.dur)
		if err != tt.err {
			t.Fatalf("unexpected error: %s", err)
		}

		if got := d.Seconds(); got != tt.out {
			t.Errorf("Parse(%q) -> d.Seconds() = %+v, want %+v", tt.dur, got, tt.out)
		}
	}
}

func ExampleParse() {
	if d, err := duration.Parse("PT1M65S"); err == nil {
		fmt.Println(d)
	}

	// Output: 2m5s
}
