package shadbala

import (
	"testing"

	"github.com/tejzpr/go-swisseph"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// Raman examples 49-51: Chesta Kendra and Cheshta Bala from the Seeghrochcha, mean and true longitude.
func TestCheshtaBalaRamanExamples(t *testing.T) {
	cases := []struct {
		planet                 string
		seeghra, mean, trueLon float64
		want                   float64
	}{
		{"Mars", 181.23, 266.34, 229.50, 22.23},   // kendra 293.31, reduced 66.69
		{"Mercury", 174.49, 181.23, 181.52, 2.30}, // 353.11 -> 6.89
		{"Jupiter", 181.23, 66.91, 84.01, 35.26},  // 105.77
		{"Venus", 158.35, 181.23, 171.16, 5.95},   // 342.15 -> 17.85
		{"Saturn", 181.23, 111.23, 124.39, 21.14}, // 63.42
	}
	for _, c := range cases {
		near(t, c.planet, cheshtaBala(c.seeghra, c.mean, c.trueLon), c.want, 0.01)
	}
}

// A mean and a true longitude on either side of 0/360 must be averaged across the line, not through
// 180 (Raman does not discuss this; the arithmetic mean of 359 and 1 is 0, not 180).
func TestCheshtaBalaAcrossZero(t *testing.T) {
	// mean 359, true 1 -> average 0; Seeghrochcha 181 -> kendra 181 -> reduced 179 -> 59.67
	near(t, "across 0", cheshtaBala(181, 359, 1), 179.0/3.0, 1e-9)
	near(t, "across 0 reversed", cheshtaBala(181, 1, 359), 179.0/3.0, 1e-9)
	// without crossing the line the same geometry gives the same value
	near(t, "plain", cheshtaBala(181+100, 99, 101), 179.0/3.0, 1e-9)
	// Kendra 0 gives 0 and 180 gives 60
	near(t, "kendra 0", cheshtaBala(100, 100, 100), 0, 1e-9)
	near(t, "kendra 180", cheshtaBala(280, 50, 150), 60, 1e-9)
}

// Raman art. 89-92: the interval from the epoch (0 h on 1 Jan 1900 at 76 E). Example 1: 1 pm on
// 8 Aug 1912 at 91 E is 4602.5 days; example 39: Raman's chart is 6862.581 days (he rounds 14 h to 0.581).
func TestRamanInterval(t *testing.T) {
	jd := swisseph.Julday(1912, 8, 8, 13.0-91.0/15.0, swisseph.GregCal)
	near(t, "example 1", ramanInterval(jd), 4602.5, 0.001)
	lon := 77.0 + 35.0/60
	jd = swisseph.Julday(1918, 10, 16, 14.0+6.0/60+16.0/3600-lon/15, swisseph.GregCal)
	near(t, "standard horoscope", ramanInterval(jd), 6862.581, 0.005)
}

// Raman examples 1-4 (8 Aug 1912, 1 pm at 91 E: interval 4602.5, t = 12) and 41-48 (standard
// horoscope: interval 6862.578, t = 18). Example 3 (Jupiter, 240.48) is not used: it adds 1.66 for
// two days where his Table VI gives 0.17 and takes 332.29 for 4,000 days where the table has 332.39.
func TestRamanMeanPositions(t *testing.T) {
	a := computeRamanMeans(4602.5, 12)
	near(t, "Sun 1912", a.Sun, 113.6930, 0.01)
	near(t, "Mars 1912", a.Mars, 162.020, 0.01)
	near(t, "Saturn 1912", a.Saturn, 35.65, 0.01)
	b := computeRamanMeans(6862.578, 18)
	near(t, "Sun", b.Sun, 181.2275, 0.01)
	near(t, "Mars", b.Mars, 266.34, 0.01)
	near(t, "Jupiter", b.Jupiter, 66.91, 0.1) // his Table VI hundreds column row 8 is 0.1 high
	near(t, "Saturn", b.Saturn, 111.23, 0.02)
	near(t, "Mercury Seeghrochcha", b.MercurySeeghra, 174.49, 0.02)
	near(t, "Venus Seeghrochcha", b.VenusSeeghra, 158.35, 0.02)
}

func TestCheshtaSunAndMoonHaveNone(t *testing.T) {
	ctx := domain.CalculationContext{JulianDayUT: 2421885.0, Input: domain.BirthInput{Timezone: 5.5}}
	lons := map[string]float64{"Sun": 180, "Moon": 311, "Mars": 229}
	near(t, "Sun", CalculateCheshtaBala("Sun", &ctx, lons), 0, 0)
	near(t, "Moon", CalculateCheshtaBala("Moon", &ctx, lons), 0, 0)
}

// Raman's standard horoscope with his true longitudes (art. 6) in his own zodiac: Cheshta Bala of
// example 51, from mean positions computed for his interval.
func TestCheshtaBalaRamanChart(t *testing.T) {
	lon := 77.0 + 35.0/60
	jd := swisseph.Julday(1918, 10, 16, 14.0+6.0/60+16.0/3600-lon/15, swisseph.GregCal)
	m := computeRamanMeans(ramanInterval(jd), 18)
	lons := map[string]float64{"Mars": dms(229, 30, 34), "Mercury": dms(181, 31, 34), "Jupiter": dms(84, 0, 49), "Venus": dms(171, 9, 56), "Saturn": dms(124, 22, 41)}
	want := map[string]float64{"Mars": 22.23, "Mercury": 2.30, "Jupiter": 35.26, "Venus": 5.95, "Saturn": 21.14}
	for p, w := range want {
		near(t, p, cheshtaFromMeans(p, m, lons[p]), w, 0.05)
	}
}

func ramanChartContext(t *testing.T, ayanamsa string) domain.CalculationContext {
	t.Helper()
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	lon := 77.0 + 35.0/60
	utc, _ := astronomyTime.ParseLocalToUTC("1918-10-16", "14:06:16", lon/15)
	return domain.CalculationContext{
		Input:       domain.BirthInput{Latitude: 13, Longitude: lon, DateOfBirth: "1918-10-16", TimeOfBirth: "14:06:16", Timezone: lon / 15},
		Config:      domain.CalculationConfig{AyanamsaMode: ephemeris.GetAyanamsaMode(ayanamsa), HouseCode: ephemeris.GetHouseSystemCode("Placidus")},
		UTCTime:     utc,
		JulianDayUT: astronomyTime.UTCToJulianDay(utc),
	}
}

// With the Raman ayanamsa the engine's own planet positions reproduce his Cheshta Balas (example 51).
// The same sky must give the same Cheshta Bala whichever ayanamsa is chosen, because the mean
// longitudes are moved into the chart's zodiac.
func TestCheshtaBalaDoesNotDependOnAyanamsa(t *testing.T) {
	want := map[string]float64{"Mars": 22.23, "Mercury": 2.30, "Jupiter": 35.26, "Venus": 5.95, "Saturn": 21.14}
	for _, ay := range []string{"raman", "lahiri", "krishnamurti", "fagan-bradley"} {
		ctx := ramanChartContext(t, ay)
		res, err := CalculateShadbala(&ctx)
		if err != nil {
			t.Fatal(err)
		}
		for _, p := range res.Planets {
			switch p.Planet {
			case "Sun", "Moon":
				near(t, ay+" "+p.Planet, p.CheshtaBala, 0, 0)
			default:
				near(t, ay+" "+p.Planet, p.CheshtaBala, want[p.Planet], 0.08)
			}
		}
	}
}
