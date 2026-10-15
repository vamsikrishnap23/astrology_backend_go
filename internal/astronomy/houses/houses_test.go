package houses

import (
	"math"
	"testing"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
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

// Every horary number must produce an ascendant at the start of its sub-division,
// regardless of how far the real time's ascendant is from it.
func TestHoraryAscendantMatchesTable(t *testing.T) {
	if err := ephemeris.Init("../../../ephe_data"); err != nil {
		t.Fatalf("ephemeris init failed: %v", err)
	}
	defer ephemeris.Close()

	locations := []struct{ lat, lon float64 }{
		{17.3850, 78.4867}, // Hyderabad
		{51.5074, -0.1278}, // London (faster/slower sign rising rates)
		{-33.87, 151.21},   // Sydney (southern hemisphere)
	}
	times := []time.Time{
		time.Date(1998, 4, 28, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 12, 21, 23, 40, 0, 0, time.UTC),
	}

	for _, loc := range locations {
		for _, utc := range times {
			jd := astronomyTime.UTCToJulianDay(utc)
			for n := 1; n <= 249; n++ {
				ctx := domain.CalculationContext{
					Input: domain.BirthInput{
						Latitude:     loc.lat,
						Longitude:    loc.lon,
						HoraryNumber: n,
					},
					JulianDayUT: jd,
					Config: domain.CalculationConfig{
						AyanamsaMode: ephemeris.GetAyanamsaMode("KP"),
						HouseCode:    ephemeris.GetHouseSystemCode("Placidus"),
					},
				}
				asc, _, _, err := CalculateHouses(&ctx)
				if err != nil {
					t.Fatalf("lat=%v n=%d: unexpected err: %v", loc.lat, n, err)
				}
				want, _ := GetHoraryAscendant(n)
				diff := math.Abs(math.Mod(asc-want+540.0, 360.0) - 180.0)
				if diff > 1e-4 {
					t.Errorf("lat=%v utc=%v n=%d: asc=%.6f want=%.6f", loc.lat, utc, n, asc, want)
				}
			}
		}
	}
}

// Beyond the polar circle the ascendant is discontinuous. A horary request there
// may fail, but must never return an ascendant that differs from the table value.
func TestHoraryPolarNeverSilentlyWrong(t *testing.T) {
	if err := ephemeris.Init("../../../ephe_data"); err != nil {
		t.Fatalf("ephemeris init failed: %v", err)
	}
	defer ephemeris.Close()

	jd := astronomyTime.UTCToJulianDay(time.Date(1998, 4, 28, 9, 0, 0, 0, time.UTC))
	for _, lat := range []float64{70, -70} {
		for n := 1; n <= 249; n++ {
			ctx := domain.CalculationContext{
				Input:       domain.BirthInput{Latitude: lat, Longitude: 25, HoraryNumber: n},
				JulianDayUT: jd,
				Config: domain.CalculationConfig{
					AyanamsaMode: ephemeris.GetAyanamsaMode("KP"),
					HouseCode:    ephemeris.GetHouseSystemCode("Placidus"),
				},
			}
			asc, _, _, err := CalculateHouses(&ctx)
			if err != nil {
				continue
			}
			want, _ := GetHoraryAscendant(n)
			if diff := math.Abs(math.Mod(asc-want+540.0, 360.0) - 180.0); diff > 1e-4 {
				t.Errorf("lat=%v n=%d: silently wrong asc=%.4f want=%.4f", lat, n, asc, want)
			}
		}
	}
}
