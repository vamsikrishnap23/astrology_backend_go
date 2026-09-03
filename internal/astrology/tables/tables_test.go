package tables

import (
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"testing"
)

func TestGenerateTables(t *testing.T) {
	planets := []domain.PlanetPosition{
		{
			Planet:            "Sun",
			Sign:              "Scorpio",
			Degree:            7,
			Minute:            15,
			Second:            0,
			SiderealLongitude: 217.25,
			Retrograde:        false,
		},
	}

	houses := []domain.HouseCusp{
		{HouseNumber: 1, Longitude: 10.0, Sign: "Aries"},
	}

	res := GenerateTables(planets, houses)

	if res.PlanetaryTable[0].SubLord == "" {
		t.Error("Expected sub lord to be generated")
	}
}
