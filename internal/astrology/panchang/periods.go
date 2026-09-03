package panchang

func calculateDailyPeriods(sunriseJD, sunsetJD float64, weekday int) ([2]float64, [2]float64, [2]float64) {
	dayDur := sunsetJD - sunriseJD
	partDur := dayDur / 8.0

	rahuMap := []int{7, 1, 6, 4, 5, 3, 2}
	yamaMap := []int{4, 3, 2, 1, 0, 6, 5}
	guliMap := []int{6, 5, 4, 3, 2, 1, 0}

	rIdx := float64(rahuMap[weekday])
	yIdx := float64(yamaMap[weekday])
	gIdx := float64(guliMap[weekday])

	rahu := [2]float64{sunriseJD + rIdx*partDur, sunriseJD + (rIdx+1)*partDur}
	yama := [2]float64{sunriseJD + yIdx*partDur, sunriseJD + (yIdx+1)*partDur}
	guli := [2]float64{sunriseJD + gIdx*partDur, sunriseJD + (gIdx+1)*partDur}

	return rahu, yama, guli
}
