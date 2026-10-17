package shadbala

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

var (
	sweepAyanamsas = []string{"lahiri", "raman", "krishnamurti", "kpnew", "fagan-bradley"}
	sweepHouses    = []string{"placidus", "koch", "porphyry", "equal", "whole sign", "regiomontanus"}
)

type sweepChart struct {
	date, clock    string
	lat, lon, tz   float64
	ayanamsa, hsys string
}

func randomCharts(n int, seed int64) []sweepChart {
	r := rand.New(rand.NewSource(seed))
	out := make([]sweepChart, n)
	start := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range out {
		d := start.AddDate(0, 0, r.Intn(73000)) // to about 2099
		lon := -180 + r.Float64()*360
		out[i] = sweepChart{
			date:     d.Format("2006-01-02"),
			clock:    fmt.Sprintf("%02d:%02d:%02d", r.Intn(24), r.Intn(60), r.Intn(60)),
			lat:      -55 + r.Float64()*110,
			lon:      lon,
			tz:       math.Round(lon/15*2) / 2,
			ayanamsa: sweepAyanamsas[r.Intn(len(sweepAyanamsas))],
			hsys:     sweepHouses[r.Intn(len(sweepHouses))],
		}
	}
	return out
}

func (c sweepChart) context(t testing.TB) domain.CalculationContext {
	utc, err := astronomyTime.ParseLocalToUTC(c.date, c.clock, c.tz)
	if err != nil {
		t.Fatal(err)
	}
	return domain.CalculationContext{
		Input:       domain.BirthInput{Latitude: c.lat, Longitude: c.lon, DateOfBirth: c.date, TimeOfBirth: c.clock, Timezone: c.tz},
		Config:      domain.CalculationConfig{AyanamsaMode: ephemeris.GetAyanamsaMode(c.ayanamsa), HouseCode: ephemeris.GetHouseSystemCode(c.hsys)},
		UTCTime:     utc,
		JulianDayUT: astronomyTime.UTCToJulianDay(utc),
	}
}

func inRange(t *testing.T, who, what string, v, lo, hi float64) {
	t.Helper()
	if math.IsNaN(v) || math.IsInf(v, 0) || v < lo-1e-6 || v > hi+1e-6 {
		t.Errorf("%s: %s = %v, want within [%v, %v]", who, what, v, lo, hi)
	}
}

func oneOf(t *testing.T, who, what string, v float64, allowed ...float64) {
	t.Helper()
	for _, a := range allowed {
		if math.Abs(v-a) < 1e-9 {
			return
		}
	}
	t.Errorf("%s: %s = %v, want one of %v", who, what, v, allowed)
}

// Properties that must hold for every chart, over several hundred random ones: ranges of every component,
// who gets the 45 / 60 / 15 / 30 lordships (exactly one planet each), the groupings of Nathonnatha and
// Paksha Bala, the sums, and that the Hindu weekday agrees with the calendar.
func TestShadbalaInvariants(t *testing.T) {
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	n := 600
	if testing.Short() {
		n = 60
	}
	for i, c := range randomCharts(n, 20261017) {
		ctx := c.context(t)
		res, err := CalculateShadbala(&ctx)
		if err != nil {
			t.Fatalf("chart %d %+v: %v", i, c, err)
		}
		if len(res.Planets) != 7 {
			t.Fatalf("chart %d: %d planets", i, len(res.Planets))
		}
		by := map[string]domain.PlanetShadbala{}
		var dina, hora, varsha, masa, tri, yud float64
		var nDina, nHora, nVarsha, nMasa, nTri int
		for _, p := range res.Planets {
			by[p.Planet] = p
			who := fmt.Sprintf("chart %d (%s %s %.1f,%.1f %s %s) %s", i, c.date, c.clock, c.lat, c.lon, c.ayanamsa, c.hsys, p.Planet)
			s, k := p.SthanaBala, p.KalaBala
			inRange(t, who, "uchcha", s.UchchaBala, 0, 60)
			inRange(t, who, "saptavargaja", s.SaptavargajaBala, 13.125, 225)
			oneOf(t, who, "ojayugma", s.OjayugmaBala, 0, 15, 30)
			oneOf(t, who, "kendradi", s.KendradiBala, 15, 30, 60)
			oneOf(t, who, "drekkana", s.DrekkanaBala, 0, 15)
			inRange(t, who, "sthana total", s.Total, 0, 60+225+30+60+15)
			near(t, who+" sthana sum", s.Total, s.UchchaBala+s.SaptavargajaBala+s.OjayugmaBala+s.KendradiBala+s.DrekkanaBala, 1e-9)
			inRange(t, who, "dig", p.DigBala, 0, 60)
			inRange(t, who, "nathonnatha", k.NathonnathaBala, 0, 60)
			if p.Planet == "Moon" {
				inRange(t, who, "paksha", k.PakshaBala, 0, 120)
			} else {
				inRange(t, who, "paksha", k.PakshaBala, 0, 60)
			}
			oneOf(t, who, "tribhaga", k.TribhagaBala, 0, 60)
			oneOf(t, who, "varsha", k.VarshaBala, 0, 15)
			oneOf(t, who, "masa", k.MasaBala, 0, 30)
			oneOf(t, who, "dina", k.DinaBala, 0, 45)
			oneOf(t, who, "hora", k.HoraBala, 0, 60)
			if p.Planet == "Sun" {
				inRange(t, who, "ayana", k.AyanaBala, 0, 120)
			} else {
				inRange(t, who, "ayana", k.AyanaBala, 0, 60)
			}
			sum := k.NathonnathaBala + k.PakshaBala + k.TribhagaBala + k.VarshaBala + k.MasaBala + k.DinaBala + k.HoraBala + k.AyanaBala + k.YuddhaBala
			near(t, who+" kala sum", k.Total, sum, 1e-9)
			inRange(t, who, "cheshta", p.CheshtaBala, 0, 60)
			if p.Planet == "Sun" || p.Planet == "Moon" {
				near(t, who+" cheshta", p.CheshtaBala, 0, 0)
			}
			near(t, who+" naisargika", p.NaisargikaBala, NaisargikaBala[p.Planet], 0.006)
			inRange(t, who, "drik", p.DrikBala, -160, 160)
			total := s.Total + p.DigBala + k.Total + p.CheshtaBala + p.NaisargikaBala + p.DrikBala
			near(t, who+" total", p.TotalShadbala, total, 0.03)
			near(t, who+" rupas", p.Rupas, p.TotalShadbala/60, 0.001)
			near(t, who+" ratio", p.StrengthRatio, p.Rupas/MinimumRequirements[p.Planet], 0.002)
			if p.MeetsMinimum != (p.Rupas >= MinimumRequirements[p.Planet]) {
				t.Errorf("%s: MeetsMinimum inconsistent", who)
			}
			dina += k.DinaBala
			hora += k.HoraBala
			varsha += k.VarshaBala
			masa += k.MasaBala
			tri += k.TribhagaBala
			yud += k.YuddhaBala
			for _, v := range []struct {
				n   *int
				val float64
			}{{&nDina, k.DinaBala}, {&nHora, k.HoraBala}, {&nVarsha, k.VarshaBala}, {&nMasa, k.MasaBala}, {&nTri, k.TribhagaBala}} {
				if v.val != 0 {
					*v.n++
				}
			}
		}
		who := fmt.Sprintf("chart %d (%s %s)", i, c.date, c.clock)
		if nDina != 1 || nHora != 1 || nVarsha != 1 || nMasa != 1 {
			t.Errorf("%s: lords of day/hora/year/month are %d/%d/%d/%d planets, want 1 each", who, nDina, nHora, nVarsha, nMasa)
		}
		if nTri != 2 || by["Jupiter"].KalaBala.TribhagaBala != 60 {
			t.Errorf("%s: Tribhaga Bala goes to %d planets (Jupiter %v), want 2 including Jupiter", who, nTri, by["Jupiter"].KalaBala.TribhagaBala)
		}
		near(t, who+" yuddha sum", yud, 0, 1e-9)
		// Nathonnatha: Sun, Jupiter, Venus share one value, Moon, Mars, Saturn the complement, Mercury 60.
		nb := func(p string) float64 { return by[p].KalaBala.NathonnathaBala }
		near(t, who+" nath Sun/Jupiter", nb("Sun"), nb("Jupiter"), 1e-9)
		near(t, who+" nath Sun/Venus", nb("Sun"), nb("Venus"), 1e-9)
		near(t, who+" nath Moon/Mars", nb("Moon"), nb("Mars"), 1e-9)
		near(t, who+" nath Moon/Saturn", nb("Moon"), nb("Saturn"), 1e-9)
		near(t, who+" nath Sun+Moon", nb("Sun")+nb("Moon"), 60, 1e-9)
		near(t, who+" nath Mercury", nb("Mercury"), 60, 0)
		// Paksha: Sun, Mars, Saturn equal; Jupiter and Venus equal and the complement; Moon doubled.
		pk := func(p string) float64 { return by[p].KalaBala.PakshaBala }
		near(t, who+" paksha Sun/Mars", pk("Sun"), pk("Mars"), 1e-9)
		near(t, who+" paksha Sun/Saturn", pk("Sun"), pk("Saturn"), 1e-9)
		near(t, who+" paksha Jupiter/Venus", pk("Jupiter"), pk("Venus"), 1e-9)
		near(t, who+" paksha Sun+Jupiter", pk("Sun")+pk("Jupiter"), 60, 1e-9)
		near(t, who+" paksha Moon", pk("Moon"), 2*pk("Jupiter"), 1e-9)
		if m := pk("Mercury"); math.Abs(m-pk("Jupiter")) > 1e-9 && math.Abs(m-pk("Sun")) > 1e-9 {
			t.Errorf("%s: Mercury's Paksha Bala %v is neither the benefic nor the malefic value", who, m)
		}
	}
}

// For a birth at least four hours after sunrise and before midnight (local clock, so the Hindu day is the
// calendar day) the Dina lord must be the weekday of the calendar date, and the Ahargana weekday agrees.
func TestDinaLordIsTheCalendarWeekday(t *testing.T) {
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	n := 400
	if testing.Short() {
		n = 40
	}
	checked := 0
	for _, c := range randomCharts(n, 7) {
		// afternoon births only: 13:00 to 17:59 clock time, well after sunrise and before any sunset
		c.clock = fmt.Sprintf("%02d:30:00", 13+checked%5)
		c.tz = math.Round(c.lon/15*2) / 2
		ctx := c.context(t)
		res, err := CalculateShadbala(&ctx)
		if err != nil {
			t.Fatal(err)
		}
		d, _ := time.Parse("2006-01-02", c.date)
		want := WeekdayLords[int(d.Weekday())]
		for _, p := range res.Planets {
			if (p.Planet == want) != (p.KalaBala.DinaBala == 45) {
				t.Errorf("%s %s at %.1f,%.1f: Dina Bala of %s is %v, weekday lord should be %s", c.date, c.clock, c.lat, c.lon, p.Planet, p.KalaBala.DinaBala, want)
			}
		}
		jdn := jdnOf(int32(d.Year()), int32(d.Month()), int32(d.Day()))
		if got := varadhipathi(aharganaOfJDN(jdn)); got != want {
			t.Errorf("%s: Ahargana weekday %s, calendar weekday %s", c.date, got, want)
		}
		checked++
	}
}

// Components that depend only on time and place must not change with the ayanamsa or the house system;
// Cheshta Bala must not depend on the ayanamsa either (the frame shift), and the results must be repeatable.
func TestComponentsIndependentOfAyanamsaAndHouses(t *testing.T) {
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	n := 60
	if testing.Short() {
		n = 10
	}
	for i, c := range randomCharts(n, 99) {
		var ref []domain.PlanetShadbala
		for j, ay := range sweepAyanamsas {
			cc := c
			cc.ayanamsa, cc.hsys = ay, sweepHouses[j%len(sweepHouses)]
			ctx := cc.context(t)
			res, err := CalculateShadbala(&ctx)
			if err != nil {
				t.Fatal(err)
			}
			if ref == nil {
				ref = res.Planets
				continue
			}
			for k, p := range res.Planets {
				r := ref[k]
				who := fmt.Sprintf("chart %d %s %s vs %s", i, p.Planet, ay, sweepAyanamsas[0])
				kb, rb := p.KalaBala, r.KalaBala
				near(t, who+" nathonnatha", kb.NathonnathaBala, rb.NathonnathaBala, 1e-9)
				near(t, who+" tribhaga", kb.TribhagaBala, rb.TribhagaBala, 0)
				near(t, who+" dina/hora/varsha/masa", kb.DinaBala+kb.HoraBala+kb.VarshaBala+kb.MasaBala, rb.DinaBala+rb.HoraBala+rb.VarshaBala+rb.MasaBala, 0)
				near(t, who+" ayana", kb.AyanaBala, rb.AyanaBala, 1e-9)
				if p.Planet == "Mercury" {
					// Whether Mercury "joins a malefic" is judged by sidereal sign, so near a sign boundary
					// two ayanamsas can disagree; either way it is the benefic or the malefic value.
					bene := res.Planets[4].KalaBala.PakshaBala // Jupiter
					mal := res.Planets[0].KalaBala.PakshaBala  // Sun
					if math.Abs(kb.PakshaBala-bene) > 0.01 && math.Abs(kb.PakshaBala-mal) > 0.01 {
						t.Errorf("%s: paksha %v is neither the benefic %v nor the malefic %v value", who, kb.PakshaBala, bene, mal)
					}
				} else {
					near(t, who+" paksha", kb.PakshaBala, rb.PakshaBala, 0.01) // the Sun-Moon arc differs by < 0.03 deg between ayanamsa formulas
				}
				near(t, who+" cheshta", p.CheshtaBala, r.CheshtaBala, 0.05)
				near(t, who+" naisargika", p.NaisargikaBala, r.NaisargikaBala, 0)
			}
		}
	}
}

// Many requests at once (the ephemeris is shared and thread-bound) must give the same answers as one at a time.
func TestShadbalaConcurrentCalculations(t *testing.T) {
	if err := ephemeris.Init(mustAbs(t, "../../../ephe_data")); err != nil {
		t.Fatal(err)
	}
	charts := randomCharts(24, 5)
	want := make([]string, len(charts))
	for i, c := range charts {
		ctx := c.context(t)
		res, err := CalculateShadbala(&ctx)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.Marshal(res)
		want[i] = string(b)
	}
	var wg sync.WaitGroup
	errs := make(chan string, 8*len(charts))
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for k := range charts {
				i := (k + g*3) % len(charts)
				ctx := charts[i].context(t)
				res, err := CalculateShadbala(&ctx)
				if err != nil {
					errs <- err.Error()
					return
				}
				b, _ := json.Marshal(res)
				if string(b) != want[i] {
					errs <- fmt.Sprintf("goroutine %d chart %d differs from the sequential result", g, i)
				}
			}
		}(g)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		t.Error(e)
	}
}
