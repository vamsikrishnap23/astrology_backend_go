package progression

import (
	"math"
	"testing"
	"time"

	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
	"path/filepath"
)

func TestCalculateSecondaryProgression(t *testing.T) {
	// Initialize Ephemeris for tests
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	birthStr := "2005-11-23"
	timeStr := "15:35:00"
	tz := 5.5

	utcTime, _ := astronomyTime.ParseLocalToUTC(birthStr, timeStr, tz)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	natalCtx := domain.CalculationContext{
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

	targetYear := 2040
	res, err := CalculateSecondaryProgression(&natalCtx, targetYear)
	if err != nil {
		t.Fatalf("Progression failed: %v", err)
	}

	if res.TargetProgressionYear != 2040 {
		t.Errorf("Expected year 2040, got %d", res.TargetProgressionYear)
	}

	expectedAge := 2040.0 - 2005.0
	if math.Abs(res.AgeInYears-expectedAge) > 1e-6 {
		t.Errorf("Expected age 35, got %f", res.AgeInYears)
	}

	expectedProgressedJD := jd + expectedAge
	if math.Abs(res.ProgressedJulianDay-expectedProgressedJD) > 1e-6 {
		t.Errorf("Expected Progressed JD %f, got %f", expectedProgressedJD, res.ProgressedJulianDay)
	}

	expectedProgressedUTC := utcTime.Add(time.Duration(expectedAge * 24 * float64(time.Hour))).Format(time.RFC3339)
	if res.ProgressedDateUTC != expectedProgressedUTC {
		t.Errorf("Expected Progressed UTC %s, got %s", expectedProgressedUTC, res.ProgressedDateUTC)
	}

	// Wait, we need to check if Sun progressed by about 35 degrees (slightly less since tropical year is ~365.24, day is 360/365.24 deg)
	// Actually Sun moves ~1 degree per day. Age = 35 days. Sun moves ~35 degrees.
	// Let's just ensure we have planets.
	if len(res.ProgressedPlanets) < 10 {
		t.Errorf("Expected at least 10 planets, got %d", len(res.ProgressedPlanets))
	}
	if len(res.ProgressedHouses) != 12 {
		t.Errorf("Expected 12 houses, got %d", len(res.ProgressedHouses))
	}
}

// Zero elapsed time test
func TestCalculateSecondaryProgression_ZeroElapsed(t *testing.T) {
	// Initialize Ephemeris for tests
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	birthStr := "2005-11-23"
	timeStr := "15:35:00"
	tz := 5.5

	utcTime, _ := astronomyTime.ParseLocalToUTC(birthStr, timeStr, tz)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	natalCtx := domain.CalculationContext{
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

	targetYear := 2005 // Zero elapsed
	res, err := CalculateSecondaryProgression(&natalCtx, targetYear)
	if err != nil {
		t.Fatalf("Progression failed: %v", err)
	}

	if res.AgeInYears != 0 {
		t.Errorf("Expected age 0, got %f", res.AgeInYears)
	}

	if res.ProgressedJulianDay != jd {
		t.Errorf("Expected Progressed JD to be equal to natal JD for year 0")
	}
}
