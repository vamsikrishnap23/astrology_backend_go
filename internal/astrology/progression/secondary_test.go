package progression

import (
	"math"
	"testing"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
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

	targetDate := "2040-11-23"
	targetUtc, _ := astronomyTime.ParseLocalToUTC(targetDate, timeStr, tz)
	targetJD := astronomyTime.UTCToJulianDay(targetUtc)

	res, err := CalculateSecondaryProgression(&natalCtx, targetDate, targetJD)
	if err != nil {
		t.Fatalf("Progression failed: %v", err)
	}

	if res.TargetProgressionDate != targetDate {
		t.Errorf("Expected date %s, got %s", targetDate, res.TargetProgressionDate)
	}

	expectedAge := (targetJD - jd) / 365.242190402
	if math.Abs(res.AgeInYears-expectedAge) > 1e-6 {
		t.Errorf("Expected age ~35, got %f", res.AgeInYears)
	}

	expectedProgressedJD := jd + expectedAge
	if math.Abs(res.ProgressedJulianDay-expectedProgressedJD) > 1e-6 {
		t.Errorf("Expected Progressed JD %f, got %f", expectedProgressedJD, res.ProgressedJulianDay)
	}

	expectedProgressedUTC := utcTime.Add(time.Duration(expectedAge * 24 * float64(time.Hour))).Format(time.RFC3339)
	if res.ProgressedDateUTC != expectedProgressedUTC {
		t.Errorf("Expected Progressed UTC %s, got %s", expectedProgressedUTC, res.ProgressedDateUTC)
	}

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

	targetDate := "2005-11-23"
	res, err := CalculateSecondaryProgression(&natalCtx, targetDate, jd)
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
