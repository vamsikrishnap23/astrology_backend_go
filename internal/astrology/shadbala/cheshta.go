package shadbala

func CalculateCheshtaBala(planet string, isRetrograde bool, ayanaBala float64) float64 {
	if planet == "Sun" || planet == "Moon" {
		return ayanaBala
	}
	if isRetrograde {
		return 60.0
	}
	// Direct motion. Classically involves complex calculations based on average speed and arcs from the Sun.
	// For baseline implementation, stationary/slow can be 30, normal direct 0.
	// We'll assign 30 as a default for direct motion (half strength) if speed isn't evaluated, or 0.
	// BPHS typically assigns 30 to somewhat slow direct, 15 to fast. Let's return 0 for now as standard direct default, or 30 for baseline.
	// Actually, many software grant 30 points for direct motion, 60 for retrograde.
	return 30.0
}
