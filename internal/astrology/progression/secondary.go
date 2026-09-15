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

			var aspectType, nature, astrologicalRule, reason string
			var exactAngle float64

			if diff <= orbMax {
				aspectType = "Conjunction"
				exactAngle = 0.0
				nature = "Variable"
				astrologicalRule = "At 0 degrees, the progressed planet and natal planet occupy the exact same point in space. This creates an intense blending of their energies where they can no longer operate independently."
			} else if math.Abs(diff-60.0) <= orbMax {
				aspectType = "Sextile"
				exactAngle = 60.0
				nature = "Harmonious"
				astrologicalRule = "At 60 degrees, the planets are in complementary elements (like Fire and Air, or Earth and Water). This creates a cooperative geometry that naturally generates favorable opportunities."
			} else if math.Abs(diff-90.0) <= orbMax {
				aspectType = "Square"
				exactAngle = 90.0
				nature = "Hard/Dynamic"
				astrologicalRule = "At 90 degrees, the planets are in conflicting elemental natures but share the same modality. They block each other's path, creating intense psychological or external friction that forces action."
			} else if math.Abs(diff-120.0) <= orbMax {
				aspectType = "Trine"
				exactAngle = 120.0
				nature = "Harmonious"
				astrologicalRule = "At exactly 120 degrees, both planets are positioned in the exact same Astrological Element (e.g., both in Fire). Their energies flow together without any resistance, generating luck and effortless harmony."
			} else if math.Abs(diff-180.0) <= orbMax {
				aspectType = "Opposition"
				exactAngle = 180.0
				nature = "Hard/Dynamic"
				astrologicalRule = "At 180 degrees, the planets are at opposite ends of the zodiac. They pull in completely opposing directions, creating a dynamic tug-of-war that requires conscious balance and compromise."
			}

			if aspectType != "" {
				reason = generateAspectReason(pPlanet.Planet, nPlanet.Planet, aspectType, nature)
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
					AstrologicalRule: astrologicalRule,
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

func generateAspectReason(progPlanet, natPlanet, aspectType, nature string) string {
	progKeywords := map[string]string{
		"Sun":     "core identity, ego, and life focus",
		"Moon":    "emotional needs, intuition, and domestic life",
		"Mercury": "communication, mindset, and daily routines",
		"Venus":   "values, romantic desires, and financial flow",
		"Mars":    "drive, ambition, and physical energy",
		"Jupiter": "desire for expansion, growth, and optimism",
		"Saturn":  "sense of duty, discipline, and restriction",
		"Uranus":  "need for radical change, freedom, and innovation",
		"Neptune": "spiritual ideals, dreams, and potential illusions",
		"Pluto":   "urge for deep transformation, power, and rebirth",
	}

	natKeywords := map[string]string{
		"Sun":     "your fundamental life purpose and vitality",
		"Moon":    "your baseline emotional security",
		"Mercury": "how you naturally process information",
		"Venus":   "your capacity for love and receiving abundance",
		"Mars":    "your natural assertiveness and conflict resolution",
		"Jupiter": "where you naturally seek luck and higher meaning",
		"Saturn":  "your deep-seated boundaries, fears, and structures",
		"Uranus":  "your authentic individuality and rebelliousness",
		"Neptune": "your inherent spiritual connection and compassion",
		"Pluto":   "your psychological depths and hidden power",
	}

	pK, ok1 := progKeywords[progPlanet]
	nK, ok2 := natKeywords[natPlanet]

	if !ok1 || !ok2 {
		return "A significant energetic exchange between these two celestial bodies, heavily influenced by the angle of the aspect."
	}

	if progPlanet == natPlanet {
		if aspectType == "Conjunction" {
			return "A major life milestone. Your " + pK + " is undergoing a powerful reset and renewal, returning to its purest form."
		} else if aspectType == "Opposition" {
			return "A profound mid-cycle crisis. Your current " + pK + " is clashing heavily with " + nK + ", demanding massive re-evaluation."
		}
	}

	switch aspectType {
	case "Conjunction":
		return "An intense merging of forces. Your current " + pK + " is powerfully activating " + nK + ". This creates a hyper-focused period where these two areas of life cannot be separated."
	case "Trine":
		return "A period of supreme ease and luck. Your evolving " + pK + " effortlessly supports and enhances " + nK + ". Doors open naturally without forcing them."
	case "Sextile":
		return "An opportunity for productive growth. Your " + pK + " is in a cooperative position with " + nK + ". If you put in the effort, you will see highly positive results."
	case "Square":
		return "A highly stressful but necessary turning point. Your current " + pK + " is creating severe friction with " + nK + ". This tension forces you to break through obstacles and make hard choices."
	case "Opposition":
		return "A major tug-of-war. Your evolving " + pK + " is directly clashing with " + nK + ". You must find a compromise between these opposing forces, often triggered by other people or external events."
	default:
		return "A significant energetic exchange."
	}
}
