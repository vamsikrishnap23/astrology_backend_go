package main

import (
	"encoding/json"
	"fmt"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/manglik"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func main() {
	ephemeris.Init("ephe_data")
	defer ephemeris.Close()

	input := domain.BirthInput{
		DateOfBirth: "1995-10-05",
		TimeOfBirth: "12:00:00",
		Latitude:    16.3938,
		Longitude:   80.1522,
		Timezone:    5.5,
		Ayanamsa:    "Lahiri",
	}

	utcTime, _ := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	ctx := &domain.CalculationContext{
		Input:       input,
		Config:      domain.CalculationConfig{AyanamsaMode: 1},
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	planetPositions, _ := planets.CalculatePlanets(ctx)
	// mock ascendant
	ascLon := 120.0

	res := manglik.CalculateManglikDosha(planetPositions, ascLon)

	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}
