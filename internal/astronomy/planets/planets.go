package planets

import (
	"fmt"
	"math"
	"strings"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

var PlanetMap = map[string]int{
	"Sun":     swisseph.Sun,
	"Moon":    swisseph.Moon,
	"Mars":    swisseph.Mars,
	"Mercury": swisseph.Mercury,
	"Jupiter": swisseph.Jupiter,
	"Venus":   swisseph.Venus,
	"Saturn":  swisseph.Saturn,
	"Uranus":  swisseph.Uranus,
	"Neptune": swisseph.Neptune,
	"Pluto":   swisseph.Pluto,
}

// CalculatePlanets calculates the positions of requested planets.
func CalculatePlanets(ctx *domain.CalculationContext) ([]domain.PlanetPosition, error) {
	var positions []domain.PlanetPosition

	// Set Ayanamsa for this calculation
	swisseph.SetSidMode(int32(ctx.Config.AyanamsaMode), 0, 0)

	// Pre-calculate ayanamsa value
	ctx.Ayanamsa = swisseph.GetAyanamsaUT(ctx.JulianDayUT)

	planetNames := []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn", "Uranus", "Neptune", "Pluto"}

	for _, pName := range planetNames {
		seID := PlanetMap[pName]
		pos, err := calculateSinglePlanet(pName, seID, ctx)
		if err != nil {
			return nil, err
		}
		positions = append(positions, pos)
	}

	// Calculate Nodes (Mean or True based on config, using Mean for now unless specified)
	rahuPos, err := calculateSinglePlanet("Rahu", swisseph.MeanNode, ctx)
	if err != nil {
		return nil, err
	}
	// Ketu is Rahu + 180 degrees
	ketuPos := rahuPos
	ketuPos.Planet = "Ketu"
	ketuPos.TropicalLongitude = math.Mod(rahuPos.TropicalLongitude+180.0, 360.0)
	ketuPos.SiderealLongitude = math.Mod(rahuPos.SiderealLongitude+180.0, 360.0)
	ketuPos.Latitude = -rahuPos.Latitude // Opposite latitude

	setDMS(&ketuPos)

	positions = append(positions, rahuPos, ketuPos)

	return positions, nil
}

func calculateSinglePlanet(name string, seID int, ctx *domain.CalculationContext) (domain.PlanetPosition, error) {
	// Calculate Tropical (flag = FlagSwieph | FlagSpeed)
	iflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed)
	res := swisseph.CalcUT(ctx.JulianDayUT, int32(seID), iflag)

	if (res.Flag & swisseph.FlagSwieph) == 0 {
		return domain.PlanetPosition{}, fmt.Errorf("failed to use Swiss Ephemeris for %s, fell back to another model (Moshier). Check EPHE_PATH", name)
	}
	if res.Error != "" && strings.Contains(strings.ToLower(res.Error), "moshier") {
		return domain.PlanetPosition{}, fmt.Errorf("Swiss Ephemeris calculation warned about Moshier fallback for %s: %s", name, res.Error)
	}

	tropicalLon := res.Data[0]
	lat := res.Data[1]
	dist := res.Data[2]
	speed := res.Data[3]

	// Calculate Sidereal
	iflagSid := int32(swisseph.FlagSwieph | swisseph.FlagSpeed | swisseph.FlagSidereal)
	resSid := swisseph.CalcUT(ctx.JulianDayUT, int32(seID), iflagSid)
	siderealLon := resSid.Data[0]

	pos := domain.PlanetPosition{
		Planet:            name,
		TropicalLongitude: tropicalLon,
		SiderealLongitude: siderealLon,
		Latitude:          lat,
		Distance:          dist,
		Speed:             speed,
		Retrograde:        speed < 0,
	}

	setDMS(&pos)
	return pos, nil
}

func setDMS(pos *domain.PlanetPosition) {
	sign, deg, min, sec := time.DecimalToDMS(pos.SiderealLongitude)
	pos.Sign = sign
	pos.Degree = deg
	pos.Minute = min
	pos.Second = sec
	pos.DegreeInSign = float64(deg) + float64(min)/60.0 + sec/3600.0
}
