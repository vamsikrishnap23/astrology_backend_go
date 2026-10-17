package shadbala

import (
	"testing"

	"github.com/tejzpr/go-swisseph"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// Raman, Graha and Bhava Balas, example 33 (Ayana Bala, standard horoscope). The kranti values are
// the ones he derives in example 32; his printed results are in the comments.
func TestAyanaBalaRamanExample(t *testing.T) {
	cases := []struct {
		planet string
		kranti float64 // north positive
		want   float64
	}{
		{"Sun", -8.75, 38.12},    // 8.75 S: (24-8.75)/48*60 = 19.06, doubled
		{"Moon", -10.75, 43.44},  // 10.75 S, south adds for the Moon
		{"Mercury", -9.0, 41.25}, // 9.0 S, always additive
		{"Jupiter", 23.5, 59.4},  // 23.5 N
		{"Venus", -5.0, 23.75},   // printed with 4.96 S
		{"Saturn", 13.0, 13.75},  // 13 N, north subtracts for Saturn
	}
	for _, c := range cases {
		near(t, c.planet, ayanaBala(c.planet, c.kranti), c.want, 0.03)
	}
	// Mars 22.45 S: the formula gives (24-22.45)*1.25 = 1.94; Raman prints 1.84 (his rounding of the kranti).
	near(t, "Mars", ayanaBala("Mars", -22.45), 1.94, 0.01)
}

// Mercury's kranti is additive whichever hemisphere it is in (previously only north was handled).
func TestAyanaBalaMercuryBothHemispheres(t *testing.T) {
	near(t, "Mercury N", ayanaBala("Mercury", 20), ayanaBala("Mercury", -20), 1e-9)
	near(t, "Mercury 20", ayanaBala("Mercury", -20), 55, 1e-9)
	near(t, "Sun on the equator", ayanaBala("Sun", 0), 60, 1e-9) // 30 doubled
	near(t, "Moon on the equator", ayanaBala("Moon", 0), 30, 1e-9)
	near(t, "limit", ayanaBala("Jupiter", 28), 60, 1e-9) // greatest kranti is 24
	near(t, "limit", ayanaBala("Jupiter", -28), 0, 1e-9)
}

// Raman example 18: Sun 180 54' (we use 180 53' 55"), Moon 311 17' 19", Mercury is with the Sun.
func TestPakshaBalaRamanExample(t *testing.T) {
	sun := 180.0 + 53.0/60 + 55.0/3600
	moon := 311.0 + 17.0/60 + 19.0/3600
	want := map[string]float64{"Sun": 16.54, "Moon": 86.92, "Mars": 16.54, "Mercury": 16.54, "Jupiter": 43.46, "Venus": 43.46, "Saturn": 16.54}
	for p, w := range want {
		near(t, p, pakshaBala(p, sun, moon, true), w, 0.01)
	}
	// Mercury not with a malefic takes the benefic value.
	near(t, "Mercury benefic", pakshaBala("Mercury", sun, moon, false), 43.46, 0.01)
	// The arc is reflected, so a Moon 130 deg behind the Sun gives the same as 130 deg ahead.
	near(t, "waning", pakshaBala("Jupiter", moon, sun, false), 43.46, 0.01)
	near(t, "full moon", pakshaBala("Moon", 0, 180, false), 120, 1e-9)
	near(t, "new moon", pakshaBala("Moon", 10, 10, false), 0, 1e-9)
}

func TestMercuryIsMalefic(t *testing.T) {
	base := func(m map[string]int) map[string]int {
		out := map[string]int{"Sun": 0, "Moon": 1, "Mars": 2, "Mercury": 5, "Jupiter": 6, "Venus": 7, "Saturn": 8, "Rahu": 9, "Ketu": 3}
		for k, v := range m {
			out[k] = v
		}
		return out
	}
	cases := []struct {
		name string
		with map[string]int
		want bool
	}{
		{"alone", nil, false},
		{"with Sun (Raman's chart)", map[string]int{"Sun": 5}, true},
		{"with Mars", map[string]int{"Mars": 5}, true},
		{"with Saturn", map[string]int{"Saturn": 5}, true},
		{"with Rahu", map[string]int{"Rahu": 5}, true},
		{"with Ketu", map[string]int{"Ketu": 5}, true},
		{"with Jupiter only", map[string]int{"Jupiter": 5}, false},
		{"with Venus only", map[string]int{"Venus": 5}, false},
		{"with the Moon only", map[string]int{"Moon": 5}, false},
	}
	for _, c := range cases {
		if got := mercuryIsMalefic(base(c.with)); got != c.want {
			t.Errorf("%s: got %v want %v", c.name, got, c.want)
		}
	}
}

// Raman art. 69: "the 25th hora, the first on the next day, is ruled by Chandra", i.e. after 24
// horas the lord is the next weekday's lord; and the sequence on Sunday is Ravi, Sukra, Budha, ...
func TestHoraSequence(t *testing.T) {
	for i, d := range WeekdayLords {
		if got, want := horaLord(d, 24.2), WeekdayLords[(i+1)%7]; got != want {
			t.Errorf("after 24 horas from %s: got %s want %s", d, got, want)
		}
		if horaLord(d, 0.5) != d {
			t.Errorf("first hora of %s must belong to %s", d, d)
		}
	}
	sunday := []string{"Sun", "Venus", "Mercury", "Moon", "Saturn", "Jupiter", "Mars", "Sun", "Venus", "Mercury", "Moon", "Saturn"}
	for n, w := range sunday {
		if got := horaLord("Sun", float64(n)+0.3); got != w {
			t.Errorf("Sunday hora %d: got %s want %s", n+1, got, w)
		}
	}
}

// Raman example 31: Wednesday, 8h 6m after sunrise, i.e. the 9th hora, is ruled by Chandra.
func TestHoraRamanExample(t *testing.T) {
	if got := horaLord("Mercury", 8.0+6.0/60); got != "Moon" {
		t.Errorf("9th hora on Wednesday: got %s want Moon", got)
	}
}

// End to end for Raman's standard horoscope: 16 Oct 1918 (Wednesday), 2:06:16 pm local mean time
// at 13 N, 77 35' E. Hora Bala goes to Chandra and Dina (Vara) Bala to Budha, as in his tables
// (examples 30, 31 and 35).
func TestHoraAndDinaRamanChart(t *testing.T) {
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	lon := 77.0 + 35.0/60
	ut := 14.0 + 6.0/60 + 16.0/3600 - lon/15
	jd := swisseph.Julday(1918, 10, 16, ut, swisseph.GregCal)
	ctx := domain.CalculationContext{
		Input:       domain.BirthInput{Latitude: 13, Longitude: lon, Timezone: 5.5},
		JulianDayUT: jd,
	}
	for _, c := range []struct {
		planet     string
		hora, dina float64
	}{{"Moon", 60, 0}, {"Mercury", 0, 45}, {"Sun", 0, 0}, {"Saturn", 0, 0}} {
		kb, err := CalculateKalaBala(c.planet, &ctx, 0, 0, 0, false)
		if err != nil {
			t.Fatal(err)
		}
		near(t, c.planet+" hora", kb.HoraBala, c.hora, 0)
		near(t, c.planet+" dina", kb.DinaBala, c.dina, 0)
	}
}

// A night birth before sunrise belongs to the previous Hindu day: 3:20 am on Wednesday 10 July 1985
// in Chennai is still Tuesday (Mars's day) for Dina Bala.
func TestDinaBeforeSunriseUsesPreviousDay(t *testing.T) {
	kb := kalaFor(t, "1985-07-10", "03:20:00", 13.0878, 80.2785, "Mars")
	near(t, "Mars dina", kb.DinaBala, 45, 0)
}
