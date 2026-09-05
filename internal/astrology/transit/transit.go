package transit

import (
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/tables"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// CalculateTransitChart calculates a complete transit chart for a given date/time.
// It mathematically superimposes the Transiting Planets directly onto the user's
// exact Natal House Cusps (D1 Chart), respecting their chosen house system (e.g., Placidus).
func CalculateTransitChart(natalCtx *domain.CalculationContext, transitCtx *domain.CalculationContext) (domain.TransitResult, error) {
	// 1. Calculate exact Natal Houses (D1)
	natalAsc, _, natalHouseCusps, err := houses.CalculateHouses(natalCtx)
	if err != nil {
		return domain.TransitResult{}, err
	}

	// 2. Calculate Transit Planets
	progPlanets, err := planets.CalculatePlanets(transitCtx)
	if err != nil {
		return domain.TransitResult{}, err
	}

	// 3. Superimpose Transiting Planets onto Natal (D1) Houses
	tblRes := tables.GenerateTables(progPlanets, natalHouseCusps)

	res := domain.TransitResult{
		TransitDateUTC: transitCtx.UTCTime.Format(time.RFC3339),
		JulianDay:      transitCtx.JulianDayUT,
		Ayanamsa:       transitCtx.Ayanamsa,
		Ascendant:      natalAsc,
		TransitData:    tblRes,
	}

	return res, nil
}
