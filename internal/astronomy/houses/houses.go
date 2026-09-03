package houses

import (
	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"math"
)

// CalculateHouses calculates the Ascendant, MC and House cusps.
func CalculateHouses(ctx *domain.CalculationContext) (float64, float64, []domain.HouseCusp, error) {

	// Lock the global ephemeris mutex to prevent C-level concurrency corruption
	ephemeris.Mu.Lock()
	defer ephemeris.Mu.Unlock()

	// IMPORTANT: Re-initialize path for the current OS thread
	swisseph.SetEphePath(ephemeris.EphePath)

	// Set Ayanamsa for Sidereal House calculation
	swisseph.SetSidMode(int32(ctx.Config.AyanamsaMode), 0, 0)

	// Calculate Tropical houses first
	// (Vedic astrology subtracts Ayanamsa directly to get Sidereal, avoiding SE's complex projection)
	res := swisseph.HousesEx(ctx.JulianDayUT, 0, ctx.Input.Latitude, ctx.Input.Longitude, ctx.Config.HouseCode)

	// Pre-calculate ayanamsa value (ensure it's available in context, or fetch here)
	ctx.Ayanamsa = swisseph.GetAyanamsaUT(ctx.JulianDayUT)

	// Convert points to Sidereal
	ascendant := res.Points[0] - ctx.Ayanamsa
	if ascendant < 0 {
		ascendant += 360.0
	}
	mc := res.Points[1] - ctx.Ayanamsa
	if mc < 0 {
		mc += 360.0
	}

	var nakshatras = []string{
		"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
		"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni",
		"Hasta", "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha",
		"Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha",
		"Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
	}
	var nakshatraLords = []string{"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury"}

	var houseCusps []domain.HouseCusp
	for i := 0; i < len(res.Houses) && i < 12; i++ {
		siderealLon := res.Houses[i] - ctx.Ayanamsa
		for siderealLon < 0 {
			siderealLon += 360.0
		}
		for siderealLon >= 360.0 {
			siderealLon -= 360.0
		}

		sign, deg, min, sec := time.DecimalToDMS(siderealLon)

		interval := 13.0 + 1.0/3.0
		nakIdx := int(math.Floor(siderealLon / interval))
		nakName := nakshatras[nakIdx%27]
		nakLord := nakshatraLords[nakIdx%9]
		nakProgress := math.Mod(siderealLon, interval) / interval * 100.0
		pada := int(math.Floor(nakProgress/25.0)) + 1

		houseCusps = append(houseCusps, domain.HouseCusp{
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

	return ascendant, mc, houseCusps, nil
}
