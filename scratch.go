package main

import (
	"fmt"
	"math"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
)

func getRiseSet(jd float64, lat float64, lon float64, rsmiFlags int32, body int) (float64, float64) {
	geopos := [3]float64{lon, lat, 0}
	epheflag := int32(swisseph.FlagSwieph)

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
	fmt.Printf("%s: %02d:%02d:%02d\n", name, sh, sm, ss)
}

func main() {
	ephemeris.Init("ephe_data")
	utcTime, _ := astronomyTime.ParseLocalToUTC("2026-09-13", "12:00:00", 5.5)
	jd := astronomyTime.UTCToJulianDay(utcTime)

	lat, lon := 17.38405, 78.45636

	fmt.Println("HYDERABAD SUNRISE/SUNSET SEP 13 2026")
	r1, s1 := getRiseSet(jd, lat, lon, 0, swisseph.Sun)                                                      // Standard
	r2, s2 := getRiseSet(jd, lat, lon, int32(swisseph.BitDiscCenter|swisseph.BitNoRefraction), swisseph.Sun) // Hindu

	printTime("Sun Standard Rise", r1)
	printTime("Sun Standard Set", s1)
	printTime("Sun Hindu Rise", r2)
	printTime("Sun Hindu Set", s2)

	r3, s3 := getRiseSet(jd, lat, lon, 0, swisseph.Moon)                                                      // Standard
	r4, s4 := getRiseSet(jd, lat, lon, int32(swisseph.BitDiscCenter|swisseph.BitNoRefraction), swisseph.Moon) // Hindu
	printTime("Moon Standard Rise", r3)
	printTime("Moon Standard Set", s3)
	printTime("Moon Hindu Rise", r4)
	printTime("Moon Hindu Set", s4)
}
