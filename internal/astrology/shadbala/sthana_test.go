package shadbala

import (
	"testing"
)

// Raman art. 36-39 and example 12 (his standard horoscope): male planets (Sun, Mars, Jupiter) in the
// first drekkana, hermaphrodite (Mercury, Saturn) in the second, female (Moon, Venus) in the last.
func TestDrekkanaBalaRamanExample12(t *testing.T) {
	cases := []struct {
		planet string
		lon    float64
		want   float64
	}{
		{"Sun", dms(180, 53, 55), 15}, {"Moon", dms(311, 17, 19), 0}, {"Mars", dms(229, 30, 34), 0},
		{"Mercury", dms(181, 31, 34), 0}, {"Jupiter", dms(84, 0, 49), 0}, {"Venus", dms(171, 9, 56), 15},
		{"Saturn", dms(124, 22, 41), 0},
	}
	for _, c := range cases {
		near(t, c.planet, calculateDrekkanaBala(c.planet, c.lon), c.want, 0)
	}
}

func TestDrekkanaBalaRules(t *testing.T) {
	want := map[string]int{"Sun": 1, "Mars": 1, "Jupiter": 1, "Mercury": 2, "Saturn": 2, "Moon": 3, "Venus": 3}
	for p, d := range want {
		for k := 1; k <= 3; k++ {
			lon := 60.0 + float64(k-1)*10 + 5 // inside the k-th drekkana of Gemini
			exp := 0.0
			if k == d {
				exp = 15
			}
			near(t, p, calculateDrekkanaBala(p, lon), exp, 0)
		}
	}
	// the joints between drekkanas
	near(t, "Sun 9.99", calculateDrekkanaBala("Sun", 9.99), 15, 0)
	near(t, "Sun 10", calculateDrekkanaBala("Sun", 10), 0, 0)
	near(t, "Saturn 10", calculateDrekkanaBala("Saturn", 10), 15, 0)
	near(t, "Moon 20", calculateDrekkanaBala("Moon", 20), 15, 0)
}

// Raman art. 30: the Moolatrikona sign of each planet (the whole sign; BPHS's degree limits are only noted in
// relationships.go).
func TestIsMoolatrikona(t *testing.T) {
	cases := []struct {
		planet, sign string
		want         bool
	}{
		{"Sun", "Leo", true}, {"Sun", "Aries", false}, {"Moon", "Taurus", true}, {"Moon", "Cancer", false},
		{"Mars", "Aries", true}, {"Mars", "Scorpio", false}, {"Mercury", "Virgo", true}, {"Mercury", "Gemini", false},
		{"Jupiter", "Sagittarius", true}, {"Jupiter", "Pisces", false}, {"Venus", "Libra", true}, {"Venus", "Taurus", false},
		{"Saturn", "Aquarius", true}, {"Saturn", "Capricorn", false},
	}
	for _, c := range cases {
		if got := isMoolatrikona(c.planet, c.sign); got != c.want {
			t.Errorf("%s in %s: got %v want %v", c.planet, c.sign, got, c.want)
		}
	}
}

// Raman art. 30: Moolatrikona counts (45) only in the Rasi; in the other vargas the same sign is the
// planet's own sign (30) -- his example is the Sun in Simha, 45 in the Rasi and not in the Navamsa.
func TestMoolatrikonaOnlyInRasi(t *testing.T) {
	signs := map[string]int{"Sun": 4, "Moon": 1, "Venus": 0, "Mars": 0, "Jupiter": 0, "Mercury": 0, "Saturn": 0}
	if got := GetPlacementStrength("Sun", "Leo", signs, true); got != Moolatrikona {
		t.Errorf("Sun in its Moolatrikona sign Leo in the Rasi: got %v", got)
	}
	if got := GetPlacementStrength("Sun", "Leo", signs, false); got != OwnHouse {
		t.Errorf("Sun in Leo in another varga: got %v, want OwnHouse", got)
	}
	// Taurus is the Moon's Moolatrikona sign but not its own: elsewhere it is judged by the lord, Venus.
	if got := GetPlacementStrength("Moon", "Taurus", signs, false); got == Moolatrikona || got == OwnHouse {
		t.Errorf("Moon in Taurus in another varga must be judged by its lord Venus, got %v", got)
	}
}

func TestMinimumRequirementsRaman(t *testing.T) {
	want := map[string]float64{"Sun": 300, "Moon": 360, "Mars": 300, "Mercury": 420, "Jupiter": 390, "Venus": 330, "Saturn": 300}
	for p, v := range want {
		near(t, p, MinimumRequirements[p], v/60, 1e-9)
	}
}

// Raman examples 9, 12 and 13, from the whole engine with the Raman ayanamsa on his standard horoscope:
// Saptavargaja Bala 90, 48.75, 90, 135, 71.25, 116.25, 97.5 (his Moon is printed 49.750 in example 9 but
// its seven parts add to 48.75, as in example 13; Mercury and Venus are printed 13.500 and 11.625,
// decimal-point slips for 135.0 and 116.25), and Sthana Bala totals of example 13.
func TestSthanaBalaRamanChart(t *testing.T) {
	ctx := ramanChartContext(t, "raman")
	res, err := CalculateShadbala(&ctx)
	if err != nil {
		t.Fatal(err)
	}
	sapta := map[string]float64{"Sun": 90, "Moon": 48.75, "Mars": 90, "Mercury": 135, "Jupiter": 71.25, "Venus": 116.25, "Saturn": 97.5}
	drek := map[string]float64{"Sun": 15, "Moon": 0, "Mars": 0, "Mercury": 0, "Jupiter": 0, "Venus": 15, "Saturn": 0}
	sthana := map[string]float64{"Sun": 198.0, "Moon": 126.5, "Mars": 172.06, "Mercury": 279.5, "Jupiter": 157.58, "Venus": 178.2, "Saturn": 177.3}
	for _, p := range res.Planets {
		near(t, p.Planet+" Saptavargaja", p.SthanaBala.SaptavargajaBala, sapta[p.Planet], 0.01)
		near(t, p.Planet+" Drekkana", p.SthanaBala.DrekkanaBala, drek[p.Planet], 0)
		near(t, p.Planet+" Sthana total", p.SthanaBala.Total, sthana[p.Planet], 0.15)
	}
}
