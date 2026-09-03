package kp

import (
	"sort"

	"github.com/vamsi/astrology_backend_go/internal/domain"
)

// CalculateFourStepSignificators derives the 4-step KP signification hierarchy for all planets.
func CalculateFourStepSignificators(tables domain.TablesResult, sigs domain.KPSignificatorsResult) domain.FourStepSignificatorsResult {
	var result domain.FourStepSignificatorsResult

	// Helper to lookup planet details from tables
	getPlanetTable := func(name string) *domain.TablePlanet {
		for _, p := range tables.PlanetaryTable {
			if p.PlanetName == name {
				return &p
			}
		}
		return nil
	}

	// Helper to get significated houses (Occupied + Owned = B + D) for any given planet
	getHouses := func(name string) []int {
		var houses []int
		for _, ps := range sigs.PlanetView {
			if ps.Planet == name {
				houses = append(houses, ps.B...)
				houses = append(houses, ps.D...)
				break
			}
		}

		uniqueMap := make(map[int]bool)
		var uniqueHouses []int
		for _, h := range houses {
			if !uniqueMap[h] {
				uniqueMap[h] = true
				uniqueHouses = append(uniqueHouses, h)
			}
		}
		sort.Ints(uniqueHouses)
		if uniqueHouses == nil {
			uniqueHouses = []int{}
		}
		return uniqueHouses
	}

	for _, p := range tables.PlanetaryTable {
		var fourStep domain.FourStepSignificator
		fourStep.Planet = p.PlanetName

		// Step 1: Planet
		fourStep.PlanetDetails = domain.StepDetails{
			Planet: p.PlanetName,
			Houses: getHouses(p.PlanetName),
		}

		// Step 2: Star Lord
		fourStep.StarLord = domain.StepDetails{
			Planet: p.NakshatraLord,
			Houses: getHouses(p.NakshatraLord),
		}

		// Step 3: Sub Lord
		fourStep.SubLord = domain.StepDetails{
			Planet: p.SubLord,
			Houses: getHouses(p.SubLord),
		}

		// Step 4: Star Lord of Sub Lord
		slOfSub := ""
		subTable := getPlanetTable(p.SubLord)
		if subTable != nil {
			slOfSub = subTable.NakshatraLord
		}
		fourStep.StarLordOfSub = domain.StepDetails{
			Planet: slOfSub,
			Houses: getHouses(slOfSub),
		}

		result.FourStepView = append(result.FourStepView, fourStep)
	}

	return result
}
