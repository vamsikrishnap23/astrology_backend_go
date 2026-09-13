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

	input := domain.BTRInput{
		BirthInput: domain.BirthInput{
			DateOfBirth: "2005-11-23",
			TimeOfBirth: "15:36:00",
			Latitude:    16.066,
			Longitude:   79.9833,
			Timezone:    5.5,
			Ayanamsa:    "Lahiri",
		},
		Gender: "Male",
		ScanMinusMinutes: 10,
		ScanPlusMinutes: 5,
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
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}
