package shadbala

import (
	"math"

	"github.com/tejzpr/go-swisseph"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// Cheshta (motional) Bala, Raman chapter VI (art. 79-107, pp. 64-79 and Tables IV-IX, pp. 112-117),
// after Sripati; BPHS 27.24-25.
//
// For Mars, Mercury, Jupiter, Venus and Saturn the Chesta Kendra is the planet's Seeghrochcha minus
// half the sum of its mean and true longitudes, reduced to 0..180 deg; Cheshta Bala is that divided
// by 3 (so 0 at a Kendra of 0 and 60 at 180). The mean longitudes come from mean-motion tables
// counted from the epoch 0 h on 1 January 1900 at Ujjain (76 deg E). The Sun and Moon have no Cheshta
// Bala in the Shadbala total (Raman example 56 leaves them blank); their Chesta Balas of art. 136-137
// are only used for Ishta and Kashta Phala.

const (
	ujjainLongitude = 76.0
	epoch1900JDUT   = 2415020.5 // 0 h UT on 1 January 1900
)

// Mean daily motions in degrees (nirayana), from Raman's Tables IV-IX. Where his printed table
// entries disagree with each other (Table IV 146.0265, Table VII rows 8-9, Table VIII 133.32) the
// rate that all the other entries of the same column imply is used.
const (
	sunMeanRate        = 0.98560265
	marsMeanRate       = 0.5240193
	jupiterMeanRate    = 0.0830963
	saturnMeanRate     = 0.033439
	mercurySeeghraRate = 4.092318
	venusSeeghraRate   = 1.6021465
)

// ramanMeans are the longitudes needed for the Chesta Kendra of the five planets.
type ramanMeans struct {
	Sun            float64 // mean Sun: also the mean longitude of Mercury and Venus, and the Seeghrochcha of Mars, Jupiter and Saturn
	Mars           float64
	Jupiter        float64
	Saturn         float64
	MercurySeeghra float64
	VenusSeeghra   float64
}

// ramanInterval is the "total interval" in days from the epoch to the birth moment (art. 89-92): the
// birth time is taken at Ujjain's longitude, so 76/360 of a day is added to Universal Time.
func ramanInterval(jdUT float64) float64 {
	return jdUT + ujjainLongitude/360.0 - epoch1900JDUT
}

func norm360(x float64) float64 {
	x = math.Mod(x, 360.0)
	if x < 0 {
		x += 360.0
	}
	return x
}

// computeRamanMeans applies the epoch constants, the mean motions and the long-term corrections of
// art. 94-103; t is the birth year minus 1900.
func computeRamanMeans(interval, t float64) ramanMeans {
	return ramanMeans{
		Sun:            norm360(257.4568 + sunMeanRate*interval),
		Mars:           norm360(270.22 + marsMeanRate*interval),
		Jupiter:        norm360(220.04 + jupiterMeanRate*interval - (3.33 + 0.0067*t)),
		Saturn:         norm360(236.74 + saturnMeanRate*interval + (5.0 + 0.001*t)),
		MercurySeeghra: norm360(164.0 + mercurySeeghraRate*interval + (6.67 - 0.00133*t)),
		VenusSeeghra:   norm360(328.51 + venusSeeghraRate*interval - (5.0 + 0.001*t)),
	}
}

// cheshtaBala is Cheshta Bala from the Seeghrochcha, the mean and the true longitude (art. 105-107).
// The true longitude is brought within 180 deg of the mean one before they are averaged, so that a
// pair straddling 0/360 deg is averaged correctly.
func cheshtaBala(seeghrochcha, mean, trueLon float64) float64 {
	if d := trueLon - mean; d > 180.0 {
		trueLon -= 360.0
	} else if d < -180.0 {
		trueLon += 360.0
	}
	kendra := norm360(seeghrochcha - (mean+trueLon)/2.0)
	if kendra > 180.0 {
		kendra = 360.0 - kendra
	}
	return kendra / 3.0
}

// cheshtaFromMeans is Cheshta Bala of planet given the mean longitudes of ramanMeans and the planet's
// true longitude, both in the same zodiac. The Sun and Moon return 0.
func cheshtaFromMeans(planet string, m ramanMeans, trueLon float64) float64 {
	switch planet {
	case "Mars":
		return cheshtaBala(m.Sun, m.Mars, trueLon)
	case "Mercury":
		return cheshtaBala(m.MercurySeeghra, m.Sun, trueLon)
	case "Jupiter":
		return cheshtaBala(m.Sun, m.Jupiter, trueLon)
	case "Venus":
		return cheshtaBala(m.VenusSeeghra, m.Sun, trueLon)
	case "Saturn":
		return cheshtaBala(m.Sun, m.Saturn, trueLon)
	}
	return 0
}

// ramanFrameShift is the angle to add to Raman's mean longitudes to express them in the zodiac of the
// chart's ayanamsa. His constants and examples are in his own zodiac (the Raman ayanamsa, as the
// match of his examples with SE's Raman mode shows), while the true longitudes follow the ayanamsa the
// user chose; averaging the two without this shift would make Cheshta Bala depend on that choice by
// half the difference of the two ayanamsas. With the Raman ayanamsa the shift is 0 and the book's
// numbers are reproduced exactly.
func ramanFrameShift(jdUT float64, ayanamsaMode int) float64 {
	if ayanamsaMode == swisseph.SidmRaman {
		return 0
	}
	ephemeris.Lock()
	defer ephemeris.Unlock()
	swisseph.SetSidMode(swisseph.SidmRaman, 0, 0)
	raman := ephemeris.AyanamsaUT(jdUT)
	swisseph.SetSidMode(int32(ayanamsaMode), 0, 0)
	return raman - ephemeris.AyanamsaUT(jdUT)
}

// shifted returns the means moved by delta degrees.
func (m ramanMeans) shifted(delta float64) ramanMeans {
	return ramanMeans{
		Sun: norm360(m.Sun + delta), Mars: norm360(m.Mars + delta), Jupiter: norm360(m.Jupiter + delta),
		Saturn: norm360(m.Saturn + delta), MercurySeeghra: norm360(m.MercurySeeghra + delta),
		VenusSeeghra: norm360(m.VenusSeeghra + delta),
	}
}

// CalculateCheshtaBala returns the Cheshta Bala of planet for the chart in ctx; lons holds the
// sidereal (true) longitudes of the planets in the chart's ayanamsa. The Sun and Moon return 0.
func CalculateCheshtaBala(planet string, ctx *domain.CalculationContext, lons map[string]float64) float64 {
	switch planet {
	case "Mars", "Mercury", "Jupiter", "Venus", "Saturn":
	default:
		return 0 // the Sun and Moon have none in the total
	}
	trueLon, ok := lons[planet]
	if !ok {
		return 0
	}
	year := swisseph.Revjul(ctx.JulianDayUT+ctx.Input.Timezone/24.0, swisseph.GregCal).Year
	m := computeRamanMeans(ramanInterval(ctx.JulianDayUT), float64(year-1900))
	m = m.shifted(ramanFrameShift(ctx.JulianDayUT, ctx.Config.AyanamsaMode))
	return cheshtaFromMeans(planet, m, trueLon)
}
