package jaimini

import (
	"path/filepath"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func TestVamsiRegression(t *testing.T) {
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	utcTime, _ := astronomyTime.ParseLocalToUTC("2005-11-23", "15:35:00", 5.5)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	ctx := domain.CalculationContext{
		Input: domain.BirthInput{
			Latitude:    16.3938,
			Longitude:   80.1522,
			DateOfBirth: "2005-11-23",
			TimeOfBirth: "15:35:00",
			Timezone:    5.5,
		},
		Config: domain.CalculationConfig{
			AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri"),
			HouseCode:    ephemeris.GetHouseSystemCode("Placidus"),
		},
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	res, err := CalculateCharaKarakas(&ctx)
	if err != nil {
		t.Fatalf("Failed to calculate jaimini karakas: %v", err)
	}

	if len(res.Planets) != 7 {
		t.Fatalf("Expected 7 karakas, got %d", len(res.Planets))
	}

	expectedKarakas := []string{"Venus", "Saturn", "Mars", "Jupiter", "Mercury", "Sun", "Moon"}
	expectedNames := []string{"AK", "AmK", "BK", "MK", "PK", "GK", "DK"}

	for i, k := range res.Planets {
		if k.Planet != expectedKarakas[i] {
			t.Errorf("Rank %d: Expected Planet %s, got %s", i+1, expectedKarakas[i], k.Planet)
		}
		if k.Karaka != expectedNames[i] {
			t.Errorf("Rank %d: Expected Karaka %s, got %s", i+1, expectedNames[i], k.Karaka)
		}
	}
}

// TODO: test for exact ties, boundaries, etc.
// The calculation relies strictly on Swiss Ephemeris exact longs, so ties only occur in synthetic data.
// We'll trust the sort logic handles ties cleanly via the fallback since it was explicitly implemented.
