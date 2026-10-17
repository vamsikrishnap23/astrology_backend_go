package shadbala

import (
	"fmt"
	"testing"
)

func dms(d, m, s float64) float64 { return d + m/60 + s/3600 }

// Raman art. 110 and 114: the shape of the aspect curve.
func TestDrishtiBaseCurve(t *testing.T) {
	cases := []struct{ dk, want float64 }{
		{0, 0}, {29.9, 0}, {30, 0}, {45, 7.5}, {60, 15}, {75, 30}, {90, 45}, {105, 37.5}, {120, 30},
		{135, 15}, {150, 0}, {165, 30}, {180, 60}, {210, 45}, {240, 30}, {270, 15}, {299.9, 0.05}, {300, 0}, {330, 0}, {359.9, 0},
	}
	for _, c := range cases {
		near(t, "DK "+ftoa(c.dk), drishtiBase(c.dk), c.want, 0.06)
	}
	// continuity at the joints, except the jump between 150 and 180 that the book describes
	for _, dk := range []float64{30, 60, 90, 120, 150, 300} {
		near(t, "joint "+ftoa(dk), drishtiBase(dk-1e-9), drishtiBase(dk), 1e-6)
	}
	// 360 deg wraps
	near(t, "wrap", drishtiBase(-180), 60, 1e-9)
}

func ftoa(f float64) string { return fmt.Sprintf("%g", f) }

// Raman art. 115: Saturn adds 45 on the 3rd and 10th, Jupiter 30 on the 5th and 9th, Mars 15 on the
// 4th and 8th; nobody else has a special aspect.
func TestVisheshaDrishti(t *testing.T) {
	cases := []struct {
		planet string
		dk     float64
		want   float64
	}{
		{"Saturn", 75, 45}, {"Saturn", 285, 45}, {"Saturn", 100, 0}, {"Saturn", 59.9, 0}, {"Saturn", 90, 0},
		{"Jupiter", 135, 30}, {"Jupiter", 255, 30}, {"Jupiter", 100, 0}, {"Jupiter", 150, 0},
		{"Mars", 105, 15}, {"Mars", 225, 15}, {"Mars", 135, 0}, {"Mars", 240, 0},
		{"Sun", 105, 0}, {"Venus", 135, 0}, {"Mercury", 75, 0}, {"Moon", 225, 0},
	}
	for _, c := range cases {
		near(t, c.planet+" "+ftoa(c.dk), visheshaDrishti(c.planet, c.dk), c.want, 0)
	}
	// "Kuja ... *15.00 + 42.75" on Jupiter in example 54: ordinary 42.75 plus the special 15.
	near(t, "Mars on Jupiter", drishtiValue("Mars", dms(214, 30, 0)), 42.75+15, 0.01)
	// "Guru ... *30.00 + 4.50" on Mars.
	near(t, "Jupiter on Mars", drishtiValue("Jupiter", dms(145, 30, 0)), 30+4.5, 0.01)
}

// Raman example 54, aspect by aspect (Drishti Kendra from example 53; the angle printed as 145 20 for
// Jupiter on Mars is a misprint for 145 30, which his own longitudes and his 4.50 require).
func TestDrishtiValuesRamanExample54(t *testing.T) {
	cases := []struct {
		from, to string
		dk       float64
		want     float64 // value printed in example 54 (without any special-aspect addition)
	}{
		{"Jupiter", "Sun", dms(96, 53, 0), 41.50}, {"Jupiter", "Moon", dms(227, 16, 0), 36.38},
		{"Jupiter", "Mars", dms(145, 30, 0), 4.50}, {"Jupiter", "Mercury", dms(97, 31, 0), 41.25},
		{"Jupiter", "Venus", dms(87, 9, 0), 42.20}, {"Jupiter", "Saturn", dms(40, 22, 0), 5.18},
		{"Moon", "Sun", dms(229, 37, 0), 35.20}, {"Moon", "Mars", dms(278, 14, 0), 10.88},
		{"Moon", "Mercury", dms(230, 15, 0), 34.88}, {"Moon", "Jupiter", dms(132, 44, 0), 17.25},
		{"Moon", "Venus", dms(219, 53, 0), 40.06}, {"Moon", "Saturn", dms(173, 6, 0), 46.20},
		{"Venus", "Moon", dms(140, 7, 0), 9.87}, {"Venus", "Mars", dms(58, 21, 0), 14.17}, {"Venus", "Jupiter", dms(272, 51, 0), 13.54},
		{"Sun", "Moon", dms(130, 23, 0), 19.62}, {"Sun", "Mars", dms(48, 37, 0), 9.31}, {"Sun", "Jupiter", dms(263, 7, 0), 18.43},
		{"Mars", "Moon", dms(81, 46, 0), 36.75}, {"Mars", "Saturn", dms(254, 52, 0), 22.56},
		{"Saturn", "Sun", dms(56, 31, 0), 13.25}, {"Saturn", "Moon", dms(186, 54, 0), 56.56}, {"Saturn", "Mars", dms(105, 8, 0), 37.43},
		{"Saturn", "Mercury", dms(57, 9, 0), 13.56}, {"Saturn", "Venus", dms(46, 47, 0), 8.38},
		{"Mercury", "Moon", dms(129, 45, 0), 20.25}, {"Mercury", "Mars", dms(47, 59, 0), 9.00}, {"Mercury", "Jupiter", dms(262, 29, 0), 18.75},
	}
	for _, c := range cases {
		// base value only: compare without the special aspect, which the book prints separately
		near(t, c.from+" on "+c.to, drishtiBase(c.dk), c.want, 0.06)
	}
	// no aspect on the Sun from Mercury: 359 22' is beyond 300
	near(t, "Mercury on Sun", drishtiBase(dms(359, 22, 0)), 0, 0)
}

// Raman's standard horoscope (nirayana longitudes of art. 6): Drishti Pindas of example 54 and Drik
// Balas of example 55 (the printed Pinda of the Moon, -86.23 on p. 87, is a misprint for -86.93, as
// his Drik Bala -21.73 and his own columns show).
func TestDrikBalaRamanExample(t *testing.T) {
	lons := map[string]float64{
		"Sun": dms(180, 53, 55), "Moon": dms(311, 17, 19), "Mars": dms(229, 30, 34), "Mercury": dms(181, 31, 34),
		"Jupiter": dms(84, 0, 49), "Venus": dms(171, 9, 56), "Saturn": dms(124, 22, 41),
	}
	signs := map[string]int{}
	for p, l := range lons {
		signs[p] = int(l / 30)
	}
	signs["Rahu"], signs["Ketu"] = int(dms(234, 23, 47)/30), int(dms(54, 23, 47)/30)
	benefic := beneficMap(lons, signs)
	if benefic["Mercury"] || !benefic["Moon"] || !benefic["Jupiter"] || !benefic["Venus"] || benefic["Sun"] || benefic["Mars"] || benefic["Saturn"] {
		t.Fatalf("classification differs from the book's: %v", benefic)
	}
	pinda := map[string]float64{"Sun": 63.45, "Moon": -86.93, "Mars": 3.81, "Mercury": 62.57, "Jupiter": -64.14, "Venus": 73.88, "Saturn": 28.82}
	drik := map[string]float64{"Sun": 15.86, "Moon": -21.73, "Mars": 0.95, "Mercury": 15.64, "Jupiter": -16.04, "Venus": 18.47, "Saturn": 7.21}
	for _, p := range TargetPlanets {
		near(t, p+" Drishti Pinda", drishtiPinda(p, lons, benefic), pinda[p], 0.2)
		near(t, p+" Drik Bala", CalculateDrikBala(p, lons, benefic), drik[p], 0.05)
	}
}

// Raman art. 53: the Moon is a benefic from the 8th day of the bright half to the 8th day of the dark
// half (84 to 264 deg elongation) and a malefic otherwise.
func TestMoonIsBenefic(t *testing.T) {
	cases := []struct {
		elong float64
		want  bool
	}{{0, false}, {40, false}, {83.9, false}, {84, true}, {130, true}, {180, true}, {263.9, true}, {264, false}, {300, false}, {359, false}}
	for _, c := range cases {
		if got := moonIsBenefic(100, 100+c.elong); got != c.want {
			t.Errorf("elongation %.1f: got %v want %v", c.elong, got, c.want)
		}
	}
	// works across the 0/360 line
	if !moonIsBenefic(350, 350+130-360+360) {
		t.Error("elongation 130 across 360 should be benefic")
	}
}
