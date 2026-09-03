package shadbala

import "math"

func getAspectValue(aspecting string, d float64) float64 {
	var v float64
	if d <= 30.0 {
		v = 0
	} else if d <= 60.0 {
		v = (d - 30.0) / 2.0
		if aspecting == "Saturn" {
			v = (d - 30.0) * 2.0
		} // 60 at 60
	} else if d <= 90.0 {
		v = (90.0 - d) / 2.0
		if aspecting == "Saturn" {
			v = (90.0 - d) * 2.0
		}
		if aspecting == "Mars" {
			v = (d - 60.0) * 2.0
		} // 60 at 90
	} else if d <= 120.0 {
		v = (d - 90.0)
		if aspecting == "Mars" {
			v = (120.0 - d) * 2.0
		}
		if aspecting == "Jupiter" {
			v = (d - 90.0) * 2.0
		} // 60 at 120
	} else if d <= 150.0 {
		v = (150.0 - d)
		if aspecting == "Jupiter" {
			v = (150.0 - d) * 2.0
		}
	} else if d <= 180.0 {
		v = (d - 150.0) * 2.0
	} else if d <= 210.0 {
		v = (300.0 - d) / 2.0
		if aspecting == "Mars" {
			v = (d - 180.0) * 2.0
		} // 60 at 210
	} else if d <= 240.0 {
		v = (300.0 - d) / 2.0
		if aspecting == "Mars" {
			v = (240.0 - d) * 2.0
		}
		if aspecting == "Jupiter" {
			v = (d - 210.0) * 2.0
		} // 60 at 240
	} else if d <= 270.0 {
		v = (300.0 - d) / 2.0
		if aspecting == "Jupiter" {
			v = (270.0 - d) * 2.0
		}
		if aspecting == "Saturn" {
			v = (d - 240.0) * 2.0
		} // 60 at 270
	} else if d <= 300.0 {
		v = (300.0 - d) / 2.0
		if aspecting == "Saturn" {
			v = (300.0 - d) * 2.0
		}
	} else {
		v = 0
	}

	if v < 0 {
		v = 0
	}
	if v > 60 {
		v = 60
	}
	return v
}

func CalculateDrikBala(planet string, targetLon float64, planets map[string]float64) float64 {
	totalV := 0.0

	for aspecting, aspectingLon := range planets {
		if aspecting == planet {
			continue
		}
		if aspecting == "Uranus" || aspecting == "Neptune" || aspecting == "Pluto" || aspecting == "Rahu" || aspecting == "Ketu" {
			continue
		}

		d := math.Mod(targetLon-aspectingLon+360.0, 360.0)
		v := getAspectValue(aspecting, d)

		isBenefic := false
		if aspecting == "Jupiter" || aspecting == "Venus" {
			isBenefic = true
		} else if aspecting == "Mercury" {
			// For Drik Bala, Mercury is usually a benefic if alone, we'll consider it benefic here.
			isBenefic = true
		} else if aspecting == "Moon" {
			isBenefic = true // Waning/Waxing rules apply, assume benefic for simplicity unless further defined
		}

		if isBenefic {
			totalV += v
		} else {
			totalV -= v
		}
	}

	return totalV / 4.0
}
