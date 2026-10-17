package planets

import (
	"math"
	"testing"

	"github.com/tejzpr/go-swisseph"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
)

// The manual sidereal longitude (apparent tropical - AyanamsaUT) must equal Swiss Ephemeris's own
// FlagSidereal result, at any date (the nutation term varies by +-17").
func TestManualSiderealMatchesFlagSidereal(t *testing.T) {
	if err := ephemeris.Init("../../../ephe_data"); err != nil {
		t.Fatal(err)
	}
	ephemeris.Lock()
	defer ephemeris.Unlock()
	swisseph.SetSidMode(swisseph.SidmLahiri, 0, 0)
	for _, y := range []int32{1950, 1985, 2000, 2024, 2040} {
		for _, body := range []int32{swisseph.Moon, swisseph.Sun, swisseph.Mars} {
			jd := swisseph.Julday(y, 6, 15, 12.0, swisseph.GregCal)
			trop := swisseph.CalcUT(jd, body, swisseph.FlagSwieph)
			sid := swisseph.CalcUT(jd, body, swisseph.FlagSwieph|swisseph.FlagSidereal)
			manual := math.Mod(trop.Data[0]-ephemeris.AyanamsaUT(jd)+360, 360)
			if d := math.Abs(manual-sid.Data[0]) * 3600; d > 0.05 {
				t.Errorf("year %d body %d: manual vs FlagSidereal differ by %.3f arcsec", y, body, d)
			}
		}
	}
}
