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
		TimeOfBirth: "15:15:00",
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
	
	progUtcTime, _ := astronomyTime.ParseLocalToUTC("2026-08-01", "00:00:00", input.Timezone)
	progJd := astronomyTime.UTCToJulianDay(progUtcTime)
	res, _ := progression.CalculateSecondaryProgression(ctx, "2026-08-01", progJd)
	
	for _, a := range res.Aspects {
		if a.NatalPlanet == "Sun" && a.Orb >= 0.3 && a.Orb <= 0.4 {
			fmt.Printf("Prog %s to Natal Sun: %s (Orb: %.3f)\n", a.ProgressedPlanet, a.AspectType, a.Orb)
		}
	}
}
