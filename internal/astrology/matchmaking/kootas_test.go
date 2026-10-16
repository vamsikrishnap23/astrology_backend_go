package matchmaking

import (
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func moon(nak, sign string, pada int) domain.MoonDetails {
	return domain.MoonDetails{Nakshatra: nak, Sign: sign, Pada: pada, RashiLord: SignLords[sign]}
}

func TestLagnaQuality(t *testing.T) {
	want := map[int]string{
		1: "Auspicious", 3: "Auspicious", 5: "Auspicious", 9: "Auspicious", 11: "Auspicious",
		2: "Moderate", 7: "Moderate",
		4: "Inauspicious", 6: "Inauspicious", 8: "Inauspicious", 10: "Inauspicious", 12: "Inauspicious",
	}
	for d := 1; d <= 12; d++ {
		if got := lagnaQuality(d); got != want[d] {
			t.Errorf("distance %d: %s want %s", d, got, want[d])
		}
	}
}

func TestVashyaKootaBookTable(t *testing.T) {
	yes := [][2]string{
		{"Aries", "Leo"}, {"Aries", "Virgo"}, {"Taurus", "Cancer"}, {"Taurus", "Libra"}, {"Gemini", "Virgo"},
		{"Cancer", "Scorpio"}, {"Cancer", "Sagittarius"}, {"Leo", "Libra"}, {"Virgo", "Pisces"},
		{"Virgo", "Gemini"}, {"Libra", "Capricorn"}, {"Libra", "Virgo"}, {"Sagittarius", "Pisces"},
		{"Capricorn", "Aries"}, {"Capricorn", "Aquarius"}, {"Aquarius", "Aries"}, {"Pisces", "Capricorn"},
	}
	for _, p := range yes {
		for _, pair := range [][2]string{{p[0], p[1]}, {p[1], p[0]}} {
			r := vashyaKoota(pair[0], pair[1])
			if r.Score != 2 || r.RawScore != 2 || r.Maximum != 2 || r.Groom != pair[0] || r.Bride != pair[1] {
				t.Errorf("groom %s bride %s: %+v", pair[0], pair[1], r)
			}
		}
	}
	// the old animal-category pairs must no longer score: Aries/Taurus were both "Chatuspada" (2 before)
	for _, p := range [][2]string{{"Aries", "Taurus"}, {"Aries", "Aries"}, {"Leo", "Taurus"}, {"Gemini", "Cancer"}, {"Aries", "Scorpio"}} {
		if r := vashyaKoota(p[0], p[1]); r.Score != 0 {
			t.Errorf("%v should score 0, got %+v", p, r)
		}
	}
}

// The gana table is read as [bride gana][groom gana].
func TestGanaScoreForOrientation(t *testing.T) {
	// Deva: Ashwini; Manushya: Bharani; Rakshasa: Magha, Krittika (not exempt); Ashlesha is an exempt bride
	cases := []struct {
		groom, bride string
		want         float64
	}{
		{"Ashwini", "Bharani", 5},  // Deva groom, Manushya bride
		{"Bharani", "Ashwini", 6},  // Manushya groom, Deva bride
		{"Magha", "Ashwini", 0},    // Rakshasa groom, Deva bride
		{"Ashwini", "Krittika", 1}, // Deva groom, Rakshasa bride
		{"Bharani", "Krittika", 0}, // Manushya groom, Rakshasa bride
		{"Magha", "Bharani", 0},    // Rakshasa groom, Manushya bride
		{"Ashwini", "Hasta", 6},    // Deva / Deva
		{"Bharani", "Rohini", 6},   // Manushya / Manushya
		{"Magha", "Krittika", 6},   // Rakshasa / Rakshasa
		{"Ashwini", "Ashlesha", 6}, // exempt Rakshasa bride gets the full score
		{"Bharani", "Ashlesha", 6},
	}
	for _, c := range cases {
		if got := ganaScoreFor(c.groom, c.bride); got != c.want {
			t.Errorf("groom %s bride %s: %v want %v", c.groom, c.bride, got, c.want)
		}
	}
}

// Book rule: a gana mismatch is greatly reduced when the count from the bride's nakshatra to the groom's
// is more than 14. It is only a suggestion; the score does not change.
func TestGanaSuggestion(t *testing.T) {
	// groom Deva (Ashwini), bride Rakshasa (Krittika): score 1; count Krittika -> Ashwini is 26
	k := ganaKoota(moon("Ashwini", "Aries", 1), moon("Krittika", "Taurus", 1))
	if k.Score != 1 || k.Suggestion == "" {
		t.Errorf("mismatch with count 26 should suggest: %+v", k)
	}
	// mismatch but count 10 (Ashwini -> Magha): no suggestion
	k = ganaKoota(moon("Magha", "Leo", 1), moon("Ashwini", "Aries", 1))
	if k.Score != 0 || k.Suggestion != "" {
		t.Errorf("count 10 must not suggest: %+v", k)
	}
	// same gana with count 16 (Hasta -> Ashwini): nothing to reduce
	k = ganaKoota(moon("Ashwini", "Aries", 1), moon("Hasta", "Virgo", 1))
	if k.Score != 6 || k.Suggestion != "" {
		t.Errorf("no mismatch must not suggest: %+v", k)
	}
}

func TestNadiKoota(t *testing.T) {
	// same nadi (Antya), both in the 21-nakshatra list, no other cancellation: score stays 0, list is a suggestion
	k := nadiKoota(moon("Magha", "Leo", 1), moon("Rohini", "Taurus", 3))
	if k.RawScore != 0 || k.Score != 0 || !k.Dosha || k.CancellationApplied || k.Suggestion == "" || k.Maximum != 8 {
		t.Errorf("list only: %+v", k)
	}

	// different nadi: full score, no dosha, nothing to suggest
	k = nadiKoota(moon("Ashwini", "Aries", 1), moon("Bharani", "Aries", 1))
	if k.RawScore != 8 || k.Score != 8 || k.Dosha || k.Suggestion != "" || k.Groom != "Adi" || k.Bride != "Madhya" {
		t.Errorf("different nadi: %+v", k)
	}

	// same nakshatra, different pada: cancelled
	k = nadiKoota(moon("Magha", "Leo", 1), moon("Magha", "Leo", 2))
	if k.RawScore != 0 || k.Score != 8 || !k.CancellationApplied || !k.Dosha {
		t.Errorf("same nakshatra different pada: %+v", k)
	}

	// same rasi, different nakshatra, same nadi (both Adi): cancelled
	k = nadiKoota(moon("Ardra", "Gemini", 2), moon("Punarvasu", "Gemini", 1))
	if k.RawScore != 0 || k.Score != 8 || !k.CancellationApplied {
		t.Errorf("same rasi different nakshatra: %+v", k)
	}

	// same nadi, nakshatras not in the list, different signs: no cancellation, no suggestion
	k = nadiKoota(moon("Hasta", "Virgo", 1), moon("Jyeshtha", "Scorpio", 1))
	if k.RawScore != 0 || k.Score != 0 || k.CancellationApplied || k.Suggestion != "" {
		t.Errorf("no exemption: %+v", k)
	}
}

func TestBhakootKoota(t *testing.T) {
	// Aries/Taurus is a 2/12 dosha; the pair is not auspicious and Mars/Venus are neutral (Graha Maitri 3)
	k := bhakootKoota(moon("Ashwini", "Aries", 1), moon("Rohini", "Taurus", 1), 3)
	if !k.Dosha || k.RawScore != 0 || k.Score != 0 || k.CancellationApplied {
		t.Errorf("uncancelled dosha: %+v", k)
	}
	// the same dosha is cancelled when Graha Maitri is 4 or more
	k = bhakootKoota(moon("Ashwini", "Aries", 1), moon("Rohini", "Taurus", 1), 4)
	if !k.Dosha || k.RawScore != 0 || k.Score != 7 || !k.CancellationApplied {
		t.Errorf("cancelled by maitri: %+v", k)
	}
	// Aries/Pisces is a 2/12 pair listed as auspicious: cancelled whatever the maitri
	k = bhakootKoota(moon("Ashwini", "Aries", 1), moon("Revati", "Pisces", 1), 0)
	if k.Score != 7 || !k.CancellationApplied {
		t.Errorf("auspicious pair: %+v", k)
	}
	// no dosha at 4/10
	k = bhakootKoota(moon("Ashwini", "Aries", 1), moon("Pushya", "Cancer", 1), 0)
	if k.Dosha || k.Score != 7 {
		t.Errorf("4/10 is not a dosha: %+v", k)
	}
	d := k.Details.(domain.BhakootDetails)
	if d.GroomToBrideDistance != 4 || d.BrideToGroomDistance != 10 || d.Relationship != "4/10" {
		t.Errorf("details: %+v", d)
	}
}

func TestSummarize(t *testing.T) {
	k := domain.Kootas{
		Varna: domain.Koota{RawScore: 1, Score: 1}, Vashya: domain.Koota{RawScore: 2, Score: 2},
		Tara: domain.Koota{RawScore: 3, Score: 3}, Yoni: domain.Koota{RawScore: 4, Score: 4},
		GrahaMaitri: domain.Koota{RawScore: 5, Score: 5}, Gana: domain.Koota{RawScore: 0, Score: 0},
		Bhakoot: domain.Koota{RawScore: 0, Score: 7}, Nadi: domain.Koota{RawScore: 0, Score: 0},
	}
	s := summarize(k)
	if s.RawTotal != 15 || s.EffectiveTotal != 22 || s.Maximum != 36 || !s.TraditionalThresholdMet {
		t.Errorf("summary: %+v", s)
	}
	if s.Percentage < 61.1 || s.Percentage > 61.2 { // 22/36
		t.Errorf("percentage: %v", s.Percentage)
	}
	if s2 := summarize(domain.Kootas{Varna: domain.Koota{RawScore: 17.5, Score: 17.5}}); s2.TraditionalThresholdMet {
		t.Error("17.5 is below the threshold of 18")
	}
}
