package ashtakavarga

import (
	"math"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

var signNames = []string{
	"Aries", "Taurus", "Gemini", "Cancer",
	"Leo", "Virgo", "Libra", "Scorpio",
	"Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

// targetBAVPlanets represents the 7 planets that get a BAV chart
var targetBAVPlanets = []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"}

// getSignIndex returns the 0-11 index of the sign for a given longitude
func getSignIndex(longitude float64) int {
	return int(math.Floor(math.Mod(longitude, 360.0) / 30.0))
}

func CalculateAshtakavarga(ctx *domain.CalculationContext) (domain.AshtakavargaResult, error) {
	// 1. Calculate positions of planets
	calculatedPlanets, err := planets.CalculatePlanets(ctx)
	if err != nil {
		return domain.AshtakavargaResult{}, err
	}

	// 2. Calculate Ascendant
	ascLongitude, _, _, err := houses.CalculateHouses(ctx)
	if err != nil {
		return domain.AshtakavargaResult{}, err
	}

	// 3. Map positions to sign indices (0-11)
	positions := make(map[string]int)
	for _, p := range calculatedPlanets {
		positions[p.Planet] = getSignIndex(p.SiderealLongitude)
	}
	positions["Ascendant"] = getSignIndex(ascLongitude)

	var bavList []domain.Bhinnashtakavarga

	// Initialize SAV signs
	savSigns := make([]domain.SignBindu, 12)
	for i := 0; i < 12; i++ {
		savSigns[i] = domain.SignBindu{
			Sign:        signNames[i],
			SignIndex:   i + 1,
			TotalBindus: 0,
		}
	}

	totalSAVBindus := 0

	// 4. Calculate BAV for each of the 7 planets
	for _, bavPlanet := range targetBAVPlanets {
		bav := domain.Bhinnashtakavarga{
			Planet:      bavPlanet,
			TotalBindus: 0,
			Signs:       make([]domain.SignBindu, 12),
		}

		for i := 0; i < 12; i++ {
			bav.Signs[i] = domain.SignBindu{
				Sign:          signNames[i],
				SignIndex:     i + 1,
				TotalBindus:   0,
				Contributions: []domain.BinduContribution{},
			}
		}

		planetRules := Rules[bavPlanet]

		for sourcePlanet, relativeHouses := range planetRules {
			sourceSignIdx := positions[sourcePlanet]

			for _, rh := range relativeHouses {
				// rh is 1-indexed relative house.
				// e.g. rh=1 means the same sign as sourcePlanet.
				targetSignIdx := (sourceSignIdx + rh - 1) % 12

				// Add to BAV
				bav.Signs[targetSignIdx].TotalBindus++
				bav.Signs[targetSignIdx].Contributions = append(bav.Signs[targetSignIdx].Contributions, domain.BinduContribution{
					SourcePlanet: sourcePlanet,
					Value:        1,
				})
				bav.TotalBindus++

				// Add to SAV
				savSigns[targetSignIdx].TotalBindus++
				totalSAVBindus++
			}
		}

		bavList = append(bavList, bav)
	}

	res := domain.AshtakavargaResult{
		CalculationTimeUTC: ctx.UTCTime.Format(time.RFC3339),
		BAV:                bavList,
		SAV:                savSigns,
		TotalSAVBindus:     totalSAVBindus,
	}

	return res, nil
}
