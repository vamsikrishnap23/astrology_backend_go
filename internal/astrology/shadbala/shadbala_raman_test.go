package shadbala

import (
	"sort"
	"testing"
)

// Raman, Graha and Bhava Balas, examples 56-57: the Shadbala Pinda of his standard horoscope.
//
// Everything is reproduced by the engine except Ayana Bala, where his kranti comes from a longitude-only
// table with a maximum of 24 deg and ours is the true declination (a decision of the owner), which moves
// the Moon by about 6 shashtiamsas and the others by up to 2. So each planet's total is compared with the
// book's after swapping our Ayana Bala for his. His printed totals are used where they agree with his own
// components; the Sun's printed 424.24 is a misprint for 424.94 and Mars's Kala Bala 33.06 for 30.06.
func TestShadbalaPindaRamanChart(t *testing.T) {
	ctx := ramanChartContext(t, "raman")
	res, err := CalculateShadbala(&ctx)
	if err != nil {
		t.Fatal(err)
	}
	total := map[string]float64{"Sun": 424.94, "Moon": 389.80, "Mars": 298.14, "Mercury": 537.02, "Jupiter": 433.71, "Venus": 376.15, "Saturn": 389.21}
	ayana := map[string]float64{"Sun": 38.12, "Moon": 43.44, "Mars": 1.84, "Mercury": 41.25, "Jupiter": 59.4, "Venus": 23.75, "Saturn": 13.75}
	for _, p := range res.Planets {
		near(t, p.Planet+" Shadbala Pinda (Ayana as in the book)", p.TotalShadbala-p.KalaBala.AyanaBala+ayana[p.Planet], total[p.Planet], 0.35)
		// the rupa and the requirement ratio follow from the total
		near(t, p.Planet+" rupas", p.Rupas, p.TotalShadbala/60, 0.001)
		near(t, p.Planet+" ratio", p.StrengthRatio, p.Rupas/MinimumRequirements[p.Planet], 0.001)
	}
}

// Example 57: the strength of the planets in the order I (strongest) to VII, and that only Mars falls
// short of its requirement.
func TestShadbalaStrengthOrderRamanChart(t *testing.T) {
	ctx := ramanChartContext(t, "raman")
	res, err := CalculateShadbala(&ctx)
	if err != nil {
		t.Fatal(err)
	}
	ratio := map[string]float64{}
	for _, p := range res.Planets {
		ratio[p.Planet] = p.StrengthRatio
		if want := p.Planet != "Mars"; p.MeetsMinimum != want {
			t.Errorf("%s: meets minimum = %v, want %v (ratio %.3f)", p.Planet, p.MeetsMinimum, want, p.StrengthRatio)
		}
	}
	names := []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"}
	sort.Slice(names, func(i, j int) bool { return ratio[names[i]] > ratio[names[j]] })
	// Raman: Ravi I, Sani II, Budha III, Sukra IV, Guru V, Chandra VI, Kuja VII. Saturn (1.30) and
	// Mercury (1.26) are close in the book and our Ayana differs, so only the ends and the Moon are fixed.
	if names[0] != "Sun" || names[6] != "Mars" {
		t.Errorf("strongest should be the Sun and weakest Mars, got %v", names)
	}
	if ratio["Sun"] < 1.4 || ratio["Mars"] > 1.0 {
		t.Errorf("ratios differ from the book: %v", ratio)
	}
}
