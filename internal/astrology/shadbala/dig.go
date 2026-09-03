package shadbala

import "math"

func CalculateDigBala(planet string, longitude, ascendant, midheaven float64) float64 {
	var maxPowerPoint float64

	switch planet {
	case "Sun", "Mars":
		maxPowerPoint = midheaven
	case "Jupiter", "Mercury":
		maxPowerPoint = ascendant
	case "Moon", "Venus":
		maxPowerPoint = math.Mod(midheaven+180.0, 360.0)
	case "Saturn":
		maxPowerPoint = math.Mod(ascendant+180.0, 360.0)
	default:
		return 0
	}

	zeroPowerPoint := math.Mod(maxPowerPoint+180.0, 360.0)

	arc := math.Abs(longitude - zeroPowerPoint)
	if arc > 180.0 {
		arc = 360.0 - arc
	}

	return arc / 3.0
}
