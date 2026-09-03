package bhava_chalit

import (
	"github.com/vamsi/astrology_backend_go/internal/astrology/tables"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsi/astrology_backend_go/internal/domain"
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
	}

	for _, th := range tblRes.HouseTable {
		bh := domain.BhavaChalitHouse{
			HouseNumber:   th.HouseNumber,
			CuspLongitude: th.CuspLongitude,
			Sign:          th.Sign,
			Degree:        th.Degree,
			Minute:        th.Minute,
			Second:        th.Second,
			Occupants:     []domain.BhavaChalitPlanet{},
		}

		// Find the occupants
		for _, occName := range th.Occupants {
			// Find planet in planetaryTable
			for _, tp := range tblRes.PlanetaryTable {
				if tp.PlanetName == occName {
					bp := domain.BhavaChalitPlanet{
						PlanetName:     tp.PlanetName,
						HouseNumber:    tp.HouseNumber,
						Sign:           tp.Sign,
						Degree:         tp.Degree,
						Minute:         tp.Minute,
						Second:         tp.Second,
						ExactLongitude: tp.ExactLongitude,
					}
					bh.Occupants = append(bh.Occupants, bp)
					break
				}
			}
		}

		res.Houses = append(res.Houses, bh)
	}

	return res, nil
}
