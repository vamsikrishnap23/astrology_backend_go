package jaimini

import (
	"sort"
	"time"

	"github.com/vamsi/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsi/astrology_backend_go/internal/domain"
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

	var candidates []domain.CharaKaraka

	for _, p := range calcPlanets {
		isJaimini := false
		for _, jp := range JaiminiPlanets {
			if jp == p.Planet {
				isJaimini = true
				break
			}
		}

		if isJaimini {
			candidates = append(candidates, domain.CharaKaraka{
				Planet:       p.Planet,
				Sign:         p.Sign,
				Degree:       p.Degree,
				Minute:       p.Minute,
				Second:       p.Second,
				DegreeInSign: p.DegreeInSign,
				Retrograde:   p.Retrograde,
			})
		}
	}

	// Sort by DegreeInSign descending
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].DegreeInSign == candidates[j].DegreeInSign {
			// Deterministic tie-breaker based on standard planetary order
			return getFallbackPriority(candidates[i].Planet) < getFallbackPriority(candidates[j].Planet)
		}
		return candidates[i].DegreeInSign > candidates[j].DegreeInSign
	})

	// Assign Karakas
	for i := range candidates {
		if i < len(KarakaNames) {
			candidates[i].Karaka = KarakaNames[i]
		}
	}

	return domain.JaiminiKarakasResult{
		CalculationTimeUTC: ctx.UTCTime.Format(time.RFC3339),
		Karakas:            candidates,
	}, nil
}
