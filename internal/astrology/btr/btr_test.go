package btr

import (
	"path/filepath"
	"testing"

	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func TestBTR(t *testing.T) {
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	input := domain.BTRInput{
		Name:             "Vamsi",
		DateOfBirth:      "2005-11-23",
		TimeOfBirth:      "15:35:00",
		Latitude:         16.3938,
		Longitude:        80.1522,
		Timezone:         5.5,
		Ayanamsa:         "Lahiri",
		HouseSystem:      "Placidus",
		Gender:           "male",
		ScanMinusMinutes: 5,
		ScanPlusMinutes:  5,
	}

	res, err := ProcessBTR(input)
	if err != nil {
		t.Fatalf("Failed to process BTR: %v", err)
	}

	if res.AstronomicalContext.DayLord != "Wednesday" && res.AstronomicalContext.DayLord != "Mercury" {
		t.Errorf("Expected Wednesday/Mercury for day lord, got %v", res.AstronomicalContext.DayLord)
	}

	if len(res.Candidates) == 0 {
		t.Errorf("Expected candidates, got 0")
	}
}
