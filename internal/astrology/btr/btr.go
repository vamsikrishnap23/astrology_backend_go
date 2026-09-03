package btr

import (
	"math"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsi/astrology_backend_go/internal/astrology/ashtakoota"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func getAscendant(jd float64, lat float64, lon float64, ayanamsa float64, houseSystem string) (string, float64) {
	ephemeris.Mu.Lock()
	defer ephemeris.Mu.Unlock()

	hsCode := byte(ephemeris.GetHouseSystemCode(houseSystem))

	resHouses := swisseph.HousesEx(jd, 0, lat, lon, hsCode)
	if len(resHouses.Points) == 0 {
		return "Unknown", 0
	}
	tropAsc := resHouses.Points[0]
	sidAsc := math.Mod(tropAsc-ayanamsa+360.0, 360.0)

	signIdx := int(math.Floor(sidAsc / 30.0))
	degInSign := math.Mod(sidAsc, 30.0)

	signs := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}

	return signs[signIdx], degInSign
}

func getSignType(sign string) string {
	switch sign {
	case "Aries", "Cancer", "Libra", "Capricorn":
		return "Movable"
	case "Taurus", "Leo", "Scorpio", "Aquarius":
		return "Fixed"
	case "Gemini", "Virgo", "Sagittarius", "Pisces":
		return "Dual"
	}
	return "Unknown"
}

func calculateLMT(utcTime time.Time, lon float64) string {
	offsetHours := lon / 15.0
	lmtTime := utcTime.Add(time.Duration(offsetHours * float64(time.Hour)))
	return lmtTime.Format(time.RFC3339)
}

func ProcessBTR(input domain.BTRInput) (domain.BTRResult, error) {
	utcTime, err := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	if err != nil {
		return domain.BTRResult{}, err
	}

	jdUT := astronomyTime.UTCToJulianDay(utcTime)

	ephemeris.Mu.Lock()
	swisseph.SetEphePath(ephemeris.EphePath)
	amode := ephemeris.GetAyanamsaMode(input.Ayanamsa)
	swisseph.SetSidMode(int32(amode), 0, 0)
	ayanamsa := swisseph.GetAyanamsaUT(jdUT)

	geopos := [3]float64{input.Longitude, input.Latitude, 0}
	resRise := swisseph.RiseTrans(jdUT, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcRise), geopos, 0, 0)
	sunriseJD := resRise.Time

	if jdUT < sunriseJD {
		resPrev := swisseph.RiseTrans(jdUT-1.0, swisseph.Sun, "", int32(swisseph.FlagSwieph), int32(swisseph.CalcRise), geopos, 0, 0)
		sunriseJD = resPrev.Time
	}
	ephemeris.Mu.Unlock()

	lmtStr := calculateLMT(utcTime, input.Longitude)

	localSunriseOffset := sunriseJD + (input.Timezone / 24.0)
	weekdayIdx := int(math.Floor(localSunriseOffset+1.5)) % 7
	if weekdayIdx < 0 {
		weekdayIdx += 7
	}
	dayLord := WeekdayLords[weekdayIdx]

	ephemeris.Mu.Lock()
	resMoon := swisseph.CalcUT(jdUT, swisseph.Moon, int32(swisseph.FlagSwieph))
	ephemeris.Mu.Unlock()

	tropLon := resMoon.Data[0]
	sidLon := math.Mod(tropLon-ayanamsa+360.0, 360.0)
	totalMinutes := sidLon * 60.0
	nakIdx := int(math.Floor(totalMinutes/800.0)) % 27
	nakshatraLords := []string{"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury"}
	moonNak := ashtakoota.Nakshatras[nakIdx]
	starLord := nakshatraLords[nakIdx%9]

	scanMin := -input.ScanMinusMinutes
	scanMax := input.ScanPlusMinutes
	stepSeconds := 12

	var candidates []domain.BTRCandidate

	for offset := scanMin * 60; offset <= scanMax*60; offset += stepSeconds {
		candTime := utcTime.Add(time.Duration(offset) * time.Second)
		candJD := astronomyTime.UTCToJulianDay(candTime)
		candLocal := candTime.Add(time.Duration(input.Timezone * float64(time.Hour)))

		elapsedMins := (candJD - sunriseJD) * 24.0 * 60.0
		if elapsedMins < 0 {
			elapsedMins += 24.0 * 60.0
		}

		mainT, antarT, isAroha, cycleIdx := CalculateTatwa(dayLord, elapsedMins)

		planetIdx := (weekdayIdx + cycleIdx) % 7
		ninetyMinPlanet := WeekdayLords[planetIdx]

		ascSign, ascDeg := getAscendant(candJD, input.Latitude, input.Longitude, ayanamsa, input.HouseSystem)
		ascType := getSignType(ascSign)

		row := domain.NadiRow{
			NinetyMinutePlanet: ninetyMinPlanet,
			AscendantSignType:  ascType,
			StarLord:           starLord,
			Gender:             input.Gender,
			ExpectedTatwa:      "Unresolved",
			ExpectedAntarTatwa: "Unresolved",
		}

		cand := domain.BTRCandidate{
			CandidateTimeUTC:   candTime.Format(time.RFC3339),
			CandidateTimeLocal: candLocal.Format(time.RFC3339),
			OffsetMinutes:      float64(offset) / 60.0,
			Tatwa: domain.TatwaDetails{
				MainTatwa:     mainT.Name,
				AntarTatwa:    antarT.Name,
				MainDuration:  mainT.Duration,
				AntarDuration: antarT.Duration,
				IsAroha:       isAroha,
				CycleIndex:    cycleIdx,
			},
			AscendantDegree:    ascDeg,
			AscendantSign:      ascSign,
			AscendantSignType:  ascType,
			NinetyMinutePlanet: ninetyMinPlanet,
			MoonNakshatra:      moonNak,
			StarLord:           starLord,
			Gender:             input.Gender,
			NadiRow:            row,
			MatchStatus:        "unresolved",
			MatchExplanation:   "The exact Tatwa Shodhana / Nadi table mapping these variables to an expected Tatwa is unavailable.",
		}

		if len(candidates) == 0 || candidates[len(candidates)-1].Tatwa.MainTatwa != mainT.Name || candidates[len(candidates)-1].Tatwa.AntarTatwa != antarT.Name {
			candidates = append(candidates, cand)
		}
	}

	warnings := []string{
		"The authoritative Nadi Rectification (Tatwa Shodana) Table relies on exact multi-variable lookups that are not legally/openly available in standard libraries.",
		"The required proprietary table must map (90m Planet, AscType, StarLord, Gender) -> Expected Tatwa/AntarTatwa.",
		"Because the authoritative lookup is unavailable, the match status is explicitly set to 'unresolved' rather than guessing or fabricating the table.",
	}

	res := domain.BTRResult{
		Status:       "partially_implemented",
		RecordedTime: utcTime.Format(time.RFC3339),
		AstronomicalContext: domain.BTRAstronomicalContext{
			Sunrise:       astronomyTime.JulianDayToUTC(sunriseJD).Format(time.RFC3339),
			LMT:           lmtStr,
			DayLord:       dayLord,
			MoonNakshatra: moonNak,
			StarLord:      starLord,
		},
		Scan: domain.BTRScanConfig{
			MinusMinutes: input.ScanMinusMinutes,
			PlusMinutes:  input.ScanPlusMinutes,
		},
		Candidates: candidates,
		Result: map[string]interface{}{
			"status": "unresolved",
			"reason": "authoritative_nadi_table_unavailable",
		},
		Warnings: warnings,
	}

	return res, nil
}
