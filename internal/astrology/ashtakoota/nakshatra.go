package ashtakoota

import (
	"math"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func getMoonDetails(jd float64, ayanamsa float64) domain.MoonDetails {
	ephemeris.Mu.Lock()
	defer ephemeris.Mu.Unlock()

	iflag := int32(swisseph.FlagSwieph)
	res := swisseph.CalcUT(jd, swisseph.Moon, iflag)
	tropLon := res.Data[0]

	sidLon := math.Mod(tropLon-ayanamsa+360.0, 360.0)

	signIdx := int(math.Floor(sidLon / 30.0))
	sign := Signs[signIdx]

	degInSign := math.Mod(sidLon, 30.0)

	// Nakshatra is 13deg 20min (13.333... deg) = 800 minutes
	// Total minutes = 360 * 60 = 21600
	totalMinutes := sidLon * 60.0
	nakIdx := int(math.Floor(totalMinutes/800.0)) % 27

	// Pada is 3deg 20min = 200 minutes
	pada := int(math.Floor(math.Mod(totalMinutes, 800.0)/200.0)) + 1

	nakshatraLords := []string{"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury"}

	return domain.MoonDetails{
		Longitude:     sidLon,
		Sign:          sign,
		Degree:        degInSign,
		Nakshatra:     Nakshatras[nakIdx],
		Pada:          pada,
		NakshatraLord: nakshatraLords[nakIdx%9],
		RashiLord:     SignLords[sign],
	}
}
