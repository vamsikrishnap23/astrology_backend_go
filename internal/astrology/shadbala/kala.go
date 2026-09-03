package shadbala

import (
	"math"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

// Weekday Lords
var WeekdayLords = []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"}

func CalculateKalaBala(planet string, ctx *domain.CalculationContext, sunLon, moonLon, declination float64) domain.KalaBala {
	ephemeris.Mu.Lock()
	swisseph.SetEphePath(ephemeris.EphePath)
	geopos := [3]float64{ctx.Input.Longitude, ctx.Input.Latitude, 0}
	searchJD := ctx.JulianDayUT - 0.5
	resRise := swisseph.RiseTrans(searchJD, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcRise), geopos, 0, 0)
	resSet := swisseph.RiseTrans(searchJD, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcSet), geopos, 0, 0)
	ephemeris.Mu.Unlock()

	sunriseJD := resRise.Time
	sunsetJD := resSet.Time
	noonJD := (sunriseJD + sunsetJD) / 2.0
	midnightJD := noonJD - 0.5
	if ctx.JulianDayUT < midnightJD {
		midnightJD = noonJD + 0.5
	}

	// 1. Nathonnatha Bala
	nathonnatha := 0.0
	diffNoon := math.Abs(ctx.JulianDayUT-noonJD) * 24.0 // in hours
	if diffNoon > 12.0 {
		diffNoon = 24.0 - diffNoon
	}
	diffMidnight := math.Abs(ctx.JulianDayUT-midnightJD) * 24.0
	if diffMidnight > 12.0 {
		diffMidnight = 24.0 - diffMidnight
	}

	if planet == "Moon" || planet == "Mars" || planet == "Saturn" {
		nathonnatha = (12.0 - diffMidnight) * 5.0
	} else if planet == "Sun" || planet == "Jupiter" || planet == "Venus" {
		nathonnatha = (12.0 - diffNoon) * 5.0
	} else if planet == "Mercury" {
		nathonnatha = 60.0
	}

	// 2. Paksha Bala
	paksha := 0.0
	angle := math.Mod(moonLon-sunLon+360.0, 360.0)
	if angle > 180.0 {
		angle = 360.0 - angle
	}
	val := angle / 3.0 // Max 60

	isBenefic := false
	if planet == "Jupiter" || planet == "Venus" || planet == "Moon" {
		isBenefic = true
	} else if planet == "Mercury" {
		isBenefic = true // Mercury usually conditionally benefic, we treat as benefic for Paksha
	}

	if isBenefic {
		paksha = val
	} else {
		paksha = 60.0 - val
	}
	if planet == "Moon" {
		paksha *= 2.0
	}

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
				ephemeris.Mu.Lock()
				resPrevSet := swisseph.RiseTrans(ctx.JulianDayUT-1.0, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcSet), geopos, 0, 0)
				ephemeris.Mu.Unlock()
				start = resPrevSet.Time
				end = sunriseJD
			} else {
				start = sunsetJD
				ephemeris.Mu.Lock()
				resNextRise := swisseph.RiseTrans(ctx.JulianDayUT, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcRise), geopos, 0, 0)
				ephemeris.Mu.Unlock()
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

	// 4. Ayana Bala
	ayana := 0.0
	// declination is passed in. Max is 24 degrees classically.
	// Formula: (24 + declination) / 48 * 60 = (24 + decl) * 1.25 for North.
	// For South: (24 - decl) * 1.25
	valNorth := (24.0 + declination) * (60.0 / 48.0)
	valSouth := (24.0 - declination) * (60.0 / 48.0)
	if valNorth > 60 {
		valNorth = 60
	}
	if valNorth < 0 {
		valNorth = 0
	}
	if valSouth > 60 {
		valSouth = 60
	}
	if valSouth < 0 {
		valSouth = 0
	}

	if planet == "Sun" || planet == "Mars" || planet == "Jupiter" || planet == "Venus" {
		ayana = valNorth
	} else if planet == "Moon" || planet == "Saturn" {
		ayana = valSouth
	} else if planet == "Mercury" {
		ayana = valNorth // Mercury often considered North, or 30 everywhere. Let's use North.
	}
	if planet == "Sun" {
		ayana *= 2.0
	}

	// 5. Dina and Hora Bala
	// Find local weekday at the relevant sunrise.
	relevantSunrise := sunriseJD
	if ctx.JulianDayUT < sunriseJD {
		relevantSunrise = sunriseJD - 1.0 // Use previous day's sunrise for Vedic day
	}

	localSunriseOffset := relevantSunrise + (ctx.Input.Timezone / 24.0)
	weekdayIdx := int(math.Floor(localSunriseOffset+1.5)) % 7 // 0=Sunday
	if weekdayIdx < 0 {
		weekdayIdx += 7
	}
	dayLord := WeekdayLords[weekdayIdx]

	dinaBala := 0.0
	if planet == dayLord {
		dinaBala = 45.0
	}

	// Hora Bala: Hour lord. 1st hour of day is ruled by day lord. Succeeding hours follow planetary order: Sun, Ven, Merc, Moon, Sat, Jup, Mars.
	// Order index: Sun(0), Ven(5), Merc(3), Moon(1), Sat(6), Jup(4), Mars(2). Step is +5 (or -2).
	horaOrder := []string{"Sun", "Venus", "Mercury", "Moon", "Saturn", "Jupiter", "Mars"}
	// Find index of dayLord in horaOrder
	startIdx := 0
	for i, v := range horaOrder {
		if v == dayLord {
			startIdx = i
			break
		}
	}

	// How many astrological hours since sunrise?
	// Day is divided into 12 horas, Night into 12 horas.
	var horaBala float64 = 0.0
	isDay := ctx.JulianDayUT >= sunriseJD && ctx.JulianDayUT < sunsetJD
	var hoursPassed int
	if isDay {
		dayDuration := sunsetJD - sunriseJD
		horaDuration := dayDuration / 12.0
		hoursPassed = int((ctx.JulianDayUT - sunriseJD) / horaDuration)
	} else {
		var activeSunset, nextSunrise float64
		if ctx.JulianDayUT < sunriseJD {
			// Born in the early morning before today's sunrise -> belongs to previous night
			ephemeris.Mu.Lock()
			resPrevSet := swisseph.RiseTrans(ctx.JulianDayUT-1.0, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcSet), geopos, 0, 0)
			ephemeris.Mu.Unlock()
			activeSunset = resPrevSet.Time
			nextSunrise = sunriseJD
		} else {
			// Born after today's sunset -> belongs to tonight
			activeSunset = sunsetJD
			ephemeris.Mu.Lock()
			resNextRise := swisseph.RiseTrans(ctx.JulianDayUT, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcRise), geopos, 0, 0)
			ephemeris.Mu.Unlock()
			nextSunrise = resNextRise.Time
		}

		nightDuration := nextSunrise - activeSunset
		horaDuration := nightDuration / 12.0
		hoursPassed = 12 + int((ctx.JulianDayUT-activeSunset)/horaDuration)
	}
	currentHoraLord := horaOrder[(startIdx+hoursPassed)%7]
	if planet == currentHoraLord {
		horaBala = 60.0
	}

	// Return KalaBala
	kb := domain.KalaBala{
		NathonnathaBala: nathonnatha,
		PakshaBala:      paksha,
		TribhagaBala:    tribhaga,
		VarshaBala:      0.0, // Simplified for now
		MasaBala:        0.0, // Simplified for now
		DinaBala:        dinaBala,
		HoraBala:        horaBala,
		AyanaBala:       ayana,
		YuddhaBala:      0.0,
	}
	kb.Total = kb.NathonnathaBala + kb.PakshaBala + kb.TribhagaBala + kb.VarshaBala + kb.MasaBala + kb.DinaBala + kb.HoraBala + kb.AyanaBala + kb.YuddhaBala
	return kb
}
