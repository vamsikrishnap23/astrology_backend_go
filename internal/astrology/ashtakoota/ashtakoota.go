package ashtakoota

import (
	"fmt"

	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func getNakshatraIndex(nak string) int {
	for i, n := range Nakshatras {
		if n == nak {
			return i
		}
	}
	return 0
}

func getSignIndex(sign string) int {
	for i, s := range Signs {
		if s == sign {
			return i
		}
	}
	return 0
}

func CalculateMatch(groomCtx, brideCtx *domain.CalculationContext) (domain.AshtakootaResult, error) {
	groomMoon := getMoonDetails(groomCtx.JulianDayUT, groomCtx.Ayanamsa)
	brideMoon := getMoonDetails(brideCtx.JulianDayUT, brideCtx.Ayanamsa)

	// 1. Varna
	gv := VarnaMapping[groomMoon.Sign]
	bv := VarnaMapping[brideMoon.Sign]
	gvr := VarnaRank[gv]
	bvr := VarnaRank[bv]

	var varnaScore float64 = 0
	if gvr >= bvr {
		varnaScore = 1
	}

	varnaRes := domain.KootaResult{
		Score:       varnaScore,
		Maximum:     1,
		GroomValue:  gv,
		BrideValue:  bv,
		Explanation: fmt.Sprintf("Groom varna %s (Rank %d) vs Bride varna %s (Rank %d)", gv, gvr, bv, bvr),
	}

	// 2. Vashya
	gvash := VashyaMapping[groomMoon.Sign]
	bvash := VashyaMapping[brideMoon.Sign]
	vashyaScore := VashyaScores[gvash][bvash]

	vashyaRes := domain.KootaResult{
		Score:       vashyaScore,
		Maximum:     2,
		GroomValue:  gvash,
		BrideValue:  bvash,
		Explanation: fmt.Sprintf("Groom vashya %s vs Bride vashya %s", gvash, bvash),
	}

	// 3. Tara
	gNakIdx := getNakshatraIndex(groomMoon.Nakshatra)
	bNakIdx := getNakshatraIndex(brideMoon.Nakshatra)

	// Count from Bride to Groom (Bride's perspective)
	countBtoG := (gNakIdx-bNakIdx+27)%27 + 1
	taraBtoGIdx := (countBtoG - 1) % 9

	// Count from Groom to Bride (Groom's perspective)
	countGtoB := (bNakIdx-gNakIdx+27)%27 + 1
	taraGtoBIdx := (countGtoB - 1) % 9

	taraNames := []string{"Janma", "Sampat", "Vipat", "Kshema", "Pratyari", "Sadhaka", "Vadha", "Mitra", "Ati-Mitra"}
	// Unfavorable: Vipat(3), Pratyari(5), Vadha(7) (0-indexed: 2, 4, 6)

	auspBtoG := !(taraBtoGIdx == 2 || taraBtoGIdx == 4 || taraBtoGIdx == 6)
	auspGtoB := !(taraGtoBIdx == 2 || taraGtoBIdx == 4 || taraGtoBIdx == 6)

	var scoreBtoG, scoreGtoB float64 = 0, 0
	if auspBtoG {
		scoreBtoG = 1.5
	}
	if auspGtoB {
		scoreGtoB = 1.5
	}

	taraRes := domain.TaraKootaResult{
		Score:   scoreBtoG + scoreGtoB,
		Maximum: 3,
		GroomToBride: domain.TaraDirection{
			NakshatraCount: countGtoB,
			Tara:           taraNames[taraGtoBIdx],
			Auspicious:     auspGtoB,
			Score:          scoreGtoB,
		},
		BrideToGroom: domain.TaraDirection{
			NakshatraCount: countBtoG,
			Tara:           taraNames[taraBtoGIdx],
			Auspicious:     auspBtoG,
			Score:          scoreBtoG,
		},
		Explanation: "Calculated cyclic distances between birth nakshatras.",
	}

	// 4. Yoni
	gYoni := NakshatraYoni[groomMoon.Nakshatra]
	bYoni := NakshatraYoni[brideMoon.Nakshatra]
	yoniScore := YoniScores[gYoni][bYoni]

	yoniRes := domain.KootaResult{
		Score:       yoniScore,
		Maximum:     4,
		GroomValue:  gYoni,
		BrideValue:  bYoni,
		Explanation: fmt.Sprintf("Groom yoni %s vs Bride yoni %s", gYoni, bYoni),
	}

	// 5. Graha Maitri
	gLord := groomMoon.RashiLord
	bLord := brideMoon.RashiLord
	gToBRel := "Same"
	bToGRel := "Same"
	var gmScore float64 = 5.0

	if gLord != bLord {
		gToBRel = GrahaNaturalMaitri[gLord][bLord]
		bToGRel = GrahaNaturalMaitri[bLord][gLord]

		if gToBRel == "Friend" && bToGRel == "Friend" {
			gmScore = 5
		}
		if (gToBRel == "Friend" && bToGRel == "Neutral") || (bToGRel == "Friend" && gToBRel == "Neutral") {
			gmScore = 4
		}
		if gToBRel == "Neutral" && bToGRel == "Neutral" {
			gmScore = 3
		}
		if (gToBRel == "Friend" && bToGRel == "Enemy") || (bToGRel == "Friend" && gToBRel == "Enemy") {
			gmScore = 1
		} // Actually 1 or 0.5 depending on rules, 1 is standard
		if (gToBRel == "Neutral" && bToGRel == "Enemy") || (bToGRel == "Neutral" && gToBRel == "Enemy") {
			gmScore = 0.5
		}
		if gToBRel == "Enemy" && bToGRel == "Enemy" {
			gmScore = 0
		}
	}

	gmRes := domain.GrahaMaitriResult{
		Score:                            gmScore,
		Maximum:                          5,
		GroomRashi:                       groomMoon.Sign,
		BrideRashi:                       brideMoon.Sign,
		GroomRashiLord:                   gLord,
		BrideRashiLord:                   bLord,
		GroomLordRelationshipToBrideLord: gToBRel,
		BrideLordRelationshipToGroomLord: bToGRel,
		Explanation:                      "Natural relationship between moon sign lords.",
	}

	// 6. Gana
	gGana := NakshatraGana[groomMoon.Nakshatra]
	bGana := NakshatraGana[brideMoon.Nakshatra]
	ganaScore := GanaScores[gGana][bGana]

	ganaRes := domain.KootaResult{
		Score:       ganaScore,
		Maximum:     6,
		GroomValue:  gGana,
		BrideValue:  bGana,
		Explanation: fmt.Sprintf("Groom gana %s vs Bride gana %s", gGana, bGana),
	}

	// 7. Bhakoot
	gSignIdx := getSignIndex(groomMoon.Sign)
	bSignIdx := getSignIndex(brideMoon.Sign)

	gToBDist := (bSignIdx-gSignIdx+12)%12 + 1
	bToGDist := (gSignIdx-bSignIdx+12)%12 + 1

	rel := fmt.Sprintf("%d/%d", gToBDist, bToGDist)
	isDosha := (rel == "2/12" || rel == "12/2" || rel == "5/9" || rel == "9/5" || rel == "6/8" || rel == "8/6")

	var rawBhakoot float64 = 7
	if isDosha {
		rawBhakoot = 0
	}

	effectiveBhakoot := rawBhakoot
	cancellation := false
	cancelReason := ""

	if isDosha {
		// Cancellation: same lord or friendly lords
		if gmScore >= 4.0 { // Lords are same or friendly
			cancellation = true
			cancelReason = "Bhakoot dosha cancelled because Rashi lords are same or friends (Graha Maitri)."
			effectiveBhakoot = 7
		}
	}

	bhakootRes := domain.BhakootResult{
		RawScore:             rawBhakoot,
		EffectiveScore:       effectiveBhakoot,
		Maximum:              7,
		GroomRashi:           groomMoon.Sign,
		BrideRashi:           brideMoon.Sign,
		GroomToBrideDistance: gToBDist,
		BrideToGroomDistance: bToGDist,
		Relationship:         rel,
		Dosha:                isDosha,
		CancellationApplied:  cancellation,
		CancellationReason:   cancelReason,
		Explanation:          "Sign distance between Moons.",
	}

	// 8. Nadi
	gNadi := NakshatraNadi[groomMoon.Nakshatra]
	bNadi := NakshatraNadi[brideMoon.Nakshatra]
	sameNadi := (gNadi == bNadi)

	var rawNadi float64 = 8
	if sameNadi {
		rawNadi = 0
	}

	effNadi := rawNadi
	nCancel := false
	nReason := ""

	if sameNadi {
		// Cancellation 1: Same nakshatra but different padas
		if groomMoon.Nakshatra == brideMoon.Nakshatra && groomMoon.Pada != brideMoon.Pada {
			nCancel = true
			nReason = "Nadi dosha cancelled due to same Nakshatra but different Padas."
			effNadi = 8
		} else if groomMoon.Sign == brideMoon.Sign && groomMoon.Nakshatra != brideMoon.Nakshatra {
			// Cancellation 2: Same sign but different nakshatras
			nCancel = true
			nReason = "Nadi dosha cancelled due to same Rashi but different Nakshatras."
			effNadi = 8
		}
	}

	nadiRes := domain.NadiResult{
		RawScore:            rawNadi,
		EffectiveScore:      effNadi,
		Maximum:             8,
		GroomNadi:           gNadi,
		BrideNadi:           bNadi,
		SameNadi:            sameNadi,
		Dosha:               sameNadi,
		CancellationApplied: nCancel,
		CancellationReason:  nReason,
		Explanation:         "Nadi mapping.",
	}

	rawTotal := varnaScore + vashyaScore + taraRes.Score + yoniScore + gmScore + ganaScore + rawBhakoot + rawNadi
	effTotal := varnaScore + vashyaScore + taraRes.Score + yoniScore + gmScore + ganaScore + effectiveBhakoot + effNadi

	res := domain.AshtakootaResult{
		RuleSet: domain.RuleSetInfo{Name: "classical-guna-milan-v1", Ayanamsa: groomCtx.Input.Ayanamsa},
		Groom:   domain.PersonMatchDetails{Name: groomCtx.Input.Name, Moon: groomMoon},
		Bride:   domain.PersonMatchDetails{Name: brideCtx.Input.Name, Moon: brideMoon},
		Kootas: domain.Kootas{
			Varna:       varnaRes,
			Vashya:      vashyaRes,
			Tara:        taraRes,
			Yoni:        yoniRes,
			GrahaMaitri: gmRes,
			Gana:        ganaRes,
			Bhakoot:     bhakootRes,
			Nadi:        nadiRes,
		},
		Summary: domain.MatchSummary{
			RawTotal:                rawTotal,
			EffectiveTotal:          effTotal,
			Maximum:                 36,
			Percentage:              (effTotal / 36.0) * 100.0,
			TraditionalThresholdMet: effTotal >= 18.0,
		},
	}

	return res, nil
}
