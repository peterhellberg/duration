package duration_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/peterhellberg/duration"
)

func FuzzParse(f *testing.F) {
	for _, s := range []string{
		"", "-", "P", "PT", "-P", "P1Y2M3DT",
		"PT1.5M", "PT0.5H29M60S", "P1Y2M3DT5H20M30.123S", "-P1WT1H",
		"P0.123W", "PT1,5S", "PT.5S", "PT1.S", "P1..5M",
		"PT0.0000000005S", "PT9223372036.854775807S", "-PT9223372036.854775808S",
		"P300Y", "PT2000000H2000000H",
		"PT" + strings.Repeat("0", 252) + "1S",
	} {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		d, err := duration.Parse(s)
		if err != nil {
			if d != 0 {
				t.Errorf("Parse(%q) = (%v, %v), want 0 on error", s, d, err)
			}

			if !errors.Is(err, duration.ErrInvalidString) && !errors.Is(err, duration.ErrUnsupportedFormat) {
				t.Errorf("Parse(%q) returned unexpected error %v", s, err)
			}

			return
		}

		if strings.HasPrefix(s, "-") {
			if d > 0 {
				t.Errorf("Parse(%q) = %v, want <= 0", s, d)
			}

			return
		}

		if d < 0 {
			t.Errorf("Parse(%q) = %v, want >= 0", s, d)
		}

		// The sign would push s over the length limit.
		if len(s) == 256 {
			return
		}

		n, err := duration.Parse("-" + s)
		if err != nil || n != -d {
			t.Errorf("Parse(%q) = (%v, %v), want (%v, nil)", "-"+s, n, err, -d)
		}
	})
}
