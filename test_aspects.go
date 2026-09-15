package main

import (
	"fmt"
	"math"
)

func main() {
	pLon := 265.006 // Prog Moon
	nLons := map[string]float64{
		"Sun":   326.185,
		"Venus": 355.527,
	}

	for name, nLon := range nLons {
		diff := math.Abs(pLon - nLon)
		diff = math.Mod(diff, 360.0)
		if diff > 180.0 {
			diff = 360.0 - diff
		}

		orbMax := 1.0
		var aspectType string
		var exactAngle float64

		if diff <= orbMax {
			aspectType = "Conjunction"
			exactAngle = 0.0
		} else if math.Abs(diff-60.0) <= orbMax {
			aspectType = "Sextile"
			exactAngle = 60.0
		} else if math.Abs(diff-90.0) <= orbMax {
			aspectType = "Square"
			exactAngle = 90.0
		}

		if aspectType != "" {
			fmt.Printf("Aspect to %s: %s (orb %.3f)\n", name, aspectType, math.Abs(diff-exactAngle))
		} else {
			fmt.Printf("No aspect to %s (diff %.3f)\n", name, diff)
		}
	}
}
