package main

import (
	"fmt"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/progression"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func main() {
	ephemeris.Init("ephe_data")

	input := domain.BirthInput{
		DateOfBirth: "1975-03-11",
		TimeOfBirth: "03:15:00",
		Latitude:    16.42593,
		Longitude:   79.54041,
		Timezone:    5.5,
		Ayanamsa:    "KP",
		HouseSystem: "Placidus",
	}

	utcTime, _ := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	jd := astronomyTime.UTCToJulianDay(utcTime)
	ctx := &domain.CalculationContext{
		Input:       input,
		Config:      domain.CalculationConfig{AyanamsaMode: ephemeris.GetAyanamsaMode(input.Ayanamsa), HouseCode: ephemeris.GetHouseSystemCode(input.HouseSystem)},
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}
	
	dateStr := "2026-08-25"
	progUtcTime, _ := astronomyTime.ParseLocalToUTC(dateStr, "00:00:00", input.Timezone)
	progJd := astronomyTime.UTCToJulianDay(progUtcTime)
	res, _ := progression.CalculateSecondaryProgression(ctx, dateStr, progJd)
	
	fmt.Println("Aspects on Aug 25 2026:")
	for _, a := range res.Aspects {
		if a.ProgressedPlanet == "Moon" && (a.NatalPlanet == "Sun" || a.NatalPlanet == "Venus") {
			fmt.Printf("To %s: %s (Orb: %.3f)\n", a.NatalPlanet, a.AspectType, a.Orb)
		}
	}
}
