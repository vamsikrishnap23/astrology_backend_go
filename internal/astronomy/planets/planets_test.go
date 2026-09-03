package planets

import (
	"testing"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func TestCalculatePlanets(t *testing.T) {
	if err := ephemeris.Init("../../../ephe_data"); err != nil {
		t.Fatalf("ephemeris init failed: %v", err)
	}
	defer ephemeris.Close()

	utc := time.Date(1998, 4, 28, 9, 0, 0, 0, time.UTC)
	jd := astronomyTime.UTCToJulianDay(utc)

	ctx := domain.CalculationContext{
		JulianDayUT: jd,
		Config: domain.CalculationConfig{
			AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri"),
		},
	}

	positions, err := CalculatePlanets(&ctx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	foundRahu := false
	foundKetu := false

	for _, p := range positions {
		if p.Planet == "Rahu" {
			foundRahu = true
			if p.Speed > 0 { // Typically nodes are retrograde (negative speed)
				t.Log("Warning: Mean Rahu speed is positive")
			}
		}
		if p.Planet == "Ketu" {
			foundKetu = true
		}

		// Retrograde detection test
		if p.Speed < 0 && !p.Retrograde {
			t.Errorf("Expected Retrograde=true for planet %s with speed %f", p.Planet, p.Speed)
		}
		if p.Speed >= 0 && p.Retrograde {
			t.Errorf("Expected Retrograde=false for planet %s with speed %f", p.Planet, p.Speed)
		}
	}

	if !foundRahu || !foundKetu {
		t.Error("Missing Nodes")
	}
}
