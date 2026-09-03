package houses

import (
	"github.com/tejzpr/go-swisseph"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

// CalculateHouses calculates the Ascendant, MC and House cusps.
func CalculateHouses(ctx *domain.CalculationContext) (float64, float64, []domain.HouseCusp, error) {

	// Set Ayanamsa for Sidereal House calculation
	swisseph.SetSidMode(int32(ctx.Config.AyanamsaMode), 0, 0)

	// Sidereal houses
	res := swisseph.HousesEx(ctx.JulianDayUT, int32(swisseph.FlagSidereal), ctx.Input.Latitude, ctx.Input.Longitude, ctx.Config.HouseCode)

	ascendant := res.Points[0]
	mc := res.Points[1]

	var houseCusps []domain.HouseCusp
	for i := 0; i < len(res.Houses) && i < 12; i++ {
		sign, _, _, _ := time.DecimalToDMS(res.Houses[i])
		houseCusps = append(houseCusps, domain.HouseCusp{
			HouseNumber: i + 1,
			Longitude:   res.Houses[i],
			Sign:        sign,
		})
	}

	return ascendant, mc, houseCusps, nil
}
