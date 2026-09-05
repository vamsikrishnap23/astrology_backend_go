package transit

import (
	"math"
	"testing"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"path/filepath"
)

func TestCalculateTransitChart(t *testing.T) {
	// Initialize Ephemeris for tests
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	transitStr := "2024-05-15"
	timeStr := "10:00:00"
	tz := 5.5

	utcTime, _ := astronomyTime.ParseLocalToUTC(transitStr, timeStr, tz)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	natalStr := "2005-11-23"
	natalTimeStr := "15:35:00"
	natalUtcTime, _ := astronomyTime.ParseLocalToUTC(natalStr, natalTimeStr, tz)
	natalJd := astronomyTime.UTCToJulianDay(natalUtcTime)

	natalCtx := domain.CalculationContext{
		Input: domain.BirthInput{
			Latitude:  16.3938,
			Longitude: 80.1522,
		},
		Config: domain.CalculationConfig{
			AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri"),
			HouseCode:    ephemeris.GetHouseSystemCode("Placidus"),
		},
		UTCTime:     natalUtcTime,
		JulianDayUT: natalJd,
	}

	transitCtx := domain.CalculationContext{
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

	res, err := CalculateTransitChart(&natalCtx, &transitCtx)
	if err != nil {
		t.Fatalf("Transit failed: %v", err)
	}

	expectedUTC := utcTime.Format(time.RFC3339)
	if res.TransitDateUTC != expectedUTC {
		t.Errorf("Expected Transit UTC %s, got %s", expectedUTC, res.TransitDateUTC)
	}

	if math.Abs(res.JulianDay-jd) > 1e-6 {
		t.Errorf("Expected Transit JD %f, got %f", jd, res.JulianDay)
	}

	if len(res.TransitData.PlanetaryTable) < 10 {
		t.Errorf("Expected at least 10 planets, got %d", len(res.TransitData.PlanetaryTable))
	}

	if len(res.TransitData.HouseTable) != 12 {
		t.Errorf("Expected 12 houses, got %d", len(res.TransitData.HouseTable))
	}
}
