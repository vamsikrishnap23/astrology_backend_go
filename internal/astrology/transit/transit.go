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

// CalculateTransitChart calculates a complete transit chart for a given date/time,
// using the natal moon (Chandra Lagna) as the reference point for Whole-Sign house cusps.
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

	// 2. Calculate Transit Planets
	progPlanets, err := planets.CalculatePlanets(transitCtx)
	if err != nil {
		return domain.TransitResult{}, err
	}

	// 3. Generate Whole Sign houses based on the Natal Moon's sign (Chandra Lagna)
	var moonHouseCusps []domain.HouseCusp
	moonSignIndex := int(math.Floor(natalMoon.SiderealLongitude / 30.0))

	for i := 0; i < 12; i++ {
		signIdx := (moonSignIndex + i) % 12
		siderealLon := float64(signIdx * 30.0)

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
