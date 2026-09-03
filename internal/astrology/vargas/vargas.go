package vargas

import (
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"math"
)

var signNames = []string{
	"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
	"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

var signLords = []string{
	"Mars", "Venus", "Mercury", "Moon", "Sun", "Mercury",
	"Venus", "Mars", "Jupiter", "Saturn", "Saturn", "Jupiter",
}

// CalculateVargas generates all registered Varga charts based on planetary sidereal longitudes and Ascendant.
func CalculateVargas(tables domain.TablesResult, cusps []domain.HouseCusp) domain.VargasResult {
	var res domain.VargasResult

	// Find the Ascendant from the House Table (House 1)
	var ascLon float64
	var ascTableHouse *domain.TableHouse
	for _, h := range tables.HouseTable {
		if h.HouseNumber == 1 {
			ascTableHouse = &h

			// Reconstruct absolute sidereal longitude of Ascendant based on sign + degree + minute + second
			signIdx := -1
			for i, s := range signNames {
				if s == h.Sign {
					signIdx = i
					break
				}
			}
			if signIdx != -1 {
				ascLon = float64(signIdx*30) + float64(h.Degree) + float64(h.Minute)/60.0 + h.Second/3600.0
			}
			break
		}
	}

	for _, rule := range Registry {
		chart := domain.VargaChart{
			Division: rule.Division(),
			Name:     rule.Name(),
		}

		// Calculate for Ascendant
		if ascTableHouse != nil {
			pos := rule.Calculate(ascLon)
			// Need to convert the pos.LongitudeInDivision (0-30) into deg/min/sec
			_, deg, min, sec := astronomyTime.DecimalToDMS(pos.LongitudeInDivision)
			nakName, nakPada, nakLord := getNakshatraInfo(pos.SignIndex, pos.LongitudeInDivision)

			chart.Ascendant = domain.VargaPlanet{
				Planet:          "Ascendant",
				SourceLongitude: ascLon,
				DivisionalSign:  signNames[pos.SignIndex],
				Degree:          deg,
				Minute:          min,
				Second:          sec,
				Nakshatra:       nakName,
				NakshatraPada:   nakPada,
				NakshatraLord:   nakLord,
				SignLord:        signLords[pos.SignIndex],
				Retrograde:      false,
			}
		}

		// Calculate for each Planet
		for _, p := range tables.PlanetaryTable {
			pos := rule.Calculate(p.ExactLongitude)
			_, deg, min, sec := astronomyTime.DecimalToDMS(pos.LongitudeInDivision)
			nakName, nakPada, nakLord := getNakshatraInfo(pos.SignIndex, pos.LongitudeInDivision)

			chart.Planets = append(chart.Planets, domain.VargaPlanet{
				Planet:          p.PlanetName,
				SourceLongitude: p.ExactLongitude,
				DivisionalSign:  signNames[pos.SignIndex],
				Degree:          deg,
				Minute:          min,
				Second:          sec,
				Nakshatra:       nakName,
				NakshatraPada:   nakPada,
				NakshatraLord:   nakLord,
				SignLord:        signLords[pos.SignIndex],
				Retrograde:      p.Retrograde,
			})
		}

		// Calculate fractional placement for all mathematical house cusps in this divisional chart
		var vargaHouses []domain.HouseCusp
		for _, hc := range cusps {
			cPos := rule.Calculate(hc.Longitude)
			_, cDeg, cMin, cSec := astronomyTime.DecimalToDMS(cPos.LongitudeInDivision)
			nakName, nakPada, nakLord := getNakshatraInfo(cPos.SignIndex, cPos.LongitudeInDivision)
			absLon := float64(cPos.SignIndex*30) + cPos.LongitudeInDivision

			vargaHouses = append(vargaHouses, domain.HouseCusp{
				HouseNumber:   hc.HouseNumber,
				Longitude:     absLon, // Pass the projected mathematical longitude instead of the source longitude!
				Sign:          signNames[cPos.SignIndex],
				Degree:        cDeg,
				Minute:        cMin,
				Second:        cSec,
				Nakshatra:     nakName,
				NakshatraPada: nakPada,
				NakshatraLord: nakLord,
			})
		}
		chart.Houses = vargaHouses

		res.Vargas = append(res.Vargas, chart)
	}

	return res
}

var nakshatras = []string{
	"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
	"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni",
	"Hasta", "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha",
	"Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha",
	"Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
}

var nakshatraLords = []string{"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury"}

func getNakshatraInfo(signIndex int, lonInDiv float64) (string, int, string) {
	absLon := float64(signIndex*30) + lonInDiv
	interval := 13.0 + 1.0/3.0
	nakIdx := int(math.Floor(absLon / interval))
	nakName := nakshatras[nakIdx%27]
	nakLord := nakshatraLords[nakIdx%9]
	nakProgress := math.Mod(absLon, interval) / interval * 100.0
	pada := int(math.Floor(nakProgress/25.0)) + 1
	return nakName, pada, nakLord
}
