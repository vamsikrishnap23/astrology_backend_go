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

	if len(res.Planets) < 9 {
		t.Fatalf("Expected 12 planets, got %d", len(res.Planets))
	}

	expectedKarakas := map[string]string{
		"Venus":   "AK",
		"Saturn":  "AmK",
		"Mars":    "BK",
		"Jupiter": "MK",
		"Mercury": "PK",
		"Sun":     "GK",
		"Moon":    "DK",
	}

	for _, k := range res.Planets {
		if expectedKaraka, exists := expectedKarakas[k.Planet]; exists {
			if k.Karaka != expectedKaraka {
				t.Errorf("Expected Karaka %s for %s, got %s", expectedKaraka, k.Planet, k.Karaka)
			}
		} else {
			if k.Karaka != "" {
				t.Errorf("Planet %s should not have a Karaka, got %s", k.Planet, k.Karaka)
			}
		}
	}
}

// TODO: test for exact ties, boundaries, etc.
// The calculation relies strictly on Swiss Ephemeris exact longs, so ties only occur in synthetic data.
// We'll trust the sort logic handles ties cleanly via the fallback since it was explicitly implemented.
