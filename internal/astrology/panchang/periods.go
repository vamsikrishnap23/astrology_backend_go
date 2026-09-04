package panchang

func calculateDailyPeriods(sunriseJD, sunsetJD float64, weekday int) ([2]float64, [2]float64, [][2]float64) {
	dayDur := sunsetJD - sunriseJD
	partDur := dayDur / 8.0

	rahuMap := []int{7, 1, 6, 4, 5, 3, 2}
	yamaMap := []int{4, 3, 2, 1, 0, 6, 5}

	rIdx := float64(rahuMap[weekday])
	yIdx := float64(yamaMap[weekday])

	rahu := [2]float64{sunriseJD + rIdx*partDur, sunriseJD + (rIdx+1)*partDur}
	yama := [2]float64{sunriseJD + yIdx*partDur, sunriseJD + (yIdx+1)*partDur}

	muhurtaDur := dayDur / 15.0
	nightDur := 1.0 - dayDur // Approximate 24h - day duration for night
	nightMuhurtaDur := nightDur / 15.0

	var durmuhurtams [][2]float64
	switch weekday {
	case 0: // Sunday
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 13*muhurtaDur, sunriseJD + 14*muhurtaDur})
	case 1: // Monday
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 8*muhurtaDur, sunriseJD + 9*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 11*muhurtaDur, sunriseJD + 12*muhurtaDur})
	case 2: // Tuesday
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 3*muhurtaDur, sunriseJD + 4*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunsetJD + 7*nightMuhurtaDur, sunsetJD + 8*nightMuhurtaDur})
	case 3: // Wednesday
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 7*muhurtaDur, sunriseJD + 8*muhurtaDur})
	case 4: // Thursday
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 5*muhurtaDur, sunriseJD + 6*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 11*muhurtaDur, sunriseJD + 12*muhurtaDur})
	case 5: // Friday
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 3*muhurtaDur, sunriseJD + 4*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 11*muhurtaDur, sunriseJD + 12*muhurtaDur})
	case 6: // Saturday
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 0*muhurtaDur, sunriseJD + 1*muhurtaDur})
		durmuhurtams = append(durmuhurtams, [2]float64{sunriseJD + 1*muhurtaDur, sunriseJD + 2*muhurtaDur})
	}

	return rahu, yama, durmuhurtams
}
