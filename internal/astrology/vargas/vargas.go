package vargas

import (
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
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

			// We need to calculate Nakshatra inside the divisional chart if requested,
			// but wait, standard Varga tooltip usually means the "original" nakshatra or the nakshatra IN the divisional chart?
			// The requirements say: "Do not overwrite the original planetary longitude. Be careful to distinguish: source natal longitude, calculated divisional sign, degree within the divisional sign".
			// Let's map the degree inside the divisional sign. We won't re-calculate Nakshatra inside the varga unless explicitly requested. We'll use the original Nakshatra for the tooltip as metadata, or we can just leave Nakshatra as original.

			chart.Ascendant = domain.VargaPlanet{
				Planet:          "Ascendant",
				SourceLongitude: ascLon,
				DivisionalSign:  signNames[pos.SignIndex],
				Degree:          deg,
				Minute:          min,
				Second:          sec,
				Nakshatra:       ascTableHouse.NakshatraLord, // We can just provide the lord or original name
				NakshatraPada:   0,                           // Asc doesn't natively have pada mapped in TableHouse
				SignLord:        signLords[pos.SignIndex],
				Retrograde:      false,
			}
		}

		// Calculate for each Planet
		for _, p := range tables.PlanetaryTable {
			pos := rule.Calculate(p.ExactLongitude)
			_, deg, min, sec := astronomyTime.DecimalToDMS(pos.LongitudeInDivision)

			chart.Planets = append(chart.Planets, domain.VargaPlanet{
				Planet:          p.PlanetName,
				SourceLongitude: p.ExactLongitude,
				DivisionalSign:  signNames[pos.SignIndex],
				Degree:          deg,
				Minute:          min,
				Second:          sec,
				Nakshatra:       p.Nakshatra,
				NakshatraPada:   p.NakshatraPada,
				SignLord:        signLords[pos.SignIndex],
				Retrograde:      p.Retrograde,
			})
		}

		// Calculate fractional placement for all mathematical house cusps in this divisional chart
		var vargaHouses []domain.HouseCusp
		for _, hc := range cusps {
			cPos := rule.Calculate(hc.Longitude)
			_, cDeg, cMin, cSec := astronomyTime.DecimalToDMS(cPos.LongitudeInDivision)

			vargaHouses = append(vargaHouses, domain.HouseCusp{
				HouseNumber:   hc.HouseNumber,
				Longitude:     hc.Longitude, // Keep original longitude or maybe they want the divisional longitude? The frontend just wants the Sign, Degree, Minute, Second
				Sign:          signNames[cPos.SignIndex],
				Degree:        cDeg,
				Minute:        cMin,
				Second:        cSec,
				Nakshatra:     hc.Nakshatra,
				NakshatraPada: hc.NakshatraPada,
				NakshatraLord: hc.NakshatraLord,
			})
		}
		chart.Houses = vargaHouses

		res.Vargas = append(res.Vargas, chart)
	}

	return res
}
