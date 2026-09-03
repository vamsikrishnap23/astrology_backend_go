package houses

import (
	"github.com/mshafiee/swephgo"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

// CalculateHouses calculates the Ascendant, MC and House cusps.
func CalculateHouses(ctx *domain.CalculationContext) (float64, float64, []domain.HouseCusp, error) {
	cusps := make([]float64, 13)
	ascmc := make([]float64, 10)

	// Set Ayanamsa for Sidereal House calculation
	swephgo.SetSidMode(ctx.Config.AyanamsaMode, 0, 0)
	
	// Need eps, armc? HousesEx calculates this
	// Sidereal houses
	swephgo.HousesEx(ctx.JulianDayUT, swephgo.SeflgSidereal, ctx.Input.Latitude, ctx.Input.Longitude, int(ctx.Config.HouseCode), cusps, ascmc)

	ascendant := ascmc[0]
	mc := ascmc[1]

	var houseCusps []domain.HouseCusp
	for i := 1; i <= 12; i++ {
		sign, _, _, _ := time.DecimalToDMS(cusps[i])
		houseCusps = append(houseCusps, domain.HouseCusp{
			HouseNumber: i,
			Longitude:   cusps[i],
			Sign:        sign,
		})
	}

	return ascendant, mc, houseCusps, nil
}
