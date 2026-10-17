package shadbala

import (
	"math"
	"sort"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// Bhava Bala (house strength), Raman chapter IX (art. 124-134, pp. 91-100). The strength of a house has
// three parts: Bhavadhipathi Bala (the Shadbala Pinda of the house lord), Bhava Digbala (from the kind of
// sign the middle of the house falls in) and Bhava Drishti Bala (the aspects on the middle of the house).
// BPHS 27.26-31 adds the occupants of the house and the rising of the sign (+/- 1 rupa for the occupying
// planets, 15 virupas by day, twilight or night); Raman does not, so they are not included.

var bhavaNames = []string{"Thanu", "Dhana", "Bhratru", "Matru", "Putra", "Satru", "Kalatra", "Ayur", "Bhagya", "Karma", "Labha", "Vyaya"}

var zodiacSigns = []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}

// signRulers are the planets ruling the signs from Aries (Raman's Manual art. 29: the traditional lords).
var signRulers = []string{"Mars", "Venus", "Mercury", "Moon", "Sun", "Mercury", "Venus", "Mars", "Jupiter", "Saturn", "Saturn", "Jupiter"}

// sripatiBhavaMadhyas returns the middle of each of the twelve houses by Sripati's method (Raman's Manual
// art. 88, used throughout the book): the Ascendant, the 4th house (the lower meridian, 180 deg from the
// Midheaven), the 7th (opposite the Ascendant) and the 10th (the Midheaven) are the middles of those houses,
// and the arc between each two is divided in three for the houses between them. asc and mc are longitudes.
func sripatiBhavaMadhyas(asc, mc float64) [12]float64 {
	kendras := [4]float64{norm360(asc), norm360(mc + 180.0), norm360(asc + 180.0), norm360(mc)} // houses 1, 4, 7, 10
	var m [12]float64
	for i := 0; i < 4; i++ {
		from, to := kendras[i], kendras[(i+1)%4]
		arc := norm360(to - from)
		m[3*i] = from
		m[3*i+1] = norm360(from + arc/3.0)
		m[3*i+2] = norm360(from + 2.0*arc/3.0)
	}
	return m
}

// bhavaSignGroup classifies the sign of a bhava madhya (art. 127-130): Nara (human) signs are Gemini, Virgo,
// Libra, Aquarius and the first half of Sagittarius; Jalachara (watery) Cancer, Pisces and the second half
// of Capricorn; Chathushpada (quadruped) Aries, Taurus, Leo, the second half of Sagittarius and the first
// half of Capricorn; Keeta (insect) Scorpio alone.
func bhavaSignGroup(lon float64) string {
	lon = norm360(lon)
	sign := int(lon / 30.0)
	firstHalf := math.Mod(lon, 30.0) < 15.0
	switch sign {
	case 7:
		return "keeta"
	case 0, 1, 4:
		return "chathushpada"
	case 2, 5, 6, 10:
		return "nara"
	case 3, 11:
		return "jalachara"
	case 8:
		if firstHalf {
			return "nara"
		}
		return "chathushpada"
	}
	if firstHalf { // Capricorn
		return "chathushpada"
	}
	return "jalachara"
}

// bhavaDigbala is Bhava Digbala in shashtiamsas (art. 131): a house is strongest (60) when its middle falls in
// a Nara sign for the 1st house, a Jalachara sign for the 4th, a Keeta sign for the 7th and a Chathushpada sign
// for the 10th, and weakest (0) in the opposite house; the value changes by 10 per house. As in the book,
// the number of the house is taken from 1 (Keeta), 4 (Chathushpada), 7 (Nara) or 10 (Jalachara), the
// difference is folded back if it exceeds 6, and multiplied by 10.
func bhavaDigbala(house int, group string) float64 {
	base := map[string]int{"keeta": 1, "chathushpada": 4, "nara": 7, "jalachara": 10}[group]
	d := ((base-house)%12 + 12) % 12
	if d > 6 {
		d = 12 - d
	}
	return float64(d) * 10.0
}

// bhavaDrishti returns the benefic (subha) and malefic (asubha) aspect strength on a bhava madhya (art. 132):
// the aspect value of each planet at the Drishti Kendra (madhya minus planet) with its special aspect, one
// fourth of it for every planet except Jupiter and Mercury, whose whole value is taken. Mercury counts as a
// full benefic here (the footnote on p. 98); the Moon is benefic by the rule of art. 53.
func bhavaDrishti(madhya float64, lons map[string]float64, benefic map[string]bool) (subha, asubha float64) {
	for _, p := range TargetPlanets {
		from, ok := lons[p]
		if !ok {
			continue
		}
		v := drishtiValue(p, norm360(madhya-from))
		if p != "Jupiter" && p != "Mercury" {
			v /= 4.0
		}
		isBenefic := benefic[p]
		if p == "Mercury" {
			isBenefic = true
		}
		if isBenefic {
			subha += v
		} else {
			asubha += v
		}
	}
	return
}

// calculateBhavaBala returns the strength of the twelve houses. totals holds each planet's Shadbala Pinda in
// shashtiamsas; the result is ranked by the total (1 = strongest).
func calculateBhavaBala(asc, mc float64, lons map[string]float64, benefic map[string]bool, totals map[string]float64) domain.BhavaBalaResult {
	madhyas := sripatiBhavaMadhyas(asc, mc)
	res := domain.BhavaBalaResult{Bhavas: make([]domain.BhavaBalaHouse, 12)}
	for i, m := range madhyas {
		house := i + 1
		sign := int(m / 30.0)
		group := bhavaSignGroup(m)
		lord := signRulers[sign]
		subha, asubha := bhavaDrishti(m, lons, benefic)
		lordBala := totals[lord]
		dig := bhavaDigbala(house, group)
		total := lordBala + dig + subha - asubha
		res.Bhavas[i] = domain.BhavaBalaHouse{
			House: house, Name: bhavaNames[i], BhavaMadhya: round2(m), Sign: zodiacSigns[sign], Lord: lord, SignGroup: group,
			BhavadhipathiBala: round2(lordBala), BhavaDigbala: dig,
			SubhaDrishtiBala: round2(subha), AsubhaDrishtiBala: round2(asubha), BhavaDrishtiBala: round2(subha - asubha),
			TotalBhavaBala: round2(total), Rupas: math.Round(total/60.0*100) / 100,
		}
	}
	idx := make([]int, 12)
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(a, b int) bool { return res.Bhavas[idx[a]].TotalBhavaBala > res.Bhavas[idx[b]].TotalBhavaBala })
	for rank, i := range idx {
		res.Bhavas[i].Rank = rank + 1
	}
	res.StrongestHouse = idx[0] + 1
	res.WeakestHouse = idx[11] + 1
	return res
}
