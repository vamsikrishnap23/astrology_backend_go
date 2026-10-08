package mrityu

import "math"

// IsMrityuBhaga checks if a planet in a given sign is at its Mrityu Bhaga (Fatal Degree).
// We use a standard +/- 1.0 degree orb.
func IsMrityuBhaga(planet, sign string, degreeInSign float64) bool {
	// Map of Planet -> Map of Sign -> Exact Degree
	mbDegrees := map[string]map[string]float64{
		"Sun": {
			"Aries": 20, "Taurus": 9, "Gemini": 12, "Cancer": 6, "Leo": 8, "Virgo": 24,
			"Libra": 16, "Scorpio": 17, "Sagittarius": 22, "Capricorn": 2, "Aquarius": 3, "Pisces": 23,
		},
		"Moon": {
			"Aries": 26, "Taurus": 12, "Gemini": 13, "Cancer": 25, "Leo": 25, "Virgo": 11,
			"Libra": 26, "Scorpio": 14, "Sagittarius": 13, "Capricorn": 25, "Aquarius": 5, "Pisces": 12,
		},
		"Mars": {
			"Aries": 19, "Taurus": 28, "Gemini": 26, "Cancer": 11, "Leo": 16, "Virgo": 14,
			"Libra": 22, "Scorpio": 10, "Sagittarius": 18, "Capricorn": 20, "Aquarius": 8, "Pisces": 18,
		},
		"Mercury": {
			"Aries": 15, "Taurus": 12, "Gemini": 14, "Cancer": 13, "Leo": 12, "Virgo": 21,
			"Libra": 10, "Scorpio": 18, "Sagittarius": 8, "Capricorn": 20, "Aquarius": 18, "Pisces": 21,
		},
		"Jupiter": {
			"Aries": 10, "Taurus": 13, "Gemini": 16, "Cancer": 17, "Leo": 10, "Virgo": 13,
			"Libra": 15, "Scorpio": 17, "Sagittarius": 12, "Capricorn": 15, "Aquarius": 14, "Pisces": 10,
		},
		"Venus": {
			"Aries": 2, "Taurus": 5, "Gemini": 11, "Cancer": 5, "Leo": 8, "Virgo": 5,
			"Libra": 21, "Scorpio": 5, "Sagittarius": 22, "Capricorn": 11, "Aquarius": 7, "Pisces": 17,
		},
		"Saturn": {
			"Aries": 10, "Taurus": 14, "Gemini": 13, "Cancer": 12, "Leo": 15, "Virgo": 15,
			"Libra": 18, "Scorpio": 8, "Sagittarius": 18, "Capricorn": 19, "Aquarius": 19, "Pisces": 21,
		},
		"Rahu": {
			"Aries": 14, "Taurus": 13, "Gemini": 12, "Cancer": 11, "Leo": 10, "Virgo": 9,
			"Libra": 8, "Scorpio": 7, "Sagittarius": 6, "Capricorn": 5, "Aquarius": 4, "Pisces": 3,
		},
		"Ketu": {
			"Aries": 18, "Taurus": 17, "Gemini": 16, "Cancer": 15, "Leo": 14, "Virgo": 13,
			"Libra": 12, "Scorpio": 11, "Sagittarius": 10, "Capricorn": 9, "Aquarius": 8, "Pisces": 7,
		},
	}

	planetMap, exists := mbDegrees[planet]
	if !exists {
		return false // Planet not in standard MB tables (such as Uranus)
	}

	targetDegree, exists := planetMap[sign]
	if !exists {
		return false
	}

	// Apply 1 degree orb
	return math.Abs(degreeInSign-targetDegree) <= 1.0
}
