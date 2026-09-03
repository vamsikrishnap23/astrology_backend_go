package main

import (
	"fmt"
	"math"
	"path/filepath"

	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/planets"
	astronomyTime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

var signNames = []string{
	"Aries", "Taurus", "Gemini", "Cancer",
	"Leo", "Virgo", "Libra", "Scorpio",
	"Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

func main() {
	absPath, _ := filepath.Abs("./ephe_data")
	ephemeris.Init(absPath)

	utcTime, _ := astronomyTime.ParseLocalToUTC("2005-11-23", "15:35:00", 5.5)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	ctx := domain.CalculationContext{
		Input: domain.BirthInput{
			Latitude:  16.3938,
			Longitude: 80.1522,
		},
		Config: domain.CalculationConfig{
			AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri"),
			HouseCode:    ephemeris.GetHouseSystemCode("Placidus"),
		},
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	calculatedPlanets, _ := planets.CalculatePlanets(&ctx)
	ascLongitude, _, _, _ := houses.CalculateHouses(&ctx)

	for _, p := range calculatedPlanets {
		idx := int(math.Floor(math.Mod(p.SiderealLongitude, 360.0) / 30.0))
		fmt.Printf("%s: %.4f -> %s\n", p.Planet, p.SiderealLongitude, signNames[idx])
	}
	ascIdx := int(math.Floor(math.Mod(ascLongitude, 360.0) / 30.0))
	fmt.Printf("Ascendant: %.4f -> %s\n", ascLongitude, signNames[ascIdx])
}
