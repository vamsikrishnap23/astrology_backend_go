package matchmaking

import (
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"math"

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

	ascLong, _, _, _ := houses.CalculateHouses(ctx)
	ascSign, _, _, _ := time.DecimalToDMS(ascLong)

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

// CalculateMatch computes the full matchmaking result for a groom and a bride.
func CalculateMatch(groomCtx, brideCtx *domain.CalculationContext) (domain.MatchmakingResult, error) {
	groom := getMoonDetails(groomCtx.JulianDayUT, groomCtx.Ayanamsa)
	bride := getMoonDetails(brideCtx.JulianDayUT, brideCtx.Ayanamsa)

	lagna := lagnaCompatibility(groomCtx, brideCtx)

	// the eight scored kootas
	var k domain.Kootas
	k.Varna = varnaKoota(groom, bride)
	k.Vashya = vashyaKoota(groom.Sign, bride.Sign)
	k.Tara = taraKoota(groom, bride)
	k.Yoni = yoniKoota(groom, bride)
	k.GrahaMaitri = grahaMaitriKoota(groom, bride)
	k.Gana = ganaKoota(groom, bride)
	k.Bhakoot = bhakootKoota(groom, bride, k.GrahaMaitri.Score)
	k.Nadi = nadiKoota(groom, bride)

	// the other checks (not scored)
	count := nakshatraCount(bride.Nakshatra, groom.Nakshatra)
	bhakootGood := k.Bhakoot.Score >= 7
	checks := domain.Checks{
		Rajju:        rajjuCheck(groom, bride, k.GrahaMaitri.Score, lagna.LordsRelationship),
		Vedha:        vedhaCheck(groom, bride),
		Mahendra:     mahendraCheck(groom, bride, count),
		StreeDeergha: streeDeerghaCheck(groom, bride, count, bhakootGood, k.GrahaMaitri.Score),
		MrugaVairam:  mrugaVairamCheck(groom, bride),
		SpecialTara:  specialTaraCheck(groom, bride),
		KalaSarpa:    kalaSarpaCheck(groomCtx, brideCtx),
	}

	return domain.MatchmakingResult{
		RuleSet:    domain.RuleSetInfo{Name: "classical-guna-milan-v2", Ayanamsa: groomCtx.Input.Ayanamsa},
		Groom:      domain.PersonMatchDetails{Name: groomCtx.Input.Name, Moon: groom, Risks: evaluateIndividualRisks(groomCtx, true)},
		Bride:      domain.PersonMatchDetails{Name: brideCtx.Input.Name, Moon: bride, Risks: evaluateIndividualRisks(brideCtx, false)},
		Lagna:      lagna,
		Ashtakoota: domain.Ashtakoota{Kootas: k, Summary: summarize(k)},
		Checks:     checks,
	}, nil
}
