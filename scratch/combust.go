package main

import (
	"fmt"
	"math"
)

func angularDistance(lon1, lon2 float64) float64 {
	dist := math.Abs(lon1 - lon2)
	if dist > 180.0 {
		dist = 360.0 - dist
	}
	return dist
}

func isCombust(planet string, dist float64, isRetrograde bool) bool {
	switch planet {
	case "Moon":
		return dist <= 12.0
	case "Mars":
		return dist <= 17.0
	case "Mercury":
		if isRetrograde {
			return dist <= 12.0
		}
		return dist <= 14.0
	case "Jupiter":
		return dist <= 11.0
	case "Venus":
		if isRetrograde {
			return dist <= 8.0
		}
		return dist <= 10.0
	case "Saturn":
		return dist <= 15.0
	default:
		return false // Rahu, Ketu, Outer planets, etc.
	}
}

func main() {
	fmt.Println(isCombust("Mars", 16.5, false))
}
