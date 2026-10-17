package shadbala

import (
	"errors"
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

// ErrNoSunriseSunset is returned when the Sun does not rise or set at the birth place around the
// birth date (polar day or night). Hora, Dina, Varsha, Masa and Tribhaga Bala all count from sunrise,
// so Shadbala cannot be given for such a chart; an error is better than silently giving them 0.
var ErrNoSunriseSunset = errors.New("shadbala cannot be calculated: the Sun does not rise or set at this latitude on this date")

// hinduDay is the Hindu day a moment falls in: it begins at Sunrise, its daytime ends at Sunset and
// the next day begins at NextSunrise (all Julian days, UT).
type hinduDay struct {
	Sunrise, Sunset, NextSunrise float64
}

// hinduDayAt finds the Hindu day containing jd for the geographic position (lon, lat, 0): the latest
// sunrise at or before jd, the sunset after it and the sunrise after that. It returns
// ErrNoSunriseSunset if any of the three cannot be found.
func hinduDayAt(jd float64, geopos [3]float64) (hinduDay, error) {
	ephemeris.Lock()
	defer ephemeris.Unlock()
	event := func(from float64, kind int) (float64, bool) {
		r := swisseph.RiseTrans(from, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(kind), geopos, 0, 0)
		return r.Time, r.Error == "" && r.Flag >= 0
	}
	sunrise, ok := event(jd-1.0, swisseph.CalcRise)
	if !ok {
		return hinduDay{}, ErrNoSunriseSunset
	}
	if sunrise > jd {
		// The next sunrise after jd-1 can fall a minute or so after jd; take the one before it.
		if sunrise, ok = event(jd-2.0, swisseph.CalcRise); !ok || sunrise > jd {
			return hinduDay{}, ErrNoSunriseSunset
		}
	}
	sunset, ok := event(sunrise, swisseph.CalcSet)
	if !ok {
		return hinduDay{}, ErrNoSunriseSunset
	}
	next, ok := event(sunset, swisseph.CalcRise)
	if !ok {
		return hinduDay{}, ErrNoSunriseSunset
	}
	return hinduDay{Sunrise: sunrise, Sunset: sunset, NextSunrise: next}, nil
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
