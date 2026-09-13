package manglik

import (
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"math"
)

var signNames = []string{
	"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
	"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

func getSignIndex(lon float64) int {
	return int(math.Floor(lon / 30.0))
}

func getHouseFromRef(planetLon, refLon float64) int {
	pSign := getSignIndex(planetLon)
	rSign := getSignIndex(refLon)
	house := (pSign-rSign+12)%12 + 1
	return house
}

func isManglikHouse(h int) bool {
	return h == 1 || h == 2 || h == 4 || h == 7 || h == 8 || h == 12
}

func CalculateManglikDosha(planets []domain.PlanetPosition, ascLon float64) domain.ManglikResult {
	var mars, moon, venus, jupiter domain.PlanetPosition
	for _, p := range planets {
		switch p.Planet {
		case "Mars":
			mars = p
		case "Moon":
			moon = p
		case "Venus":
			venus = p
		case "Jupiter":
			jupiter = p
		}
	}

	res := domain.ManglikResult{
		BaseAnalysis:  make(map[string]domain.ManglikAnalysis),
		Cancellations: []domain.ManglikCancellation{},
		Remedies:      []string{},
	}

	ascHouse := getHouseFromRef(mars.SiderealLongitude, ascLon)
	moonHouse := getHouseFromRef(mars.SiderealLongitude, moon.SiderealLongitude)
	venusHouse := getHouseFromRef(mars.SiderealLongitude, venus.SiderealLongitude)

	hasAsc := isManglikHouse(ascHouse)
	hasMoon := isManglikHouse(moonHouse)
	hasVenus := isManglikHouse(venusHouse)

	res.BaseAnalysis["from_ascendant"] = domain.ManglikAnalysis{IsPresent: hasAsc, House: ascHouse}
	res.BaseAnalysis["from_moon"] = domain.ManglikAnalysis{IsPresent: hasMoon, House: moonHouse}
	res.BaseAnalysis["from_venus"] = domain.ManglikAnalysis{IsPresent: hasVenus, House: venusHouse}

	marsSignIdx := getSignIndex(mars.SiderealLongitude)
	marsSign := signNames[marsSignIdx]
	ascSignIdx := getSignIndex(ascLon)

	isCancelled := false

	// Exception 1: Benefic Ascendants (Leo, Cancer)
	if ascSignIdx == 3 || ascSignIdx == 4 { // Cancer or Leo
		res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is a Yogakaraka for Cancer/Leo Ascendant, nullifying the dosha."})
		isCancelled = true
	}

	// Exception 2: Own Sign, Exalted, Debilitated
	if marsSign == "Aries" || marsSign == "Scorpio" {
		res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is placed in its own sign (" + marsSign + ")."})
		isCancelled = true
	} else if marsSign == "Capricorn" {
		res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is Exalted in Capricorn."})
		isCancelled = true
	} else if marsSign == "Cancer" {
		res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is Debilitated in Cancer."})
		isCancelled = true
	}

	// Exception 3: Conjunction with Jupiter or Moon
	if getHouseFromRef(mars.SiderealLongitude, jupiter.SiderealLongitude) == 1 {
		res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is conjunct with Jupiter."})
		isCancelled = true
	}
	if getHouseFromRef(mars.SiderealLongitude, moon.SiderealLongitude) == 1 {
		res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is conjunct with the Moon (Chandra-Mangala Yoga)."})
		isCancelled = true
	}

	// Exception 4: Specific House/Sign rules (From Ascendant)
	if hasAsc {
		if ascHouse == 2 && (marsSign == "Gemini" || marsSign == "Virgo") {
			res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is in the 2nd house in a Mercury sign (" + marsSign + ")."})
			isCancelled = true
		} else if ascHouse == 4 && (marsSign == "Aries" || marsSign == "Scorpio") {
			res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is in the 4th house in its own sign (" + marsSign + ")."})
			isCancelled = true
		} else if ascHouse == 7 && (marsSign == "Cancer" || marsSign == "Capricorn") {
			res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is in the 7th house in " + marsSign + "."})
			isCancelled = true
		} else if ascHouse == 8 && (marsSign == "Sagittarius" || marsSign == "Pisces") {
			res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is in the 8th house in a Jupiter sign (" + marsSign + ")."})
			isCancelled = true
		} else if ascHouse == 12 && (marsSign == "Taurus" || marsSign == "Libra") {
			res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is in the 12th house in a Venus sign (" + marsSign + ")."})
			isCancelled = true
		}
	}

	// Exception 5: Nakshatra Bhanga
	exemptNakshatras := map[string]bool{
		"Ashwini": true, "Mrigashira": true, "Punarvasu": true, "Pushya": true,
		"Ashlesha": true, "Uttara Phalguni": true, "Swati": true, "Anuradha": true,
		"Purva Ashadha": true, "Uttara Ashadha": true, "Shravana": true,
		"Uttara Bhadrapada": true, "Revati": true,
	}
	if exemptNakshatras[mars.Nakshatra] {
		res.Cancellations = append(res.Cancellations, domain.ManglikCancellation{Rule: "Mars is in an exempt Nakshatra (" + mars.Nakshatra + "), completely nullifying the dosha."})
		isCancelled = true
	}

	if !hasAsc && !hasMoon && !hasVenus {
		res.IsManglik = false
		res.Status = "Non-Manglik"
	} else if isCancelled {
		res.IsManglik = false
		res.Status = "Cancelled Manglik"
	} else if !hasAsc && (hasMoon || hasVenus) {
		res.IsManglik = true
		res.Status = "Partial Manglik"
	} else {
		res.IsManglik = true
		res.Status = "High Manglik"
	}

	if res.Status == "High Manglik" || res.Status == "Partial Manglik" {
		res.Remedies = []string{
			"Perform Kumbh Vivah (symbolic marriage) before actual marriage.",
			"Worship Lord Hanuman and recite Hanuman Chalisa daily.",
			"Manglik should ideally marry another Manglik.",
			"Chant the Mangal Mantra ('Om Kraam Kreem Kroum Sah Bhaumaya Namah').",
			"Donate red lentils (masoor dal) or red cloth on Tuesdays.",
		}
	}

	return res
}
