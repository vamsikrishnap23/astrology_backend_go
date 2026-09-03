package ashtakavarga

import (
	"path/filepath"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func TestCalculateAshtakavarga(t *testing.T) {
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	utcTime, _ := astronomyTime.ParseLocalToUTC("2005-11-23", "15:35:00", 5.5)
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

	res, err := CalculateAshtakavarga(&ctx)
	if err != nil {
		t.Fatalf("Failed to calculate ashtakavarga: %v", err)
	}

	// SAV Total Check (Must be 337 exactly)
	if res.TotalSAVBindus != 337 {
		t.Errorf("Expected exactly 337 Total SAV Bindus, got %d", res.TotalSAVBindus)
	}

	// SAV slice check
	if len(res.SAV) != 12 {
		t.Errorf("Expected 12 signs in SAV, got %d", len(res.SAV))
	}

	var manualSAVSum int
	for _, sign := range res.SAV {
		manualSAVSum += sign.TotalBindus
	}
	if manualSAVSum != 337 {
		t.Errorf("Sum of SAV signs %d != 337", manualSAVSum)
	}

	// BAV total check
	expectedBAVTotals := map[string]int{
		"Sun":     48,
		"Moon":    49,
		"Mars":    39,
		"Mercury": 54,
		"Jupiter": 56,
		"Venus":   52,
		"Saturn":  39,
	}

	if len(res.BAV) != 7 {
		t.Errorf("Expected 7 planets for BAV, got %d", len(res.BAV))
	}

	for _, bav := range res.BAV {
		expectedTotal, ok := expectedBAVTotals[bav.Planet]
		if !ok {
			t.Errorf("Unexpected planet %s in BAV", bav.Planet)
			continue
		}

		if bav.TotalBindus != expectedTotal {
			t.Errorf("Expected %s BAV total to be %d, got %d", bav.Planet, expectedTotal, bav.TotalBindus)
		}

		var manualBAVSum int
		for _, sign := range bav.Signs {
			manualBAVSum += sign.TotalBindus
		}
		if manualBAVSum != expectedTotal {
			t.Errorf("Sum of %s signs %d != expected %d", bav.Planet, manualBAVSum, expectedTotal)
		}
	}
}
