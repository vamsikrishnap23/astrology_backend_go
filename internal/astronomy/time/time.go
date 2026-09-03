package time

import (
	"fmt"
	"github.com/mshafiee/swephgo"
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

	jd := swephgo.Julday(year, month, day, hour, swephgo.SeGregCal)
	return jd
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
