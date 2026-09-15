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
		DateOfBirth: "2005-11-23", // Guessing from a previous log
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

	// We don't know the exact progression date. Let's just try to parse what the user sent
	// Wait, the chart says 08/01/2026.
	progUtcTime, _ := astronomyTime.ParseLocalToUTC("2026-01-08", "00:00:00", input.Timezone) // Could be Jan 8 or Aug 1
	progJd := astronomyTime.UTCToJulianDay(progUtcTime)

	res, _ := progression.CalculateSecondaryProgression(ctx, "2026-01-08", progJd)
	
	// Just print the longitudes of Sun and Moon to see if they match the image
	for _, p := range res.ProgressedPlanets {
		if p.Planet == "Moon" || p.Planet == "Sun" || p.Planet == "Venus" {
			fmt.Printf("Prog %s: %.3f\n", p.Planet, p.SiderealLongitude)
		}
	}

	b, _ := json.MarshalIndent(res.Aspects, "", "  ")
	fmt.Println(string(b))
}
