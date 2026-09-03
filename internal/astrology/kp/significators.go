package kp

import (
	"sort"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// CalculateSignificators derives the KP significators (A, B, C, D) based on the planetary and house tables.
func CalculateSignificators(tables domain.TablesResult) domain.KPSignificatorsResult {
	// 1. Build lookup maps
	planetHouseMap := make(map[string]int)       // Planet -> House it occupies
	houseOwnersMap := make(map[string][]int)     // Planet -> Houses it owns
	planetStarLordMap := make(map[string]string) // Planet -> Star Lord

	for _, p := range tables.PlanetaryTable {
		planetHouseMap[p.PlanetName] = p.HouseNumber
		planetStarLordMap[p.PlanetName] = p.NakshatraLord
	}

	for _, h := range tables.HouseTable {
		houseOwnersMap[h.SignLord] = append(houseOwnersMap[h.SignLord], h.HouseNumber)
	}

	// 2. Build Planet View
	var planetView []domain.PlanetSignificator

	// Maintain a consistent order for planets (same as PlanetaryTable input)
	for _, p := range tables.PlanetaryTable {
		planetName := p.PlanetName
		starLord := planetStarLordMap[planetName]

		// A = Star Lord's occupant Bhava
		var a []int
		if slHouse, exists := planetHouseMap[starLord]; exists {
			a = []int{slHouse}
		} else {
			a = []int{}
		}

		// B = Planet's occupant Bhava
		b := []int{p.HouseNumber}

		// C = Houses owned by Star Lord
		c := houseOwnersMap[starLord]
		if c == nil {
			c = []int{}
		} else {
			// Copy to avoid modifying map slice
			cCopy := make([]int, len(c))
			copy(cCopy, c)
			c = cCopy
			sort.Ints(c)
		}

		// D = Houses owned by the Planet
		d := houseOwnersMap[planetName]
		if d == nil {
			d = []int{}
		} else {
			dCopy := make([]int, len(d))
			copy(dCopy, d)
			d = dCopy
			sort.Ints(d)
		}

		planetView = append(planetView, domain.PlanetSignificator{
			Planet: planetName,
			A:      a,
			B:      b,
			C:      c,
			D:      d,
		})
	}

	// 3. Build House View
	var houseView []domain.HouseSignificator
	for i := 1; i <= 12; i++ {
		var a, b, c, d []string

		for _, ps := range planetView {
			if contains(ps.A, i) {
				a = append(a, ps.Planet)
			}
			if contains(ps.B, i) {
				b = append(b, ps.Planet)
			}
			if contains(ps.C, i) {
				c = append(c, ps.Planet)
			}
			if contains(ps.D, i) {
				d = append(d, ps.Planet)
			}
		}

		if a == nil {
			a = []string{}
		}
		if b == nil {
			b = []string{}
		}
		if c == nil {
			c = []string{}
		}
		if d == nil {
			d = []string{}
		}

		houseView = append(houseView, domain.HouseSignificator{
			House: i,
			A:     a,
			B:     b,
			C:     c,
			D:     d,
		})
	}

	return domain.KPSignificatorsResult{
		PlanetView: planetView,
		HouseView:  houseView,
	}
}

func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
