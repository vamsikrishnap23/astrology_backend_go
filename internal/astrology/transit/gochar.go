package transit

import (
	"fmt"
	"math"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func getSignIndex(lon float64) int {
	return int(math.Floor(lon / 30.0))
}

func getHouseFromRef(planetLon, refLon float64) int {
	pSign := getSignIndex(planetLon)
	rSign := getSignIndex(refLon)
	return (pSign-rSign+12)%12 + 1
}

// Define Vedha pairs: goodHouse -> vedhaHouse
var vedhaRules = map[string]map[int]int{
	"Sun":     {3: 9, 6: 12, 10: 4, 11: 5},
	"Moon":    {1: 5, 3: 9, 6: 12, 7: 2, 10: 4, 11: 8},
	"Mars":    {3: 12, 6: 9, 11: 5},
	"Mercury": {2: 5, 4: 3, 6: 9, 8: 1, 10: 8, 11: 12},
	"Jupiter": {2: 12, 5: 4, 7: 3, 9: 10, 11: 8},
	"Venus":   {1: 8, 2: 7, 3: 1, 4: 10, 5: 9, 8: 5, 9: 11, 11: 6, 12: 3},
	"Saturn":  {3: 12, 6: 9, 11: 5},
	"Rahu":    {3: 12, 6: 9, 11: 5},
	"Ketu":    {3: 12, 6: 9, 11: 5},
}

var vamaVedhaRules = map[string]map[int]int{}

func init() {
	// Dynamically build Vama Vedha by reversing Vedha rules
	for planet, rules := range vedhaRules {
		vamaVedhaRules[planet] = make(map[int]int)
		for goodH, badH := range rules {
			vamaVedhaRules[planet][badH] = goodH
		}
	}
}

func getTransitOccupants(transitingPlanets []domain.PlanetPosition, natalMoonLon float64) map[int][]string {
	occupants := make(map[int][]string)
	for _, p := range transitingPlanets {
		h := getHouseFromRef(p.SiderealLongitude, natalMoonLon)
		occupants[h] = append(occupants[h], p.Planet)
	}
	return occupants
}

func hasObstructingPlanet(house int, occupants map[int][]string, excludedPlanet string) (bool, string) {
	for _, p := range occupants[house] {
		if p != excludedPlanet {
			return true, p
		}
	}
	return false, ""
}

func CalculateGochar(natalMoonLon float64, transitingPlanets []domain.PlanetPosition) []domain.GocharStatus {
	var results []domain.GocharStatus
	occupants := getTransitOccupants(transitingPlanets, natalMoonLon)

	planetList := []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn", "Rahu", "Ketu"}

	for _, pName := range planetList {
		var tp domain.PlanetPosition
		found := false
		for _, p := range transitingPlanets {
			if p.Planet == pName {
				tp = p
				found = true
				break
			}
		}
		if !found {
			continue
		}

		house := getHouseFromRef(tp.SiderealLongitude, natalMoonLon)

		status := "Bad"
		reason := fmt.Sprintf("%s is transiting house %d from Natal Moon, which yields inauspicious results.", pName, house)

		// Check if it's a Good house (i.e. has a Vedha rule defined for it)
		vedhaHouse, isGoodHouse := vedhaRules[pName][house]

		// Set excluded planet for Vedha exceptions
		excludedPlanet := ""
		if pName == "Sun" {
			excludedPlanet = "Saturn"
		} else if pName == "Saturn" {
			excludedPlanet = "Sun"
		} else if pName == "Moon" {
			excludedPlanet = "Mercury"
		} else if pName == "Mercury" {
			excludedPlanet = "Moon"
		}

		if isGoodHouse {
			status = "Good"
			reason = fmt.Sprintf("%s is transiting house %d from Natal Moon, which yields auspicious results.", pName, house)

			// Check Vedha
			isObstructed, obstructingPlanet := hasObstructingPlanet(vedhaHouse, occupants, excludedPlanet)
			if isObstructed {
				status = "Obstructed (Vedha)"
				reason = fmt.Sprintf("%s is in an auspicious transit house (%d), but its good results are blocked by %s transiting the Vedha house (%d).", pName, house, obstructingPlanet, vedhaHouse)
			}
		} else {
			// Check Vama Vedha
			vamaVedhaHouse, isVamaVedhaHouse := vamaVedhaRules[pName][house]
			if isVamaVedhaHouse {
				isObstructed, obstructingPlanet := hasObstructingPlanet(vamaVedhaHouse, occupants, excludedPlanet)
				if isObstructed {
					status = "Obstructed (Vama Vedha)"
					reason = fmt.Sprintf("%s is in an inauspicious transit house (%d), but its bad results are blocked by %s transiting the Vama Vedha house (%d).", pName, house, obstructingPlanet, vamaVedhaHouse)
				}
			}
		}

		results = append(results, domain.GocharStatus{
			Planet:        pName,
			HouseFromMoon: house,
			Status:        status,
			Reason:        reason,
		})
	}

	return results
}
