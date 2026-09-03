package kp

import (
	"time"

	"github.com/vamsi/astrology_backend_go/internal/domain"
)

var weekdayRulers = []string{
	"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn",
}

// CalculateRulingPlanets derives the KP Ruling Planets from the calculated tables and input date.
//
// Nodes as Agents Convention:
// Rahu and Ketu act as agents for a planet if they are:
// 1. Conjoined with that planet (occupying the same sign).
// 2. In the sign owned by that planet (Sign Lord).
// 3. In the nakshatra owned by that planet (Star Lord).
// If a true planet is a ruling planet, any Node acting as its agent is appended as a ruling planet.
//
// Deduplication Decision:
// We do NOT deduplicate planets. A planet appearing multiple times (e.g., as Moon Sign Lord AND Day Lord)
// is considered a stronger ruling planet in KP astrology. We preserve all occurrences and their sources.
func CalculateRulingPlanets(input domain.BirthInput, tables domain.TablesResult) domain.KPRulingPlanetsResult {
	var rps []domain.RulingPlanet

	// 1. Ascendant Lords
	var ascHouse domain.TableHouse
	for _, h := range tables.HouseTable {
		if h.HouseNumber == 1 {
			ascHouse = h
			break
		}
	}
	if ascHouse.HouseNumber == 1 {
		rps = append(rps, domain.RulingPlanet{Planet: ascHouse.NakshatraLord, Source: domain.SourceAscendantStarLord})
		rps = append(rps, domain.RulingPlanet{Planet: ascHouse.SignLord, Source: domain.SourceAscendantSignLord})
		rps = append(rps, domain.RulingPlanet{Planet: ascHouse.SubLord, Source: domain.SourceAscendantSubLord})
	}

	// 2. Moon Lords
	var moon domain.TablePlanet
	for _, p := range tables.PlanetaryTable {
		if p.PlanetName == "Moon" {
			moon = p
			break
		}
	}
	if moon.PlanetName == "Moon" {
		rps = append(rps, domain.RulingPlanet{Planet: moon.NakshatraLord, Source: domain.SourceMoonStarLord})
		rps = append(rps, domain.RulingPlanet{Planet: moon.SignLord, Source: domain.SourceMoonSignLord})
		rps = append(rps, domain.RulingPlanet{Planet: moon.SubLord, Source: domain.SourceMoonSubLord})
	}

	// 3. Day Lord (Using the local calendar day)
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, input.DateOfBirth)
	if err == nil {
		wd := int(parsedDate.Weekday())
		dayLord := weekdayRulers[wd]
		rps = append(rps, domain.RulingPlanet{Planet: dayLord, Source: domain.SourceDayLord})
	}

	// 4. Nodes as Agents
	var rahu, ketu domain.TablePlanet
	for _, p := range tables.PlanetaryTable {
		if p.PlanetName == "Rahu" {
			rahu = p
		}
		if p.PlanetName == "Ketu" {
			ketu = p
		}
	}

	// Helper to find if a Node acts as an agent for a given target planet
	isAgentFor := func(node domain.TablePlanet, target string) bool {
		// Rule 1: Sign Lord
		if node.SignLord == target {
			return true
		}
		// Rule 2: Nakshatra Lord
		if node.NakshatraLord == target {
			return true
		}
		// Rule 3: Conjunction (in the same sign)
		for _, p := range tables.PlanetaryTable {
			if p.PlanetName == target && p.Sign == node.Sign {
				return true
			}
		}
		return false
	}

	// Extract current unique true planets in ruling planets
	trueRPs := make(map[string]bool)
	for _, rp := range rps {
		if rp.Planet != "Rahu" && rp.Planet != "Ketu" {
			trueRPs[rp.Planet] = true
		}
	}

	// Create a unique list of RPs to avoid duplicating agent entries if a planet appears twice in RPs
	uniqueTargets := []string{}
	seenTargets := make(map[string]bool)
	for _, rp := range rps {
		if !seenTargets[rp.Planet] && rp.Planet != "Rahu" && rp.Planet != "Ketu" {
			uniqueTargets = append(uniqueTargets, rp.Planet)
			seenTargets[rp.Planet] = true
		}
	}

	for _, target := range uniqueTargets {
		if rahu.PlanetName != "" && isAgentFor(rahu, target) {
			rps = append(rps, domain.RulingPlanet{Planet: "Rahu", Source: domain.SourceNodeAgent, AgentFor: target})
		}
		if ketu.PlanetName != "" && isAgentFor(ketu, target) {
			rps = append(rps, domain.RulingPlanet{Planet: "Ketu", Source: domain.SourceNodeAgent, AgentFor: target})
		}
	}

	return domain.KPRulingPlanetsResult{
		RulingPlanets: rps,
	}
}
