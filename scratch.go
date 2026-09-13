package main

import (
	"fmt"
	"math"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
)

func getRiseSet(jd float64, lat float64, lon float64, body int, rsmiFlags int32, topo bool) (float64, float64) {
	geopos := [3]float64{lon, lat, 0}
	epheflag := int32(swisseph.FlagSwieph)
	if topo {
		epheflag |= swisseph.FlagTopoctr
		swisseph.SetTopo(lon, lat, 0)
	}

	rsmiRise := int32(swisseph.CalcRise) | rsmiFlags
	resRise := swisseph.RiseTrans(jd, int32(body), "", epheflag, rsmiRise, geopos, 0, 0)
	fracRise := resRise.Time + 0.5 - math.Floor(resRise.Time+0.5)

	rsmiSet := int32(swisseph.CalcSet) | rsmiFlags
	resSet := swisseph.RiseTrans(jd, int32(body), "", epheflag, rsmiSet, geopos, 0, 0)
	fracSet := resSet.Time + 0.5 - math.Floor(resSet.Time+0.5)

	return fracRise * 24.0, fracSet * 24.0
}

func printTime(name string, tUTC float64) {
	tLocal := tUTC + 5.5
	sh := int(tLocal)
	sm := int((tLocal - float64(sh)) * 60)
	ss := int((tLocal - float64(sh) - float64(sm)/60.0) * 3600)

	ampm := "AM"
	if sh >= 12 {
		ampm = "PM"
		if sh > 12 {
			sh -= 12
		}
	}
	if sh == 0 {
		sh = 12
	}

	fmt.Printf("%s: %02d:%02d:%02d %s\n", name, sh, sm, ss, ampm)
}

func main() {
	ephemeris.Init("ephe_data")
	utcTime, _ := astronomyTime.ParseLocalToUTC("2005-11-23", "00:00:00", 5.5)
	jd := astronomyTime.UTCToJulianDay(utcTime)
	lat, lon := 16.3900, 80.1500

	fmt.Println("MOON")
	mr1, ms1 := getRiseSet(jd, lat, lon, swisseph.Moon, 0, false)
	printTime("Standard Moon (0 flags)", mr1)
	printTime("Standard Set", ms1)

	mr2, ms2 := getRiseSet(jd, lat, lon, swisseph.Moon, int32(swisseph.BitDiscCenter|swisseph.BitNoRefraction), false)
	printTime("Center+NoRefrac (Geocentric/Default)", mr2)
	printTime("Center+NoRefrac Set", ms2)

	mr3, ms3 := getRiseSet(jd, lat, lon, swisseph.Moon, int32(swisseph.BitDiscCenter|swisseph.BitNoRefraction), true)
	printTime("Center+NoRefrac (Topocentric)", mr3)
	printTime("Center+NoRefrac Set", ms3)

	mr4, ms4 := getRiseSet(jd, lat, lon, swisseph.Moon, int32(swisseph.BitHinduRising), false)
	printTime("Hindu Rising Flag", mr4)
	printTime("Hindu Rising Set", ms4)
}
