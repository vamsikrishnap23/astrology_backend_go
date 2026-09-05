package progression

import (
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
	}

	return res, nil
}
