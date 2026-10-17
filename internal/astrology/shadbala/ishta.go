package shadbala

import (
	"math"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// Ishta and Kashta Phala (benefic and malefic tendency of a planet), Raman chapter X (art. 135-140, pp. 101-104).
// Ishta Phala is the square root of Uchcha Bala times Chesta Bala; Kashta Phala is the square root of
// (60 - Uchcha Bala) times (60 - Chesta Bala). A planet with more Ishta than Kashta is inclined to do good
// in its dasa and bhukti, one with more Kashta to do harm. BPHS ch. 28 gives a longer method with the
// Uchcha and Chesta Rasmis; Raman's is followed.

// sunChestaBala is the Sun's Chesta Bala for Ishta and Kashta Phala (art. 136): 90 deg are added to the
// Sun's Sayana (tropical) longitude to get its Chesta Kendra, which is reduced to 0..180 and divided by 3.
func sunChestaBala(tropicalLon float64) float64 {
	k := norm360(tropicalLon + 90.0)
	if k > 180.0 {
		k = 360.0 - k
	}
	return k / 3.0
}

// moonChestaBala is the Moon's Chesta Bala (art. 137): the arc from the Sun, reduced to 0..180, divided by 3.
func moonChestaBala(sunLon, moonLon float64) float64 {
	k := norm360(moonLon - sunLon)
	if k > 180.0 {
		k = 360.0 - k
	}
	return k / 3.0
}

// ishtaKashta returns the two phalas from Uchcha Bala and Chesta Bala (art. 138-139).
func ishtaKashta(uchcha, chesta float64) (ishta, kashta float64) {
	return math.Sqrt(uchcha * chesta), math.Sqrt((60.0 - uchcha) * (60.0 - chesta))
}

// calculateIshtaKashta gives the phalas of the seven planets; uchcha and cheshta are their Uchcha and Cheshta
// Balas from the Shadbala, sunTropical the Sun's tropical longitude and sunLon, moonLon the sidereal ones.
func calculateIshtaKashta(uchcha, cheshta map[string]float64, sunTropical, sunLon, moonLon float64) []domain.IshtaKashtaPlanet {
	out := make([]domain.IshtaKashtaPlanet, 0, len(TargetPlanets))
	for _, p := range TargetPlanets {
		chesta := cheshta[p]
		switch p {
		case "Sun":
			chesta = sunChestaBala(sunTropical)
		case "Moon":
			chesta = moonChestaBala(sunLon, moonLon)
		}
		ishta, kashta := ishtaKashta(uchcha[p], chesta)
		incl := "kashta"
		if ishta > kashta {
			incl = "ishta"
		}
		out = append(out, domain.IshtaKashtaPlanet{
			Planet: p, UchchaBala: round2(uchcha[p]), ChestaBala: round2(chesta),
			IshtaPhala: round2(ishta), KashtaPhala: round2(kashta), Inclination: incl,
		})
	}
	return out
}
