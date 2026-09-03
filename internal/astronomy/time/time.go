package time

import (
	"fmt"
	"github.com/tejzpr/go-swisseph"
	"math"
	"time"
)

// ParseLocalToUTC parses the date and time strings and the decimal timezone offset,
// returning the normalized UTC time.
func ParseLocalToUTC(dateStr, timeStr string, tzOffset float64) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	localStr := fmt.Sprintf("%s %s", dateStr, timeStr)

	// Parse as if it's UTC just to get the parts
	t, err := time.Parse(layout, localStr)
	if err != nil {
		return time.Time{}, err
	}

	// tzOffset is hours. e.g. 5.5 hours = 5 hours 30 mins
	// local = UTC + offset
	// UTC = local - offset
	offsetDuration := time.Duration(tzOffset * float64(time.Hour))

	utcTime := t.Add(-offsetDuration)
	return utcTime, nil
}

// UTCToJulianDay converts UTC time to Julian Day (UT).
func UTCToJulianDay(utc time.Time) float64 {
	year := utc.Year()
	month := int(utc.Month())
	day := utc.Day()
	hour := float64(utc.Hour()) + float64(utc.Minute())/60.0 + float64(utc.Second())/3600.0

	jd := swisseph.Julday(int32(year), int32(month), int32(day), hour, swisseph.GregCal)
	return jd
}

// JulianDayToUTC converts Julian Day (UT) back to a time.Time UTC object.
func JulianDayToUTC(jd float64) time.Time {
	res := swisseph.Revjul(jd, swisseph.GregCal)

	year := int(res.Year)
	month := time.Month(res.Month)
	day := int(res.Day)

	hourFloat := res.Hour
	h := int(hourFloat)
	rem := (hourFloat - float64(h)) * 60.0
	m := int(rem)
	s := int((rem - float64(m)) * 60.0)

	return time.Date(year, month, day, h, m, s, 0, time.UTC)
}

// ConvertDecimalToDegree minute second
func DecimalToDMS(decimal float64) (sign string, deg int, min int, sec float64) {
	// Normalize between 0 and 360
	decimal = math.Mod(decimal, 360.0)
	if decimal < 0 {
		decimal += 360.0
	}

	signs := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}

	signIdx := int(decimal / 30.0)
	sign = signs[signIdx]

	degInSign := decimal - float64(signIdx*30)

	deg = int(degInSign)
	rem := (degInSign - float64(deg)) * 60.0
	min = int(rem)
	sec = (rem - float64(min)) * 60.0

	return
}
