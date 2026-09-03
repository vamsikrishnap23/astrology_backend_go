package planets

import (
	"math"

	"github.com/mshafiee/swephgo"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

var PlanetMap = map[string]int{
	"Sun":     swephgo.SeSun,
	"Moon":    swephgo.SeMoon,
	"Mars":    swephgo.SeMars,
	"Mercury": swephgo.SeMercury,
	"Jupiter": swephgo.SeJupiter,
	"Venus":   swephgo.SeVenus,
	"Saturn":  swephgo.SeSaturn,
	"Uranus":  swephgo.SeUranus,
	"Neptune": swephgo.SeNeptune,
	"Pluto":   swephgo.SePluto,
}

// CalculatePlanets calculates the positions of requested planets.
func CalculatePlanets(ctx *domain.CalculationContext) ([]domain.PlanetPosition, error) {
	var positions []domain.PlanetPosition

	// Set Ayanamsa for this calculation
	swephgo.SetSidMode(ctx.Config.AyanamsaMode, 0, 0)

	// Pre-calculate ayanamsa value
	ctx.Ayanamsa = swephgo.GetAyanamsaUt(ctx.JulianDayUT)

	planetNames := []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn", "Uranus", "Neptune", "Pluto"}

	for _, pName := range planetNames {
		seID := PlanetMap[pName]
		pos := calculateSinglePlanet(pName, seID, ctx)
		positions = append(positions, pos)
	}

	// Calculate Nodes (Mean or True based on config, using Mean for now unless specified)
	rahuPos := calculateSinglePlanet("Rahu", swephgo.SeMeanNode, ctx)
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

func calculateSinglePlanet(name string, seID int, ctx *domain.CalculationContext) domain.PlanetPosition {
	xx := make([]float64, 6)
	serr := make([]byte, 256)

	// Calculate Tropical (flag = SeflgSwieph | SeflgSpeed)
	iflag := swephgo.SeflgSwieph | swephgo.SeflgSpeed
	swephgo.CalcUt(ctx.JulianDayUT, seID, iflag, xx, serr)
	tropicalLon := xx[0]
	lat := xx[1]
	dist := xx[2]
	speed := xx[3]

	// Calculate Sidereal
	iflagSid := swephgo.SeflgSwieph | swephgo.SeflgSpeed | swephgo.SeflgSidereal
	swephgo.CalcUt(ctx.JulianDayUT, seID, iflagSid, xx, serr)
	siderealLon := xx[0]

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
	return pos
}

func setDMS(pos *domain.PlanetPosition) {
	sign, deg, min, sec := time.DecimalToDMS(pos.SiderealLongitude)
	pos.Sign = sign
	pos.Degree = deg
	pos.Minute = min
	pos.Second = sec
	pos.DegreeInSign = float64(deg) + float64(min)/60.0 + sec/3600.0
}
