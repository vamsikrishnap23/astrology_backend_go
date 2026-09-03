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
	}

	for i, th := range tblRes.HouseTable {
		hc := houseCusps[i] // They align 1:1
		bh := domain.BhavaChalitHouse{
			HouseNumber:   th.HouseNumber,
			CuspLongitude: th.CuspLongitude,
			Sign:          th.Sign,
			Degree:        th.Degree,
			Minute:        th.Minute,
			Second:        th.Second,
			Nakshatra:     hc.Nakshatra,
			NakshatraPada: hc.NakshatraPada,
			NakshatraLord: hc.NakshatraLord,
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
						Nakshatra:      tp.Nakshatra,
						NakshatraPada:  tp.NakshatraPada,
						NakshatraLord:  tp.NakshatraLord,
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
