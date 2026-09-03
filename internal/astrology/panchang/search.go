package panchang

import "math"

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

func findElementBoundaries(jd, currentPos, interval float64, calcFunc func(float64) float64) (float64, float64) {
	currentIndex := math.Floor(currentPos / interval)
	startBoundary := currentIndex * interval
	endBoundary := math.Mod((currentIndex+1)*interval, 360)

	startSearch := jd - 2.5
	endSearch := jd + 2.5

	startJD := bisectionSearch(startSearch, jd, startBoundary, calcFunc)
	endJD := bisectionSearch(jd, endSearch, endBoundary, calcFunc)

	return startJD, endJD
}
