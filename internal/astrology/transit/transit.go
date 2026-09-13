package transit

import (
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/tables"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// CalculateTransitChart calculates a complete transit chart for a given date/time.
func CalculateTransitChart(ctx *domain.CalculationContext, natalMoonLon float64) (domain.TransitResult, error) {
	progPlanets, err := planets.CalculatePlanets(ctx)
	if err != nil {
		return domain.TransitResult{}, err
	}

	progAsc, _, progHouseCusps, err := houses.CalculateHouses(ctx)
	if err != nil {
		return domain.TransitResult{}, err
	}

	tblRes := tables.GenerateTables(progPlanets, progHouseCusps)

	gocharResults := CalculateGochar(natalMoonLon, progPlanets)

	res := domain.TransitResult{
		TransitDateUTC: ctx.UTCTime.Format(time.RFC3339),
		JulianDay:      ctx.JulianDayUT,
		Ayanamsa:       ctx.Ayanamsa,
		Ascendant:      progAsc,
		TransitData:    tblRes,
		Gochar:         gocharResults,
	}

	return res, nil
}
