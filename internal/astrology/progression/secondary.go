package progression

import (
	"math"
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

	var aspects []domain.ProgressedAspect

	for _, pPlanet := range progPlanets {
		// Only check major planets/luminaries
		if pPlanet.Planet == "Rahu" || pPlanet.Planet == "Ketu" {
			continue
		}

		for _, nPlanet := range natalPlanets {
			if nPlanet.Planet == "Rahu" || nPlanet.Planet == "Ketu" {
				continue
			}

			// Don't check outer planets against outer planets usually, but let's check all for completeness
			diff := math.Abs(pPlanet.SiderealLongitude - nPlanet.SiderealLongitude)
			diff = math.Mod(diff, 360.0)
			if diff > 180.0 {
				diff = 360.0 - diff
			}

			// Define orb for progressions
			orbMax := 1.0

			var aspectType, nature, reason string
			var exactAngle float64

			if diff <= orbMax {
				aspectType = "Conjunction"
				exactAngle = 0.0
				nature = "Variable"
				reason = "Intensely powerful merging of energies. Effect depends heavily on the planets involved."
			} else if math.Abs(diff-60.0) <= orbMax {
				aspectType = "Sextile"
				exactAngle = 60.0
				nature = "Harmonious"
				reason = "Brings positive opportunities and favorable circumstances. A door is unlocked, but action is required."
			} else if math.Abs(diff-90.0) <= orbMax {
				aspectType = "Square"
				exactAngle = 90.0
				nature = "Hard/Dynamic"
				reason = "Creates friction, stress, or tension that forces necessary action and growth."
			} else if math.Abs(diff-120.0) <= orbMax {
				aspectType = "Trine"
				exactAngle = 120.0
				nature = "Harmonious"
				reason = "Brings effortless flow, luck, and natural talents manifesting without friction."
			} else if math.Abs(diff-180.0) <= orbMax {
				aspectType = "Opposition"
				exactAngle = 180.0
				nature = "Hard/Dynamic"
				reason = "Indicates a tug-of-war or conflict requiring balance and compromise between opposing forces."
			}

			if aspectType != "" {
				orb := math.Abs(diff - exactAngle)
				// Round orb to 2 decimal places
				orb = math.Round(orb*100) / 100

				aspects = append(aspects, domain.ProgressedAspect{
					ProgressedPlanet: pPlanet.Planet,
					NatalPlanet:      nPlanet.Planet,
					Angle:            exactAngle,
					Orb:              orb,
					AspectType:       aspectType,
					Nature:           nature,
					Reason:           reason,
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
