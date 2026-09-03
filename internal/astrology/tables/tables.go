package tables

import (
	"math"

	astrotime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

var signLords = map[string]string{
	"Aries":       "Mars",
	"Taurus":      "Venus",
	"Gemini":      "Mercury",
	"Cancer":      "Moon",
	"Leo":         "Sun",
	"Virgo":       "Mercury",
	"Libra":       "Venus",
	"Scorpio":     "Mars",
	"Sagittarius": "Jupiter",
	"Capricorn":   "Saturn",
	"Aquarius":    "Saturn",
	"Pisces":      "Jupiter",
}

var nakshatraNames = []string{
	"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
	"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni",
	"Uttara Phalguni", "Hasta", "Chitra", "Swati", "Vishakha", "Anuradha",
	"Jyeshtha", "Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana",
	"Dhanishta", "Shatabhisha", "Purva Bhadrapada", "Uttara Bhadrapada",
	"Revati",
}

var nakshatraLords = []string{
	"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
	"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
	"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
}

func GenerateTables(planets []domain.PlanetPosition, houses []domain.HouseCusp) domain.TablesResult {
	var planetaryTable []domain.TablePlanet
	var houseTable []domain.TableHouse

	// Build House Table first to know the house boundaries
	houseLen := len(houses)

	// Helper to find which house a degree falls into
	getHouseNumber := func(degree float64) int {
		for i := 0; i < houseLen; i++ {
			nextI := (i + 1) % houseLen

			h1 := houses[i].Longitude
			h2 := houses[nextI].Longitude

			// Handle 360 wrap around
			if h1 < h2 {
				if degree >= h1 && degree < h2 {
					return houses[i].HouseNumber
				}
			} else {
				if degree >= h1 || degree < h2 {
					return houses[i].HouseNumber
				}
			}
		}
		return 1 // Fallback
	}

	// 1. Planetary Table
	for _, p := range planets {
		// Calculate Nakshatra from Sidereal Longitude
		interval := 13.0 + 1.0/3.0
		nakIdx := int(math.Floor(p.SiderealLongitude / interval))
		nakProgress := math.Mod(p.SiderealLongitude, interval) / interval * 100.0
		pada := int(math.Floor(nakProgress/25.0)) + 1

		if nakIdx >= 27 {
			nakIdx = 0
		} // wrap just in case

		nakName := nakshatraNames[nakIdx]
		hNum := getHouseNumber(p.SiderealLongitude)

		sl, nl, ssl, sssl, ssssl := GetKPLords(p.SiderealLongitude)

		tp := domain.TablePlanet{
			PlanetName:     p.Planet,
			Sign:           p.Sign,
			Degree:         p.Degree,
			Minute:         p.Minute,
			Second:         p.Second,
			ExactLongitude: p.SiderealLongitude,
			Retrograde:     p.Retrograde,
			Speed:          p.Speed,
			HouseNumber:    hNum,
			Nakshatra:      nakName,
			NakshatraPada:  pada,
			SignLord:       sl,
			NakshatraLord:  nl,
			SubLord:        ssl,
			SubSubLord:     sssl,
			SubSubSubLord:  ssssl,
		}
		planetaryTable = append(planetaryTable, tp)
	}

	// 2. House Table
	for _, h := range houses {
		// Format cusp degrees
		sign, deg, min, sec := astrotime.DecimalToDMS(h.Longitude)

		// Find occupants
		var occupants []string
		for _, p := range planetaryTable {
			if p.HouseNumber == h.HouseNumber {
				// Avoid duplicate formatting or just push name
				// Let's just put planet names
				occupants = append(occupants, p.PlanetName)
			}
		}

		if occupants == nil {
			occupants = []string{}
		}

		sl, nl, ssl, sssl, ssssl := GetKPLords(h.Longitude)

		th := domain.TableHouse{
			HouseNumber:   h.HouseNumber,
			CuspLongitude: h.Longitude,
			Sign:          sign,
			Degree:        deg,
			Minute:        min,
			Second:        sec,
			Nakshatra:     h.Nakshatra,
			NakshatraPada: h.NakshatraPada,
			SignLord:      sl,
			NakshatraLord: nl,
			SubLord:       ssl,
			SubSubLord:    sssl,
			SubSubSubLord: ssssl,
			Occupants:     occupants,
		}
		houseTable = append(houseTable, th)
	}

	return domain.TablesResult{
		PlanetaryTable: planetaryTable,
		HouseTable:     houseTable,
	}
}

// Helper to determine if planet is near boundary?
func isNearBoundary(degree float64) bool {
	// Not strictly required for the output table but user asked to test it
	mod := math.Mod(degree, 30.0)
	return mod < 1.0 || mod > 29.0
}
