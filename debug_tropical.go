package main

import (
	"fmt"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
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
	}

	utcTime, _ := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	jd := astronomyTime.UTCToJulianDay(utcTime)
	ctx := &domain.CalculationContext{
		Input:       input,
		Config:      domain.CalculationConfig{AyanamsaMode: ephemeris.GetAyanamsaMode(input.Ayanamsa)},
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	natPlanets, _ := planets.CalculatePlanets(ctx)
	for _, p := range natPlanets {
		if p.Planet == "Sun" || p.Planet == "Venus" {
			fmt.Printf("Natal %s Tropical: %.3f\n", p.Planet, p.TropicalLongitude)
		}
	}

	progUtcTime, _ := astronomyTime.ParseLocalToUTC("2026-08-01", "00:00:00", input.Timezone)
	progJd := astronomyTime.UTCToJulianDay(progUtcTime)
	
	// tropical year
	daysAlive := progJd - jd
	ageInYears := daysAlive / 365.242190402
	pJD := jd + ageInYears
	
	pCtx := &domain.CalculationContext{
		Input:       input,
		Config:      ctx.Config,
		JulianDayUT: pJD,
	}
	progPlanets, _ := planets.CalculatePlanets(pCtx)
	for _, p := range progPlanets {
		if p.Planet == "Moon" {
			fmt.Printf("Prog Moon Tropical: %.3f\n", p.TropicalLongitude)
		}
	}
}
