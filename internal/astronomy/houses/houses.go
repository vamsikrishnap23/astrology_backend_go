package houses

import (
	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
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

	var houseCusps []domain.HouseCusp
	for i := 0; i < len(res.Houses) && i < 12; i++ {
		siderealLon := res.Houses[i] - ctx.Ayanamsa
		if siderealLon < 0 {
			siderealLon += 360.0
		}

		sign, _, _, _ := time.DecimalToDMS(siderealLon)
		houseCusps = append(houseCusps, domain.HouseCusp{
			HouseNumber: i + 1,
			Longitude:   siderealLon,
			Sign:        sign,
		})
	}

	return ascendant, mc, houseCusps, nil
}
