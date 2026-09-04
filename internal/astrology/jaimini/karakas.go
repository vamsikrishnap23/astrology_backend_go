package jaimini

import (
	"sort"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

var JaiminiPlanets = []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"}

var KarakaNames = []string{"AK", "AmK", "BK", "MK", "PK", "GK", "DK"}

func getFallbackPriority(planet string) int {
	for i, p := range JaiminiPlanets {
		if p == planet {
			return i
		}
	}
	return 99
}

// CalculateCharaKarakas computes the 7 Jaimini Chara Karakas.
func CalculateCharaKarakas(ctx *domain.CalculationContext) (domain.JaiminiKarakasResult, error) {
	calcPlanets, err := planets.CalculatePlanets(ctx)
	if err != nil {
		return domain.JaiminiKarakasResult{}, err
	}

	var allPlanets []domain.CharaKaraka

	for _, p := range calcPlanets {
		allPlanets = append(allPlanets, domain.CharaKaraka{
			Planet:          p.Planet,
			SourceLongitude: p.SiderealLongitude,
			DivisionalSign:  p.Sign,
			Degree:          p.Degree,
			Minute:          p.Minute,
			Second:          p.Second,
			Nakshatra:       p.Nakshatra,
			NakshatraPada:   p.NakshatraPada,
			NakshatraLord:   p.NakshatraLord,
			DegreeInSign:    p.DegreeInSign,
			Retrograde:      p.Retrograde,
		})
	}

	// Filter out the 7 core planets for Jaimini ranking
	var jaiminiRefs []*domain.CharaKaraka
	for i := range allPlanets {
		for _, jp := range JaiminiPlanets {
			if jp == allPlanets[i].Planet {
				jaiminiRefs = append(jaiminiRefs, &allPlanets[i])
				break
			}
		}
	}

	// Sort by DegreeInSign descending
	sort.Slice(jaiminiRefs, func(i, j int) bool {
		if jaiminiRefs[i].DegreeInSign == jaiminiRefs[j].DegreeInSign {
			// Deterministic tie-breaker based on standard planetary order
			return getFallbackPriority(jaiminiRefs[i].Planet) < getFallbackPriority(jaiminiRefs[j].Planet)
		}
		return jaiminiRefs[i].DegreeInSign > jaiminiRefs[j].DegreeInSign
	})

	// Assign Karakas
	for i := range jaiminiRefs {
		if i < len(KarakaNames) {
			jaiminiRefs[i].Karaka = KarakaNames[i]
		}
	}

	return domain.JaiminiKarakasResult{
		CalculationTimeUTC: ctx.UTCTime.Format(time.RFC3339),
		Planets:            allPlanets,
	}, nil
}
