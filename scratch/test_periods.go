package main

import (
	"fmt"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
)

func calculateDailyPeriods(sunriseJD, sunsetJD float64, weekday int) ([][2]float64, [][2]float64, [][2]float64) {
	dayDur := sunsetJD - sunriseJD
	partDur := dayDur / 8.0

	rahuMap := []int{7, 1, 6, 4, 5, 3, 2}
	yamaMap := []int{4, 3, 2, 1, 0, 6, 5}

	rIdx := float64(rahuMap[weekday])
	yIdx := float64(yamaMap[weekday])

	rahu := [][2]float64{{sunriseJD + rIdx*partDur, sunriseJD + (rIdx + 1)*partDur}}
	yama := [][2]float64{{sunriseJD + yIdx*partDur, sunriseJD + (yIdx + 1)*partDur}}

	muhurtaDur := dayDur / 15.0
	nightDur := 1.0 - dayDur // Approximate 24h - day duration for night
	nightMuhurtaDur := nightDur / 15.0

	var durmuhurtams [][2]float64
	switch weekday {
	case 0:
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 13*muhurtaDur, sunriseJD + 14*muhurtaDur})
	case 1:
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 8*muhurtaDur, sunriseJD + 9*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 11*muhurtaDur, sunriseJD + 12*muhurtaDur})
	case 2:
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 3*muhurtaDur, sunriseJD + 4*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunsetJD + 7*nightMuhurtaDur, sunsetJD + 8*nightMuhurtaDur})
	case 3:
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 7*muhurtaDur, sunriseJD + 8*muhurtaDur})
	case 4:
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 5*muhurtaDur, sunriseJD + 6*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 11*muhurtaDur, sunriseJD + 12*muhurtaDur})
	case 5:
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 3*muhurtaDur, sunriseJD + 4*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 11*muhurtaDur, sunriseJD + 12*muhurtaDur})
	case 6:
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 0*muhurtaDur, sunriseJD + 1*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 1*muhurtaDur, sunriseJD + 2*muhurtaDur})
	}

	return rahu, yama, durmuhurtams
}

func jdToUTC(jd float64) time.Time {
	res := swisseph.Revjul(jd, swisseph.GregCal)
	return time.Date(int(res.Year), time.Month(res.Month), int(res.Day), int(res.Hour), int((res.Hour-float64(int(res.Hour)))*60), 0, 0, time.UTC)
}

func main() {
	ephemeris.Init("ephe_data")
	defer ephemeris.Close()
	
	loc, _ := time.LoadLocation("Asia/Kolkata")
	
	// Exact Sunrise/Sunset for Sattenapalle on Sep 9 2026
	lat, lon := 16.3938, 80.1522
	
	// Start of Day
	utcTime := time.Date(2026, 9, 9, 0, 0, 0, 0, loc).UTC()
	h := float64(utcTime.Hour()) + float64(utcTime.Minute())/60.0 + float64(utcTime.Second())/3600.0
	startOfDayJD := swisseph.Julday(int32(utcTime.Year()), int32(utcTime.Month()), int32(utcTime.Day()), h, swisseph.GregCal)
	
	rsflag := int32(swisseph.FlagSwieph | swisseph.BitDiscCenter | swisseph.BitNoRefraction)
	resRise := swisseph.RiseTrans(startOfDayJD, swisseph.Sun, "", rsflag, int32(swisseph.CalcRise), [3]float64{lon, lat, 0}, 0, 0)
	resSet := swisseph.RiseTrans(startOfDayJD, swisseph.Sun, "", rsflag, int32(swisseph.CalcSet), [3]float64{lon, lat, 0}, 0, 0)
	
	rahu, yama, dur := calculateDailyPeriods(resRise.Time, resSet.Time, 3) // Wednesday is 3
	
	fmt.Printf("Sunrise: %s\n", jdToUTC(resRise.Time).In(loc).Format("03:04 PM"))
	fmt.Printf("Sunset:  %s\n", jdToUTC(resSet.Time).In(loc).Format("03:04 PM"))
	fmt.Printf("Rahu Kalam: %s to %s\n", jdToUTC(rahu[0][0]).In(loc).Format("03:04 PM"), jdToUTC(rahu[0][1]).In(loc).Format("03:04 PM"))
	fmt.Printf("Yamaganda:  %s to %s\n", jdToUTC(yama[0][0]).In(loc).Format("03:04 PM"), jdToUTC(yama[0][1]).In(loc).Format("03:04 PM"))
	fmt.Printf("Durmuhurtam: %s to %s\n", jdToUTC(dur[0][0]).In(loc).Format("03:04 PM"), jdToUTC(dur[0][1]).In(loc).Format("03:04 PM"))
}
