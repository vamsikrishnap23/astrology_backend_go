package transit

import (
	"math"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/planets"
	astrotime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

var signNames = []string{
	"Aries", "Taurus", "Gemini", "Cancer",
	"Leo", "Virgo", "Libra", "Scorpio",
	"Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

var targetPlanets = []string{
	"Sun", "Moon", "Mars", "Mercury",
	"Jupiter", "Venus", "Saturn", "Rahu", "Ketu",
}

func getSiderealLon(jd float64, planet string, ayanamsaMode int) float64 {
	ephemeris.Mu.Lock()
	defer ephemeris.Mu.Unlock()
	swisseph.SetEphePath(ephemeris.EphePath)

	swisseph.SetSidMode(int32(ayanamsaMode), 0, 0)
	ayanamsa := swisseph.GetAyanamsaUT(jd)

	var seID int
	var isKetu bool
	if planet == "Rahu" {
		seID = swisseph.MeanNode
	} else if planet == "Ketu" {
		seID = swisseph.MeanNode
		isKetu = true
	} else {
		seID = planets.PlanetMap[planet]
	}

	iflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed)
	res := swisseph.CalcUT(jd, int32(seID), iflag)

	tropicalLon := res.Data[0]
	if isKetu {
		tropicalLon = math.Mod(tropicalLon+180.0, 360.0)
	}

	siderealLon := math.Mod(tropicalLon-ayanamsa, 360.0)
	if siderealLon < 0 {
		siderealLon += 360.0
	}
	return siderealLon
}

func findNextIngress(startJD float64, planet string, ayanamsaMode int) (float64, int) {
	// Step size based on planet speed to not miss retrograde boundary jumps
	step := 1.0 // 1 day
	if planet == "Moon" {
		step = 0.1 // Moon moves ~13 deg/day, so 0.1 day is safe
	}

	startLon := getSiderealLon(startJD, planet, ayanamsaMode)
	startSign := int(math.Floor(startLon / 30.0))

	currentJD := startJD
	var boundaryLon float64
	var destSign int

	// Fast forward until sign changes
	for {
		currentJD += step
		lon := getSiderealLon(currentJD, planet, ayanamsaMode)
		sign := int(math.Floor(lon / 30.0))

		if sign != startSign {
			destSign = sign

			// Find the boundary we crossed
			if sign == (startSign+1)%12 {
				// Forward cross
				boundaryLon = float64(sign * 30)
			} else if sign == (startSign-1+12)%12 {
				// Retrograde cross
				boundaryLon = float64(startSign * 30)
			} else {
				// Large jump? Adjust boundary (e.g. Moon fast movement)
				if (sign-startSign+12)%12 < 6 {
					boundaryLon = float64(sign * 30) // Forward
				} else {
					boundaryLon = float64(startSign * 30) // Backward
				}
			}

			break
		}

		// Guard against infinite loop (Saturn takes ~2.5 years, so 1000 days is max. Loop limit to 2000)
		if currentJD-startJD > 2000 {
			break
		}
	}

	// Bisection to find exact time of crossing
	low := currentJD - step
	high := currentJD
	for i := 0; i < 50; i++ {
		mid := (low + high) / 2.0
		lon := getSiderealLon(mid, planet, ayanamsaMode)

		diff := math.Mod(lon-boundaryLon+540.0, 360.0) - 180.0

		if math.Abs(diff) < 1e-6 {
			return mid, destSign
		}

		// We need to know the direction of movement to adjust low/high.
		// A simple way is to check the sign at mid.
		midSign := int(math.Floor(lon / 30.0))
		if midSign == startSign {
			low = mid
		} else {
			high = mid
		}
	}

	return (low + high) / 2.0, destSign
}

func CalculateUpcomingTransits(startJD float64, ayanamsaMode int) []domain.UpcomingTransitPlanet {
	var results []domain.UpcomingTransitPlanet

	for _, p := range targetPlanets {
		ingressJD, destSign := findNextIngress(startJD, p, ayanamsaMode)

		utcTime := astrotime.JulianDayToUTC(ingressJD)

		results = append(results, domain.UpcomingTransitPlanet{
			Planet:             p,
			DestinationSign:    signNames[destSign],
			TransitionDateTime: utcTime.Format(time.RFC3339),
		})
	}

	return results
}
