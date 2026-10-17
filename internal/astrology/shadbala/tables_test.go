package shadbala

import (
	"encoding/json"
	"math"
	"sort"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func tablesFor(t *testing.T, ayanamsa string) domain.ShadbalaResult {
	t.Helper()
	ctx := ramanChartContext(t, ayanamsa)
	res, err := CalculateShadbala(&ctx)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

// Every table has the seven planets in every row and in the total, the parts add up to the total, and the
// final table agrees with the per-planet results.
func TestShadbalaTablesAreConsistent(t *testing.T) {
	res := tablesFor(t, "lahiri")
	if res.RuleSet != RuleSet || res.Unit == "" {
		t.Errorf("rule set %q unit %q", res.RuleSet, res.Unit)
	}
	tb := res.Tables
	if len(tb.Planets) != 7 || tb.Planets[0] != "Sun" || tb.Planets[6] != "Saturn" {
		t.Fatalf("planets: %v", tb.Planets)
	}
	for name, table := range map[string]domain.ShadbalaTable{
		"sthana": tb.SthanaBala, "dig": tb.DigBala, "kala": tb.KalaBala, "cheshta": tb.CheshtaBala, "naisargika": tb.NaisargikaBala, "drik": tb.DrikBala,
	} {
		if table.Title == "" || len(table.Rows) == 0 {
			t.Errorf("%s: empty table", name)
		}
		for _, p := range tb.Planets {
			sum := 0.0
			for _, r := range table.Rows {
				v, ok := r.Values[p]
				if !ok || r.Label == "" || r.Component == "" {
					t.Errorf("%s %s: missing value or label in row %q", name, p, r.Component)
				}
				sum += v
			}
			tot, ok := table.Total[p]
			if !ok {
				t.Errorf("%s: no total for %s", name, p)
			}
			near(t, name+" "+p+" total", tot, sum, 0.05) // each part is rounded to two decimals
		}
	}
	if len(tb.SthanaBala.Rows) != 5 || len(tb.KalaBala.Rows) != 9 {
		t.Errorf("Sthana has %d parts and Kala %d, want 5 and 9", len(tb.SthanaBala.Rows), len(tb.KalaBala.Rows))
	}
	pd := tb.ShadbalaPinda
	if len(pd.Rows) != 6 {
		t.Fatalf("pinda table has %d rows", len(pd.Rows))
	}
	ranks := []int{}
	for _, p := range res.Planets {
		sum := 0.0
		for _, r := range pd.Rows {
			sum += r.Values[p.Planet]
		}
		near(t, p.Planet+" pinda sum", sum, p.TotalShadbala, 0.05)
		near(t, p.Planet+" pinda total", pd.Total[p.Planet], p.TotalShadbala, 0.006)
		near(t, p.Planet+" rupas", pd.Rupas[p.Planet], p.Rupas, 0)
		near(t, p.Planet+" minimum", pd.MinimumRequiredRupas[p.Planet], p.MinimumRequired, 0)
		near(t, p.Planet+" ratio", pd.StrengthRatio[p.Planet], p.StrengthRatio, 0)
		if pd.Rank[p.Planet] != p.Rank {
			t.Errorf("%s: table rank %d, planet rank %d", p.Planet, pd.Rank[p.Planet], p.Rank)
		}
		ranks = append(ranks, p.Rank)
	}
	sort.Ints(ranks)
	for i, r := range ranks {
		if r != i+1 {
			t.Fatalf("ranks must be 1..7 once each, got %v", ranks)
		}
	}
	// rank 1 has the highest ratio
	for _, a := range res.Planets {
		for _, b := range res.Planets {
			if a.Rank < b.Rank && a.StrengthRatio < b.StrengthRatio {
				t.Errorf("%s (rank %d, %.3f) ranked above %s (rank %d, %.3f)", a.Planet, a.Rank, a.StrengthRatio, b.Planet, b.Rank, b.StrengthRatio)
			}
		}
	}
	if _, err := json.Marshal(res); err != nil {
		t.Fatal(err)
	}
}

// Raman's printed tables for his standard horoscope (Raman ayanamsa): example 13 (Sthana Bala), example 35
// (Kala Bala parts, where he prints them) and example 57 (rank).
func TestShadbalaTablesRamanChart(t *testing.T) {
	res := tablesFor(t, "raman")
	tb := res.Tables
	sthana := map[string]float64{"Sun": 198.0, "Moon": 126.5, "Mars": 172.06, "Mercury": 279.5, "Jupiter": 157.58, "Venus": 178.2, "Saturn": 177.3}
	for p, w := range sthana {
		near(t, p+" Sthana total", tb.SthanaBala.Total[p], w, 0.15)
	}
	row := func(table domain.ShadbalaTable, comp string) map[string]float64 {
		for _, r := range table.Rows {
			if r.Component == comp {
				return r.Values
			}
		}
		t.Fatalf("no row %s", comp)
		return nil
	}
	for p, w := range map[string]float64{"Sun": 90, "Moon": 48.75, "Mars": 90, "Mercury": 135, "Jupiter": 71.25, "Venus": 116.25, "Saturn": 97.5} {
		near(t, p+" Saptavargaja", row(tb.SthanaBala, "saptavargaja_bala")[p], w, 0.01)
	}
	for p, w := range map[string]float64{"Sun": 48.32, "Moon": 11.68, "Mars": 11.68, "Mercury": 60, "Jupiter": 48.32, "Venus": 48.32, "Saturn": 11.68} {
		near(t, p+" Nathonnatha", row(tb.KalaBala, "nathonnatha_bala")[p], w, 0.05)
	}
	near(t, "Moon Hora", row(tb.KalaBala, "hora_bala")["Moon"], 60, 0)
	near(t, "Saturn Abda", row(tb.KalaBala, "varsha_bala")["Saturn"], 15, 0)
	near(t, "Mercury Masa", row(tb.KalaBala, "masa_bala")["Mercury"], 30, 0)
	near(t, "Mercury Vara", row(tb.KalaBala, "dina_bala")["Mercury"], 45, 0)
	for p, w := range map[string]float64{"Mars": 22.23, "Mercury": 2.30, "Jupiter": 35.26, "Venus": 5.95, "Saturn": 21.14} {
		near(t, p+" Cheshta", tb.CheshtaBala.Total[p], w, 0.05)
	}
	for p, w := range map[string]float64{"Sun": 15.86, "Moon": -21.73, "Mars": 0.95, "Mercury": 15.64, "Jupiter": -16.04, "Venus": 18.47, "Saturn": 7.21} {
		near(t, p+" Drik", tb.DrikBala.Total[p], w, 0.06)
	}
	for p, w := range map[string]float64{"Sun": 60, "Moon": 51.43, "Venus": 42.85, "Jupiter": 34.28, "Mercury": 25.71, "Mars": 17.14, "Saturn": 8.57} {
		near(t, p+" Naisargika", tb.NaisargikaBala.Total[p], w, 0.01)
	}
	// Example 57: the Sun is strongest, Mars weakest and the only planet below its requirement.
	pd := tb.ShadbalaPinda
	if pd.Rank["Sun"] != 1 || pd.Rank["Mars"] != 7 {
		t.Errorf("ranks: %v", pd.Rank)
	}
	for p, ratio := range pd.StrengthRatio {
		if (ratio >= 1) == (p == "Mars") {
			t.Errorf("%s: ratio %.3f", p, ratio)
		}
	}
	for _, p := range res.Planets {
		if math.IsNaN(p.Rupas) {
			t.Errorf("%s: NaN", p.Planet)
		}
	}
}
