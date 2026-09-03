package houses

import (
	"testing"
	"time"

	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func TestCalculateHouses(t *testing.T) {
	if err := ephemeris.Init("../../../ephe_data"); err != nil {
		t.Fatalf("ephemeris init failed: %v", err)
	}
	defer ephemeris.Close()

	utc := time.Date(1998, 4, 28, 9, 0, 0, 0, time.UTC)
	jd := astronomyTime.UTCToJulianDay(utc)

	ctx := domain.CalculationContext{
		Input: domain.BirthInput{
			Latitude:  17.3850,
			Longitude: 78.4867,
		},
		JulianDayUT: jd,
		Config: domain.CalculationConfig{
			AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri"),
			HouseCode:    ephemeris.GetHouseSystemCode("Placidus"),
		},
	}

	asc, mc, cusps, err := CalculateHouses(&ctx)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if asc == 0 || mc == 0 {
		t.Error("Invalid asc or mc")
	}

	if len(cusps) != 12 {
		t.Errorf("Expected 12 cusps, got %d", len(cusps))
	}

	// Ascendant usually matches 1st house cusp in quadrant systems
	if asc != cusps[0].Longitude {
		t.Errorf("Ascendant (%f) does not match 1st house cusp (%f)", asc, cusps[0].Longitude)
	}
}
