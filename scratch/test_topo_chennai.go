package main

import (
	"fmt"
	"math"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
)

func calcTithiAngleTopocentric(jd, lat, lon float64) float64 {
	tflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed | swisseph.FlagTopoctr)
	swisseph.SetTopo(lon, lat, 0)
	sun := swisseph.CalcUT(jd, swisseph.Sun, tflag).Data[0]
	moon := swisseph.CalcUT(jd, swisseph.Moon, tflag).Data[0]
	return math.Mod(moon-sun+360.0, 360.0)
}

func bisectionSearch(startJD, endJD, targetVal float64, calcFunc func(float64) float64) float64 {
	low := startJD
	high := endJD
	for i := 0; i < 50; i++ {
		mid := (low + high) / 2
		val := calcFunc(mid)
		diff := math.Mod(val-targetVal+540, 360) - 180
		if math.Abs(diff) < 0.000001 {
			return mid
		}
		if diff < 0 {
			low = mid
		} else {
			high = mid
		}
	}
	return (low + high) / 2
}

func jdToUTC(jd float64) time.Time {
	res := swisseph.Revjul(jd, swisseph.GregCal)
	y, m, d, h := res.Year, res.Month, res.Day, res.Hour
	hr := int(h)
	min := int((h - float64(hr)) * 60)
	sec := int(math.Round((h - float64(hr) - float64(min)/60) * 3600))
	if sec >= 60 { sec -= 60; min++ }
	if min >= 60 { min -= 60; hr++ }
	return time.Date(int(y), time.Month(m), int(d), hr, min, sec, 0, time.UTC)
}

func main() {
	ephemeris.Init("ephe_data")
	defer ephemeris.Close()
	loc, _ := time.LoadLocation("Asia/Kolkata")
	utcTime := time.Date(2026, 9, 9, 6, 30, 0, 0, time.UTC)
	h := float64(utcTime.Hour()) + float64(utcTime.Minute())/60.0 + float64(utcTime.Second())/3600.0
	jd := swisseph.Julday(int32(utcTime.Year()), int32(utcTime.Month()), int32(utcTime.Day()), h, swisseph.GregCal)
	
	lat, lon := 13.0827, 80.2707 // Chennai
	
	endBoundaryTithi := 336.0
	endJDTithiTopo := bisectionSearch(jd-1.0, jd+1.0, endBoundaryTithi, func(j float64) float64 { return calcTithiAngleTopocentric(j, lat, lon) })
	fmt.Printf("Topocentric Chennai: %s\n", jdToUTC(endJDTithiTopo).In(loc).Format("03:04:05 PM"))
}
