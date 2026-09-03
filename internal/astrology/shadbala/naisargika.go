package shadbala

var NaisargikaBala = map[string]float64{
	"Sun":     60.00,
	"Moon":    51.43,
	"Venus":   42.85,
	"Jupiter": 34.28,
	"Mercury": 25.71,
	"Mars":    17.14,
	"Saturn":  8.57,
}

func CalculateNaisargikaBala(planet string) float64 {
	if val, ok := NaisargikaBala[planet]; ok {
		return val
	}
	return 0.0
}
