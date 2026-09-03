package shadbala

import (
	"path/filepath"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func TestCalculateShadbala(t *testing.T) {
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

	res, err := CalculateShadbala(&ctx)
	if err != nil {
		t.Fatalf("Failed to calculate shadbala: %v", err)
	}

	if len(res.Planets) != 7 {
		t.Errorf("Expected 7 planets, got %d", len(res.Planets))
	}

	for _, p := range res.Planets {
		if p.TotalShadbala <= 0 {
			t.Errorf("Planet %s has 0 or negative total shadbala", p.Planet)
		}
		if p.Rupas <= 0 {
			t.Errorf("Planet %s has 0 or negative rupas", p.Planet)
		}
		if p.NaisargikaBala <= 0 {
			t.Errorf("Planet %s has 0 or negative naisargika bala", p.Planet)
		}
	}
}
