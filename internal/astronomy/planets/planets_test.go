package planets

import (
	"testing"
	"time"

	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func TestCalculatePlanets(t *testing.T) {
	ephemeris.Init("")
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
