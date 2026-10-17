package shadbala

import (
	"math"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func CalculateKalaBala(planet string, ctx *domain.CalculationContext, sunLon, moonLon, declination float64, mercuryMalefic bool) domain.KalaBala {
	ephemeris.Lock()
	geopos := [3]float64{ctx.Input.Longitude, ctx.Input.Latitude, 0}
	searchJD := ctx.JulianDayUT - 0.5
	resRise := swisseph.RiseTrans(searchJD, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcRise), geopos, 0, 0)
	resSet := swisseph.RiseTrans(searchJD, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcSet), geopos, 0, 0)
	ephemeris.Unlock()

	sunriseJD := resRise.Time
	sunsetJD := resSet.Time

	// 1. Nathonnatha Bala (from local apparent time; see nathonnatha.go)
	nathonnatha := 0.0
	ephemeris.Lock()
	apparentHours, errLAT := localApparentHours(ctx.JulianDayUT, ctx.Input.Longitude)
	ephemeris.Unlock()
	if errLAT == nil {
		nathonnatha = nathonnathaBala(planet, apparentHours)
	}

	// 2. Paksha Bala (see paksha.go)
	paksha := pakshaBala(planet, sunLon, moonLon, mercuryMalefic)

	// 3. Tribhaga Bala
	tribhaga := 0.0
	if planet == "Jupiter" {
		tribhaga = 60.0
	} else {
		isDay := false
		var start, end float64
		if ctx.JulianDayUT >= sunriseJD && ctx.JulianDayUT < sunsetJD {
			isDay = true
			start = sunriseJD
			end = sunsetJD
		} else {
			isDay = false
			if ctx.JulianDayUT < sunriseJD {
				ephemeris.Lock()
				resPrevSet := swisseph.RiseTrans(ctx.JulianDayUT-1.0, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcSet), geopos, 0, 0)
				ephemeris.Unlock()
				start = resPrevSet.Time
				end = sunriseJD
			} else {
				start = sunsetJD
				ephemeris.Lock()
				resNextRise := swisseph.RiseTrans(ctx.JulianDayUT, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcRise), geopos, 0, 0)
				ephemeris.Unlock()
				end = resNextRise.Time
			}
		}

		duration := end - start
		third := duration / 3.0
		part := 1
		if ctx.JulianDayUT >= start+2*third {
			part = 3
		} else if ctx.JulianDayUT >= start+third {
			part = 2
		}

		if isDay {
			if part == 1 && planet == "Mercury" {
				tribhaga = 60.0
			}
			if part == 2 && planet == "Sun" {
				tribhaga = 60.0
			}
			if part == 3 && planet == "Saturn" {
				tribhaga = 60.0
			}
		} else {
			if part == 1 && planet == "Moon" {
				tribhaga = 60.0
			}
			if part == 2 && planet == "Venus" {
				tribhaga = 60.0
			}
			if part == 3 && planet == "Mars" {
				tribhaga = 60.0
			}
		}
	}

	// 4. Ayana Bala (see ayana.go)
	ayana := ayanaBala(planet, declination)

	// 5. Dina Bala and Hora Bala: the Hindu day starts at sunrise; the weekday lord of that sunrise
	// gets 45, and the planet ruling the running 60-minute hora gets 60 (see hora.go).
	dinaBala, horaBala, varshaBala, masaBala := 0.0, 0.0, 0.0, 0.0
	if sunrise, ok := lastSunrise(ctx.JulianDayUT, geopos); ok {
		dayLord := weekdayLordOfSunrise(sunrise, ctx.Input.Timezone)
		if planet == dayLord {
			dinaBala = 45.0
		}
		// Varsha and Masa Bala: lords of the 360-day year and 30-day month of the Hindu day of birth
		// (see ahargana.go); the Hindu day is the one whose sunrise began it.
		ahargana := aharganaOfJDN(int64(math.Floor(sunrise + ctx.Input.Timezone/24.0 + 0.5)))
		if planet == abdadhipathi(ahargana) {
			varshaBala = 15.0
		}
		if planet == masadhipathi(ahargana) {
			masaBala = 30.0
		}
		if planet == horaLord(dayLord, (ctx.JulianDayUT-sunrise)*24.0) {
			horaBala = 60.0
		}
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
	return kb
}
