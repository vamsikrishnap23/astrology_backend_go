package matchmaking

import (
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// The book's worked example: groom Mrigashira (Taurus), bride Dhanishta. The book states the count is 10,
// stree deergha is satisfied (beyond 9) and Mrigashira and Dhanishta are a vedha pair. 10 is also one of
// the mahendra counts.
func TestChecksBookExample(t *testing.T) {
	groom := moon("Mrigashira", "Taurus", 1)
	bride := moon("Dhanishta", "Capricorn", 1)
	count := nakshatraCount(bride.Nakshatra, groom.Nakshatra)
	if count != 10 {
		t.Fatalf("count = %d, want 10", count)
	}
	if m := mahendraCheck(groom, bride, count); !m.Satisfied {
		t.Errorf("mahendra should hold for count 10: %+v", m)
	}
	s := streeDeerghaCheck(groom, bride, count, true, 3)
	if d := s.Details.(domain.StreeDeerghaDetails); !s.Satisfied || d.Navaka != 2 {
		t.Errorf("stree deergha: %+v", s)
	}
	if v := vedhaCheck(groom, bride); v.Satisfied {
		t.Errorf("vedha: %+v", v)
	}
}

func TestNakshatraCountMatchesBookExamples(t *testing.T) {
	cases := []struct {
		from, to string
		want     int
	}{
		{"Ashwini", "Jyeshtha", 18}, // book: 18th from Ashwini
		{"Bharani", "Anuradha", 16}, // book: 16th from Bharani
		{"Dhanishta", "Mrigashira", 10},
		{"Ashwini", "Ashwini", 1},
		{"Revati", "Ashwini", 2},
	}
	for _, c := range cases {
		if got := nakshatraCount(c.from, c.to); got != c.want {
			t.Errorf("count %s->%s = %d, want %d", c.from, c.to, got, c.want)
		}
	}
}

func TestMahendraCheck(t *testing.T) {
	good := map[int]bool{4: true, 7: true, 10: true, 13: true, 16: true, 19: true, 22: true, 25: true}
	g, b := moon("Ashwini", "Aries", 1), moon("Bharani", "Aries", 1)
	for count := 1; count <= 27; count++ {
		c := mahendraCheck(g, b, count)
		if c.Satisfied != good[count] {
			t.Errorf("count %d: got %v want %v", count, c.Satisfied, good[count])
		}
		if c.Details.(domain.CountDetails).BrideToGroomCount != count {
			t.Errorf("count %d: details %+v", count, c.Details)
		}
	}
}

func TestStreeDeerghaCheck(t *testing.T) {
	g, b := moon("Ashwini", "Aries", 1), moon("Bharani", "Aries", 1)
	for count := 1; count <= 27; count++ {
		c := streeDeerghaCheck(g, b, count, false, 0)
		d := c.Details.(domain.StreeDeerghaDetails)
		if d.Navaka != (count-1)/9+1 {
			t.Errorf("count %d: navaka %d", count, d.Navaka)
		}
		if c.Satisfied != (count >= 10) {
			t.Errorf("count %d: satisfied=%v", count, c.Satisfied)
		}
		if d.SatisfiedByAlternateOpinion != (count >= 8) {
			t.Errorf("count %d: alternate opinion=%v", count, d.SatisfiedByAlternateOpinion)
		}
	}
	// exception: only when it failed AND bhakoot is good AND graha maitri is good
	exception := func(count int, bhakootGood bool, maitri float64) bool {
		return streeDeerghaCheck(g, b, count, bhakootGood, maitri).Details.(domain.StreeDeerghaDetails).ExceptionApplies
	}
	if !exception(5, true, 4) {
		t.Error("exception should apply (bhakoot good, maitri 4)")
	}
	if exception(5, false, 5) || exception(5, true, 3) {
		t.Error("exception must need both bhakoot and graha maitri good")
	}
	if exception(12, true, 5) {
		t.Error("exception is only for a failed stree deergha")
	}
}

func TestVedhaPairs(t *testing.T) {
	if len(vedhaPairs) != 13 {
		t.Fatalf("expected 13 pairs, got %d", len(vedhaPairs))
	}
	seen := map[string]bool{}
	for _, p := range vedhaPairs {
		for _, n := range p {
			if seen[n] {
				t.Errorf("%s appears in more than one pair", n)
			}
			seen[n] = true
			if getNakshatraIndex(n) == 0 && n != "Ashwini" {
				t.Errorf("%s is not a known nakshatra name", n)
			}
		}
		if vedhaCheck(moon(p[0], "Aries", 1), moon(p[1], "Aries", 1)).Satisfied || vedhaCheck(moon(p[1], "Aries", 1), moon(p[0], "Aries", 1)).Satisfied {
			t.Errorf("pair %v must be symmetric", p)
		}
	}
	if len(seen) != 26 || seen["Chitra"] {
		t.Errorf("26 distinct nakshatras expected, Chitra excluded; got %d", len(seen))
	}
	if !vedhaCheck(moon("Ashwini", "Aries", 1), moon("Bharani", "Aries", 1)).Satisfied || !vedhaCheck(moon("Chitra", "Virgo", 1), moon("Chitra", "Virgo", 1)).Satisfied {
		t.Error("unrelated nakshatras must not have vedha")
	}
}

func TestRajjuCheck(t *testing.T) {
	// Ashwini and Ashlesha are both Pada rajju: a dosha. Mars (Aries) and Moon (Cancer) lords are friend/neutral.
	g, b := moon("Ashwini", "Aries", 1), moon("Ashlesha", "Cancer", 1)
	c := rajjuCheck(g, b, 3, "Inimical")
	if c.Satisfied || c.CancellationApplied || c.Groom != "Pada" || c.Bride != "Pada" {
		t.Errorf("uncancelled rajju dosha: %+v", c)
	}
	// cancelled by Graha Maitri of 4 or more
	if c := rajjuCheck(g, b, 4, "Inimical"); !c.Satisfied || !c.CancellationApplied || c.CancellationReason == "" {
		t.Errorf("cancel by maitri: %+v", c)
	}
	// cancelled by friendly lagna lords
	if c := rajjuCheck(g, b, 3, "Friendly"); !c.Satisfied || !c.CancellationApplied {
		t.Errorf("cancel by lagna lords: %+v", c)
	}
	// different rajju groups: satisfied, nothing cancelled
	c = rajjuCheck(moon("Ashwini", "Aries", 1), moon("Rohini", "Taurus", 1), 0, "Inimical")
	if !c.Satisfied || c.CancellationApplied || c.Groom == c.Bride {
		t.Errorf("different rajju: %+v", c)
	}
}

func TestMrugaVairamAndSpecialTara(t *testing.T) {
	// Ashwini (Horse) and Hasta (Buffalo) are natural enemies
	if c := mrugaVairamCheck(moon("Ashwini", "Aries", 1), moon("Hasta", "Virgo", 1)); c.Satisfied || c.Groom != "Horse" || c.Bride != "Buffalo" {
		t.Errorf("mruga vairam: %+v", c)
	}
	if c := mrugaVairamCheck(moon("Ashwini", "Aries", 1), moon("Bharani", "Aries", 1)); !c.Satisfied {
		t.Errorf("no mruga vairam: %+v", c)
	}
	// bride Krittika with groom Bharani is one of the three special rejections; bride Krittika/groom Ashlesha is in the 15-pair list
	if c := specialTaraCheck(moon("Bharani", "Aries", 1), moon("Krittika", "Taurus", 1)); c.Satisfied {
		t.Errorf("special tara: %+v", c)
	}
	if c := specialTaraCheck(moon("Ashlesha", "Cancer", 1), moon("Krittika", "Taurus", 1)); c.Satisfied {
		t.Errorf("ashubha kurpu: %+v", c)
	}
	if c := specialTaraCheck(moon("Ashwini", "Aries", 1), moon("Krittika", "Taurus", 1)); !c.Satisfied {
		t.Errorf("no special tara: %+v", c)
	}
}
