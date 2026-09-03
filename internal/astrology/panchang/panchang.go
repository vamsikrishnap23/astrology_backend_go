package panchang

import (
	"math"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func CalculatePanchang(ctx *domain.CalculationContext) (domain.PanchangResult, error) {
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

	// Calculate Solar and Lunar Times
	sunriseJD, sunsetJD := calculateRiseSet(ctx.JulianDayUT, ctx.Input.Latitude, ctx.Input.Longitude, swisseph.Sun)
	moonriseJD, moonsetJD := calculateRiseSet(ctx.JulianDayUT, ctx.Input.Latitude, ctx.Input.Longitude, swisseph.Moon)
	noonJD := (sunriseJD + sunsetJD) / 2.0

	tithi := calculateTithi(ctx.JulianDayUT, sunTrop, moonTrop)
	vara := calculateVara(ctx.Input.DateOfBirth, ctx.Input.Timezone)
	nakshatra := calculateNakshatra(ctx.JulianDayUT, moonSid)
	yoga := calculateYoga(ctx.JulianDayUT, sunSid, moonSid)
	karana := calculateKarana(ctx.JulianDayUT, sunTrop, moonTrop)

	rahu, yama, gulika := calculateDailyPeriods(sunriseJD, sunsetJD, vara.Number)

	tzOff := time.Duration(ctx.Input.Timezone * float64(time.Hour))
	loc := time.FixedZone("Local", int(tzOff.Seconds()))

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
		Gulika:    domain.DailyPeriod{Start: formatTime(gulika[0]), End: formatTime(gulika[1])},
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

	// Start search from roughly 12 hours before to ensure we catch today's sunrise/sunset
	// even if the birth time is late in the day.
	searchJD := jd - 0.5

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
