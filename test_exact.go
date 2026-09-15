package main

import (
	"fmt"
	"math"
)

func main() {
	pLon := 265.00611 // Prog Moon 25:00:22 Sagittarius
	nSun := 326.18527 // Natal Sun 26:11:07 Aquarius
	nVe  := 355.52722 // Natal Venus 25:31:38 Pisces

	diffSun := math.Abs(pLon - nSun)
	diffSun = math.Mod(diffSun, 360.0)
	if diffSun > 180.0 {
		diffSun = 360.0 - diffSun
	}

	diffVe := math.Abs(pLon - nVe)
	diffVe = math.Mod(diffVe, 360.0)
	if diffVe > 180.0 {
		diffVe = 360.0 - diffVe
	}

	fmt.Printf("Diff Sun: %.3f\n", diffSun)
	fmt.Printf("Diff Venus: %.3f\n", diffVe)
}
