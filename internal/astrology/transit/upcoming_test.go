package transit

import (
	"path/filepath"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
)

func TestCalculateUpcomingTransits(t *testing.T) {
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	utcTime, _ := astronomyTime.ParseLocalToUTC("2024-05-15", "10:00:00", 5.5)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	ayanamsaMode := ephemeris.GetAyanamsaMode("Lahiri")

	res := CalculateUpcomingTransits(jd, ayanamsaMode)

	if len(res) != 9 {
		t.Errorf("Expected 9 transiting planets, got %d", len(res))
	}

	for _, p := range res {
		if p.Planet == "" {
			t.Errorf("Empty planet name in transits")
		}
		if p.DestinationSign == "" {
			t.Errorf("Empty destination sign for %s", p.Planet)
		}
		if p.TransitionDateTime == "" {
			t.Errorf("Empty transition date for %s", p.Planet)
		}
	}
}
