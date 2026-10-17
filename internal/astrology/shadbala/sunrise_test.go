package shadbala

import (
	"errors"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// Raman art. 56-57: the day and the night are each cut in three; Mercury, the Sun and Saturn rule the
// parts of the day, the Moon, Venus and Mars those of the night; Jupiter always has 60.
func TestTribhagaBalaThirds(t *testing.T) {
	day := hinduDay{Sunrise: 100.0, Sunset: 100.5, NextSunrise: 101.0} // thirds of 1/6 day each
	cases := []struct {
		jd   float64
		lord string
	}{
		{100.0, "Mercury"}, {100.16, "Mercury"}, {100.17, "Sun"}, {100.33, "Sun"}, {100.34, "Saturn"}, {100.49, "Saturn"},
		{100.5, "Moon"}, {100.66, "Moon"}, {100.67, "Venus"}, {100.83, "Venus"}, {100.84, "Mars"}, {100.99, "Mars"},
	}
	for _, c := range cases {
		for _, p := range []string{"Sun", "Moon", "Mars", "Mercury", "Venus", "Saturn"} {
			want := 0.0
			if p == c.lord {
				want = 60
			}
			near(t, p, tribhagaBala(p, day, c.jd), want, 0)
		}
		near(t, "Jupiter", tribhagaBala("Jupiter", day, c.jd), 60, 0)
	}
}

// The Hindu day found for a moment starts at the latest sunrise: sunrise <= jd < next sunrise, with
// the sunset in between.
func TestHinduDayAt(t *testing.T) {
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	geopos := [3]float64{80.2785, 13.0878, 0}
	base := 2446256.0 // 10 July 1985, 12 h UT is 2446257.0
	var prevSunrise float64
	for h := 0.0; h < 48; h += 0.37 {
		jd := base + h/24
		d, err := hinduDayAt(jd, geopos)
		if err != nil {
			t.Fatal(err)
		}
		if !(d.Sunrise <= jd && jd < d.NextSunrise && d.Sunrise < d.Sunset && d.Sunset < d.NextSunrise) {
			t.Fatalf("jd %.4f: day %+v does not bracket it", jd, d)
		}
		if l := d.NextSunrise - d.Sunrise; l < 0.99 || l > 1.01 {
			t.Errorf("a Hindu day should last about 24 hours, got %.4f", l)
		}
		if d.Sunrise < prevSunrise {
			t.Errorf("sunrise went backwards at %.4f", jd)
		}
		prevSunrise = d.Sunrise
	}
}

// Regression: Tribhaga Bala used to treat a birth in the first hour or so after sunrise as night (the
// sunset it looked up was the previous evening's). Chennai, 10 July 1985: sunrise is about 05:50 IST.
func TestTribhagaJustAfterSunrise(t *testing.T) {
	cases := []struct {
		clock string
		lord  string
	}{
		{"05:20:00", "Mars"}, {"05:40:00", "Mars"}, // last third of the night
		{"06:00:00", "Mercury"}, {"06:10:00", "Mercury"}, {"06:20:00", "Mercury"}, {"06:40:00", "Mercury"}, {"09:00:00", "Mercury"},
		{"11:00:00", "Sun"}, {"14:00:00", "Sun"}, {"16:00:00", "Saturn"}, {"18:00:00", "Saturn"},
		{"19:00:00", "Moon"}, {"22:00:00", "Moon"}, {"23:30:00", "Venus"}, {"01:30:00", "Venus"},
	}
	for _, c := range cases {
		got := ""
		for _, p := range []string{"Sun", "Moon", "Mars", "Mercury", "Venus", "Saturn"} {
			if kalaFor(t, "1985-07-10", c.clock, 13.0878, 80.2785, p).TribhagaBala == 60 {
				got += p + " "
			}
		}
		if got != c.lord+" " {
			t.Errorf("%s: Tribhaga Bala goes to %q, want %q", c.clock, got, c.lord)
		}
	}
}

func polarContext(t *testing.T, date, clock string, lat float64) domain.CalculationContext {
	t.Helper()
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	utc, _ := astronomyTime.ParseLocalToUTC(date, clock, 0)
	return domain.CalculationContext{
		Input:       domain.BirthInput{Latitude: lat, Longitude: 25.0, DateOfBirth: date, TimeOfBirth: clock},
		Config:      domain.CalculationConfig{AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri"), HouseCode: ephemeris.GetHouseSystemCode("Placidus")},
		UTCTime:     utc,
		JulianDayUT: astronomyTime.UTCToJulianDay(utc),
	}
}

// Where the Sun does not rise or set there is no sunrise to count Hora, Dina, Varsha, Masa and Tribhaga
// Bala from; the calculation must say so instead of returning zeros.
func TestNoSunriseGivesAnError(t *testing.T) {
	for _, c := range []struct{ name, date string }{{"midnight sun", "2023-06-21"}, {"polar night", "2023-12-21"}} {
		ctx := polarContext(t, c.name+"", "12:00:00", 80)
		ctx = polarContext(t, c.date, "12:00:00", 80)
		_, err := CalculateShadbala(&ctx)
		if !errors.Is(err, ErrNoSunriseSunset) {
			t.Errorf("%s at 80 N: got %v, want ErrNoSunriseSunset", c.name, err)
		}
	}
	// A high but ordinary latitude still works, including a short midsummer night and a short midwinter day.
	for _, c := range []struct {
		date string
		lat  float64
	}{{"2023-06-21", 60}, {"2023-12-21", 60}, {"2023-03-21", 60}, {"2023-06-21", 64}, {"2023-12-21", 64}, {"2023-12-21", -64}, {"2023-06-21", -64}} {
		date := c.date
		ctx := polarContext(t, date, "12:00:00", c.lat)
		res, err := CalculateShadbala(&ctx)
		if err != nil {
			t.Fatalf("%.0f N on %s: %v", c.lat, date, err)
		}
		if len(res.Planets) != 7 {
			t.Errorf("%.0f N on %s: %d planets", c.lat, date, len(res.Planets))
		}
	}
}
