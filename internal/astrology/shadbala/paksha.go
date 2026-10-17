package shadbala

import "math"

// mercuryIsMalefic reports whether Mercury counts as a malefic for Paksha Bala: it does when it
// joins a malefic (BPHS 3.11; Raman art. 53 "afflicted Budha"), i.e. is in the same sign as the
// Sun, Mars, Saturn, Rahu or Ketu. signs maps planet name to sidereal sign index 0..11.
func mercuryIsMalefic(signs map[string]int) bool {
	m, ok := signs["Mercury"]
	if !ok {
		return false
	}
	for _, p := range []string{"Sun", "Mars", "Saturn", "Rahu", "Ketu"} {
		if s, ok := signs[p]; ok && s == m {
			return true
		}
	}
	return false
}

// pakshaBala is Paksha Bala in shashtiamsas (Raman art. 55; BPHS 27.10-11). The arc between the
// Sun and the Moon (reflected into 0..180) divided by 3 is the strength of every benefic; a malefic
// gets 60 minus that. The Moon's value is doubled. The Moon takes the benefic value, as in Raman's
// table (Chandra is "Subha" and the value is "always to be doubled"); Jupiter and Venus are
// benefic; Sun, Mars and Saturn are malefic; Mercury follows mercuryMalefic.
func pakshaBala(planet string, sunLon, moonLon float64, mercuryMalefic bool) float64 {
	arc := math.Mod(moonLon-sunLon+720.0, 360.0)
	if arc > 180.0 {
		arc = 360.0 - arc
	}
	benefic := arc / 3.0
	var v float64
	switch planet {
	case "Moon":
		return benefic * 2
	case "Jupiter", "Venus":
		v = benefic
	case "Mercury":
		if mercuryMalefic {
			v = 60.0 - benefic
		} else {
			v = benefic
		}
	case "Sun", "Mars", "Saturn":
		v = 60.0 - benefic
	}
	return v
}
