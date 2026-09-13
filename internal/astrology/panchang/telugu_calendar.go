package panchang

import (
	"math"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

var samvatsaraNames = []string{
	"Prabhava", "Vibhava", "Shukla", "Pramodoota", "Prajothpatti", "Angeerasa",
	"Sreemukha", "Bhava", "Yuva", "Dhata", "Eeswara", "Bahudhanya", "Pramadi",
	"Vikrama", "Vrusha", "Chitrabhanu", "Svabhanu", "Tarana", "Parthiva", "Vyaya",
	"Sarvajit", "Sarvadhari", "Virodhi", "Vikruti", "Khara", "Nandana", "Vijaya",
	"Jaya", "Manmadha", "Durmukhi", "Hevilambi", "Vilambi", "Vikari", "Sarvari",
	"Plava", "Shubhakritu", "Shobhakritu", "Krodhi", "Viswavasu", "Parabhava",
	"Plavanga", "Keelaka", "Saumya", "Sadharana", "Virodhikritu", "Paridhavi",
	"Pramadeecha", "Ananda", "Rakshasa", "Nala", "Pingala", "Kalayukti",
	"Siddharthi", "Raudri", "Durmathi", "Dundubhi", "Rudhirodgari", "Raktakshi",
	"Krodhana", "Akshaya",
}

var rutuNames = []string{"Vasantha", "Greeshma", "Varsha", "Sharad", "Hemantha", "Sisira"}
var masamNames = []string{"Chaitra", "Vaishakha", "Jyeshtha", "Ashadha", "Shravana", "Bhadrapada", "Ashwayuja", "Kartika", "Margashira", "Pushya", "Magha", "Phalguna"}

func calcSunSidereal(jd float64) float64 {
	tflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed | swisseph.FlagSidereal)
	sun := swisseph.CalcUT(jd, swisseph.Sun, tflag).Data[0]
	return math.Mod(sun+360.0, 360.0)
}

func CalculateTeluguCalendar(jd float64) domain.TeluguCalendar {
	// 1. Ayana
	sunSid := calcSunSidereal(jd)
	ayana := "Dakshinayana"
	if sunSid >= 270.0 || sunSid < 90.0 {
		ayana = "Uttarayana"
	}

	// 2. Paksha
	tithiAngle := calcTithiAngle(jd)
	paksha := "Shukla"
	if tithiAngle >= 180.0 {
		paksha = "Krishna"
	}

	// 3. Find Masam
	// Find previous and next Amavasya (0 degrees tithi)
	guessPrev := jd - (tithiAngle / 12.1907)
	guessNext := jd + ((360.0 - tithiAngle) / 12.1907)

	prevAmavasya := bisectionSearch(guessPrev-2.0, guessPrev+2.0, 0.0, calcTithiAngle)
	nextAmavasya := bisectionSearch(guessNext-2.0, guessNext+2.0, 0.0, calcTithiAngle)

	sunStart := calcSunSidereal(prevAmavasya)
	sunEnd := calcSunSidereal(nextAmavasya)

	signStart := int(math.Floor(sunStart / 30.0))
	signEnd := int(math.Floor(sunEnd / 30.0))

	isAdhika := false
	var masamIndex int

	if signStart == signEnd {
		isAdhika = true
		masamIndex = (signStart + 1) % 12
	} else {
		masamIndex = signEnd
	}

	masamName := masamNames[masamIndex]
	if isAdhika {
		masamName = "Adhika " + masamName
	}

	// 4. Rutu
	rutuIndex := (masamIndex / 2) % 6
	rutu := rutuNames[rutuIndex]

	// 5. Samvatsara
	t := jdToUTC(jd)
	gregYear := t.Year()
	gregMonth := t.Month()

	var shakaYear int
	if gregMonth < time.March {
		shakaYear = gregYear - 79
	} else if gregMonth > time.April {
		shakaYear = gregYear - 78
	} else {
		if masamIndex >= 10 {
			shakaYear = gregYear - 79
		} else {
			shakaYear = gregYear - 78
		}
	}

	samvatsaraIndex := (shakaYear + 11) % 60
	if samvatsaraIndex < 0 {
		samvatsaraIndex += 60
	}
	samvatsara := samvatsaraNames[samvatsaraIndex]

	return domain.TeluguCalendar{
		Samvatsara: samvatsara,
		Ayana:      ayana,
		Rutu:       rutu,
		Masam:      masamName,
		Paksha:     paksha,
	}
}
