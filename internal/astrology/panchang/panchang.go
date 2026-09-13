package panchang

import (
	"math"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func CalculatePanchang(ctx *domain.CalculationContext) (domain.PanchangResult, error) {
	ephemeris.Mu.Lock()
	defer ephemeris.Mu.Unlock()
	swisseph.SetEphePath(ephemeris.EphePath)

	swisseph.SetSidMode(int32(ctx.Config.AyanamsaMode), 0, 0)

	sflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed | swisseph.FlagSidereal)
	moonSidRes := swisseph.CalcUT(ctx.JulianDayUT, swisseph.Moon, sflag)
	moonSid := moonSidRes.Data[0]

	// Find the Julian day for 00:00:00 of the input date in the input timezone
	tzOff := time.Duration(ctx.Input.Timezone * float64(time.Hour))
	loc := time.FixedZone("Local", int(tzOff.Seconds()))

	t, _ := time.Parse("2006-01-02", ctx.Input.DateOfBirth)
	localStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	utcStart := localStart.UTC()

	// Convert UTC time to JD manually (or use swisseph if time is not imported)
	// swisseph.Julday takes year, month, day, hour (decimal)
	startOfDayJD := swisseph.Julday(int32(utcStart.Year()), int32(utcStart.Month()), int32(utcStart.Day()), float64(utcStart.Hour())+float64(utcStart.Minute())/60.0+float64(utcStart.Second())/3600.0, swisseph.GregCal)

	// Calculate Solar and Lunar Times from Start of Day
	sunriseJD, sunsetJD := calculateRiseSet(startOfDayJD, ctx.Input.Latitude, ctx.Input.Longitude, swisseph.Sun)
	moonriseJD, moonsetJD := calculateRiseSet(startOfDayJD, ctx.Input.Latitude, ctx.Input.Longitude, swisseph.Moon)
	noonJD := (sunriseJD + sunsetJD) / 2.0

	tithi := calculateTithi(ctx.JulianDayUT)
	vara := calculateVara(ctx.Input.DateOfBirth, ctx.Input.Timezone)
	nakshatra := calculateNakshatra(ctx.JulianDayUT)
	yoga := calculateYoga(ctx.JulianDayUT)
	karana := calculateKarana(ctx.JulianDayUT)

	rahu, yama, durmuhurtams := calculateDailyPeriods(sunriseJD, sunsetJD, vara.Number)

	formatTime := func(jd float64) string {
		utc := jdToUTC(jd)
		return utc.In(loc).Format("2006-01-02T15:04:05-07:00")
	}

	zodiacSigns := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
	moonSignIdx := int(math.Floor(moonSid/30.0)) % 12
	if moonSignIdx < 0 {
		moonSignIdx += 12
	}
	rasiName := zodiacSigns[moonSignIdx]

	teluguCal := CalculateTeluguCalendar(ctx.JulianDayUT)

	res := domain.PanchangResult{
		Date:           ctx.Input.DateOfBirth,
		LocalTime:      ctx.Input.TimeOfBirth,
		Timezone:       ctx.Input.Timezone,
		Sunrise:        formatTime(sunriseJD),
		Sunset:         formatTime(sunsetJD),
		SolarNoon:      formatTime(noonJD),
		Moonrise:       formatTime(moonriseJD),
		Moonset:        formatTime(moonsetJD),
		Rasi:           rasiName,
		Vara:           vara,
		Tithi:          formatTithi(tithi, formatTime),
		Nakshatra:      formatNakshatra(nakshatra, formatTime),
		Yoga:           formatYoga(yoga, formatTime),
		Karana:         formatKarana(karana, formatTime),
		TeluguCalendar: teluguCal,
	}

	for _, r := range rahu {
		res.RahuKalam = append(res.RahuKalam, domain.DailyPeriod{Start: formatTime(r[0]), End: formatTime(r[1])})
	}
	for _, y := range yama {
		res.Yamaganda = append(res.Yamaganda, domain.DailyPeriod{Start: formatTime(y[0]), End: formatTime(y[1])})
	}

	for _, d := range durmuhurtams {
		res.Durmuhurtam = append(res.Durmuhurtam, domain.DailyPeriod{Start: formatTime(d[0]), End: formatTime(d[1])})
	}

	// Standard Varjyam (Visha Ghatis) starting ghatis for 27 Nakshatras
	varjyamGhatis := []float64{50, 24, 30, 40, 14, 21, 30, 20, 32, 30, 20, 18, 21, 20, 14, 14, 10, 14, 56, 24, 20, 10, 10, 18, 16, 24, 30}

	// Standard Amrita Kalam starting ghatis for 27 Nakshatras (Varjyam + 24 ghatis rule, wrapped at 60)
	amruthaGhatis := []float64{14, 48, 54, 4, 38, 45, 54, 44, 56, 54, 44, 42, 45, 44, 38, 38, 34, 38, 20, 48, 44, 34, 34, 42, 40, 48, 54}

	dayNakshatras := getTimeline(startOfDayJD, startOfDayJD+1.0, calculateNakshatra)

	for _, nak := range dayNakshatras {
		nNum := nak.Number - 1 // 0-indexed Nakshatra
		if nNum >= 0 && nNum < 27 && nak.StartJD > 0 && nak.EndJD > 0 {
			nakDur := nak.EndJD - nak.StartJD

			// Calculate Varjyam (duration is 4 ghatis)
			vStartGhati := varjyamGhatis[nNum]
			vStartFrac := vStartGhati / 60.0
			vEndFrac := (vStartGhati + 4.0) / 60.0
			vStartJD := nak.StartJD + vStartFrac*nakDur
			vEndJD := nak.StartJD + vEndFrac*nakDur

			// Only append if it overlaps with our calendar day or just append all?
			// Usually we append all that correspond to the day's Nakshatras.
			res.Varjyam = append(res.Varjyam, domain.DailyPeriod{Start: formatTime(vStartJD), End: formatTime(vEndJD)})

			// Calculate Amrutha Ghadiyalu (duration is 4 ghatis)
			aStartGhati := amruthaGhatis[nNum]
			aStartFrac := aStartGhati / 60.0
			aEndFrac := (aStartGhati + 4.0) / 60.0
			aStartJD := nak.StartJD + aStartFrac*nakDur
			aEndJD := nak.StartJD + aEndFrac*nakDur
			res.AmruthaGhadiyalu = append(res.AmruthaGhadiyalu, domain.DailyPeriod{Start: formatTime(aStartJD), End: formatTime(aEndJD)})
		}
	}

	return res, nil
}

func CalculateDailyPanchang(ctx *domain.CalculationContext) (domain.DailyPanchangResult, error) {
	ephemeris.Mu.Lock()
	defer ephemeris.Mu.Unlock()
	swisseph.SetEphePath(ephemeris.EphePath)
	swisseph.SetSidMode(int32(ctx.Config.AyanamsaMode), 0, 0)

	tzOff := time.Duration(ctx.Input.Timezone * float64(time.Hour))
	loc := time.FixedZone("Local", int(tzOff.Seconds()))

	t, _ := time.Parse("2006-01-02", ctx.Input.DateOfBirth)
	localStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	utcStart := localStart.UTC()

	startOfDayJD := swisseph.Julday(int32(utcStart.Year()), int32(utcStart.Month()), int32(utcStart.Day()), float64(utcStart.Hour())+float64(utcStart.Minute())/60.0+float64(utcStart.Second())/3600.0, swisseph.GregCal)
	endOfDayJD := startOfDayJD + 1.0

	sunriseJD, sunsetJD := calculateRiseSet(startOfDayJD, ctx.Input.Latitude, ctx.Input.Longitude, swisseph.Sun)
	moonriseJD, moonsetJD := calculateRiseSet(startOfDayJD, ctx.Input.Latitude, ctx.Input.Longitude, swisseph.Moon)
	noonJD := (sunriseJD + sunsetJD) / 2.0

	vara := calculateVara(ctx.Input.DateOfBirth, ctx.Input.Timezone)
	rahu, yama, durmuhurtams := calculateDailyPeriods(sunriseJD, sunsetJD, vara.Number)

	formatTime := func(jd float64) string {
		utc := jdToUTC(jd)
		return utc.In(loc).Format("2006-01-02T15:04:05-07:00")
	}

	teluguCal := CalculateTeluguCalendar(ctx.JulianDayUT)

	res := domain.DailyPanchangResult{
		Date:           ctx.Input.DateOfBirth,
		Timezone:       ctx.Input.Timezone,
		Sunrise:        formatTime(sunriseJD),
		Sunset:         formatTime(sunsetJD),
		SolarNoon:      formatTime(noonJD),
		Moonrise:       formatTime(moonriseJD),
		Moonset:        formatTime(moonsetJD),
		Vara:           vara,
		TeluguCalendar: teluguCal,
	}

	dayTithis := getTimeline(startOfDayJD, endOfDayJD, calculateTithi)
	for _, dt := range dayTithis {
		res.Tithis = append(res.Tithis, formatTithi(dt, formatTime))
	}

	dayNakshatras := getTimeline(startOfDayJD, endOfDayJD, calculateNakshatra)
	for _, dn := range dayNakshatras {
		res.Nakshatras = append(res.Nakshatras, formatNakshatra(dn, formatTime))
	}

	dayYogas := getTimeline(startOfDayJD, endOfDayJD, calculateYoga)
	for _, dy := range dayYogas {
		res.Yogas = append(res.Yogas, formatYoga(dy, formatTime))
	}

	dayKaranas := getTimeline(startOfDayJD, endOfDayJD, calculateKarana)
	for _, dk := range dayKaranas {
		res.Karanas = append(res.Karanas, formatKarana(dk, formatTime))
	}

	for _, r := range rahu {
		res.RahuKalam = append(res.RahuKalam, domain.DailyPeriod{Start: formatTime(r[0]), End: formatTime(r[1])})
	}
	for _, y := range yama {
		res.Yamaganda = append(res.Yamaganda, domain.DailyPeriod{Start: formatTime(y[0]), End: formatTime(y[1])})
	}

	for _, d := range durmuhurtams {
		res.Durmuhurtam = append(res.Durmuhurtam, domain.DailyPeriod{Start: formatTime(d[0]), End: formatTime(d[1])})
	}

	varjyamGhatis := []float64{50, 24, 30, 40, 14, 21, 30, 20, 32, 30, 20, 18, 21, 20, 14, 14, 10, 14, 56, 24, 20, 10, 10, 18, 16, 24, 30}
	amruthaGhatis := []float64{14, 48, 54, 4, 38, 45, 54, 44, 56, 54, 44, 42, 45, 44, 38, 38, 34, 38, 20, 48, 44, 34, 34, 42, 40, 48, 54}

	for _, nak := range dayNakshatras {
		nNum := nak.Number - 1
		if nNum >= 0 && nNum < 27 && nak.StartJD > 0 && nak.EndJD > 0 {
			nakDur := nak.EndJD - nak.StartJD

			vStartGhati := varjyamGhatis[nNum]
			vStartFrac := vStartGhati / 60.0
			vEndFrac := (vStartGhati + 4.0) / 60.0
			vStartJD := nak.StartJD + vStartFrac*nakDur
			vEndJD := nak.StartJD + vEndFrac*nakDur
			res.Varjyam = append(res.Varjyam, domain.DailyPeriod{Start: formatTime(vStartJD), End: formatTime(vEndJD)})

			aStartGhati := amruthaGhatis[nNum]
			aStartFrac := aStartGhati / 60.0
			aEndFrac := (aStartGhati + 4.0) / 60.0
			aStartJD := nak.StartJD + aStartFrac*nakDur
			aEndJD := nak.StartJD + aEndFrac*nakDur
			res.AmruthaGhadiyalu = append(res.AmruthaGhadiyalu, domain.DailyPeriod{Start: formatTime(aStartJD), End: formatTime(aEndJD)})
		}
	}

	return res, nil
}

func jdToUTC(jd float64) time.Time {
	res := swisseph.Revjul(jd, swisseph.GregCal)
	y, m, d, h := res.Year, res.Month, res.Day, res.Hour
	hr := int(h)
	min := int((h - float64(hr)) * 60)
	sec := int(math.Round((h - float64(hr) - float64(min)/60) * 3600))
	if sec >= 60 {
		sec -= 60
		min++
	}
	if min >= 60 {
		min -= 60
		hr++
	}
	return time.Date(int(y), time.Month(m), int(d), hr, min, sec, 0, time.UTC)
}

func calculateRiseSet(jd, lat, lon float64, body int32) (float64, float64) {
	geopos := [3]float64{lon, lat, 0}

	// Use Hindu Sunrise definitions: Center of Sun's disk, without atmospheric refraction.
	// This is the standard rule for South Indian Panchangams (Nithra, Butte, etc.)
	searchJD := jd
	rsflag := int32(swisseph.FlagSwieph | swisseph.BitDiscCenter | swisseph.BitNoRefraction)

	resRise := swisseph.RiseTrans(searchJD, body, "", rsflag, int32(swisseph.CalcRise), geopos, 0, 0)
	resSet := swisseph.RiseTrans(searchJD, body, "", rsflag, int32(swisseph.CalcSet), geopos, 0, 0)

	return resRise.Time, resSet.Time
}

type elementData struct {
	Number   int
	Progress float64
	StartJD  float64
	EndJD    float64
}
