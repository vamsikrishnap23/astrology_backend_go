package main

import (
	"encoding/json"
	"fmt"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/btr"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func main() {
	ephemeris.Init("ephe_data")
	defer ephemeris.Close()

	input := domain.BTRInput{
		BirthInput: domain.BirthInput{
			DateOfBirth: "2026-09-13",
			TimeOfBirth: "01:44:00",
			Latitude:    17.38405,
			Longitude:   78.45636,
			Timezone:    5.5,
			Ayanamsa:    "Lahiri",
		},
		Gender: "Male",
	}

	utcTime, _ := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	ctx := &domain.CalculationContext{
		Input:       input.BirthInput,
		Config:      domain.CalculationConfig{AyanamsaMode: 1, HouseCode: byte('P')},
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	res, _ := btr.CalculateBTR(input, ctx)
	b, _ := json.MarshalIndent(res.InputAnalysis, "", "  ")
	fmt.Println(string(b))
}
