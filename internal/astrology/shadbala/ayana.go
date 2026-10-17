package shadbala

import "math"

// ayanaBala is Ayana Bala in shashtiamsas (Raman art. 75, after Keshava Daivajna):
// (24 deg + kranti) / 48 x 60, where a northern kranti is added for the Sun, Mars, Jupiter and Venus
// and a southern one for the Moon and Saturn (the opposite hemisphere subtracts). For Mercury the
// kranti is always added, north or south. The Sun's value is doubled. Hindu astrology takes 24 deg
// as the greatest kranti, so the result is limited to 0..60 (120 for the Sun).
// declination is the planet's declination in degrees, north positive.
func ayanaBala(planet string, declination float64) float64 {
	var k float64
	switch planet {
	case "Sun", "Mars", "Jupiter", "Venus":
		k = declination
	case "Moon", "Saturn":
		k = -declination
	case "Mercury":
		k = math.Abs(declination)
	default:
		return 0
	}
	v := (24.0 + k) * (60.0 / 48.0)
	v = math.Max(0, math.Min(60, v))
	if planet == "Sun" {
		v *= 2
	}
	return v
}
