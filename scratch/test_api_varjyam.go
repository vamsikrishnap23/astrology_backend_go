package main

import (
	"encoding/json"
	"fmt"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/panchang"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func main() {
	ephemeris.Init("ephe_data")
	defer ephemeris.Close()
	
	input := domain.BirthInput{
		DateOfBirth: "2026-09-09",
		TimeOfBirth: "12:00:00",
		Latitude:    16.3938,
		Longitude:   80.1522,
		Timezone:    5.5,
		Ayanamsa:    "Lahiri",
	}

	utcTime, _ := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	config := domain.CalculationConfig{
		AyanamsaMode: ephemeris.GetAyanamsaMode(input.Ayanamsa),
		HouseCode:    ephemeris.GetHouseSystemCode(input.HouseSystem),
	}

	ctx := domain.CalculationContext{
		Input:       input,
		Config:      config,
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	res, _ := panchang.CalculatePanchang(&ctx)

	b, _ := json.MarshalIndent(res.Varjyam, "", "  ")
	fmt.Println("Varjyam:")
	fmt.Println(string(b))

	b2, _ := json.MarshalIndent(res.AmruthaGhadiyalu, "", "  ")
	fmt.Println("Amrutha Ghadiyalu:")
	fmt.Println(string(b2))
}
