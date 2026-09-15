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
		DateOfBirth: "1975-03-11",
		TimeOfBirth: "03:15:00",
		Latitude:    16.42593,
		Longitude:   79.54041,
		Timezone:    5.5,
		Ayanamsa:    "Lahiri",
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

	progUtcTime, _ := astronomyTime.ParseLocalToUTC("2026-01-08", "00:00:00", input.Timezone)
	progJd := astronomyTime.UTCToJulianDay(progUtcTime)

	res, _ := progression.CalculateSecondaryProgression(ctx, "2026-01-08", progJd)
	
	// Print Natal Positions
	fmt.Println("--- Natal Positions ---")
	fmt.Printf("Sun: %.3f\n", 326.185) // We'll see real values
	
	for _, a := range res.Aspects {
		if a.ProgressedPlanet == "Moon" && (a.NatalPlanet == "Sun" || a.NatalPlanet == "Venus") {
			b, _ := json.MarshalIndent(a, "", "  ")
			fmt.Println(string(b))
		}
	}
	
	// Let's print all Moon aspects to see what it found
	fmt.Println("--- All Prog Moon Aspects ---")
	for _, a := range res.Aspects {
		if a.ProgressedPlanet == "Moon" {
			fmt.Printf("To %s: %s (Orb: %.3f)\n", a.NatalPlanet, a.AspectType, a.Orb)
		}
	}

	// Just to verify manual math, print exact longitudes
	fmt.Println("--- Exact Longitudes ---")
	var moonLon float64
	for _, p := range res.ProgressedPlanets {
		if p.Planet == "Moon" { moonLon = p.SiderealLongitude }
	}
	fmt.Printf("Prog Moon: %.3f\n", moonLon)
	
	// Note: res only returns ProgressedPlanets. We need to calculate NatalPlanets manually here if we want to print them.
}
