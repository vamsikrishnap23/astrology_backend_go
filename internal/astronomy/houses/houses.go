package houses

import (
	"fmt"
	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"math"
)

// CalculateHouses calculates the Ascendant, MC and House cusps.
func CalculateHouses(ctx *domain.CalculationContext) (float64, float64, []domain.HouseCusp, error) {

	// Serialise access to the C library and pin this goroutine to one OS thread (see ephemeris.Lock)
	ephemeris.Lock()
	defer ephemeris.Unlock()

	// Set Ayanamsa for Sidereal House calculation
	swisseph.SetSidMode(int32(ctx.Config.AyanamsaMode), 0, 0)

	// Pre-calculate ayanamsa value
	ctx.Ayanamsa = ephemeris.AyanamsaUT(ctx.JulianDayUT)

	calcJD := ctx.JulianDayUT

	if ctx.Input.HoraryNumber > 0 && ctx.Input.HoraryNumber <= 249 {
		targetSiderealAsc, _ := GetHoraryAscendant(ctx.Input.HoraryNumber)
		targetTropAsc := math.Mod(targetSiderealAsc+ctx.Ayanamsa, 360.0)

		phantomJD, err := findPhantomJD(ctx.JulianDayUT, targetTropAsc, ctx.Input.Latitude, ctx.Input.Longitude, ctx.Config.HouseCode)
		if err != nil {
			return 0, 0, nil, err
		}
		calcJD = phantomJD
	}

	// Calculate Tropical houses using either real time or Phantom time
	res := swisseph.HousesEx(calcJD, 0, ctx.Input.Latitude, ctx.Input.Longitude, ctx.Config.HouseCode)

	// Convert points to Sidereal
	ascendant := res.Points[0] - ctx.Ayanamsa
	if ascendant < 0 {
		ascendant += 360.0
	}
	mc := res.Points[1] - ctx.Ayanamsa
	if mc < 0 {
		mc += 360.0
	}

	var nakshatras = []string{
		"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
		"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni",
		"Hasta", "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha",
		"Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha",
		"Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
	}
	var nakshatraLords = []string{"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury"}

	var houseCusps []domain.HouseCusp
	for i := 0; i < len(res.Houses) && i < 12; i++ {
		siderealLon := res.Houses[i] - ctx.Ayanamsa
		for siderealLon < 0 {
			siderealLon += 360.0
		}
		for siderealLon >= 360.0 {
			siderealLon -= 360.0
		}

		sign, deg, min, sec := time.DecimalToDMS(siderealLon)

		interval := 13.0 + 1.0/3.0
		nakIdx := int(math.Floor(siderealLon / interval))
		nakName := nakshatras[nakIdx%27]
		nakLord := nakshatraLords[nakIdx%9]
		nakProgress := math.Mod(siderealLon, interval) / interval * 100.0
		pada := int(math.Floor(nakProgress/25.0)) + 1

		houseCusps = append(houseCusps, domain.HouseCusp{
			HouseNumber:   i + 1,
			Longitude:     siderealLon,
			Sign:          sign,
			Degree:        deg,
			Minute:        min,
			Second:        sec,
			Nakshatra:     nakName,
			NakshatraPada: pada,
			NakshatraLord: nakLord,
		})
	}

	return ascendant, mc, houseCusps, nil
}

// findPhantomJD finds the Julian Day nearest to startJD at which the tropical
// ascendant equals targetTropAsc (used for KP horary numbers). The ascendant
// advances monotonically through 360 degrees about once per sidereal day but at
// a rate that varies with sign and latitude, so the root is bracketed by a
// coarse scan over +/- half a day before bisecting. The caller must hold ephemeris.Lock.
func findPhantomJD(startJD, targetTropAsc, lat, lon float64, houseCode byte) (float64, error) {
	const (
		window = 0.5  // days either side of startJD; one full ascendant cycle spans ~1 day
		step   = 0.02 // ~29 minutes; the ascendant moves well under 180 degrees per step
	)
	// signed angular offset of the ascendant from the target, in (-180, 180]
	offset := func(jd float64) float64 {
		asc := swisseph.HousesEx(jd, 0, lat, lon, houseCode).Points[0]
		return math.Mod(asc-targetTropAsc+540.0, 360.0) - 180.0
	}

	var (
		bestLow, bestHigh float64
		bestDist          = math.Inf(1)
		found             bool
	)

	prevJD := startJD - window
	prev := offset(prevJD)
	for jd := prevJD + step; jd <= startJD+window+1e-9; jd += step {
		cur := offset(jd)
		// A genuine crossing goes from negative to non-negative with a small jump;
		// the +180 -> -180 wrap-around is rejected by prev < 0 and the span check.
		if prev < 0 && cur >= 0 && cur-prev < 180.0 {
			dist := math.Abs((prevJD+jd)/2.0 - startJD)
			if dist < bestDist {
				bestDist, bestLow, bestHigh, found = dist, prevJD, jd, true
			}
		}
		prevJD, prev = jd, cur
	}
	if !found {
		return 0, fmt.Errorf("could not resolve ascendant %.4f for horary number at lat %.4f", targetTropAsc, lat)
	}

	low, high := bestLow, bestHigh
	for i := 0; i < 60; i++ {
		mid := (low + high) / 2.0
		if offset(mid) < 0 {
			low = mid
		} else {
			high = mid
		}
	}
	root := (low + high) / 2.0
	// Beyond the polar circle the ascendant is discontinuous, and the bracket can
	// collapse onto the jump instead of a true root. Never return that silently.
	if math.Abs(offset(root)) > 1e-4 {
		return 0, fmt.Errorf("horary ascendant %.4f is not attainable at lat %.4f (ascendant discontinuous)", targetTropAsc, lat)
	}
	return root, nil
}
