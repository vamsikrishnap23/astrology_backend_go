package shadbala

import (
	"math"

	"github.com/tejzpr/go-swisseph"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
)

// WeekdayLords is indexed 0 = Sunday.
var WeekdayLords = []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"}

// horaOrder is the order in which the horas follow one another: each hora is ruled by the planet
// that comes next in descending order of orbit (Raman art. 69; BPHS 27.13).
var horaOrder = []string{"Sun", "Venus", "Mercury", "Moon", "Saturn", "Jupiter", "Mars"}

// horaLord returns the ruler of the hora that is running hoursSinceSunrise (60-minute horas, the
// first one ruled by the weekday lord, Raman art. 68-70 / BPHS 27.13).
func horaLord(dayLord string, hoursSinceSunrise float64) string {
	start := 0
	for i, p := range horaOrder {
		if p == dayLord {
			start = i
			break
		}
	}
	n := int(math.Floor(hoursSinceSunrise))
	return horaOrder[((start+n)%7+7)%7]
}

// lastSunrise returns the Julian day (UT) of the latest sunrise at or before jd for the given
// geographic position (lon, lat, 0). It is the start of the Hindu day that jd falls in.
func lastSunrise(jd float64, geopos [3]float64) (float64, bool) {
	ephemeris.Lock()
	defer ephemeris.Unlock()
	rise := func(from float64) (float64, bool) {
		r := swisseph.RiseTrans(from, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcRise), geopos, 0, 0)
		return r.Time, r.Error == "" && r.Flag >= 0
	}
	t, ok := rise(jd - 1.0)
	if !ok {
		return 0, false
	}
	if t > jd {
		// Sunrise is a minute or so later than 24 h after the previous one; take the one before.
		t, ok = rise(jd - 2.0)
		if !ok || t > jd {
			return 0, false
		}
	}
	return t, true
}

// weekdayLordOfSunrise is the lord of the weekday whose sunrise is sunriseJD (UT), using the civil
// date at the given UTC offset in hours.
func weekdayLordOfSunrise(sunriseJD, tzHours float64) string {
	idx := int(math.Floor(sunriseJD+tzHours/24.0+1.5)) % 7 // 0 = Sunday
	if idx < 0 {
		idx += 7
	}
	return WeekdayLords[idx]
}
