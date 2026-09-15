package main

import (
	"encoding/json"
	"fmt"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/progression"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func main() {
	ephemeris.Init("ephe_data")

	input := domain.BirthInput{
		DateOfBirth: "2005-11-23",
		TimeOfBirth: "15:36:00",
		Latitude:    16.066,
		Longitude:   79.9833,
		Timezone:    5.5,
		Ayanamsa:    "Lahiri",
	}

	utcTime, _ := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	ctx := &domain.CalculationContext{
		Input:       input,
		Config:      domain.CalculationConfig{AyanamsaMode: 1, HouseCode: byte('P')},
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	// target date around 2026
	progUtcTime, _ := astronomyTime.ParseLocalToUTC("2026-09-15", input.TimeOfBirth, input.Timezone)
	progJd := astronomyTime.UTCToJulianDay(progUtcTime)

	res, _ := progression.CalculateSecondaryProgression(ctx, "2026-09-15", progJd)
	b, _ := json.MarshalIndent(res.Aspects, "", "  ")
	fmt.Println(string(b))
}
