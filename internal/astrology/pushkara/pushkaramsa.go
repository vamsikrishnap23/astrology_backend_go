package pushkara

// IsPushkaramsa determines if a planet is in a Pushkara Navamsha based on its sign and degree.
func IsPushkaramsa(sign string, degreeInSign float64) bool {
	// 3 degrees 20 minutes is exactly 3.3333333333333335 degrees
	// Helper to check if degree is within a range [start, end)
	inRange := func(d, start, end float64) bool {
		return d >= start && d < end
	}

	switch sign {
	case "Aries", "Leo", "Sagittarius": // Fire Signs
		// 7th Navamsha: 20° 00' to 23° 20' (20.0 to 23.333333...)
		if inRange(degreeInSign, 20.0, 23.333333333333336) {
			return true
		}
		// 9th Navamsha: 26° 40' to 30° 00' (26.666666... to 30.0)
		if inRange(degreeInSign, 26.666666666666668, 30.0) {
			return true
		}
	case "Taurus", "Virgo", "Capricorn": // Earth Signs
		// 3rd Navamsha: 06° 40' to 10° 00' (6.666666... to 10.0)
		if inRange(degreeInSign, 6.666666666666667, 10.0) {
			return true
		}
		// 5th Navamsha: 13° 20' to 16° 40' (13.333333... to 16.666666...)
		if inRange(degreeInSign, 13.333333333333334, 16.666666666666668) {
			return true
		}
	case "Gemini", "Libra", "Aquarius": // Air Signs
		// 6th Navamsha: 16° 40' to 20° 00' (16.666666... to 20.0)
		if inRange(degreeInSign, 16.666666666666668, 20.0) {
			return true
		}
		// 8th Navamsha: 23° 20' to 26° 40' (23.333333... to 26.666666...)
		if inRange(degreeInSign, 23.333333333333336, 26.666666666666668) {
			return true
		}
	case "Cancer", "Scorpio", "Pisces": // Water Signs
		// 1st Navamsha: 00° 00' to 03° 20' (0.0 to 3.333333...)
		if inRange(degreeInSign, 0.0, 3.3333333333333335) {
			return true
		}
		// 3rd Navamsha: 06° 40' to 10° 00' (6.666666... to 10.0)
		if inRange(degreeInSign, 6.666666666666667, 10.0) {
			return true
		}
	}

	return false
}

// IsPushkaraBhaga determines if a planet is at the exact Pushkara Bhaga degree (within a 1-degree orb).
func IsPushkaraBhaga(sign string, degreeInSign float64) bool {
	var targetDegree float64

	switch sign {
	case "Aries", "Leo", "Sagittarius":
		targetDegree = 21.0
	case "Taurus", "Virgo", "Capricorn":
		targetDegree = 14.0
	case "Gemini", "Libra", "Aquarius":
		targetDegree = 24.0
	case "Cancer", "Scorpio", "Pisces":
		targetDegree = 7.0
	default:
		return false
	}

	// We apply a strict 1-degree orb (meaning +/- 1.0 degree from the exact target)
	// Example: For Fire (21.0), anything between 20.0 and 22.0 will flag as true.
	diff := degreeInSign - targetDegree
	if diff < 0 {
		diff = -diff
	}

	return diff <= 1.0
}
