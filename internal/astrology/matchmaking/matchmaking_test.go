package matchmaking

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func TestAshtakoota(t *testing.T) {
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	utcTime1, _ := astronomyTime.ParseLocalToUTC("2005-11-23", "15:35:00", 5.5)
	jd1 := astronomyTime.UTCToJulianDay(utcTime1)
	swisseph.SetEphePath(absPath)
	swisseph.SetSidMode(swisseph.SidmLahiri, 0, 0)

	ctx1 := domain.CalculationContext{
		Input:       domain.BirthInput{Name: "Groom", Ayanamsa: "Lahiri"},
		UTCTime:     utcTime1,
		JulianDayUT: jd1,
		Ayanamsa:    swisseph.GetAyanamsaUT(jd1),
	}

	utcTime2, _ := astronomyTime.ParseLocalToUTC("2002-05-14", "10:00:00", 5.5)
	jd2 := astronomyTime.UTCToJulianDay(utcTime2)
	ctx2 := domain.CalculationContext{
		Input:       domain.BirthInput{Name: "Bride", Ayanamsa: "Lahiri"},
		UTCTime:     utcTime2,
		JulianDayUT: jd2,
		Ayanamsa:    swisseph.GetAyanamsaUT(jd2),
	}

	res, err := CalculateMatch(&ctx1, &ctx2)
	if err != nil {
		t.Fatalf("Failed to calculate ashtakoota: %v", err)
	}

	if res.Ashtakoota.Summary.Maximum != 36 {
		t.Errorf("Expected max 36, got %f", res.Ashtakoota.Summary.Maximum)
	}

	if res.Ashtakoota.Kootas.Nadi.Maximum != 8 {
		t.Errorf("Expected Nadi max 8")
	}
}

// newTestContext builds a calculation context with real Moon positions (Lahiri, Placidus, IST).
func newTestContext(name, date, tm string, lat, lon float64) *domain.CalculationContext {
	in := domain.BirthInput{Name: name, DateOfBirth: date, TimeOfBirth: tm, Latitude: lat, Longitude: lon, Timezone: 5.5, Ayanamsa: "Lahiri", HouseSystem: "Placidus"}
	utc, _ := astronomyTime.ParseLocalToUTC(date, tm, 5.5)
	jd := astronomyTime.UTCToJulianDay(utc)
	cfg := domain.CalculationConfig{AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri"), HouseCode: ephemeris.GetHouseSystemCode("Placidus")}
	ephemeris.Lock()
	swisseph.SetSidMode(int32(cfg.AyanamsaMode), 0, 0)
	ay := swisseph.GetAyanamsaUT(jd)
	ephemeris.Unlock()
	return &domain.CalculationContext{Input: in, Config: cfg, UTCTime: utc, JulianDayUT: jd, Ayanamsa: ay}
}

// End to end with real Moon positions: groom 2005-11-23 15:35 (Moon Leo/Magha), bride 2002-05-14 10:00
// (Moon Taurus/Rohini). Expected values are derived by hand from the rules.
func TestMatchEndToEnd(t *testing.T) {
	absPath, _ := filepath.Abs("../../../ephe_data")
	if err := ephemeris.Init(absPath); err != nil {
		t.Fatalf("ephemeris init: %v", err)
	}
	g := newTestContext("Groom", "2005-11-23", "15:35:00", 16.3938, 80.1522)
	b := newTestContext("Bride", "2002-05-14", "10:00:00", 12.9716, 77.5946)

	res, err := CalculateMatch(g, b)
	if err != nil {
		t.Fatal(err)
	}
	if res.Groom.Moon.Nakshatra != "Magha" || res.Bride.Moon.Nakshatra != "Rohini" {
		t.Fatalf("unexpected moons: %s / %s", res.Groom.Moon.Nakshatra, res.Bride.Moon.Nakshatra)
	}
	if res.RuleSet.Ayanamsa != "Lahiri" {
		t.Errorf("ayanamsa: %s", res.RuleSet.Ayanamsa)
	}

	k := res.Ashtakoota.Kootas
	if k.Vashya.Score != 0 { // Leo and Taurus have no vashya relation in the book table
		t.Errorf("vashya: %+v", k.Vashya)
	}
	if k.Gana.Score != 0 || k.Gana.Suggestion != "" { // bride Manushya (Rohini), groom Rakshasa (Magha); count 7
		t.Errorf("gana: %+v", k.Gana)
	}
	// both Antya; Magha and Rohini are both in the 21-nakshatra list: the score stays 0, the rule is only suggested
	n := k.Nadi
	if n.RawScore != 0 || n.Score != 0 || !n.Dosha || n.CancellationApplied || n.Suggestion == "" {
		t.Errorf("nadi: %+v", n)
	}
	// Aries groom lagna is the 11th from the Gemini bride lagna => auspicious
	if res.Lagna.Distance != 11 || res.Lagna.Quality != "Auspicious" || !res.Lagna.IsAuspicious {
		t.Errorf("lagna: %+v", res.Lagna)
	}
	// varna 1 + vashya 0 + tara 1.5 + yoni 1 + maitri 0 + gana 0 + bhakoot 7 + nadi 0
	s := res.Ashtakoota.Summary
	if s.RawTotal != 10.5 || s.EffectiveTotal != 10.5 || s.Maximum != 36 || s.TraditionalThresholdMet {
		t.Errorf("summary: %+v", s)
	}

	// checks: Rohini (4th) -> Magha (10th) is count 7
	c := res.Checks
	if !c.Mahendra.Satisfied {
		t.Errorf("mahendra (count 7): %+v", c.Mahendra)
	}
	sd := c.StreeDeergha.Details.(domain.StreeDeerghaDetails)
	if c.StreeDeergha.Satisfied || sd.Navaka != 1 || sd.SatisfiedByAlternateOpinion || sd.ExceptionApplies {
		t.Errorf("stree deergha: %+v", c.StreeDeergha)
	}
	if !c.Vedha.Satisfied { // Magha pairs with Revati, Rohini with Swati
		t.Errorf("vedha: %+v", c.Vedha)
	}
	if !c.Rajju.Satisfied || c.Rajju.Groom != "Pada" || c.Rajju.Bride != "Kantha" {
		t.Errorf("rajju: %+v", c.Rajju)
	}
}

// The response must stay organised: fixed sections, and every koota and check has the same shape.
func TestResponseShape(t *testing.T) {
	absPath, _ := filepath.Abs("../../../ephe_data")
	if err := ephemeris.Init(absPath); err != nil {
		t.Fatalf("ephemeris init: %v", err)
	}
	res, err := CalculateMatch(
		newTestContext("Groom", "2005-11-23", "15:35:00", 16.3938, 80.1522),
		newTestContext("Bride", "2002-05-14", "10:00:00", 12.9716, 77.5946))
	if err != nil {
		t.Fatal(err)
	}

	raw, _ := json.Marshal(res)
	var top map[string]map[string]any
	if err := json.Unmarshal(raw, &top); err != nil {
		// groom/bride/lagna/... are all objects, so this must decode
		t.Fatal(err)
	}
	wantTop := []string{"rule_set", "groom", "bride", "lagna", "ashtakoota", "checks"}
	if len(top) != len(wantTop) {
		t.Errorf("top-level keys = %v, want exactly %v", keys(top), wantTop)
	}
	for _, k := range wantTop {
		if _, ok := top[k]; !ok {
			t.Errorf("missing top-level section %q", k)
		}
	}
	for _, gone := range []string{"kootas", "doshas", "summary", "book_rules", "lagna_match"} {
		if _, ok := top[gone]; ok {
			t.Errorf("old top-level key %q must not come back", gone)
		}
	}

	ashta := top["ashtakoota"]
	if len(ashta) != 2 || ashta["kootas"] == nil || ashta["summary"] == nil {
		t.Errorf("ashtakoota must hold exactly kootas and summary: %v", keys(map[string]map[string]any{"": ashta}))
	}
	kootas, _ := ashta["kootas"].(map[string]any)
	for _, name := range []string{"varna", "vashya", "tara", "yoni", "graha_maitri", "gana", "bhakoot", "nadi"} {
		obj, ok := kootas[name].(map[string]any)
		if !ok {
			t.Errorf("missing koota %q", name)
			continue
		}
		for _, f := range []string{"name", "maximum", "raw_score", "score", "groom", "bride", "dosha", "cancellation_applied", "explanation"} {
			if _, ok := obj[f]; !ok {
				t.Errorf("koota %s is missing field %q", name, f)
			}
		}
	}
	if len(kootas) != 8 {
		t.Errorf("expected 8 kootas, got %d", len(kootas))
	}

	checks := top["checks"]
	for _, name := range []string{"rajju", "vedha", "mahendra", "stree_deergha", "mruga_vairam", "special_tara", "kala_sarpa"} {
		obj, ok := checks[name].(map[string]any)
		if !ok {
			t.Errorf("missing check %q", name)
			continue
		}
		for _, f := range []string{"name", "satisfied", "explanation"} {
			if _, ok := obj[f]; !ok {
				t.Errorf("check %s is missing field %q", name, f)
			}
		}
	}
	if len(checks) != 7 {
		t.Errorf("expected 7 checks, got %d", len(checks))
	}
}

func keys(m map[string]map[string]any) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
}
