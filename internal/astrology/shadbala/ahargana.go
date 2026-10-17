package shadbala

import (
	"math"

	"github.com/tejzpr/go-swisseph"
)

// Varsha (year) and Masa (month) Bala come from a count of days, the Ahargana (Raman art. 58-67,
// BPHS 27.13). Hindu astrology uses a "year" of 360 days and a "month" of 30 days. Raman counts days
// from 2 May 1827 (a Wednesday), which is day 1 (his condensed Ahargana); the full Ahargana from the
// creation is that plus ahargana1827 (BPHS gives 714,404,108,573 for 1 Jan 1860, which agrees).
const ahargana1827 int64 = 714404096640

// weekdaysFromWednesday lists the weekday lords in the order Raman counts the remainder: remainder 1
// is Wednesday, ... remainder 0 (= 7) is Tuesday (art. 60-63).
var weekdaysFromWednesday = []string{"Mercury", "Jupiter", "Venus", "Saturn", "Sun", "Moon", "Mars"}

// aharganaEpochJDN is the Julian day number of 2 May 1827 (civil date).
func aharganaEpochJDN() int64 {
	return int64(math.Floor(swisseph.Julday(1827, 5, 2, 12.0, swisseph.GregCal) + 0.5))
}

// aharganaOfJDN returns Raman's condensed Ahargana for a civil date given by its Julian day number.
func aharganaOfJDN(jdn int64) int64 {
	return jdn - aharganaEpochJDN() + 1
}

// lordFromWednesdayRemainder maps a remainder (mod 7) counted from Wednesday to its weekday lord.
func lordFromWednesdayRemainder(rem int64) string {
	return weekdaysFromWednesday[((rem-1)%7+7)%7]
}

// abdadhipathi is the lord of the Hindu year of birth: the lord of the weekday on which that 360-day
// year began. quotient = Ahargana / 360, multiplied by 3 plus 1, remainder after dividing by 7 is
// counted from Wednesday (Raman art. 60; the year starts 3 days later each year since 360 = 7 x 51 + 3).
func abdadhipathi(ahargana int64) string {
	return lordFromWednesdayRemainder((ahargana/360)*3 + 1)
}

// masadhipathi is the lord of the Hindu month of birth (30 days = 7 x 4 + 2): quotient = Ahargana /
// 30, times 2 plus 1, remainder counted from Wednesday (Raman art. 61).
func masadhipathi(ahargana int64) string {
	return lordFromWednesdayRemainder((ahargana/30)*2 + 1)
}

// varadhipathi is the lord of the weekday: Ahargana mod 7 counted from Wednesday (Raman art. 62-63).
func varadhipathi(ahargana int64) string {
	return lordFromWednesdayRemainder(ahargana)
}
