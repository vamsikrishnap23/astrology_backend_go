package transit

import (
	"time"

	"github.com/vamsi/astrology_backend_go/internal/astrology/tables"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

// CalculateTransitChart calculates a complete transit chart for a given date/time,
// using the natal configuration for location, house system, and ayanamsa.
func CalculateTransitChart(ctx *domain.CalculationContext) (domain.TransitResult, error) {
	progPlanets, err := planets.CalculatePlanets(ctx)
	if err != nil {
		return domain.TransitResult{}, err
	}

	progAsc, _, progHouseCusps, err := houses.CalculateHouses(ctx)
	if err != nil {
		return domain.TransitResult{}, err
	}

	tblRes := tables.GenerateTables(progPlanets, progHouseCusps)

	res := domain.TransitResult{
		TransitDateUTC: ctx.UTCTime.Format(time.RFC3339),
		JulianDay:      ctx.JulianDayUT,
		Ayanamsa:       ctx.Ayanamsa,
		Ascendant:      progAsc,
		TransitData:    tblRes,
	}

	return res, nil
}
