package shadbala

import (
	"math"
	"testing"
)

// Raman, example 8 (p. 21) and 59 (p. 98-99): the middles of the houses of the standard horoscope, from the
// Ascendant 298 27' and the Midheaven 216 30' by Sripati's trisection. The book prints the 4th and 10th as
// 36 36' and 216 36' (6' above the lower and upper meridian), so the others are within 0.11 deg of his.
func TestSripatiBhavaMadhyas(t *testing.T) {
	m := sripatiBhavaMadhyas(dms(298, 27, 0), dms(216, 30, 0))
	book := []float64{dms(298, 27, 0), dms(331, 10, 0), dms(3, 53, 0), dms(36, 36, 0), dms(63, 53, 0), dms(91, 10, 0),
		dms(118, 27, 0), dms(151, 10, 0), dms(183, 53, 0), dms(216, 36, 0), dms(243, 53, 0), dms(271, 10, 0)}
	for i := range book {
		d := math.Abs(math.Mod(m[i]-book[i]+540, 360) - 180)
		if d > 0.11 {
			t.Errorf("house %d: middle %.3f, book %.3f", i+1, m[i], book[i])
		}
	}
	// kendras are exactly the angles; each house is a third of the arc between two kendras
	asc, mc := 123.0, 33.0
	m = sripatiBhavaMadhyas(asc, mc)
	near(t, "1st", m[0], 123, 1e-9)
	near(t, "4th", m[3], 213, 1e-9)
	near(t, "7th", m[6], 303, 1e-9)
	near(t, "10th", m[9], 33, 1e-9)
	near(t, "2nd", m[1], 123+30, 1e-9)
	near(t, "11th", m[10], math.Mod(33+30, 360), 1e-9)
}

// Art. 127-130.
func TestBhavaSignGroup(t *testing.T) {
	cases := []struct {
		lon  float64
		want string
	}{
		{1, "chathushpada"}, {31, "chathushpada"}, {61, "nara"}, {91, "jalachara"}, {121, "chathushpada"}, {151, "nara"},
		{181, "nara"}, {211, "keeta"}, {241, "nara"}, {254.9, "nara"}, {255, "chathushpada"}, {271, "chathushpada"}, {284.9, "chathushpada"},
		{285, "jalachara"}, {299, "jalachara"}, {311, "nara"}, {331, "jalachara"},
	}
	for _, c := range cases {
		if got := bhavaSignGroup(c.lon); got != c.want {
			t.Errorf("%.1f: got %s want %s", c.lon, got, c.want)
		}
	}
}

// The table on p. 93 ("Table showing the Digbala of Bhavas"): the strength of each of the twelve houses
// for each kind of sign.
func TestBhavaDigbalaTable(t *testing.T) {
	table := map[string][]float64{
		"nara":         {60, 50, 40, 30, 20, 10, 0, 10, 20, 30, 40, 50},
		"jalachara":    {30, 40, 50, 60, 50, 40, 30, 20, 10, 0, 10, 20},
		"chathushpada": {30, 20, 10, 0, 10, 20, 30, 40, 50, 60, 50, 40},
		"keeta":        {0, 10, 20, 30, 40, 50, 60, 50, 40, 30, 20, 10},
	}
	for group, vals := range table {
		for h, want := range vals {
			near(t, group, bhavaDigbala(h+1, group), want, 0)
		}
	}
}

// Example 59: Bhava Digbala of the twelve houses of the standard horoscope (I 30, II 40, III 10, IV nil, V 20,
// VI 40, VII 30, VIII 10, IX 20, X 30, XI 40, XII 40), from the middles printed in example 8.
func TestBhavaDigbalaRamanChart(t *testing.T) {
	madhyas := []float64{dms(298, 27, 0), dms(331, 10, 0), dms(3, 53, 0), dms(36, 36, 0), dms(63, 53, 0), dms(91, 10, 0),
		dms(118, 27, 0), dms(151, 10, 0), dms(183, 53, 0), dms(216, 36, 0), dms(243, 53, 0), dms(271, 10, 0)}
	want := []float64{30, 40, 10, 0, 20, 40, 30, 10, 20, 30, 40, 40}
	for i, m := range madhyas {
		near(t, "house", bhavaDigbala(i+1, bhavaSignGroup(m)), want[i], 0)
	}
	// the 8th house in a Nara sign, worked in art. 131: 7 - 8 = 11, folded to 1, times 10
	near(t, "art. 131", bhavaDigbala(8, "nara"), 10, 0)
}

// Example 59: Subhadrishtibala and Asubhadrishtibala on each house middle (pp. 98-99). The 4th house's benefic
// aspects are printed 62.89 but its parts (42.50 + 9.31 + 10.08) add to 61.89.
func TestBhavaDrishtiRamanChart(t *testing.T) {
	lons := map[string]float64{
		"Sun": dms(180, 54, 0), "Moon": dms(311, 17, 0), "Mars": dms(229, 31, 0), "Mercury": dms(181, 32, 0),
		"Jupiter": dms(84, 1, 0), "Venus": dms(171, 10, 0), "Saturn": dms(124, 23, 0),
	}
	signs := map[string]int{}
	for p, l := range lons {
		signs[p] = int(l / 30)
	}
	signs["Rahu"], signs["Ketu"] = 7, 1
	benefic := beneficMap(lons, signs)
	madhyas := []float64{dms(298, 27, 0), dms(331, 10, 0), dms(3, 53, 0), dms(36, 36, 0), dms(63, 53, 0), dms(91, 10, 0),
		dms(118, 27, 0), dms(151, 10, 0), dms(183, 53, 0), dms(216, 36, 0), dms(243, 53, 0), dms(271, 10, 0)}
	subha := []float64{79.99, 61.77, 85.05, 61.89, 43.22, 20.19, 12.39, 34.81, 48.49, 56.20, 45.18, 111.04}
	asubha := []float64{25.81, 25.33, 26.06, 34.79, 31.65, 17.27, 6.69, 2.30, 3.69, 11.69, 12.05, 13.49}
	for i, m := range madhyas {
		s, a := bhavaDrishti(m, lons, benefic)
		near(t, "subha "+bhavaNames[i], s, subha[i], 1.0)
		near(t, "asubha "+bhavaNames[i], a, asubha[i], 1.0)
	}
}

// Example 59 (footnote, p. 98): Mercury counts as a full benefic for the aspects on houses even when it is with
// the Sun, and Jupiter's and Mercury's aspects are taken whole while the others are taken by a quarter.
func TestBhavaDrishtiMercuryAndQuarters(t *testing.T) {
	signs := map[string]int{"Sun": 0, "Mercury": 0, "Mars": 3, "Jupiter": 6, "Venus": 1, "Saturn": 1, "Moon": 0, "Rahu": 3, "Ketu": 9}
	away := 111.0 // a house middle at 101 is 350 deg from a planet at 111: no aspect
	lons := map[string]float64{"Sun": away, "Moon": away, "Mars": away, "Mercury": away, "Jupiter": away, "Venus": away, "Saturn": away}
	benefic := beneficMap(map[string]float64{"Sun": 0, "Moon": 10, "Mars": 90, "Mercury": 1, "Jupiter": 200, "Venus": 40, "Saturn": 50}, signs)
	if benefic["Mercury"] {
		t.Fatal("Mercury is with the Sun here, so it is a malefic for the aspects planets cast")
	}
	// only Mercury aspects, from 100 deg: its whole value counts as a benefic
	lons["Mercury"] = 1
	s, a := bhavaDrishti(101, lons, benefic)
	near(t, "Mercury whole and benefic", s, drishtiValue("Mercury", 100), 1e-9)
	near(t, "nothing malefic", a, 0, 1e-9)
	// Jupiter whole (benefic), Venus and the Sun by a quarter (benefic and malefic)
	lons["Mercury"], lons["Jupiter"], lons["Venus"], lons["Sun"] = away, 1, 1, 1
	s, a = bhavaDrishti(101, lons, benefic)
	near(t, "Jupiter whole + Venus quarter", s, drishtiValue("Jupiter", 100)+drishtiValue("Venus", 100)/4, 1e-9)
	near(t, "Sun quarter", a, drishtiValue("Sun", 100)/4, 1e-9)
}

// Example 59: Bhava Bala of the standard horoscope with the Shadbala Pindas of example 56, and the closing
// remark that the 8th house is the strongest and the 3rd the weakest.
func TestBhavaBalaRamanChart(t *testing.T) {
	ctx := ramanChartContext(t, "raman")
	res, err := CalculateShadbala(&ctx)
	if err != nil {
		t.Fatal(err)
	}
	b := res.BhavaBala
	if len(b.Bhavas) != 12 {
		t.Fatalf("%d houses", len(b.Bhavas))
	}
	lords := []string{"Saturn", "Jupiter", "Mars", "Venus", "Mercury", "Moon", "Moon", "Mercury", "Venus", "Mars", "Jupiter", "Saturn"}
	dig := []float64{30, 40, 10, 0, 20, 40, 30, 10, 20, 30, 40, 40}
	// The book's printed totals of the 4th, 7th and 12th houses do not agree with their parts (404.15, 7.00 and
	// 426.76 / 7.09 rupas): from the parts they are 403.25 (6.72), 425.50 (7.09) and 526.76 (8.78).
	rupas := []float64{7.89, 8.51, 6.11, 6.72, 9.49, 7.21, 7.09, 9.67, 7.35, 6.21, 8.44, 8.78}
	for i, h := range b.Bhavas {
		if h.House != i+1 || h.Name != bhavaNames[i] || h.Lord != lords[i] {
			t.Errorf("house %d: %s lord %s, want %s", i+1, h.Name, h.Lord, lords[i])
		}
		near(t, h.Name+" Digbala", h.BhavaDigbala, dig[i], 0)
		near(t, h.Name+" total", h.TotalBhavaBala, h.BhavadhipathiBala+h.BhavaDigbala+h.BhavaDrishtiBala, 0.02)
		near(t, h.Name+" drishti", h.BhavaDrishtiBala, h.SubhaDrishtiBala-h.AsubhaDrishtiBala, 0.011)
		// our planets differ slightly from his printed longitudes and Ayana Bala uses the true declination
		near(t, h.Name+" rupas", h.Rupas, rupas[i], 0.4)
	}
	if b.StrongestHouse != 8 || b.WeakestHouse != 3 {
		t.Errorf("strongest house %d, weakest %d; Raman: 8th and 3rd", b.StrongestHouse, b.WeakestHouse)
	}
	for _, h := range b.Bhavas {
		if h.Rank < 1 || h.Rank > 12 {
			t.Errorf("rank %d", h.Rank)
		}
	}
	if b.Bhavas[7].Rank != 1 || b.Bhavas[2].Rank != 12 {
		t.Errorf("ranks of the 8th and 3rd: %d, %d", b.Bhavas[7].Rank, b.Bhavas[2].Rank)
	}
}
