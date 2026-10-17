package shadbala

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/tejzpr/go-swisseph"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func near(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("%s = %.3f, want %.3f (+-%.3f)", name, got, want, tol)
	}
}

// Raman, Graha and Bhava Balas, example 17 (method B): birth at 2:20 pm local apparent time gives
// Diva Bala 48.32 (Sun, Jupiter, Venus), Ratri Bala 11.68 (Moon, Mars, Saturn), Mercury 60.
func TestNathonnathaRamanExample(t *testing.T) {
	h := 14.0 + 20.0/60.0
	for _, p := range []string{"Sun", "Jupiter", "Venus"} {
		near(t, p, nathonnathaBala(p, h), 48.32, 0.05)
	}
	for _, p := range []string{"Moon", "Mars", "Saturn"} {
		near(t, p, nathonnathaBala(p, h), 11.68, 0.05)
	}
	near(t, "Mercury", nathonnathaBala("Mercury", h), 60, 0)
}

// Noon is full strength for the Sun, Jupiter and Venus, midnight for the Moon, Mars and Saturn.
func TestNathonnathaNoonAndMidnight(t *testing.T) {
	cases := []struct {
		hour      float64
		sun, moon float64
	}{{12, 60, 0}, {0, 0, 60}, {24, 0, 60}, {6, 30, 30}, {18, 30, 30}, {23 + 10.0/60, 4.17, 55.83}}
	for _, c := range cases {
		near(t, "Sun", nathonnathaBala("Sun", c.hour), c.sun, 0.01)
		near(t, "Moon", nathonnathaBala("Moon", c.hour), c.moon, 0.01)
		// the two groups always add up to 60
		near(t, "Sun+Moon", nathonnathaBala("Sun", c.hour)+nathonnathaBala("Moon", c.hour), 60, 1e-9)
	}
}

// The equation of time: Raman's birth was 2:06:16 pm LMT at 77 deg 35' E on 16 Oct 1918, which he
// converts to 2:20 pm apparent time (EoT about +14 min).
func TestLocalApparentHoursRamanChart(t *testing.T) {
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	ephemeris.Lock()
	defer ephemeris.Unlock()
	lon := 5.0*15 + 10.0*0.25 + 20.0*(0.25/60) // 5h10m20s east
	utHours := 14.0 + 6.0/60 + 16.0/3600 - lon/15
	jd := swisseph.Julday(1918, 10, 16, utHours, swisseph.GregCal)
	got, err := localApparentHours(jd, lon)
	if err != nil {
		t.Fatal(err)
	}
	near(t, "apparent hours", got, 14.0+20.0/60, 0.03)
}

func mustAbs(t *testing.T, p string) string {
	t.Helper()
	a, err := filepath.Abs(p)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func kalaFor(t *testing.T, date, clock string, lat, lon float64, planet string) domain.KalaBala {
	t.Helper()
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	utc, _ := astronomyTime.ParseLocalToUTC(date, clock, 5.5)
	ctx := domain.CalculationContext{
		Input:       domain.BirthInput{Latitude: lat, Longitude: lon, DateOfBirth: date, TimeOfBirth: clock, Timezone: 5.5},
		Config:      domain.CalculationConfig{AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri")},
		UTCTime:     utc,
		JulianDayUT: astronomyTime.UTCToJulianDay(utc),
	}
	kb, err := CalculateKalaBala(planet, &ctx, 0, 0, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	return kb
}

// Regression: the old code picked the wrong sunrise for night births and swapped noon and midnight
// (Sun 52.7, Moon 7.3 for a birth at 23:10 in Mumbai).
func TestNathonnathaNightBirthEndToEnd(t *testing.T) {
	sun := kalaFor(t, "2012-09-03", "23:10:00", 19.076, 72.8777, "Sun")
	moon := kalaFor(t, "2012-09-03", "23:10:00", 19.076, 72.8777, "Moon")
	if sun.NathonnathaBala > 12 || moon.NathonnathaBala < 48 {
		t.Errorf("night birth: Sun %.2f (want small), Moon %.2f (want large)", sun.NathonnathaBala, moon.NathonnathaBala)
	}
	// Afternoon birth is the reverse; pre-dawn birth is moon-strong.
	if s := kalaFor(t, "2005-11-23", "13:00:00", 16.39381, 80.15221, "Sun"); s.NathonnathaBala < 50 {
		t.Errorf("afternoon Sun = %.2f, want > 50", s.NathonnathaBala)
	}
	// 03:20 IST clock in Chennai (80.28 E; IST meridian is 82.5 E) is about 3:06 apparent time
	// (-9 min longitude, -5 min equation of time): 46.5 deg from midnight, so Mars = 60 - 46.5/3.
	if m := kalaFor(t, "1985-07-10", "03:20:00", 13.0878, 80.2785, "Mars"); math.Abs(m.NathonnathaBala-44.5) > 0.7 {
		t.Errorf("3:20 am Mars = %.2f, want about 44.5", m.NathonnathaBala)
	}
}
