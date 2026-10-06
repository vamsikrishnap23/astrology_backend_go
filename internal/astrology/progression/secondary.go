package progression

import (
	"math"
	"strings"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// CalculateSecondaryProgression calculates the progressed chart for a given date.
// Rule: 1 day after birth = 1 tropical year of life.
func CalculateSecondaryProgression(natalCtx *domain.CalculationContext, targetDate string, targetJD float64) (domain.ProgressionResult, error) {
	// Tropical year length in days
	tropicalYear := 365.242190402

	// How many days have they been alive?
	daysAlive := targetJD - natalCtx.JulianDayUT

	// 1 day = 1 year, so the fraction of days to add to natal JD is daysAlive / tropicalYear
	ageInYears := daysAlive / tropicalYear
	progressedJD := natalCtx.JulianDayUT + ageInYears

	// Add days to the UTC time as well
	progressedUTC := natalCtx.UTCTime.Add(time.Duration(ageInYears * 24 * float64(time.Hour)))

	// Create a new context for the progressed date
	progressedCtx := &domain.CalculationContext{
		Input:       natalCtx.Input,  // Same natal location and settings
		Config:      natalCtx.Config, // Same ayanamsa and house settings
		UTCTime:     progressedUTC,
		JulianDayUT: progressedJD,
	}

	// Calculate progressed planets
	progPlanets, err := planets.CalculatePlanets(progressedCtx)
	if err != nil {
		return domain.ProgressionResult{}, err
	}

	// Calculate progressed houses
	progAsc, progMC, progHouseCusps, err := houses.CalculateHouses(progressedCtx)
	if err != nil {
		return domain.ProgressionResult{}, err
	}

	// Calculate natal planets to check for aspects
	natalPlanets, err := planets.CalculatePlanets(natalCtx)
	if err != nil {
		return domain.ProgressionResult{}, err
	}

	// Calculate natal houses to check for aspects to Ascendant and MC (Bhavas)
	_, _, natHouseCusps, err := houses.CalculateHouses(natalCtx)
	if err != nil {
		return domain.ProgressionResult{}, err
	}

	var aspects []domain.ProgressedAspect

	// Combine planets and key Bhavas (Asc/MC) into a unified list
	type AstrologicalPoint struct {
		Name      string
		Longitude float64
	}

	var progPoints []AstrologicalPoint
	for _, p := range progPlanets {
		if p.Planet != "Neptune" && p.Planet != "Pluto" {
			progPoints = append(progPoints, AstrologicalPoint{Name: p.Planet, Longitude: p.SiderealLongitude})
		}
	}
	for _, cusp := range progHouseCusps {
		name := "Bhava 1" // fallback
		if cusp.HouseNumber == 1 {
			name = "Bhava 1 (Ascendant)"
		}
		if cusp.HouseNumber == 2 {
			name = "Bhava 2"
		}
		if cusp.HouseNumber == 3 {
			name = "Bhava 3"
		}
		if cusp.HouseNumber == 4 {
			name = "Bhava 4"
		}
		if cusp.HouseNumber == 5 {
			name = "Bhava 5"
		}
		if cusp.HouseNumber == 6 {
			name = "Bhava 6"
		}
		if cusp.HouseNumber == 7 {
			name = "Bhava 7 (Descendant)"
		}
		if cusp.HouseNumber == 8 {
			name = "Bhava 8"
		}
		if cusp.HouseNumber == 9 {
			name = "Bhava 9"
		}
		if cusp.HouseNumber == 10 {
			name = "Bhava 10 (MC)"
		}
		if cusp.HouseNumber == 11 {
			name = "Bhava 11"
		}
		if cusp.HouseNumber == 12 {
			name = "Bhava 12"
		}
		progPoints = append(progPoints, AstrologicalPoint{Name: name, Longitude: cusp.Longitude})
	}

	var natPoints []AstrologicalPoint
	for _, n := range natalPlanets {
		if true {
			natPoints = append(natPoints, AstrologicalPoint{Name: n.Planet, Longitude: n.SiderealLongitude})
		}
	}
	for _, cusp := range natHouseCusps {
		name := "Bhava 1" // fallback
		if cusp.HouseNumber == 1 {
			name = "Bhava 1 (Ascendant)"
		}
		if cusp.HouseNumber == 2 {
			name = "Bhava 2"
		}
		if cusp.HouseNumber == 3 {
			name = "Bhava 3"
		}
		if cusp.HouseNumber == 4 {
			name = "Bhava 4"
		}
		if cusp.HouseNumber == 5 {
			name = "Bhava 5"
		}
		if cusp.HouseNumber == 6 {
			name = "Bhava 6"
		}
		if cusp.HouseNumber == 7 {
			name = "Bhava 7 (Descendant)"
		}
		if cusp.HouseNumber == 8 {
			name = "Bhava 8"
		}
		if cusp.HouseNumber == 9 {
			name = "Bhava 9"
		}
		if cusp.HouseNumber == 10 {
			name = "Bhava 10 (MC)"
		}
		if cusp.HouseNumber == 11 {
			name = "Bhava 11"
		}
		if cusp.HouseNumber == 12 {
			name = "Bhava 12"
		}
		natPoints = append(natPoints, AstrologicalPoint{Name: name, Longitude: cusp.Longitude})
	}

	for _, pPoint := range progPoints {
		for _, nPoint := range natPoints {
			// Skip self-aspects for identical Bhavas (e.g. Progressed Bhava 1 to Natal Bhava 1)
			isProgBhava := strings.Contains(pPoint.Name, "Bhava")
			isNatBhava := strings.Contains(nPoint.Name, "Bhava")
			if isProgBhava && isNatBhava {
				continue
			}

			diff := math.Abs(pPoint.Longitude - nPoint.Longitude)
			diff = math.Mod(diff, 360.0)
			if diff > 180.0 {
				diff = 360.0 - diff
			}

			orbMax := 1.0
			var aspectType, nature string
			var exactAngle float64

			if diff <= orbMax {
				aspectType = "Conjunction"
				exactAngle = 0.0
				nature = "Variable"
			} else if math.Abs(diff-60.0) <= orbMax {
				aspectType = "Sextile"
				exactAngle = 60.0
				nature = "Harmonious"
			} else if math.Abs(diff-90.0) <= orbMax {
				aspectType = "Square"
				exactAngle = 90.0
				nature = "Hard/Dynamic"
			} else if math.Abs(diff-120.0) <= orbMax {
				aspectType = "Trine"
				exactAngle = 120.0
				nature = "Harmonious"
			} else if math.Abs(diff-150.0) <= orbMax {
				aspectType = "Quincunx"
				exactAngle = 150.0
				nature = "Mixed"
			} else if math.Abs(diff-180.0) <= orbMax {
				aspectType = "Opposition"
				exactAngle = 180.0
				nature = "Hard/Dynamic"
			}

			if aspectType != "" {

				orb := math.Abs(diff - exactAngle)
				orb = math.Round(orb*100) / 100

				aspects = append(aspects, domain.ProgressedAspect{
					ProgressedPlanet: pPoint.Name,
					NatalPlanet:      nPoint.Name,
					Angle:            exactAngle,
					Orb:              orb,
					AspectType:       aspectType,
					Nature:           nature,
				})
			}
		}
	}

	res := domain.ProgressionResult{
		NatalDateUTC:          natalCtx.UTCTime.Format(time.RFC3339),
		TargetProgressionDate: targetDate,
		AgeInYears:            ageInYears,
		ProgressedDateUTC:     progressedUTC.Format(time.RFC3339),
		ProgressedJulianDay:   progressedJD,
		ProgressedAyanamsa:    progressedCtx.Ayanamsa,
		Ascendant:             progAsc,
		MC:                    progMC,
		ProgressedPlanets:     progPlanets,
		ProgressedHouses:      progHouseCusps,
		Aspects:               aspects,
	}

	return res, nil
}
