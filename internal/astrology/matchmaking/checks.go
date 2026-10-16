package matchmaking

import (
	"fmt"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// Compatibility checks that do not add to the 36 points. Each returns a domain.Check whose Satisfied
// field is true when there is no outstanding problem.

// vedhaPairs are the nakshatra pairs that obstruct each other (Telugu muhurta book p.71).
// NOTE: the book prints Pushya's partner as "Punarvasu", which already pairs with Uttara Ashadha;
// Purva Ashadha (the standard partner) is used here.
var vedhaPairs = [][2]string{
	{"Ashwini", "Jyeshtha"},
	{"Bharani", "Anuradha"},
	{"Krittika", "Vishakha"},
	{"Rohini", "Swati"},
	{"Ardra", "Shravana"},
	{"Punarvasu", "Uttara Ashadha"},
	{"Pushya", "Purva Ashadha"},
	{"Ashlesha", "Mula"},
	{"Magha", "Revati"},
	{"Purva Phalguni", "Uttara Bhadrapada"},
	{"Uttara Phalguni", "Purva Bhadrapada"},
	{"Hasta", "Shatabhisha"},
	{"Mrigashira", "Dhanishta"},
}

var vedhaPartner = func() map[string]string {
	m := make(map[string]string, len(vedhaPairs)*2)
	for _, p := range vedhaPairs {
		m[p[0]] = p[1]
		m[p[1]] = p[0]
	}
	return m
}()

// nakshatraCount counts from the `from` nakshatra to the `to` nakshatra, both inclusive (1..27).
func nakshatraCount(from, to string) int {
	return (getNakshatraIndex(to)-getNakshatraIndex(from)+27)%27 + 1
}

// rajjuCheck: both in the same rajju group is a dosha, cancelled when the Moon-sign lords are the same
// or Graha Maitri is 4 or more, or when the lagna lords are the same or friendly.
func rajjuCheck(groom, bride domain.MoonDetails, grahaMaitriScore float64, lagnaLordsRelationship string) domain.Check {
	gRajju := NakshatraRajju[groom.Nakshatra]
	bRajju := NakshatraRajju[bride.Nakshatra]
	dosha := gRajju == bRajju

	cancelled := false
	reason := ""
	if dosha {
		if groom.RashiLord == bride.RashiLord || grahaMaitriScore >= 4.0 {
			cancelled = true
			reason = "Rajju dosha cancelled because Rasi Lords are same or friendly."
		} else if lagnaLordsRelationship == "Same" || lagnaLordsRelationship == "Friendly" {
			cancelled = true
			reason = "Rajju dosha cancelled because Lagna Lords are same or friendly."
		}
	}

	explanation := fmt.Sprintf("Groom Rajju is %s. Bride Rajju is %s.", gRajju, bRajju)
	if dosha && !cancelled {
		explanation += " FATAL DOSHA: Both belong to the same Rajju. This is highly inauspicious."
	}
	return domain.Check{
		Name:                "Rajju",
		Satisfied:           !dosha || cancelled,
		Groom:               gRajju,
		Bride:               bRajju,
		CancellationApplied: cancelled,
		CancellationReason:  reason,
		Explanation:         explanation,
	}
}

// vedhaCheck: certain nakshatra pairs obstruct each other.
func vedhaCheck(groom, bride domain.MoonDetails) domain.Check {
	has := vedhaPartner[groom.Nakshatra] == bride.Nakshatra
	explanation := "The two nakshatras do not obstruct each other (no vedha)."
	if has {
		explanation = fmt.Sprintf("%s and %s obstruct each other (vedha); marriage should proceed only with remedial indications.", groom.Nakshatra, bride.Nakshatra)
	}
	return domain.Check{Name: "Vedha", Satisfied: !has, Groom: groom.Nakshatra, Bride: bride.Nakshatra, Explanation: explanation}
}

// mahendraCheck: the groom's nakshatra is the 4th, 7th, 10th, 13th, 16th, 19th, 22nd or 25th from the bride's.
func mahendraCheck(groom, bride domain.MoonDetails, count int) domain.Check {
	ok := count >= 4 && count <= 25 && (count-4)%3 == 0
	explanation := fmt.Sprintf("Groom's nakshatra is the %dth from the bride's, not one of 4, 7, 10, 13, 16, 19, 22, 25.", count)
	if ok {
		explanation = fmt.Sprintf("Groom's nakshatra is the %dth from the bride's, one of 4, 7, 10, 13, 16, 19, 22, 25: auspicious, increases married life.", count)
	}
	return domain.Check{
		Name: "Mahendra", Satisfied: ok, Groom: groom.Nakshatra, Bride: bride.Nakshatra,
		Explanation: explanation, Details: domain.CountDetails{BrideToGroomCount: count},
	}
}

// streeDeerghaCheck: the 27 nakshatras counted from the bride's are split into three navakas of 9.
// The groom in the 1st navaka is quarrelsome, the 2nd average, the 3rd auspicious. It is satisfied when
// the count is above 9 (some opinions: above 7). If Bhakoot and Graha Maitri are good, a shortfall is
// tolerable ("good" Graha Maitri is taken as a score of 4 or more).
func streeDeerghaCheck(groom, bride domain.MoonDetails, count int, bhakootGood bool, grahaMaitriScore float64) domain.Check {
	navaka := (count-1)/9 + 1
	quality := map[int]string{1: "Kalaham (quarrel)", 2: "Samanyam (average)", 3: "Shubham (auspicious)"}[navaka]
	satisfied := count > 9
	exception := !satisfied && bhakootGood && grahaMaitriScore >= 4

	explanation := fmt.Sprintf("Groom's nakshatra is the %dth from the bride's, in navaka %d: %s.", count, navaka, quality)
	if exception {
		explanation += " Rasi koota and graha maitri are good, so the shortfall can be tolerated."
	}
	return domain.Check{
		Name: "Stree Deergha", Satisfied: satisfied, Groom: groom.Nakshatra, Bride: bride.Nakshatra,
		Explanation: explanation,
		Details: domain.StreeDeerghaDetails{
			BrideToGroomCount: count, Navaka: navaka, Quality: quality,
			SatisfiedByAlternateOpinion: count > 7, ExceptionApplies: exception,
		},
	}
}

// mrugaVairamCheck: the yoni animals of the couple are natural enemies.
func mrugaVairamCheck(groom, bride domain.MoonDetails) domain.Check {
	gYoni := NakshatraYoni[groom.Nakshatra]
	bYoni := NakshatraYoni[bride.Nakshatra]
	has := MrugaEnmity[gYoni] == bYoni

	explanation := "Checks if the Yoni (Animal) of the couple are natural enemies."
	if has {
		explanation = fmt.Sprintf("MRUGA VAIRAM DOSHA: Groom's Yoni animal is %s and Bride's Yoni animal is %s. These are natural enemies, indicating extreme incompatibility.", gYoni, bYoni)
	}
	return domain.Check{Name: "Mruga Vairam", Satisfied: !has, Groom: gYoni, Bride: bYoni, Explanation: explanation}
}

// specialTaraCheck: specifically rejected nakshatra pairings (Ashubha Kurpu and three more).
func specialTaraCheck(groom, bride domain.MoonDetails) domain.Check {
	has := false
	if bride.Nakshatra == "Krittika" && groom.Nakshatra == "Bharani" {
		has = true
	} else if bride.Nakshatra == "Ashlesha" && groom.Nakshatra == "Pushya" {
		has = true
	} else if bride.Nakshatra == "Shatabhisha" && groom.Nakshatra == "Dhanishta" {
		has = true
	}
	if AshubhaKurpu[bride.Nakshatra] == groom.Nakshatra {
		has = true
	}

	explanation := "Checks for specifically rejected Tara combinations (Ashubha Kurpu)."
	if has {
		explanation = fmt.Sprintf("ASHUBHA KURPU DOSHA: The pairing of Bride's %s and Groom's %s is explicitly forbidden in the Telugu matchmaking text.", bride.Nakshatra, groom.Nakshatra)
	}
	return domain.Check{Name: "Special Tara", Satisfied: !has, Groom: groom.Nakshatra, Bride: bride.Nakshatra, Explanation: explanation}
}

// kalaSarpaCheck: the match is rejected when both persons have Kala Sarpa dosha.
func kalaSarpaCheck(groomCtx, brideCtx *domain.CalculationContext) domain.Check {
	g := checkKalaSarpa(groomCtx)
	b := checkKalaSarpa(brideCtx)
	both := g && b

	explanation := "Checks if both Bride and Groom have Kala Sarpa Dosha. (If both have it, the match is rejected)."
	if both {
		explanation = "FATAL DOSHA: Both the Groom and the Bride suffer from Kala Sarpa Dosha in their natal charts. The match is explicitly rejected."
	}
	present := func(has bool) string {
		if has {
			return "Present"
		}
		return "Absent"
	}
	return domain.Check{Name: "Kala Sarpa", Satisfied: !both, Groom: present(g), Bride: present(b), Explanation: explanation}
}
