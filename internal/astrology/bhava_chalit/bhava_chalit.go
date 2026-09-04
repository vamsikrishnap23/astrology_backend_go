package bhava_chalit

import (
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/tables"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// CalculateBhavaChalit calculates the Bhava Chalit chart, returning exactly the requested fields.
func CalculateBhavaChalit(ctx *domain.CalculationContext) (domain.BhavaChalitResult, error) {
	planetsList, err := planets.CalculatePlanets(ctx)
	if err != nil {
		return domain.BhavaChalitResult{}, err
	}

	ascendant, _, houseCusps, err := houses.CalculateHouses(ctx)
	if err != nil {
		return domain.BhavaChalitResult{}, err
	}

	// We can leverage tables.GenerateTables to reuse the logic that assigns planets to houses
	tblRes := tables.GenerateTables(planetsList, houseCusps)

	res := domain.BhavaChalitResult{
		Ascendant: ascendant,
		Planets:   []domain.BhavaChalitPlanet{},
		Houses:    houseCusps, // Direct assignment of the full array
	}

	// Build the flattened planets array
	for _, tp := range tblRes.PlanetaryTable {
		bp := domain.BhavaChalitPlanet{
			Planet:          tp.PlanetName,
			SourceLongitude: tp.ExactLongitude,
			DivisionalSign:  tp.Sign,
			Degree:          tp.Degree,
			Minute:          tp.Minute,
			Second:          tp.Second,
			Nakshatra:       tp.Nakshatra,
			NakshatraPada:   tp.NakshatraPada,
			NakshatraLord:   tp.NakshatraLord,
			SignLord:        tp.SignLord,
			Retrograde:      tp.Retrograde,
			HouseNumber:     tp.HouseNumber,
		}
		res.Planets = append(res.Planets, bp)
	}

	return res, nil
}
