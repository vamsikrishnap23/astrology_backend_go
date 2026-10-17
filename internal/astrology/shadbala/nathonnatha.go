package shadbala

import (
	"math"

	"github.com/tejzpr/go-swisseph"
)

// localApparentHours returns the local apparent (true solar) time of day in hours, 0..24, for a UT
// Julian day and a geographic longitude (east positive). Hindu astrology measures Nathonnatha Bala
// from apparent noon and midnight (BPHS 27.8-9, Raman art. 48), so the equation of time is applied
// to local mean time. Must be called under ephemeris.Lock.
func localApparentHours(jdUT, lonDeg float64) (float64, error) {
	eot, err := swisseph.TimeEqu(jdUT) // days; apparent minus mean time
	if err != nil {
		return 0, err
	}
	jdLAT := jdUT + lonDeg/360.0 + eot
	h := math.Mod((jdLAT+0.5-math.Floor(jdLAT+0.5))*24.0, 24.0)
	if h < 0 {
		h += 24.0
	}
	return h, nil
}

// nathonnathaBala is the strength from the time of day (Raman art. 47-51, method B; BPHS 27.8-9).
// The Moon, Mars and Saturn are strongest at midnight and powerless at noon; the Sun, Jupiter and
// Venus the reverse; Mercury always has 60. apparentHours is local apparent time (0..24).
func nathonnathaBala(planet string, apparentHours float64) float64 {
	// Angle from midnight, reflected into 0..180 (180 = noon).
	deg := math.Mod(apparentHours*15.0, 360.0)
	if deg < 0 {
		deg += 360.0
	}
	if deg > 180.0 {
		deg = 360.0 - deg
	}
	day := deg / 3.0 // 0 at midnight, 60 at noon
	switch planet {
	case "Sun", "Jupiter", "Venus":
		return day
	case "Moon", "Mars", "Saturn":
		return 60.0 - day
	case "Mercury":
		return 60.0
	}
	return 0
}
