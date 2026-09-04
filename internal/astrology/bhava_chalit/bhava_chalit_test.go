package bhava_chalit

import (
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"path/filepath"
)

func TestCalculateBhavaChalit(t *testing.T) {
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	birthStr := "2005-11-23"
	timeStr := "15:35:00"
	tz := 5.5

	utcTime, _ := astronomyTime.ParseLocalToUTC(birthStr, timeStr, tz)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	ctx := domain.CalculationContext{
		Input: domain.BirthInput{
			Latitude:  16.3938,
			Longitude: 80.1522,
		},
		Config: domain.CalculationConfig{
			AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri"),
			HouseCode:    ephemeris.GetHouseSystemCode("Placidus"),
		},
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	res, err := CalculateBhavaChalit(&ctx)
	if err != nil {
		t.Fatalf("Bhava Chalit calculation failed: %v", err)
	}

	if len(res.Houses) != 12 {
		t.Errorf("Expected 12 houses, got %d", len(res.Houses))
	}

	// Test if planets are correctly flattened and contain house numbers
	if len(res.Planets) < 10 { // Sun through Pluto + Rahu/Ketu = 12 total usually
		t.Errorf("Expected at least 10 planets mapped, got %d", len(res.Planets))
	}

	for _, p := range res.Planets {
		if p.Planet == "" || p.DivisionalSign == "" {
			t.Errorf("Planet missing fields: %+v", p)
		}
		if p.HouseNumber < 1 || p.HouseNumber > 12 {
			t.Errorf("Invalid house number for planet %s: %d", p.Planet, p.HouseNumber)
		}
	}
}
