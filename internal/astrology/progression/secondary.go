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
		if p.Planet == "Sun" || p.Planet == "Moon" || p.Planet == "Mercury" || p.Planet == "Venus" || p.Planet == "Mars" {
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
		if n.Planet != "Rahu" && n.Planet != "Ketu" {
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
			if pPoint.Name == nPoint.Name && strings.Contains(pPoint.Name, "Bhava") {
				continue
			}

			diff := math.Abs(pPoint.Longitude - nPoint.Longitude)
			diff = math.Mod(diff, 360.0)
			if diff > 180.0 {
				diff = 360.0 - diff
			}

			orbMax := 1.0
			var aspectType, nature, astrologicalRule string
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
			} else if math.Abs(diff-150.0) <= orbMax {
				aspectType = "Quincunx"
				exactAngle = 150.0
				nature = "Mixed"
				astrologicalRule = "At 150 degrees, the planets have nothing in common (different element, different modality, different polarity). This creates an awkward, irritating energy that requires constant adjustment."
			} else if math.Abs(diff-180.0) <= orbMax {
				aspectType = "Opposition"
				exactAngle = 180.0
				nature = "Hard/Dynamic"
				astrologicalRule = "At 180 degrees, the planets are at opposite ends of the zodiac. They pull in completely opposing directions, creating a dynamic tug-of-war that requires conscious balance and compromise."
			}

			if aspectType != "" {
				progKeywords := map[string]string{
					"Bhava 1 (Ascendant)":  "physical body, outward personality, and life path",
					"Bhava 2":              "personal wealth, values, and family lineage",
					"Bhava 3":              "courage, siblings, and communication skills",
					"Bhava 4":              "inner peace, mother, and foundational home life",
					"Bhava 5":              "creativity, intellect, and children",
					"Bhava 6":              "health, daily service, and ability to overcome enemies",
					"Bhava 7 (Descendant)": "marriage, partnerships, and public dealings",
					"Bhava 8":              "longevity, hidden wealth, and deep transformations",
					"Bhava 9":              "dharma, higher learning, and fortune",
					"Bhava 10 (MC)":        "career, public reputation, and highest ambitions",
					"Bhava 11":             "major gains, aspirations, and social networks",
					"Bhava 12":             "spirituality, isolation, and foreign connections",
					"Sun":                  "core identity, ego, and life focus",
					"Moon":                 "emotional needs, intuition, and domestic life",
					"Mercury":              "communication, mindset, and daily routines",
					"Venus":                "values, romantic desires, and financial flow",
					"Mars":                 "drive, ambition, and physical energy",
					"Jupiter":              "desire for expansion, growth, and optimism",
					"Saturn":               "sense of duty, discipline, and restriction",
					"Uranus":               "need for radical change, freedom, and innovation",
					"Neptune":              "spiritual ideals, dreams, and potential illusions",
					"Pluto":                "urge for deep transformation, power, and rebirth",
				}

				natKeywords := map[string]string{
					"Bhava 1 (Ascendant)":  "your physical presence, self-image, and approach to life",
					"Bhava 2":              "your personal wealth, values, and family lineage",
					"Bhava 3":              "your courage, siblings, and communication skills",
					"Bhava 4":              "your inner peace, mother, and foundational home life",
					"Bhava 5":              "your creativity, intellect, and children",
					"Bhava 6":              "your health, daily service, and ability to overcome enemies",
					"Bhava 7 (Descendant)": "your marriage, partnerships, and public dealings",
					"Bhava 8":              "your longevity, hidden wealth, and deep transformations",
					"Bhava 9":              "your dharma, higher learning, and fortune",
					"Bhava 10 (MC)":        "your ultimate career goals, social standing, and legacy",
					"Bhava 11":             "your major gains, aspirations, and social networks",
					"Bhava 12":             "your spirituality, isolation, and foreign connections",
					"Sun":                  "your fundamental life purpose and vitality",
					"Moon":                 "your baseline emotional security",
					"Mercury":              "how you naturally process information",
					"Venus":                "your capacity for love and receiving abundance",
					"Mars":                 "your natural assertiveness and conflict resolution",
					"Jupiter":              "where you naturally seek luck and higher meaning",
					"Saturn":               "your deep-seated boundaries, fears, and structures",
					"Uranus":               "your authentic individuality and rebelliousness",
					"Neptune":              "your inherent spiritual connection and compassion",
					"Pluto":                "your psychological depths and hidden power",
				}

				pK := progKeywords[pPoint.Name]
				nK := natKeywords[nPoint.Name]
				if pK == "" {
					pK = pPoint.Name
				}
				if nK == "" {
					nK = nPoint.Name
				}

				orb := math.Abs(diff - exactAngle)
				orb = math.Round(orb*100) / 100

				aspects = append(aspects, domain.ProgressedAspect{
					ProgressedPlanet: pPoint.Name,
					NatalPlanet:      nPoint.Name,
					Angle:            exactAngle,
					Orb:              orb,
					AspectType:       aspectType,
					Nature:           nature,
					AstrologicalRule: astrologicalRule,
					ProgKeyword:      pK,
					NatKeyword:       nK,
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
