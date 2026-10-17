package shadbala

import "math"

// Drik Bala (aspect strength), Raman chapter VIII (art. 109-120, pp. 82-88), after Sripati; BPHS 27.19.

// drishtiBase is the aspect value, in shashtiamsas, for a Drishti Kendra dk: the longitude of the
// aspected body minus the longitude of the aspecting planet, 0..360 (art. 113). It rises from 0 at
// 30 deg to 15 at 60, 45 at 90, falls to 30 at 120 and to 0 at 150, jumps up to 60 at 180 and falls to
// 0 at 300 (art. 110, 114).
func drishtiBase(dk float64) float64 {
	dk = math.Mod(dk, 360.0)
	if dk < 0 {
		dk += 360.0
	}
	switch {
	case dk < 30.0 || dk >= 300.0:
		return 0
	case dk < 60.0:
		return (dk - 30.0) / 2.0
	case dk < 90.0:
		return (dk - 60.0) + 15.0
	case dk < 120.0:
		return (120.0-dk)/2.0 + 30.0
	case dk < 150.0:
		return 150.0 - dk
	case dk < 180.0:
		return (dk - 150.0) * 2.0
	default: // 180..300
		return (300.0 - dk) / 2.0
	}
}

// visheshaDrishti is the extra aspect value added to the ordinary one for the special aspects
// (art. 115): Saturn on the 3rd and 10th (60-90 and 270-300 deg) adds 45, Jupiter on the 5th and 9th
// (120-150 and 240-270) adds 30, Mars on the 4th and 8th (90-120 and 210-240) adds 15.
func visheshaDrishti(planet string, dk float64) float64 {
	dk = math.Mod(dk, 360.0)
	if dk < 0 {
		dk += 360.0
	}
	in := func(lo, hi float64) bool { return dk >= lo && dk < hi }
	switch planet {
	case "Saturn":
		if in(60, 90) || in(270, 300) {
			return 45
		}
	case "Jupiter":
		if in(120, 150) || in(240, 270) {
			return 30
		}
	case "Mars":
		if in(90, 120) || in(210, 240) {
			return 15
		}
	}
	return 0
}

// drishtiValue is the full aspect value that planet casts at a Drishti Kendra dk.
func drishtiValue(planet string, dk float64) float64 {
	return drishtiBase(dk) + visheshaDrishti(planet, dk)
}

// moonIsBenefic reports whether the Moon is a benefic (Vriddhi Chandra) for aspects: "he is a benefic
// from the 8th day of the bright half of the lunar month to the 8th day of the dark half of the lunar
// month and a malefic in the rest of the days" (Raman art. 53). The 8th tithi of the bright half starts
// at 84 deg of Moon-Sun elongation and the 8th of the dark half at 264 deg.
func moonIsBenefic(sunLon, moonLon float64) bool {
	e := math.Mod(moonLon-sunLon+720.0, 360.0)
	return e >= 84.0 && e < 264.0
}

// beneficMap says which of the seven planets cast a benefic (+) or malefic (-) aspect (art. 116-118):
// Jupiter and Venus are benefic; the Sun, Mars and Saturn malefic; the Moon follows moonIsBenefic;
// Mercury is a malefic when it is afflicted (see mercuryIsMalefic) and a benefic otherwise.
func beneficMap(lons map[string]float64, signs map[string]int) map[string]bool {
	return map[string]bool{
		"Sun":     false,
		"Mars":    false,
		"Saturn":  false,
		"Jupiter": true,
		"Venus":   true,
		"Moon":    moonIsBenefic(lons["Sun"], lons["Moon"]),
		"Mercury": !mercuryIsMalefic(signs),
	}
}

// drishtiPinda is the net aspect on planet: the sum of the aspect values of the benefics on it minus
// those of the malefics (art. 119). Only the seven planets Sun to Saturn aspect.
func drishtiPinda(planet string, lons map[string]float64, benefic map[string]bool) float64 {
	target, ok := lons[planet]
	if !ok {
		return 0
	}
	total := 0.0
	for _, aspecting := range TargetPlanets {
		if aspecting == planet {
			continue
		}
		from, ok := lons[aspecting]
		if !ok {
			continue
		}
		v := drishtiValue(aspecting, math.Mod(target-from+720.0, 360.0))
		if benefic[aspecting] {
			total += v
		} else {
			total -= v
		}
	}
	return total
}

// CalculateDrikBala is one fourth of the Drishti Pinda, positive or negative (art. 120). It is added
// to or subtracted from the other balas in the Shadbala total (art. 121).
func CalculateDrikBala(planet string, lons map[string]float64, benefic map[string]bool) float64 {
	return drishtiPinda(planet, lons, benefic) / 4.0
}
