package matchmaking

import (
	"fmt"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// The eight scored kootas. Each takes the two Moon positions and returns a domain.Koota.

func newKoota(name string, maximum, score float64, groom, bride, explanation string) domain.Koota {
	return domain.Koota{
		Name:        name,
		Maximum:     maximum,
		RawScore:    score,
		Score:       score,
		Groom:       groom,
		Bride:       bride,
		Explanation: explanation,
	}
}

// varnaKoota: 1 point if the groom's varna ranks equal to or above the bride's.
func varnaKoota(groom, bride domain.MoonDetails) domain.Koota {
	gv := VarnaMapping[groom.Sign]
	bv := VarnaMapping[bride.Sign]
	gr := VarnaRank[gv]
	br := VarnaRank[bv]

	var score float64
	if gr >= br {
		score = 1
	}
	return newKoota("Varna", 1, score, gv, bv,
		fmt.Sprintf("Groom varna %s (Rank %d) vs Bride varna %s (Rank %d)", gv, gr, bv, br))
}

// vashyaKoota applies the book's rasi-pair table: 2 points if either rasi lists the other as vashya.
func vashyaKoota(groomSign, brideSign string) domain.Koota {
	contains := func(list []string, s string) bool {
		for _, x := range list {
			if x == s {
				return true
			}
		}
		return false
	}
	var score float64
	var explanation string
	switch {
	case contains(VashyaRasis[groomSign], brideSign):
		score = 2
		explanation = fmt.Sprintf("Bride's rasi %s is vashya to groom's rasi %s.", brideSign, groomSign)
	case contains(VashyaRasis[brideSign], groomSign):
		score = 2
		explanation = fmt.Sprintf("Groom's rasi %s is vashya to bride's rasi %s.", groomSign, brideSign)
	default:
		explanation = fmt.Sprintf("No vashya relation is listed between %s and %s.", groomSign, brideSign)
	}
	return newKoota("Vashya", 2, score, groomSign, brideSign, explanation)
}

// taraKoota counts the nakshatras both ways; each favourable direction scores 1.5.
func taraKoota(groom, bride domain.MoonDetails) domain.Koota {
	gNakIdx := getNakshatraIndex(groom.Nakshatra)
	bNakIdx := getNakshatraIndex(bride.Nakshatra)

	// Count from Bride to Groom (Bride's perspective)
	countBtoG := (gNakIdx-bNakIdx+27)%27 + 1
	taraBtoGIdx := (countBtoG - 1) % 9

	// Count from Groom to Bride (Groom's perspective)
	countGtoB := (bNakIdx-gNakIdx+27)%27 + 1
	taraGtoBIdx := (countGtoB - 1) % 9

	taraNames := []string{"Janma", "Sampat", "Vipat", "Kshema", "Pratyari", "Sadhaka", "Vadha", "Mitra", "Ati-Mitra"}

	// Unfavorable: Janma(1), Vipat(3), Pratyari(5), Vadha(7) -> 0-indexed: 0, 2, 4, 6
	auspBtoG := !(taraBtoGIdx == 0 || taraBtoGIdx == 2 || taraBtoGIdx == 4 || taraBtoGIdx == 6)
	auspGtoB := !(taraGtoBIdx == 0 || taraGtoBIdx == 2 || taraGtoBIdx == 4 || taraGtoBIdx == 6)

	// Exception 1A: Anu Janma (2nd series: 10 to 18 distance) Pada exemptions (Page 3)
	if countBtoG >= 10 && countBtoG <= 18 {
		if taraBtoGIdx == 2 && groom.Pada != 1 {
			auspBtoG = true
		}
		if taraBtoGIdx == 4 && groom.Pada != 4 {
			auspBtoG = true
		}
		if taraBtoGIdx == 6 && groom.Pada != 3 {
			auspBtoG = true
		}
	}
	if countGtoB >= 10 && countGtoB <= 18 {
		if taraGtoBIdx == 2 && bride.Pada != 1 {
			auspGtoB = true
		}
		if taraGtoBIdx == 4 && bride.Pada != 4 {
			auspGtoB = true
		}
		if taraGtoBIdx == 6 && bride.Pada != 3 {
			auspGtoB = true
		}
	}

	// Exception 1B: Eka Nakshatra (Same Nakshatra)
	if countBtoG == 1 {
		if BadEkaNakshatras[bride.Nakshatra] {
			auspBtoG = false
			auspGtoB = false
		} else {
			auspBtoG = true
			auspGtoB = true
		}
	}

	// Exception 2: Revati is supremely auspicious
	if bride.Nakshatra == "Revati" || groom.Nakshatra == "Revati" {
		auspBtoG = true
		auspGtoB = true
	}

	// Exception 3: Vipat Tara (Index 2) is cancelled if Moon signs are Sama-Saptama (1/7 axis)
	if (getSignIndex(groom.Sign)-getSignIndex(bride.Sign)+12)%12+1 == 7 {
		if taraBtoGIdx == 2 {
			auspBtoG = true
		}
		if taraGtoBIdx == 2 {
			auspGtoB = true
		}
	}

	var scoreBtoG, scoreGtoB float64
	if auspBtoG {
		scoreBtoG = 1.5
	}
	if auspGtoB {
		scoreGtoB = 1.5
	}

	k := newKoota("Tara", 3, scoreBtoG+scoreGtoB, groom.Nakshatra, bride.Nakshatra,
		"Calculated cyclic distances between birth nakshatras.")
	k.Details = domain.TaraDetails{
		GroomToBride: domain.TaraDirection{NakshatraCount: countGtoB, Tara: taraNames[taraGtoBIdx], Auspicious: auspGtoB, Score: scoreGtoB},
		BrideToGroom: domain.TaraDirection{NakshatraCount: countBtoG, Tara: taraNames[taraBtoGIdx], Auspicious: auspBtoG, Score: scoreBtoG},
	}
	return k
}

// yoniKoota scores the animal compatibility of the two birth nakshatras.
func yoniKoota(groom, bride domain.MoonDetails) domain.Koota {
	gYoni := NakshatraYoni[groom.Nakshatra]
	bYoni := NakshatraYoni[bride.Nakshatra]
	return newKoota("Yoni", 4, YoniScores[gYoni][bYoni], gYoni, bYoni,
		fmt.Sprintf("Groom yoni %s vs Bride yoni %s", gYoni, bYoni))
}

// grahaMaitriKoota scores the natural friendship of the two Moon-sign lords.
func grahaMaitriKoota(groom, bride domain.MoonDetails) domain.Koota {
	gLord := groom.RashiLord
	bLord := bride.RashiLord
	gToBRel := "Same"
	bToGRel := "Same"
	var score float64 = 5.0

	if gLord != bLord {
		gToBRel = GrahaNaturalMaitri[gLord][bLord]
		bToGRel = GrahaNaturalMaitri[bLord][gLord]

		if gToBRel == "Friend" && bToGRel == "Friend" {
			score = 5
		}
		if (gToBRel == "Friend" && bToGRel == "Neutral") || (bToGRel == "Friend" && gToBRel == "Neutral") {
			score = 4
		}
		if gToBRel == "Neutral" && bToGRel == "Neutral" {
			score = 3
		}
		if (gToBRel == "Friend" && bToGRel == "Enemy") || (bToGRel == "Friend" && gToBRel == "Enemy") {
			score = 1
		}
		if (gToBRel == "Neutral" && bToGRel == "Enemy") || (bToGRel == "Neutral" && gToBRel == "Enemy") {
			score = 0.5
		}
		if gToBRel == "Enemy" && bToGRel == "Enemy" {
			score = 0
		}
	}

	k := newKoota("Graha Maitri", 5, score, gLord, bLord, "Natural relationship between moon sign lords.")
	k.Details = domain.GrahaMaitriDetails{
		GroomRashi:                       groom.Sign,
		BrideRashi:                       bride.Sign,
		GroomLordRelationshipToBrideLord: gToBRel,
		BrideLordRelationshipToGroomLord: bToGRel,
	}
	return k
}

// ganaKoota scores gana compatibility. When the gana is mismatched but the count from the bride's
// nakshatra to the groom's is more than 14, the book says the harm is greatly reduced: that is
// reported in Suggestion and does not change the score.
func ganaKoota(groom, bride domain.MoonDetails) domain.Koota {
	gGana := NakshatraGana[groom.Nakshatra]
	bGana := NakshatraGana[bride.Nakshatra]
	score := ganaScoreFor(groom.Nakshatra, bride.Nakshatra)

	k := newKoota("Gana", 6, score, gGana, bGana, fmt.Sprintf("Groom gana %s vs Bride gana %s", gGana, bGana))
	if count := nakshatraCount(bride.Nakshatra, groom.Nakshatra); score < 6 && count > 14 {
		k.Suggestion = fmt.Sprintf("The count from the bride's nakshatra to the groom's is %d (more than 14), which greatly reduces the harm of a gana mismatch (Telugu muhurta book, p. 65). This is a suggestion only and is not counted in the score.", count)
	}
	return k
}

// ganaScoreFor scores gana compatibility from the published table, which is read as
// [bride gana][groom gana]. A Rakshasa bride whose nakshatra is in GanaExemptBrideNakshatras
// gets the full 6.
func ganaScoreFor(groomNakshatra, brideNakshatra string) float64 {
	gGana := NakshatraGana[groomNakshatra]
	bGana := NakshatraGana[brideNakshatra]
	score := GanaScores[bGana][gGana]
	if bGana == "Rakshasa" && GanaExemptBrideNakshatras[brideNakshatra] {
		score = 6
	}
	return score
}

// bhakootKoota scores the Moon-sign distances. Doshas (2/12, 5/9, 6/8, either order) score 0 unless
// cancelled by an auspicious sign pair or by Graha Maitri of 4 or more.
func bhakootKoota(groom, bride domain.MoonDetails, grahaMaitriScore float64) domain.Koota {
	gSignIdx := getSignIndex(groom.Sign)
	bSignIdx := getSignIndex(bride.Sign)

	gToBDist := (bSignIdx-gSignIdx+12)%12 + 1
	bToGDist := (gSignIdx-bSignIdx+12)%12 + 1

	rel := fmt.Sprintf("%d/%d", gToBDist, bToGDist)
	isDosha := (rel == "2/12" || rel == "12/2" || rel == "5/9" || rel == "9/5" || rel == "6/8" || rel == "8/6")

	var raw float64 = 7
	if isDosha {
		raw = 0
	}

	score := raw
	cancelled := false
	reason := ""
	if isDosha {
		pair1 := groom.Sign + "-" + bride.Sign
		pair2 := bride.Sign + "-" + groom.Sign

		// Telugu book auspicious exceptions
		if AuspiciousBhakootPairs[pair1] || AuspiciousBhakootPairs[pair2] {
			cancelled = true
			reason = "Bhakoot dosha cancelled because this specific sign combination is highly auspicious."
			score = 7
		} else if grahaMaitriScore >= 4.0 { // same lord or friendly lords
			cancelled = true
			reason = "Bhakoot dosha cancelled because Rashi lords are same or friends (Graha Maitri)."
			score = 7
		}
	}

	explanation := "Sign distance between Moons."
	if bToGDist == 2 {
		explanation += " WARNING: Groom's Moon is 2nd from Bride's Moon (Risk of Vaidhavya/Poverty per Telugu rules)."
	} else if bToGDist == 5 {
		explanation += " WARNING: Groom's Moon is 5th from Bride's Moon (Risk of Vaidhavya/Loss of children per Telugu rules)."
	} else if bToGDist == 3 {
		explanation += " WARNING: Groom's Moon is 3rd from Bride's Moon (Risk of Sorrows/Hardships)."
	} else if bToGDist == 6 {
		explanation += " WARNING: Groom's Moon is 6th from Bride's Moon (Risk of Loss of children/Quarrels)."
	}

	k := newKoota("Bhakoot", 7, score, groom.Sign, bride.Sign, explanation)
	k.RawScore = raw
	k.Dosha = isDosha
	k.CancellationApplied = cancelled
	k.CancellationReason = reason
	k.Details = domain.BhakootDetails{GroomToBrideDistance: gToBDist, BrideToGroomDistance: bToGDist, Relationship: rel}
	return k
}

// nadiKoota scores Nadi. Same nadi scores 0. Cancellations that change the score: same nakshatra with a
// different pada, or same rasi with a different nakshatra. The Jaimini Maharshi list of 21 nakshatras does
// NOT change the score; when it applies it is reported in Suggestion only.
func nadiKoota(groom, bride domain.MoonDetails) domain.Koota {
	gNadi := NakshatraNadi[groom.Nakshatra]
	bNadi := NakshatraNadi[bride.Nakshatra]
	same := gNadi == bNadi

	k := newKoota("Nadi", 8, 8, gNadi, bNadi, "Nadi mapping.")
	if !same {
		return k
	}

	k.Dosha = true
	k.RawScore = 0
	k.Score = 0
	if groom.Nakshatra == bride.Nakshatra && groom.Pada != bride.Pada {
		k.CancellationApplied = true
		k.CancellationReason = "Nadi dosha cancelled due to same Nakshatra but different Padas."
		k.Score = 8
	} else if groom.Sign == bride.Sign && groom.Nakshatra != bride.Nakshatra {
		k.CancellationApplied = true
		k.CancellationReason = "Nadi dosha cancelled due to same Rashi but different Nakshatras."
		k.Score = 8
	}
	if JaiminiNadiExceptions[groom.Nakshatra] && JaiminiNadiExceptions[bride.Nakshatra] {
		k.Suggestion = fmt.Sprintf("Per the Jaimini Maharshi exemption list, Nadi dosha does not apply when both nakshatras (%s and %s) are in that list. This is a suggestion only and is not counted in the score.", groom.Nakshatra, bride.Nakshatra)
	}
	return k
}

// summarize adds up the eight kootas (raw and counted) against the 36-point maximum.
func summarize(k domain.Kootas) domain.MatchSummary {
	raw := k.Varna.RawScore + k.Vashya.RawScore + k.Tara.RawScore + k.Yoni.RawScore + k.GrahaMaitri.RawScore + k.Gana.RawScore + k.Bhakoot.RawScore + k.Nadi.RawScore
	eff := k.Varna.Score + k.Vashya.Score + k.Tara.Score + k.Yoni.Score + k.GrahaMaitri.Score + k.Gana.Score + k.Bhakoot.Score + k.Nadi.Score
	return domain.MatchSummary{
		RawTotal:                raw,
		EffectiveTotal:          eff,
		Maximum:                 36,
		Percentage:              (eff / 36.0) * 100.0,
		TraditionalThresholdMet: eff >= 18.0,
	}
}
