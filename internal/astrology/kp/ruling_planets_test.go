package kp

import (
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func TestCalculateRulingPlanets(t *testing.T) {
	input := domain.BirthInput{
		DateOfBirth: "2005-11-23", // Wednesday (Day Lord = Mercury)
	}

	tables := domain.TablesResult{
		PlanetaryTable: []domain.TablePlanet{
			{PlanetName: "Moon", NakshatraLord: "Venus", SignLord: "Sun", SubLord: "Jupiter"},
			{PlanetName: "Mercury", Sign: "Scorpio"},
			{PlanetName: "Venus", Sign: "Sagittarius"},
			{PlanetName: "Jupiter", Sign: "Libra"},
			{PlanetName: "Rahu", Sign: "Pisces", SignLord: "Jupiter", NakshatraLord: "Saturn"},
			{PlanetName: "Ketu", Sign: "Virgo", SignLord: "Mercury", NakshatraLord: "Moon"}, // Conjoined with nobody, but sign lord is Mercury
		},
		HouseTable: []domain.TableHouse{
			{HouseNumber: 1, NakshatraLord: "Ketu", SignLord: "Mars", SubLord: "Saturn"},
		},
	}

	res := CalculateRulingPlanets(input, tables)

	expectedSources := []domain.RulingPlanet{
		{Planet: "Ketu", Source: domain.SourceAscendantStarLord},
		{Planet: "Mars", Source: domain.SourceAscendantSignLord},
		{Planet: "Saturn", Source: domain.SourceAscendantSubLord},
		{Planet: "Venus", Source: domain.SourceMoonStarLord},
		{Planet: "Sun", Source: domain.SourceMoonSignLord},
		{Planet: "Jupiter", Source: domain.SourceMoonSubLord},
		{Planet: "Mercury", Source: domain.SourceDayLord},
	}

	// Saturn triggers Rahu
	expectedSources = append(expectedSources, domain.RulingPlanet{Planet: "Rahu", Source: domain.SourceNodeAgent, AgentFor: "Saturn"})

	// Jupiter triggers Rahu
	expectedSources = append(expectedSources, domain.RulingPlanet{Planet: "Rahu", Source: domain.SourceNodeAgent, AgentFor: "Jupiter"})

	// Mercury triggers Ketu
	expectedSources = append(expectedSources, domain.RulingPlanet{Planet: "Ketu", Source: domain.SourceNodeAgent, AgentFor: "Mercury"})

	if len(res.RulingPlanets) != len(expectedSources) {
		t.Fatalf("Expected %d ruling planets, got %d. Got: %+v", len(expectedSources), len(res.RulingPlanets), res.RulingPlanets)
	}

	for i, rp := range expectedSources {
		if res.RulingPlanets[i].Planet != rp.Planet || res.RulingPlanets[i].Source != rp.Source || res.RulingPlanets[i].AgentFor != rp.AgentFor {
			t.Errorf("Mismatch at index %d: expected %+v, got %+v", i, rp, res.RulingPlanets[i])
		}
	}
}
