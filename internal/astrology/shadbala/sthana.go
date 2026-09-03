package shadbala

import (
	"math"

	"github.com/vamsi/astrology_backend_go/internal/domain"
)

var ExaltationPoints = map[string]float64{
	"Sun":     10.0,
	"Moon":    33.0,  // Taurus 3 = 30 + 3
	"Mars":    298.0, // Capricorn 28 = 270 + 28
	"Mercury": 165.0, // Virgo 15 = 150 + 15
	"Jupiter": 95.0,  // Cancer 5 = 90 + 5
	"Venus":   357.0, // Pisces 27 = 330 + 27
	"Saturn":  200.0, // Libra 20 = 180 + 20
}

func calculateUchchaBala(planet string, longitude float64) float64 {
	exaltation, ok := ExaltationPoints[planet]
	if !ok {
		return 0
	}
	debilitation := math.Mod(exaltation+180.0, 360.0)

	arc := math.Abs(longitude - debilitation)
	if arc > 180.0 {
		arc = 360.0 - arc
	}
	return arc / 3.0
}

func calculateSaptavargajaBala(planet string, vargas domain.VargasResult, natalSigns map[string]int) float64 {
	total := 0.0
	targetDivisions := []int{1, 2, 3, 7, 9, 12, 30}

	for _, div := range targetDivisions {
		for _, varga := range vargas.Vargas {
			if varga.Division == div {
				// Find planet in this varga
				for _, vp := range varga.Planets {
					if vp.Planet == planet {
						rel := GetPlacementStrength(planet, vp.DivisionalSign, natalSigns)
						switch rel {
						case Moolatrikona:
							total += 45.0
						case OwnHouse:
							total += 30.0
						case AdhiMitra:
							total += 22.5
						case Mitra:
							total += 15.0
						case Sama:
							total += 7.5
						case Shatru:
							total += 3.75
						case AdhiShatru:
							total += 1.875
						}
					}
				}
			}
		}
	}
	return total
}

func calculateOjayugmaBala(planet string, rasiSign, navamshaSign int) float64 {
	isRasiEven := (rasiSign % 2) != 0 // 0=Aries(Odd), 1=Taurus(Even)
	isNavamshaEven := (navamshaSign % 2) != 0

	bala := 0.0
	if planet == "Venus" || planet == "Moon" {
		if isRasiEven {
			bala += 15.0
		}
		if isNavamshaEven {
			bala += 15.0
		}
	} else {
		if !isRasiEven {
			bala += 15.0
		}
		if !isNavamshaEven {
			bala += 15.0
		}
	}
	return bala
}

func calculateKendradiBala(planet string, rasiSign, ascSign int) float64 {
	house := (rasiSign-ascSign+12)%12 + 1 // 1 to 12
	if house == 1 || house == 4 || house == 7 || house == 10 {
		return 60.0
	} else if house == 2 || house == 5 || house == 8 || house == 11 {
		return 30.0
	}
	return 15.0
}

func calculateDrekkanaBala(planet string, longitude float64) float64 {
	degreeInSign := math.Mod(longitude, 30.0)
	drekkana := 1
	if degreeInSign >= 20.0 {
		drekkana = 3
	} else if degreeInSign >= 10.0 {
		drekkana = 2
	}

	if (planet == "Sun" || planet == "Mars" || planet == "Jupiter") && drekkana == 1 {
		return 15.0
	} else if (planet == "Moon" || planet == "Venus") && drekkana == 2 {
		return 15.0
	} else if (planet == "Mercury" || planet == "Saturn") && drekkana == 3 {
		return 15.0
	}
	return 0.0
}
