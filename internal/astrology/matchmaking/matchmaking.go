package matchmaking

import (
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"math"

	"fmt"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
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

func checkKalaSarpa(ctx *domain.CalculationContext) bool {
	pls, err := planets.CalculatePlanets(ctx)
	if err != nil {
		return false
	}

	var rahuLon, ketuLon float64
	var otherLons []float64

	for _, p := range pls {
		if p.Planet == "Rahu" {
			rahuLon = p.SiderealLongitude
		} else if p.Planet == "Ketu" {
			ketuLon = p.SiderealLongitude
		} else if p.Planet != "Ascendant" && p.Planet != "MC" && p.Planet != "Uranus" && p.Planet != "Neptune" && p.Planet != "Pluto" {
			otherLons = append(otherLons, p.SiderealLongitude)
		}
	}

	rahuSide := true
	ketuSide := true

	for _, lon := range otherLons {
		distRahu := math.Mod(lon-rahuLon+360.0, 360.0)
		if distRahu > 180.0 {
			rahuSide = false
		}

		distKetu := math.Mod(lon-ketuLon+360.0, 360.0)
		if distKetu > 180.0 {
			ketuSide = false
		}
	}

	return rahuSide || ketuSide
}

func evaluateIndividualRisks(ctx *domain.CalculationContext, isGroom bool) domain.IndividualRisks {
	res := domain.IndividualRisks{
		VaidhavyaDosham:   domain.IndividualRiskCheck{HasRisk: false, RiskFactors: []string{}},
		DwikalatraDosham:  domain.IndividualRiskCheck{HasRisk: false, RiskFactors: []string{}},
		Napumsakatvam:     domain.IndividualRiskCheck{HasRisk: false, RiskFactors: []string{}},
		DampatyaMarakatva: domain.IndividualRiskCheck{HasRisk: false, RiskFactors: []string{}},
	}
	pls, err := planets.CalculatePlanets(ctx)
	if err != nil {
		return res
	}

	var ascSign string
	for _, p := range pls {
		if p.Planet == "Ascendant" {
			ascSign = p.Sign
			break
		}
	}

	ascIdx := getSignIndex(ascSign)
	h7Idx := (ascIdx + 6) % 12
	h8Idx := (ascIdx + 7) % 12

	h7Sign := Signs[h7Idx]
	h8Sign := Signs[h8Idx]

	h1Lord := SignLords[ascSign]
	h7Lord := SignLords[h7Sign]
	h8Lord := SignLords[h8Sign]

	planetsInSign := make(map[string][]string)
	var venusSign, h1LordSign, h7LordSign, h8LordSign string

	for _, p := range pls {
		if p.Planet != "Ascendant" && p.Planet != "MC" && p.Planet != "Uranus" && p.Planet != "Neptune" && p.Planet != "Pluto" {
			planetsInSign[p.Sign] = append(planetsInSign[p.Sign], p.Planet)
			if p.Planet == "Venus" {
				venusSign = p.Sign
			}
			if p.Planet == h1Lord {
				h1LordSign = p.Sign
			}
			if p.Planet == h7Lord {
				h7LordSign = p.Sign
			}
			if p.Planet == h8Lord {
				h8LordSign = p.Sign
			}
		}
	}

	// Helper to get House Number
	getHouseNum := func(sign string) int {
		sIdx := getSignIndex(sign)
		return (sIdx-ascIdx+12)%12 + 1
	}

	isMalefic := func(p string) bool {
		return p == "Sun" || p == "Mars" || p == "Saturn" || p == "Rahu" || p == "Ketu"
	}

	isPapaKartari := func(s string) bool {
		sIdx := getSignIndex(s)
		prevIdx := (sIdx + 11) % 12
		nextIdx := (sIdx + 1) % 12
		hasPrevMalefic, hasNextMalefic := false, false
		for _, p := range planetsInSign[Signs[prevIdx]] {
			if isMalefic(p) {
				hasPrevMalefic = true
			}
		}
		for _, p := range planetsInSign[Signs[nextIdx]] {
			if isMalefic(p) {
				hasNextMalefic = true
			}
		}
		return hasPrevMalefic && hasNextMalefic
	}

	var saturnSign string
	for _, p := range pls {
		if p.Planet == "Saturn" {
			saturnSign = p.Sign
			break
		}
	}

	// 1. Vaidhavya Dosham / Basic House Checks
	for _, p := range planetsInSign[h8Sign] {
		if p == "Saturn" {
			// Telugu Book Exception: Saturn in 8th is not a dosha for Taurus and Libra Lagnas
			if ascSign != "Taurus" && ascSign != "Libra" {
				res.VaidhavyaDosham.HasRisk = true
				res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Saturn is in the 8th house (risk of Alpayush/Vaidhavya).")
			}
		}
		if p == "Sun" {
			res.VaidhavyaDosham.HasRisk = true
			res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Sun is in the 8th house.")
		}
	}
	hasMars8 := false
	hasRahu8 := false
	for _, p := range planetsInSign[h8Sign] {
		if p == "Mars" {
			hasMars8 = true
		}
		if p == "Rahu" {
			hasRahu8 = true
		}
	}
	if hasMars8 && hasRahu8 {
		res.VaidhavyaDosham.HasRisk = true
		res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Mars and Rahu conjunct in the 8th house.")
	}

	if h7LordSign == h8Sign && h8LordSign == h7Sign {
		res.VaidhavyaDosham.HasRisk = true
		res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Exchange (Parivartana) of 7th and 8th lords.")
	}

	// New PDF Rules for Lords Placements
	h1LordHouse := getHouseNum(h1LordSign)
	if h1LordHouse == 6 || h1LordHouse == 8 || h1LordHouse == 12 {
		res.VaidhavyaDosham.HasRisk = true
		res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Lagnadhipati (Lagna Lord) is in the 6th, 8th, or 12th house.")
	}

	h7LordHouse := getHouseNum(h7LordSign)
	if h7LordHouse == 1 || h7LordHouse == 2 || h7LordHouse == 6 || h7LordHouse == 8 {
		res.DwikalatraDosham.HasRisk = true
		res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Saptamadhipati (7th Lord) is in the 1st, 2nd, 6th, or 8th house.")
	}

	h8LordHouse := getHouseNum(h8LordSign)
	if h8LordHouse == 1 || h8LordHouse == 2 || h8LordHouse == 3 || h8LordHouse == 4 || h8LordHouse == 6 || h8LordHouse == 7 || h8LordHouse == 9 {
		res.VaidhavyaDosham.HasRisk = true
		res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Ashtamadhipati (8th Lord) is in an inauspicious house.")
	}

	// Extreme Conjunction Checks from Telugu Book
	if h7LordSign == h8LordSign {
		conjHouse := getHouseNum(h7LordSign)
		if conjHouse == 1 || conjHouse == 4 || conjHouse == 7 || conjHouse == 10 {
			res.VaidhavyaDosham.HasRisk = true
			res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "7th and 8th lords are conjunct in a Kendra (1,4,7,10).")
		} else if conjHouse == 5 {
			res.VaidhavyaDosham.HasRisk = true
			res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "7th and 8th lords are conjunct in the 5th house.")
		}
	}

	// Kutumba Nasanam (Mercury + Rahu in 2nd house)
	h2Idx := (ascIdx + 1) % 12
	h2Sign := Signs[h2Idx]
	hasMerc2, hasRahu2 := false, false
	for _, p := range planetsInSign[h2Sign] {
		if p == "Mercury" {
			hasMerc2 = true
		}
		if p == "Rahu" {
			hasRahu2 = true
		}
	}
	if hasMerc2 && hasRahu2 {
		res.DwikalatraDosham.HasRisk = true
		res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Mercury and Rahu conjunct in the 2nd house (Kutumba Nasanam).")
	}

	// 2. Dwikalatra Dosham Check
	isDual := func(s string) bool {
		return s == "Gemini" || s == "Virgo" || s == "Sagittarius" || s == "Pisces"
	}
	// Fixed: The book specifically says they must be "kalisi" (conjunct) in a dual sign.
	if h1LordSign == h7LordSign && h7LordSign == venusSign && isDual(venusSign) {
		res.DwikalatraDosham.HasRisk = true
		res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Lagna Lord, 7th Lord, and Venus are conjunct together in a Dual Sign.")
	}

	if venusSign != "Taurus" && venusSign != "Libra" {
		for _, p := range planetsInSign[venusSign] {
			if p == "Mars" {
				res.DwikalatraDosham.HasRisk = true
				res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Mars and Venus conjunction outside of own signs.")
			}
		}
	}

	// 3. Extended Vaidhavya Rules
	for _, pList := range planetsInSign {
		hasM, hasV, hasR := false, false, false
		for _, p := range pList {
			if p == "Mars" {
				hasM = true
			}
			if p == "Venus" {
				hasV = true
			}
			if p == "Rahu" {
				hasR = true
			}
		}
		if hasM && hasV && hasR {
			res.VaidhavyaDosham.HasRisk = true
			res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Mars, Venus, and Rahu conjunction.")
		}
	}
	for _, p := range planetsInSign["Aquarius"] {
		if p == "Ketu" {
			res.VaidhavyaDosham.HasRisk = true
			res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Ketu in Aquarius (Madhyamayuvu).")
		}
	}

	// 4. Extended Dwikalatra Rules
	h2Idx = (ascIdx + 1) % 12
	h12Idx := (ascIdx + 11) % 12
	h2Lord := SignLords[Signs[h2Idx]]
	h12Lord := SignLords[Signs[h12Idx]]
	var h2LordSign, h12LordSign string
	for _, p := range pls {
		if p.Planet == h2Lord {
			h2LordSign = p.Sign
		}
		if p.Planet == h12Lord {
			h12LordSign = p.Sign
		}
	}
	if h2LordSign == Signs[h12Idx] && h12LordSign == Signs[h2Idx] {
		res.DwikalatraDosham.HasRisk = true
		res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Exchange (Parivartana) between 2nd and 12th lords.")
	}

	hasMerc7, hasSat7, hasMars7, hasVenus7 := false, false, false, false
	for _, p := range planetsInSign[h7Sign] {
		if p == "Mercury" {
			hasMerc7 = true
		}
		if p == "Saturn" {
			hasSat7 = true
		}
		if p == "Mars" {
			hasMars7 = true
		}
		if p == "Venus" {
			hasVenus7 = true
		}
	}
	if isGroom {
		if hasMerc7 && hasSat7 {
			res.DwikalatraDosham.HasRisk = true
			res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Mercury and Saturn conjunct in 7th house (Groom chart).")
		}
		if hasMars7 && hasVenus7 {
			res.DwikalatraDosham.HasRisk = true
			res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Mars and Venus conjunct in 7th house (Groom chart).")
		}
	}

	// 5. Corrected Napumsakatvam (Impotency indicators - Male only)
	if isGroom {
		var saturnSign string
		for _, p := range pls {
			if p.Planet == "Saturn" {
				saturnSign = p.Sign
				break
			}
		}
		hasMaleficLagna := false
		for _, p := range planetsInSign[ascSign] {
			if p == "Sun" || p == "Mars" || p == "Saturn" || p == "Rahu" || p == "Ketu" {
				hasMaleficLagna = true
				break
			}
		}
		if hasMaleficLagna {
			satHouse := getHouseNum(saturnSign)
			isWater := (saturnSign == "Cancer" || saturnSign == "Scorpio" || saturnSign == "Pisces")
			if (satHouse == 2 || satHouse == 6 || satHouse == 12) && isWater {
				res.Napumsakatvam.HasRisk = true
				res.Napumsakatvam.RiskFactors = append(res.Napumsakatvam.RiskFactors, "Saturn in 2nd/6th/12th in a water sign with malefic in Lagna.")
			}
		}
	}

	// 6. Dampatya Marakatva Dosham (Lethal Marital Flaws)
	for sign, pList := range planetsInSign {
		hasMoon, hasVen := false, false
		for _, p := range pList {
			if p == "Moon" {
				hasMoon = true
			}
			if p == "Venus" {
				hasVen = true
			}
		}
		if hasMoon && hasVen {
			oppIdx := (getSignIndex(sign) + 6) % 12
			hasOppMars, hasOppSat := false, false
			for _, p := range planetsInSign[Signs[oppIdx]] {
				if p == "Mars" {
					hasOppMars = true
				}
				if p == "Saturn" {
					hasOppSat = true
				}
			}
			if hasOppMars && hasOppSat {
				res.DampatyaMarakatva.HasRisk = true
				res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "Moon and Venus conjunct with Mars and Saturn in opposition.")
			}
		}
	}

	// Dampatya Marakatva Rule 8 & 18 (House 5, 7, 8 lord interlocks)
	h5Idx := (ascIdx + 4) % 12
	h5Sign := Signs[h5Idx]
	h5Lord := SignLords[h5Sign]
	var h5LordSign string
	for _, p := range pls {
		if p.Planet == h5Lord {
			h5LordSign = p.Sign
			break
		}
	}

	if h5LordSign == h7Sign || h7LordSign == h5Sign {
		res.DampatyaMarakatva.HasRisk = true
		res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "5th and 7th lords connection/exchange.")
	}
	if h7LordSign == h8Sign || h8LordSign == h7Sign {
		res.DampatyaMarakatva.HasRisk = true
		res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "7th and 8th lords placement/exchange.")
	}

	// Dampatya Marakatva Rule 11 (Saturn + Venus conjunction outside Libra)
	for sign, pList := range planetsInSign {
		if sign != "Libra" {
			hasSat, hasVen := false, false
			for _, p := range pList {
				if p == "Saturn" {
					hasSat = true
				}
				if p == "Venus" {
					hasVen = true
				}
			}
			if hasSat && hasVen {
				res.DampatyaMarakatva.HasRisk = true
				res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "Saturn and Venus conjunct outside Libra.")
			}
		}
	}

	// Dampatya Marakatva Rule 27 (Mars + Saturn in Lagna or 7th)
	checkMarsSat := func(signName, houseLabel string) {
		hasM, hasS := false, false
		for _, p := range planetsInSign[signName] {
			if p == "Mars" {
				hasM = true
			}
			if p == "Saturn" {
				hasS = true
			}
		}
		if hasM && hasS {
			res.DampatyaMarakatva.HasRisk = true
			res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "Mars and Saturn conjunct in "+houseLabel+".")
		}
	}
	checkMarsSat(ascSign, "Lagna")
	checkMarsSat(h7Sign, "7th house")

	// Papa Kartari Checks for 7th house, 7th lord, Venus
	if isPapaKartari(h7Sign) {
		res.DampatyaMarakatva.HasRisk = true
		res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "7th house is in Papa Kartari (hemmed between malefics).")
	}
	if isPapaKartari(h7LordSign) {
		res.DampatyaMarakatva.HasRisk = true
		res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "7th Lord is in Papa Kartari (hemmed between malefics).")
	}
	if isPapaKartari(venusSign) {
		res.DampatyaMarakatva.HasRisk = true
		res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "Venus is in Papa Kartari (hemmed between malefics).")
	}

	// Specific Lagna Marakatva Rules (14, 20-22)
	if ascSign == "Taurus" && venusSign == "Scorpio" {
		res.DampatyaMarakatva.HasRisk = true
		res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "Taurus Lagna with Venus in 7th (Scorpio).")
	}
	if ascSign == "Scorpio" {
		hasMerc7 := false
		for _, p := range planetsInSign["Taurus"] {
			if p == "Mercury" {
				hasMerc7 = true
			}
		}
		if hasMerc7 {
			res.DampatyaMarakatva.HasRisk = true
			res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "Scorpio Lagna with Mercury in 7th (Taurus).")
		}
	}
	if ascSign == "Cancer" {
		hasJup7 := false
		for _, p := range planetsInSign["Capricorn"] {
			if p == "Jupiter" {
				hasJup7 = true
			}
		}
		if hasJup7 {
			res.DampatyaMarakatva.HasRisk = true
			res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "Cancer Lagna with Jupiter in 7th (Capricorn).")
		}
	}
	if ascSign == "Virgo" && saturnSign == "Pisces" {
		hasSunLagna := false
		for _, p := range planetsInSign["Virgo"] {
			if p == "Sun" {
				hasSunLagna = true
			}
		}
		if hasSunLagna {
			res.DampatyaMarakatva.HasRisk = true
			res.DampatyaMarakatva.RiskFactors = append(res.DampatyaMarakatva.RiskFactors, "Virgo Lagna with Saturn in 7th and Sun in 1st.")
		}
	}

	// High-Impact Vaidhavya & Dwikalatra Conditions
	h10Idx := (ascIdx + 9) % 12
	h10Sign := Signs[h10Idx]
	h9Idx := (ascIdx + 8) % 12
	h9Sign := Signs[h9Idx]

	for _, p := range planetsInSign[h10Sign] {
		if p == "Saturn" {
			res.DwikalatraDosham.HasRisk = true
			res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Saturn in the 10th house.")
		}
		if p == "Mars" {
			res.DwikalatraDosham.HasRisk = true
			res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Mars in the 10th house.")
		}
	}

	h1LordHouse = getHouseNum(h1LordSign)
	if h1LordHouse == 3 || h1LordHouse == 6 {
		res.DwikalatraDosham.HasRisk = true
		res.DwikalatraDosham.RiskFactors = append(res.DwikalatraDosham.RiskFactors, "Lagna lord placed in the 3rd or 6th house.")
	}

	for _, p := range planetsInSign[h9Sign] {
		if p == "Mars" {
			res.VaidhavyaDosham.HasRisk = true
			res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Mars in the 9th house (Alpayush).")
		}
		if p == "Saturn" {
			res.VaidhavyaDosham.HasRisk = true
			res.VaidhavyaDosham.RiskFactors = append(res.VaidhavyaDosham.RiskFactors, "Saturn in the 9th house.")
		}
	}

	return res
}

func CalculateMatch(groomCtx, brideCtx *domain.CalculationContext) (domain.MatchmakingResult, error) {
	groomMoon := getMoonDetails(groomCtx.JulianDayUT, groomCtx.Ayanamsa)
	brideMoon := getMoonDetails(brideCtx.JulianDayUT, brideCtx.Ayanamsa)

	// --- Lagna Kootami (From PDF Rule 1, 2, 3) ---
	gPls, _ := planets.CalculatePlanets(groomCtx)
	bPls, _ := planets.CalculatePlanets(brideCtx)

	var gLagna, bLagna string
	for _, p := range gPls {
		if p.Planet == "Ascendant" {
			gLagna = p.Sign
			break
		}
	}
	for _, p := range bPls {
		if p.Planet == "Ascendant" {
			bLagna = p.Sign
			break
		}
	}

	gLagnaIdx := getSignIndex(gLagna)
	bLagnaIdx := getSignIndex(bLagna)
	lagnaDist := (gLagnaIdx-bLagnaIdx+12)%12 + 1

	isLagnaAusp := (lagnaDist == 1 || lagnaDist == 3 || lagnaDist == 5 || lagnaDist == 9 || lagnaDist == 11)

	gLagnaLord := SignLords[gLagna]
	bLagnaLord := SignLords[bLagna]

	lagnaRel := "Same"
	if gLagnaLord != bLagnaLord {
		// Use existing GrahaNaturalMaitri. It's a map map[string]map[string]string.
		// Groom lord from Bride lord perspective
		rel1 := GrahaNaturalMaitri[bLagnaLord][gLagnaLord]
		rel2 := GrahaNaturalMaitri[gLagnaLord][bLagnaLord]

		if rel1 == "Friend" && rel2 == "Friend" {
			lagnaRel = "Friendly"
		} else if rel1 == "Enemy" || rel2 == "Enemy" {
			lagnaRel = "Inimical"
		} else {
			lagnaRel = "Neutral"
		}
	}

	lagnaRes := domain.LagnaCompatibility{
		BrideLagna:        bLagna,
		GroomLagna:        gLagna,
		Distance:          lagnaDist,
		IsAuspicious:      isLagnaAusp,
		BrideLagnaLord:    bLagnaLord,
		GroomLagnaLord:    gLagnaLord,
		LordsRelationship: lagnaRel,
	}

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

	// Unfavorable: Janma(1), Vipat(3), Pratyari(5), Vadha(7) -> 0-indexed: 0, 2, 4, 6
	auspBtoG := !(taraBtoGIdx == 0 || taraBtoGIdx == 2 || taraBtoGIdx == 4 || taraBtoGIdx == 6)
	auspGtoB := !(taraGtoBIdx == 0 || taraGtoBIdx == 2 || taraGtoBIdx == 4 || taraGtoBIdx == 6)

	// Exception 1A: Anu Janma (2nd series: 10 to 18 distance) Pada exemptions (Page 3)
	if countBtoG >= 10 && countBtoG <= 18 {
		if taraBtoGIdx == 2 && groomMoon.Pada != 1 {
			auspBtoG = true
		}
		if taraBtoGIdx == 4 && groomMoon.Pada != 4 {
			auspBtoG = true
		}
		if taraBtoGIdx == 6 && groomMoon.Pada != 3 {
			auspBtoG = true
		}
	}
	if countGtoB >= 10 && countGtoB <= 18 {
		if taraGtoBIdx == 2 && brideMoon.Pada != 1 {
			auspGtoB = true
		}
		if taraGtoBIdx == 4 && brideMoon.Pada != 4 {
			auspGtoB = true
		}
		if taraGtoBIdx == 6 && brideMoon.Pada != 3 {
			auspGtoB = true
		}
	}

	// Exception 1B: Eka Nakshatra (Same Nakshatra)
	if countBtoG == 1 {
		if BadEkaNakshatras[brideMoon.Nakshatra] {
			auspBtoG = false
			auspGtoB = false
		} else {
			auspBtoG = true
			auspGtoB = true
		}
	}

	// Exception 2: Revati is supremely auspicious
	if brideMoon.Nakshatra == "Revati" || groomMoon.Nakshatra == "Revati" {
		auspBtoG = true
		auspGtoB = true
	}

	// Exception 3: Vipat Tara (Index 2) is cancelled if Moon signs are Sama-Saptama (1/7 axis)
	gSignIdxTemp := getSignIndex(groomMoon.Sign)
	bSignIdxTemp := getSignIndex(brideMoon.Sign)
	if (gSignIdxTemp-bSignIdxTemp+12)%12+1 == 7 {
		if taraBtoGIdx == 2 {
			auspBtoG = true
		}
		if taraGtoBIdx == 2 {
			auspGtoB = true
		}
	}

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

	// Exception for specific Rakshasa Gana brides
	if bGana == "Rakshasa" {
		if brideMoon.Nakshatra == "Ashlesha" || brideMoon.Nakshatra == "Chitra" || brideMoon.Nakshatra == "Vishakha" ||
			brideMoon.Nakshatra == "Jyeshtha" || brideMoon.Nakshatra == "Mula" || brideMoon.Nakshatra == "Shatabhisha" {
			ganaScore = 6 // Full score due to exemption
		}
	}

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
		pairStr1 := groomMoon.Sign + "-" + brideMoon.Sign
		pairStr2 := brideMoon.Sign + "-" + groomMoon.Sign

		// New Check: Telugu Book Auspicious Exceptions
		if AuspiciousBhakootPairs[pairStr1] || AuspiciousBhakootPairs[pairStr2] {
			cancellation = true
			cancelReason = "Bhakoot dosha cancelled because this specific sign combination is highly auspicious."
			effectiveBhakoot = 7
		} else if gmScore >= 4.0 { // Existing Cancellation: same lord or friendly lords
			cancellation = true
			cancelReason = "Bhakoot dosha cancelled because Rashi lords are same or friends (Graha Maitri)."
			effectiveBhakoot = 7
		}
	}

	bExpl := "Sign distance between Moons."
	if bToGDist == 2 {
		bExpl += " WARNING: Groom's Moon is 2nd from Bride's Moon (Risk of Vaidhavya/Poverty per Telugu rules)."
	} else if bToGDist == 5 {
		bExpl += " WARNING: Groom's Moon is 5th from Bride's Moon (Risk of Vaidhavya/Loss of children per Telugu rules)."
	} else if bToGDist == 3 {
		bExpl += " WARNING: Groom's Moon is 3rd from Bride's Moon (Risk of Sorrows/Hardships)."
	} else if bToGDist == 6 {
		bExpl += " WARNING: Groom's Moon is 6th from Bride's Moon (Risk of Loss of children/Quarrels)."
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
		Explanation:          bExpl,
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
		// Check Jaimini Maharshi Exemptions First
		if JaiminiNadiExceptions[groomMoon.Nakshatra] && JaiminiNadiExceptions[brideMoon.Nakshatra] {
			nCancel = true
			nReason = "Nadi dosha cancelled per Jaimini Maharshi exemptions for these specific Nakshatras."
			effNadi = 8
		} else if groomMoon.Nakshatra == brideMoon.Nakshatra && groomMoon.Pada != brideMoon.Pada {
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

	// 9. Rajju
	gRajju := NakshatraRajju[groomMoon.Nakshatra]
	bRajju := NakshatraRajju[brideMoon.Nakshatra]
	isRajjuDosha := (gRajju == bRajju)

	rCancel := false
	rReason := ""

	if isRajjuDosha {
		if gLord == bLord || gmScore >= 4.0 {
			rCancel = true
			rReason = "Rajju dosha cancelled because Rasi Lords are same or friendly."
		} else if lagnaRel == "Same" || lagnaRel == "Friendly" {
			rCancel = true
			rReason = "Rajju dosha cancelled because Lagna Lords are same or friendly."
		}
	}

	rajjuRes := domain.RajjuResult{
		GroomRajju:          gRajju,
		BrideRajju:          bRajju,
		Dosha:               isRajjuDosha,
		CancellationApplied: rCancel,
		CancellationReason:  rReason,
		Explanation:         "Body part (Rajju) matching. Same Rajju is a Dosha.",
	}

	rawTotal := varnaScore + vashyaScore + taraRes.Score + yoniScore + gmScore + ganaScore + rawBhakoot + rawNadi
	effTotal := varnaScore + vashyaScore + taraRes.Score + yoniScore + gmScore + ganaScore + effectiveBhakoot + effNadi

	// --- New Match Doshas (Telugu Book) ---

	// A. Kala Sarpa Match Dosham
	gKalaSarpa := checkKalaSarpa(groomCtx)
	bKalaSarpa := checkKalaSarpa(brideCtx)
	hasKalaSarpaMatch := (gKalaSarpa && bKalaSarpa)

	kalaSarpaRes := domain.DoshaCheck{
		HasDosha:    hasKalaSarpaMatch,
		Explanation: "If both Bride and Groom have Kala Sarpa Dosha (all planets between Rahu and Ketu), the match is strictly rejected.",
	}

	// B. Mruga Vairam (Animal Enmity)
	hasMrugaVairam := (MrugaEnmity[gYoni] == bYoni)
	mrugaVairamRes := domain.DoshaCheck{
		HasDosha:    hasMrugaVairam,
		Explanation: "Checks if the Yoni (Animal) of the couple are natural enemies (e.g. Cow-Tiger, Snake-Mongoose).",
	}

	// C. Special Tara Rejections
	hasSpecialTara := false
	// Rule from Page 2
	if brideMoon.Nakshatra == "Krittika" && groomMoon.Nakshatra == "Bharani" {
		hasSpecialTara = true
	} else if brideMoon.Nakshatra == "Ashlesha" && groomMoon.Nakshatra == "Pushya" {
		hasSpecialTara = true
	} else if brideMoon.Nakshatra == "Shatabhisha" && groomMoon.Nakshatra == "Dhanishta" {
		hasSpecialTara = true
	}
	// Rule from Page 3 (Ashubha Kurpu - 15 Pairs)
	if AshubhaKurpu[brideMoon.Nakshatra] == groomMoon.Nakshatra {
		hasSpecialTara = true
	}

	specialTaraRes := domain.DoshaCheck{
		HasDosha:    hasSpecialTara,
		Explanation: "Checks for specifically rejected Tara combinations (e.g., Bride Krittika + Groom Bharani).",
	}

	matchDoshas := domain.MatchDoshas{
		KalaSarpaMatch: kalaSarpaRes,
		MrugaVairam:    mrugaVairamRes,
		SpecialTara:    specialTaraRes,
	}

	res := domain.MatchmakingResult{
		RuleSet:    domain.RuleSetInfo{Name: "classical-guna-milan-v1", Ayanamsa: groomCtx.Input.Ayanamsa},
		Groom:      domain.PersonMatchDetails{Name: groomCtx.Input.Name, Moon: groomMoon, Risks: evaluateIndividualRisks(groomCtx, true)},
		Bride:      domain.PersonMatchDetails{Name: brideCtx.Input.Name, Moon: brideMoon, Risks: evaluateIndividualRisks(brideCtx, false)},
		LagnaMatch: lagnaRes,
		Kootas: domain.Kootas{
			Varna:       varnaRes,
			Vashya:      vashyaRes,
			Tara:        taraRes,
			Yoni:        yoniRes,
			GrahaMaitri: gmRes,
			Gana:        ganaRes,
			Bhakoot:     bhakootRes,
			Nadi:        nadiRes,
			Rajju:       rajjuRes,
		},
		Doshas: matchDoshas,
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
