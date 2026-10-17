package shadbala

import (
	"math"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// Raman art. 77: the diameters of the planetary discs in seconds of arc.
func TestBimbaParimana(t *testing.T) {
	want := map[string]float64{"Mars": 9.4, "Mercury": 6.6, "Jupiter": 190.4, "Venus": 16.6, "Saturn": 158.0}
	for p, w := range want {
		near(t, p, bimbaParimana[p], w, 0)
	}
	if len(bimbaParimana) != 5 {
		t.Errorf("the Sun and Moon never fight, but the table has %d planets", len(bimbaParimana))
	}
}

// Art. 76: the planet with the smaller longitude wins; the difference of the aggregates divided by the
// difference of the diameters is added to the winner and taken from the loser.
func TestYuddhaWinnerAndAmount(t *testing.T) {
	lons := map[string]float64{"Mars": 100.2, "Venus": 100.7}
	agg := map[string]float64{"Mars": 300, "Venus": 252}
	got := yuddhaBala(lons, agg)
	want := 48.0 / (16.6 - 9.4) // 6.667
	near(t, "Mars wins", got["Mars"], want, 1e-9)
	near(t, "Venus loses", got["Venus"], -want, 1e-9)
	// The same pair with the longitudes swapped: now Venus has the smaller longitude and wins.
	got = yuddhaBala(map[string]float64{"Mars": 100.7, "Venus": 100.2}, agg)
	near(t, "Venus wins", got["Venus"], want, 1e-9)
	near(t, "Mars loses", got["Mars"], -want, 1e-9)
	// A weaker winner still wins: the amount depends only on the size of the difference.
	got = yuddhaBala(lons, map[string]float64{"Mars": 252, "Venus": 300})
	near(t, "weaker winner", got["Mars"], want, 1e-9)
}

// Art. 76: the planets must be less than one degree apart.
func TestYuddhaOrb(t *testing.T) {
	agg := map[string]float64{"Mars": 300, "Mercury": 200, "Saturn": 250}
	if got := yuddhaBala(map[string]float64{"Mars": 10, "Mercury": 10.99}, agg); len(got) != 2 {
		t.Errorf("0.99 deg apart should fight: %v", got)
	}
	if got := yuddhaBala(map[string]float64{"Mars": 10, "Mercury": 11.0}, agg); len(got) != 0 {
		t.Errorf("exactly 1 deg apart should not fight: %v", got)
	}
	if got := yuddhaBala(map[string]float64{"Mars": 10, "Mercury": 12}, agg); len(got) != 0 {
		t.Errorf("2 deg apart should not fight: %v", got)
	}
	// The distance is measured across 0/360: the planet behind (359.8) has the smaller longitude.
	got := yuddhaBala(map[string]float64{"Mars": 359.8, "Saturn": 0.3}, agg)
	if got["Mars"] <= 0 || got["Saturn"] >= 0 {
		t.Errorf("across 0/360 the planet behind should win: %v", got)
	}
}

// The Sun and Moon never enter war (art. 76), whatever their distance from a planet.
func TestYuddhaIgnoresSunAndMoon(t *testing.T) {
	lons := map[string]float64{"Sun": 50.1, "Moon": 50.2, "Mars": 50.0}
	agg := map[string]float64{"Sun": 400, "Moon": 100, "Mars": 250}
	if got := yuddhaBala(lons, agg); len(got) != 0 {
		t.Errorf("Sun and Moon must not fight: %v", got)
	}
}

// What the winner gains the loser loses, so all the Yuddha Balas add up to 0, also with several wars.
func TestYuddhaConservation(t *testing.T) {
	lons := map[string]float64{"Mars": 10.0, "Mercury": 10.4, "Jupiter": 10.8, "Venus": 200, "Saturn": 200.5}
	agg := map[string]float64{"Mars": 310, "Mercury": 180, "Jupiter": 260, "Venus": 150, "Saturn": 290}
	got := yuddhaBala(lons, agg)
	sum := 0.0
	for _, v := range got {
		sum += v
	}
	near(t, "sum", sum, 0, 1e-9)
	if len(got) != 5 {
		t.Errorf("all five planets are at war here: %v", got)
	}
	// Mars has the smallest longitude of the first three, so it wins both of its wars.
	if got["Mars"] <= 0 {
		t.Errorf("Mars should come out ahead: %v", got)
	}
}

func shadbalaOn(t *testing.T, date string) []domain.PlanetShadbala {
	t.Helper()
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	utc, _ := astronomyTime.ParseLocalToUTC(date, "12:00:00", 0)
	ctx := domain.CalculationContext{
		Input:       domain.BirthInput{Latitude: 17.385, Longitude: 78.4867, DateOfBirth: date, TimeOfBirth: "12:00:00"},
		Config:      domain.CalculationConfig{AyanamsaMode: ephemeris.GetAyanamsaMode("Lahiri"), HouseCode: ephemeris.GetHouseSystemCode("Placidus")},
		UTCTime:     utc,
		JulianDayUT: astronomyTime.UTCToJulianDay(utc),
	}
	res, err := CalculateShadbala(&ctx)
	if err != nil {
		t.Fatal(err)
	}
	return res.Planets
}

func byName(ps []domain.PlanetShadbala) map[string]domain.PlanetShadbala {
	m := map[string]domain.PlanetShadbala{}
	for _, p := range ps {
		m[p.Planet] = p
	}
	return m
}

// 1 March 2023, noon UT: Venus (347.07) and Jupiter (347.80) are 0.73 deg apart and nobody else is
// near another planet. Venus, with the smaller longitude, wins.
func TestYuddhaEngineJupiterVenus(t *testing.T) {
	p := byName(shadbalaOn(t, "2023-03-01"))
	for _, name := range []string{"Sun", "Moon", "Mars", "Mercury", "Saturn"} {
		near(t, name, p[name].KalaBala.YuddhaBala, 0, 0)
	}
	v, j := p["Venus"], p["Jupiter"]
	if v.KalaBala.YuddhaBala <= 0 || j.KalaBala.YuddhaBala >= 0 {
		t.Fatalf("Venus should win and Jupiter lose: Venus %.3f Jupiter %.3f", v.KalaBala.YuddhaBala, j.KalaBala.YuddhaBala)
	}
	near(t, "gain = loss", v.KalaBala.YuddhaBala, -j.KalaBala.YuddhaBala, 0.011)
	// the amount is the difference of the aggregates over the difference of the diameters
	agg := func(x domain.PlanetShadbala) float64 {
		return x.SthanaBala.Total + x.DigBala + x.KalaBala.Total - x.KalaBala.AyanaBala - x.KalaBala.YuddhaBala
	}
	want := math.Abs(agg(v)-agg(j)) / (190.4 - 16.6)
	near(t, "amount", v.KalaBala.YuddhaBala, want, 0.02)
	// Kala Bala total includes it
	for _, x := range []domain.PlanetShadbala{v, j} {
		k := x.KalaBala
		sum := k.NathonnathaBala + k.PakshaBala + k.TribhagaBala + k.VarshaBala + k.MasaBala + k.DinaBala + k.HoraBala + k.AyanaBala + k.YuddhaBala
		near(t, x.Planet+" Kala total", k.Total, sum, 0.05)
	}
}

// 2 March 2023: two wars at once, Jupiter-Venus and Mercury-Saturn; 3 March: 1.26 deg apart, none.
func TestYuddhaEngineTwoWarsAndNone(t *testing.T) {
	p := byName(shadbalaOn(t, "2023-03-02"))
	fighting := 0
	sum := 0.0
	for _, x := range p {
		if x.KalaBala.YuddhaBala != 0 {
			fighting++
		}
		sum += x.KalaBala.YuddhaBala
	}
	if fighting != 4 {
		t.Errorf("expected four planets at war on 2 March 2023, got %d", fighting)
	}
	near(t, "sum of Yuddha Bala", sum, 0, 0.011)
	for _, x := range shadbalaOn(t, "2023-03-03") {
		near(t, x.Planet+" on 3 March", x.KalaBala.YuddhaBala, 0, 0)
	}
}
