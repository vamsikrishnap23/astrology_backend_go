package shadbala

import (
	"math"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// CalculateKalaBala returns the Kala Bala of planet for the chart in ctx, except Yuddha Bala (see
// yuddha.go). It returns ErrNoSunriseSunset where there is no sunrise or sunset to count from.
func CalculateKalaBala(planet string, ctx *domain.CalculationContext, sunLon, moonLon, declination float64, mercuryMalefic bool) (domain.KalaBala, error) {
	geopos := [3]float64{ctx.Input.Longitude, ctx.Input.Latitude, 0}
	day, err := hinduDayAt(ctx.JulianDayUT, geopos)
	if err != nil {
		return domain.KalaBala{}, err
	}

	// 1. Nathonnatha Bala (from local apparent time; see nathonnatha.go)
	ephemeris.Lock()
	apparentHours, err := localApparentHours(ctx.JulianDayUT, ctx.Input.Longitude)
	ephemeris.Unlock()
	if err != nil {
		return domain.KalaBala{}, err
	}
	nathonnatha := nathonnathaBala(planet, apparentHours)

	// 2. Paksha Bala (see paksha.go)
	paksha := pakshaBala(planet, sunLon, moonLon, mercuryMalefic)

	// 3. Tribhaga Bala (see tribhaga.go)
	tribhaga := tribhagaBala(planet, day, ctx.JulianDayUT)

	// 4. Ayana Bala (see ayana.go)
	ayana := ayanaBala(planet, declination)

	// 5. Dina, Hora, Varsha and Masa Bala: the Hindu day starts at sunrise; the weekday lord of that
	// sunrise gets 45, the planet ruling the running 60-minute hora gets 60 (see hora.go), and the lords
	// of the 360-day year and 30-day month get 15 and 30 (see ahargana.go).
	dinaBala, horaBala, varshaBala, masaBala := 0.0, 0.0, 0.0, 0.0
	dayLord := weekdayLordOfSunrise(day.Sunrise, ctx.Input.Timezone)
	if planet == dayLord {
		dinaBala = 45.0
	}
	if planet == horaLord(dayLord, (ctx.JulianDayUT-day.Sunrise)*24.0) {
		horaBala = 60.0
	}
	ahargana := aharganaOfJDN(int64(math.Floor(day.Sunrise + ctx.Input.Timezone/24.0 + 0.5)))
	if planet == abdadhipathi(ahargana) {
		varshaBala = 15.0
	}
	if planet == masadhipathi(ahargana) {
		masaBala = 30.0
	}

	// Return KalaBala
	kb := domain.KalaBala{
		NathonnathaBala: nathonnatha,
		PakshaBala:      paksha,
		TribhagaBala:    tribhaga,
		VarshaBala:      varshaBala,
		MasaBala:        masaBala,
		DinaBala:        dinaBala,
		HoraBala:        horaBala,
		AyanaBala:       ayana,
		YuddhaBala:      0.0,
	}
	kb.Total = kb.NathonnathaBala + kb.PakshaBala + kb.TribhagaBala + kb.VarshaBala + kb.MasaBala + kb.DinaBala + kb.HoraBala + kb.AyanaBala + kb.YuddhaBala
	return kb, nil
}
