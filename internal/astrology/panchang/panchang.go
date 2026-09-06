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

	tflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed)
	sunRes := swisseph.CalcUT(ctx.JulianDayUT, swisseph.Sun, tflag)
	moonRes := swisseph.CalcUT(ctx.JulianDayUT, swisseph.Moon, tflag)
	sunTrop := sunRes.Data[0]
	moonTrop := moonRes.Data[0]

	sflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed | swisseph.FlagSidereal)
	sunSidRes := swisseph.CalcUT(ctx.JulianDayUT, swisseph.Sun, sflag)
	moonSidRes := swisseph.CalcUT(ctx.JulianDayUT, swisseph.Moon, sflag)
	sunSid := sunSidRes.Data[0]
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

	tithi := calculateTithi(ctx.JulianDayUT, sunTrop, moonTrop)
	vara := calculateVara(ctx.Input.DateOfBirth, ctx.Input.Timezone)
	nakshatra := calculateNakshatra(ctx.JulianDayUT, moonSid)
	yoga := calculateYoga(ctx.JulianDayUT, sunSid, moonSid)
	karana := calculateKarana(ctx.JulianDayUT, sunTrop, moonTrop)

	rahu, yama, durmuhurtams := calculateDailyPeriods(sunriseJD, sunsetJD, vara.Number)

	formatTime := func(jd float64) string {
		utc := jdToUTC(jd)
		return utc.In(loc).Format("2006-01-02T15:04:05-07:00")
	}

	res := domain.PanchangResult{
		Date:      ctx.Input.DateOfBirth,
		LocalTime: ctx.Input.TimeOfBirth,
		Timezone:  ctx.Input.Timezone,
		Sunrise:   formatTime(sunriseJD),
		Sunset:    formatTime(sunsetJD),
		SolarNoon: formatTime(noonJD),
		Moonrise:  formatTime(moonriseJD),
		Moonset:   formatTime(moonsetJD),
		Vara:      vara,
		Tithi:     formatTithi(tithi, formatTime),
		Nakshatra: formatNakshatra(nakshatra, formatTime),
		Yoga:      formatYoga(yoga, formatTime),
		Karana:    formatKarana(karana, formatTime),
		RahuKalam: domain.DailyPeriod{Start: formatTime(rahu[0]), End: formatTime(rahu[1])},
		Yamaganda: domain.DailyPeriod{Start: formatTime(yama[0]), End: formatTime(yama[1])},
	}

	for _, d := range durmuhurtams {
		res.Durmuhurtam = append(res.Durmuhurtam, domain.DailyPeriod{Start: formatTime(d[0]), End: formatTime(d[1])})
	}

	// Calculate Varjyam (Tyajyam)
	varjyamGhatis := []float64{50, 4, 30, 40, 14, 21, 30, 21, 33, 30, 20, 1, 21, 20, 14, 14, 17, 21, 20, 20, 20, 10, 10, 18, 16, 16, 50}
	nNum := nakshatra.Number - 1 // 0-indexed Nakshatra
	if nNum >= 0 && nNum < 27 && nakshatra.StartJD > 0 && nakshatra.EndJD > 0 {
		vStartGhati := varjyamGhatis[nNum]
		vStartFrac := vStartGhati / 60.0
		vEndFrac := (vStartGhati + 4.0) / 60.0
		nakDur := nakshatra.EndJD - nakshatra.StartJD
		vStartJD := nakshatra.StartJD + vStartFrac*nakDur
		vEndJD := nakshatra.StartJD + vEndFrac*nakDur
		res.Varjyam = append(res.Varjyam, domain.DailyPeriod{Start: formatTime(vStartJD), End: formatTime(vEndJD)})
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

	// We start the search from exactly the given jd.
	// If jd is 00:00 Local Time, the next rise/set will accurately be today's events.
	searchJD := jd

	// Removed BitDiscCenter and BitNoRefraction to calculate apparent visual upper limb sunrise/sunset
	resRise := swisseph.RiseTrans(searchJD, body, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcRise), geopos, 0, 0)
	resSet := swisseph.RiseTrans(searchJD, body, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcSet), geopos, 0, 0)

	return resRise.Time, resSet.Time
}

type elementData struct {
	Number   int
	Progress float64
	StartJD  float64
	EndJD    float64
}
