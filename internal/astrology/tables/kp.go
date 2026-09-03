package tables

import (
	"math"
)

var vimshottariLords = []string{"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury"}
var vimshottariYears = []float64{7, 20, 6, 10, 7, 18, 16, 19, 17}

func indexOfLord(lord string) int {
	for i, l := range vimshottariLords {
		if l == lord {
			return i
		}
	}
	return 0
}

// GetKPLords returns the Sign Lord, Nakshatra Lord, Sub Lord, Sub-Sub Lord, and Sub-Sub-Sub Lord
// using the precise proportional Vimshottari mathematical subdivision method.
func GetKPLords(longitude float64) (sl, nl, ssl, sssl, ssssl string) {
	longitude = math.Mod(longitude, 360.0)
	if longitude < 0 {
		longitude += 360.0
	}

	// 1. Sign Lord (SL)
	signs := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
	signIdx := int(math.Floor(longitude/30.0 + 1e-12))
	if signIdx >= 12 {
		signIdx = 0
	}
	sign := signs[signIdx]
	sl = signLords[sign]

	// 2. Nakshatra Lord (NL)
	nakshatraLen := 13.0 + 1.0/3.0
	nakIdx := int(math.Floor(longitude/nakshatraLen + 1e-12))
	nl = vimshottariLords[nakIdx%9]

	rem := longitude - float64(nakIdx)*nakshatraLen
	if rem < 0 {
		rem = 0 // Just in case epsilon pushed it over slightly early
	}

	// Helper for subdivision using exact cumulative proportions to avoid floating point subtraction accumulation
	findSubLord := func(startLord string, totalSpan float64, currentRem float64) (string, float64, float64) {
		idx := indexOfLord(startLord)
		fraction := currentRem / totalSpan
		cumulative := 0.0

		for i := 0; i < 9; i++ {
			years := vimshottariYears[idx]
			cumulative += years / 120.0

			// If the fraction falls into this lord's proportion segment
			if fraction < cumulative-1e-12 {
				// Calculate the precise remainder inside this specific lord's segment
				prevCumulative := cumulative - (years / 120.0)
				newRem := (fraction - prevCumulative) * totalSpan
				lordSpan := (years / 120.0) * totalSpan
				return vimshottariLords[idx], newRem, lordSpan
			}
			idx = (idx + 1) % 9
		}

		// Fallback for extreme precision edge cases directly at the boundary
		lastIdx := (indexOfLord(startLord) + 8) % 9
		return vimshottariLords[lastIdx], 0, (vimshottariYears[lastIdx] / 120.0) * totalSpan
	}

	// 3. Sub Lord (SSL)
	ssl, rem, sslSpan := findSubLord(nl, nakshatraLen, rem)

	// 4. Sub-Sub Lord (SSSL)
	sssl, rem, ssslSpan := findSubLord(ssl, sslSpan, rem)

	// 5. Sub-Sub-Sub Lord (SSSSL)
	ssssl, _, _ = findSubLord(sssl, ssslSpan, rem)

	return sl, nl, ssl, sssl, ssssl
}
