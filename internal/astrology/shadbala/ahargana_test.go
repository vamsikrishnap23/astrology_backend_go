package shadbala

import (
	"fmt"
	"math"
	"testing"

	"github.com/tejzpr/go-swisseph"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func jdnOf(y, m, d int32) int64 {
	return int64(math.Floor(swisseph.Julday(y, m, d, 12.0, swisseph.GregCal) + 0.5))
}

// Raman, Table I (page 110-111): the Ahargana on 31 December of each year, counted from 2 May 1827
// (day 1), and Example 23 (16 Oct 1918 = 33,405). The table's last row prints 63,431 for 2000, but the
// year is a leap year (starred in the table) and every other leap year adds 366, so it should read
// 63,432; the rows used here are the ones that are consistent.
func TestAharganaMatchesRamanTableI(t *testing.T) {
	rows := []struct {
		year int32
		want int64
	}{{1827, 244}, {1828, 610}, {1832, 2071}, {1860, 12298}, {1900, 26907}, {1917, 33116}, {1918, 33481},
		{1950, 45169}, {1964, 50283}, {1988, 59049}, {1999, 63066}}
	for _, r := range rows {
		if got := aharganaOfJDN(jdnOf(r.year, 12, 31)); got != r.want {
			t.Errorf("Ahargana on 31 Dec %d = %d, want %d", r.year, got, r.want)
		}
	}
	if got := aharganaOfJDN(jdnOf(1918, 10, 16)); got != 33405 {
		t.Errorf("Ahargana on 16 Oct 1918 = %d, want 33405", got)
	}
}

// Raman examples 25-26, 24: year lord Saturn, month lord Budha, weekday lord Budha for 33,405.
func TestLordsRamanExample(t *testing.T) {
	if got := abdadhipathi(33405); got != "Saturn" {
		t.Errorf("Abdadhipathi = %s, want Saturn", got)
	}
	if got := masadhipathi(33405); got != "Mercury" {
		t.Errorf("Masadhipathi = %s, want Mercury", got)
	}
	if got := varadhipathi(33405); got != "Mercury" {
		t.Errorf("Varadhipathi = %s, want Mercury", got)
	}
}

// BPHS 27.13 / Raman examples 20-22 count from the creation with remainders counted from Sunday:
// year lord = (Ahargana/360 x 3 + 1) mod 7, month lord = (Ahargana/30 x 2 + 1) mod 7, weekday lord =
// Ahargana mod 7; remainder 1 = Sunday ... 0 = Saturday. Raman's condensed count (from Wednesday) must
// give the same lords for every day.
func TestWednesdayCountEqualsSundayCountOfFullAhargana(t *testing.T) {
	sundayFirst := []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"}
	lord := func(rem int64) string { return sundayFirst[((rem-1)%7+7)%7] }
	for a := int64(1); a < 60000; a++ {
		full := a + ahargana1827
		if got, want := abdadhipathi(a), lord((full/360)*3+1); got != want {
			t.Fatalf("year lord differs at condensed %d: %s vs %s", a, got, want)
		}
		if got, want := masadhipathi(a), lord((full/30)*2+1); got != want {
			t.Fatalf("month lord differs at condensed %d: %s vs %s", a, got, want)
		}
		if got, want := varadhipathi(a), lord(full); got != want {
			t.Fatalf("weekday lord differs at condensed %d: %s vs %s", a, got, want)
		}
	}
	// Raman's example 20-22: full Ahargana 714,404,130,045 gives Saturn, Budha and Budha.
	full := int64(714404130045)
	if aharganaOfJDN(jdnOf(1918, 10, 16))+ahargana1827 != full {
		t.Errorf("full Ahargana of 16 Oct 1918 = %d, want %d", aharganaOfJDN(jdnOf(1918, 10, 16))+ahargana1827, full)
	}
	// BPHS 27.13: 1 January 1860 is 714,404,108,573 days from the creation.
	if got := aharganaOfJDN(jdnOf(1860, 1, 1)) + ahargana1827; got != 714404108573 {
		t.Errorf("full Ahargana of 1 Jan 1860 = %d, want 714404108573 (BPHS)", got)
	}
}

// The weekday lord from the Ahargana must equal the real weekday (Raman's own check, art. 62).
func TestVaradhipathiIsTheRealWeekday(t *testing.T) {
	for y := int32(1900); y <= 2030; y += 7 {
		jdn := jdnOf(y, 3, 9)
		idx := int((jdn + 1) % 7) // 0 = Sunday
		if got, want := varadhipathi(aharganaOfJDN(jdn)), WeekdayLords[idx]; got != want {
			t.Errorf("9 Mar %d: Ahargana weekday %s, calendar weekday %s", y, got, want)
		}
	}
}

// A 360-day year and a 30-day month: the lord is constant inside one and changes at the boundary.
func TestYearAndMonthLordBoundaries(t *testing.T) {
	for a := int64(360); a < 360*30; a++ {
		same := a/360 == (a-1)/360
		if same && abdadhipathi(a) != abdadhipathi(a-1) {
			t.Fatalf("year lord changed inside a year at %d", a)
		}
		if !same && abdadhipathi(a) == abdadhipathi(a-1) {
			t.Fatalf("year lord did not change at the year boundary %d", a)
		}
	}
	for a := int64(30); a < 30*400; a++ {
		same := a/30 == (a-1)/30
		if same && masadhipathi(a) != masadhipathi(a-1) {
			t.Fatalf("month lord changed inside a month at %d", a)
		}
		if !same && masadhipathi(a) == masadhipathi(a-1) {
			t.Fatalf("month lord did not change at the month boundary %d", a)
		}
	}
}

func balaFor(t *testing.T, date, clock string, planet string) domain.KalaBala {
	t.Helper()
	return kalaFor(t, date, clock, 13.0878, 80.2785, planet)
}

// Raman's standard horoscope end to end: Saturn gets Varsha Bala 15 and Mercury Masa Bala 30
// (examples 28-29, 35); nobody else gets either.
func TestVarshaMasaRamanChart(t *testing.T) {
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	lon := 77.0 + 35.0/60
	jd := swisseph.Julday(1918, 10, 16, 14.0+6.0/60+16.0/3600-lon/15, swisseph.GregCal)
	ctx := domain.CalculationContext{Input: domain.BirthInput{Latitude: 13, Longitude: lon, Timezone: 5.5}, JulianDayUT: jd}
	for _, p := range TargetPlanets {
		kb := CalculateKalaBala(p, &ctx, 0, 0, 0, false)
		wantV, wantM := 0.0, 0.0
		if p == "Saturn" {
			wantV = 15
		}
		if p == "Mercury" {
			wantM = 30
		}
		near(t, p+" varsha", kb.VarshaBala, wantV, 0)
		near(t, p+" masa", kb.MasaBala, wantM, 0)
	}
}

// The Hindu day begins at sunrise, so a birth before sunrise belongs to the previous day's Ahargana:
// on the first civil day of a new 30-day month, a birth at 3 am still has the old month's lord.
func TestMasaUsesHinduDayBeforeSunrise(t *testing.T) {
	var dateStr string
	var a int64
	for jdn := jdnOf(2001, 1, 1); jdn < jdnOf(2003, 1, 1) && dateStr == ""; jdn++ {
		if aharganaOfJDN(jdn)%30 == 0 && aharganaOfJDN(jdn)%360 != 0 {
			r := swisseph.Revjul(float64(jdn), swisseph.GregCal)
			dateStr = fmt.Sprintf("%04d-%02d-%02d", r.Year, r.Month, r.Day)
			a = aharganaOfJDN(jdn)
		}
	}
	if dateStr == "" {
		t.Fatal("no month boundary found")
	}
	before := masadhipathi(a - 1)
	after := masadhipathi(a)
	if before == after {
		t.Fatal("test setup: lords should differ across the boundary")
	}
	gotBefore := balaFor(t, dateStr, "03:00:00", before).MasaBala
	gotAfter := balaFor(t, dateStr, "12:00:00", after).MasaBala
	near(t, "3 am: old month's lord", gotBefore, 30, 0)
	near(t, "noon: new month's lord", gotAfter, 30, 0)
}
