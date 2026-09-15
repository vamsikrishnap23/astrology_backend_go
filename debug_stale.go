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
	
	// Test a range of dates in 2026
	for m := 1; m <= 12; m++ {
		dateStr := fmt.Sprintf("2026-%02d-01", m)
		progUtcTime, _ := astronomyTime.ParseLocalToUTC(dateStr, "00:00:00", input.Timezone)
		progJd := astronomyTime.UTCToJulianDay(progUtcTime)
		res, _ := progression.CalculateSecondaryProgression(ctx, dateStr, progJd)
		
		for _, a := range res.Aspects {
			if a.ProgressedPlanet == "Moon" && a.NatalPlanet == "Sun" {
				fmt.Printf("Date %s -> %s to Sun, Orb: %.3f\n", dateStr, a.AspectType, a.Orb)
			}
		}
	}
}
