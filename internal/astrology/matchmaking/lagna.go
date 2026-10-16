package matchmaking

import (
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// lagnaQuality classifies the distance counted from the bride's lagna to the groom's lagna:
// 1, 3, 5, 9, 11 are auspicious; 2 and 7 moderate; 4, 6, 8, 10, 12 inauspicious.
func lagnaQuality(distance int) string {
	switch distance {
	case 1, 3, 5, 9, 11:
		return "Auspicious"
	case 2, 7:
		return "Moderate"
	default:
		return "Inauspicious"
	}
}

// lagnaCompatibility compares the two ascendants. The lords' relationship is: Same (one lord),
// Friendly (mutual friends), Inimical (either is an enemy of the other), otherwise Neutral.
func lagnaCompatibility(groomCtx, brideCtx *domain.CalculationContext) domain.LagnaCompatibility {
	gAscLong, _, _, _ := houses.CalculateHouses(groomCtx)
	gLagna, _, _, _ := time.DecimalToDMS(gAscLong)

	bAscLong, _, _, _ := houses.CalculateHouses(brideCtx)
	bLagna, _, _, _ := time.DecimalToDMS(bAscLong)

	dist := (getSignIndex(gLagna)-getSignIndex(bLagna)+12)%12 + 1

	gLord := SignLords[gLagna]
	bLord := SignLords[bLagna]

	rel := "Same"
	if gLord != bLord {
		rel1 := GrahaNaturalMaitri[bLord][gLord]
		rel2 := GrahaNaturalMaitri[gLord][bLord]
		if rel1 == "Friend" && rel2 == "Friend" {
			rel = "Friendly"
		} else if rel1 == "Enemy" || rel2 == "Enemy" {
			rel = "Inimical"
		} else {
			rel = "Neutral"
		}
	}

	return domain.LagnaCompatibility{
		GroomLagna:        gLagna,
		BrideLagna:        bLagna,
		GroomLagnaLord:    gLord,
		BrideLagnaLord:    bLord,
		Distance:          dist,
		Quality:           lagnaQuality(dist),
		IsAuspicious:      dist == 1 || dist == 3 || dist == 5 || dist == 9 || dist == 11,
		LordsRelationship: rel,
	}
}
