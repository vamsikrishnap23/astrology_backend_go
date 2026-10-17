package shadbala

import "math"

// Yuddha Bala (planetary war), Raman art. 76-78 (pp. 60-62) and summary (p. 107); BPHS 27.20 differs
// (it adds the difference of the two Shadbalas), and Raman's rule is followed here.
//
// Two planets are at war when they are in conjunction less than one degree apart; the Sun and Moon
// never fight. The one with the smaller longitude wins. The aggregate of Sthana Bala, Dig Bala and
// Kala Bala (up to Hora Bala) of each is found, and the difference between the two aggregates divided
// by the difference of the two planets' disc diameters is the Yuddha Bala: it is added to the Kala
// Bala of the winner and subtracted from that of the loser.

// bimbaParimana is the apparent diameter of each planet's disc in seconds of arc (art. 77).
var bimbaParimana = map[string]float64{
	"Mars":    9.4,
	"Mercury": 6.6,
	"Jupiter": 190.4,
	"Venus":   16.6,
	"Saturn":  158.0,
}

// yuddhaOrbit is the largest distance, in degrees, at which two planets are at war (art. 76).
const yuddhaOrbit = 1.0

// yuddhaBala returns the Yuddha Bala (shashtiamsas) of each planet at war: positive for a winner,
// negative for a loser, summed over all the wars a planet is in. aggregate holds each planet's
// Sthana + Dig + Kala (up to Hora) Bala. Planets not at war are absent from the result.
func yuddhaBala(lons map[string]float64, aggregate map[string]float64) map[string]float64 {
	fighters := []string{"Mars", "Mercury", "Jupiter", "Venus", "Saturn"}
	res := map[string]float64{}
	for i := 0; i < len(fighters); i++ {
		for j := i + 1; j < len(fighters); j++ {
			a, b := fighters[i], fighters[j]
			la, okA := lons[a]
			lb, okB := lons[b]
			if !okA || !okB {
				continue
			}
			// signed distance a - b in (-180, 180]
			d := math.Mod(la-lb+540.0, 360.0) - 180.0
			if math.Abs(d) >= yuddhaOrbit {
				continue
			}
			winner, loser := a, b
			if d > 0 { // a is ahead, so b has the smaller longitude
				winner, loser = b, a
			}
			gap := math.Abs(aggregate[a]-aggregate[b]) / math.Abs(bimbaParimana[a]-bimbaParimana[b])
			res[winner] += gap
			res[loser] -= gap
		}
	}
	return res
}
