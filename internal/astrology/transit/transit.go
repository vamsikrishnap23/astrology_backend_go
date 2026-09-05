package transit

import (
	"math"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/tables"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

var nakshatras = []string{
	"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
	"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni",
	"Hasta", "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha",
	"Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha",
	"Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
}
var nakshatraLords = []string{"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury"}

// CalculateTransitChart calculates a complete transit chart for a given date/time.
// It implements the "Equal House from Moon" (Chandra Lagna) method, which is the
// standard Vedic astrology technique for Gochara (Transits). The 1st house cusp is
// pegged exactly to the Natal Moon's Sidereal Longitude.
func CalculateTransitChart(natalCtx *domain.CalculationContext, transitCtx *domain.CalculationContext) (domain.TransitResult, error) {
	// 1. Calculate Natal Planets to find the Natal Moon
	natalPlanets, err := planets.CalculatePlanets(natalCtx)
	if err != nil {
		return domain.TransitResult{}, err
	}

	var natalMoon domain.PlanetPosition
	for _, p := range natalPlanets {
		if p.Planet == "Moon" {
			natalMoon = p
			break
		}
	}

	// 2. Generate "Equal House from Moon" Cusps
	var moonHouseCusps []domain.HouseCusp
	baseLongitude := natalMoon.SiderealLongitude

	for i := 0; i < 12; i++ {
		siderealLon := math.Mod(baseLongitude+float64(i*30), 360.0)

		sign, deg, min, sec := astronomyTime.DecimalToDMS(siderealLon)

		interval := 13.0 + 1.0/3.0
		nakIdx := int(math.Floor(siderealLon / interval))
		nakName := nakshatras[nakIdx%27]
		nakLord := nakshatraLords[nakIdx%9]
		nakProgress := math.Mod(siderealLon, interval) / interval * 100.0
		pada := int(math.Floor(nakProgress/25.0)) + 1

		moonHouseCusps = append(moonHouseCusps, domain.HouseCusp{
			HouseNumber:   i + 1,
			Longitude:     siderealLon,
			Sign:          sign,
			Degree:        deg,
			Minute:        min,
			Second:        sec,
			Nakshatra:     nakName,
			NakshatraPada: pada,
			NakshatraLord: nakLord,
		})
	}

	// 3. Calculate Transit Planets
	progPlanets, err := planets.CalculatePlanets(transitCtx)
	if err != nil {
		return domain.TransitResult{}, err
	}

	// 4. Superimpose Transiting Planets onto Moon Houses
	tblRes := tables.GenerateTables(progPlanets, moonHouseCusps)

	res := domain.TransitResult{
		TransitDateUTC: transitCtx.UTCTime.Format(time.RFC3339),
		JulianDay:      transitCtx.JulianDayUT,
		Ayanamsa:       transitCtx.Ayanamsa,
		Ascendant:      natalMoon.SiderealLongitude,
		TransitData:    tblRes,
	}

	return res, nil
}
