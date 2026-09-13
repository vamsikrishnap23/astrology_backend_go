package btr

import (
	"fmt"
	"math"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func getSunrise(jd float64, lat float64, lon float64) float64 {
	ephemeris.Mu.Lock()
	defer ephemeris.Mu.Unlock()
	swisseph.SetEphePath(ephemeris.EphePath)

	geopos := [3]float64{lon, lat, 0}
	epheflag := int32(swisseph.FlagSwieph)
	rsmiRise := int32(swisseph.CalcRise | swisseph.BitHinduRising)
	resRise := swisseph.RiseTrans(jd, swisseph.Sun, "", epheflag, rsmiRise, geopos, 0, 0)
	frac := resRise.Time + 0.5 - math.Floor(resRise.Time+0.5)
	return frac * 24.0
}

func getWeekday(jd float64, timezone float64) int {
	res := swisseph.Revjul(jd, swisseph.GregCal)
	t := time.Date(int(res.Year), time.Month(res.Month), int(res.Day), int(res.Hour), 0, 0, 0, time.UTC)
	loc := time.FixedZone("Local", int(timezone*3600))
	localTime := t.In(loc)
	return int(localTime.Weekday())
}

var planetNames = []string{"Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury", "Ketu", "Venus"}

func getNadiPlanet(row int, ascType string, baseOffset int) string {
	offset := baseOffset
	if ascType == "Fixed" {
		offset += 2
	} else if ascType == "Dual" {
		offset += 4
	}
	idx := ((row - 1) + offset) % 9
	return planetNames[idx]
}

type Tatwa struct {
	Name   string
	Gender string
	Rows   int
}

var tatwas = []Tatwa{
	{"Prithvi", "Male", 2},
	{"Jala", "Female", 4},
	{"Tejo", "Male", 6},
	{"Vayu", "Female", 8},
	{"Akash", "Male", 10},
}

func getTatwaInfo(row int, weekday int) (string, string) {
	cycleRow := (row-121+480)%480 + 1
	halfCycleIdx := (cycleRow - 1) / 30
	isReverse := (halfCycleIdx%2 != 0)
	rowInHalfCycle := (cycleRow-1)%30 + 1

	startIdx := 0
	switch weekday {
	case 3:
		startIdx = 0
	case 1, 5:
		startIdx = 1
	case 0, 2:
		startIdx = 2
	case 6:
		startIdx = 3
	case 4:
		startIdx = 4
	}

	var seq []Tatwa
	idx := startIdx
	for i := 0; i < 5; i++ {
		seq = append(seq, tatwas[idx])
		idx = (idx + 1) % 5
	}

	if isReverse {
		for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
			seq[i], seq[j] = seq[j], seq[i]
		}
	}

	current := 0
	for _, t := range seq {
		if rowInHalfCycle <= current+t.Rows {
			return t.Name, t.Gender
		}
		current += t.Rows
	}

	return "Unknown", "Unknown"
}

func getAscendantType(ascLon float64) string {
	signIdx := int(math.Floor(ascLon / 30.0))
	modality := signIdx % 3
	if modality == 0 {
		return "Movable"
	} else if modality == 1 {
		return "Fixed"
	}
	return "Dual"
}

func formatTatwa(gender, tatwa string) string {
	prefix := "Fem"
	if gender == "Male" {
		prefix = "Male"
	}
	return prefix + tatwa
}

func getWeekdayLordOffset(weekday int) int {
	// Sun=0, Mon=1, Tue=2, Wed=6(Merc), Thu=4(Jup), Fri=8(Ven), Sat=5(Sat)
	offsets := []int{0, 1, 2, 6, 4, 8, 5}
	if weekday >= 0 && weekday < 7 {
		return offsets[weekday]
	}
	return 0
}

func evaluateRow(row int, weekday int, ascType string, userGender string, actualStarLord string) (bool, string, string, string) {
	calcTatwa, calcGender := getTatwaInfo(row, weekday)
	calcPlanet := getNadiPlanet(row, ascType, getWeekdayLordOffset(weekday))
	match := (calcGender == userGender) && (calcPlanet == actualStarLord)
	return match, calcTatwa, calcGender, calcPlanet
}

func getAntarTatwa(row int, weekday int) string {
	tatwas := []string{"Prithvi", "Jala", "Tejo", "Vayu", "Akash"}
	durations := []float64{6.0, 12.0, 18.0, 24.0, 30.0}

	// Sync with getTatwaInfo's global row shift
	cycleRow := (row-121+480)%480 + 1
	halfCycleIdx := (cycleRow - 1) / 30
	isReverse := (halfCycleIdx%2 != 0) // Odd index is Avaroha (Reverse)
	rowInHalfCycle := (cycleRow-1)%30 + 1

	startIdx := 0
	switch weekday {
	case 3:
		startIdx = 0
	case 1, 5:
		startIdx = 1
	case 0, 2:
		startIdx = 2
	case 6:
		startIdx = 3
	case 4:
		startIdx = 4
	}

	// Build the sequence of Maha Tatwas for this cycle
	var seq []int
	idx := startIdx
	for i := 0; i < 5; i++ {
		seq = append(seq, idx)
		idx = (idx + 1) % 5
	}

	if isReverse {
		for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
			seq[i], seq[j] = seq[j], seq[i]
		}
	}

	mahaStartRow := 1
	var mahaDuration float64
	var mahaTatwaIdx int

	currentTotal := 0
	for _, tIdx := range seq {
		d := durations[tIdx]
		rowsForTatwa := int(d / 3.0)
		if rowInHalfCycle <= currentTotal+rowsForTatwa {
			mahaDuration = d
			mahaTatwaIdx = tIdx
			mahaStartRow = currentTotal + 1
			break
		}
		currentTotal += rowsForTatwa
	}

	minutesElapsed := float64(rowInHalfCycle-mahaStartRow)*3.0 + 1.5

	currentAntarIdx := mahaTatwaIdx
	accumulatedAntarMins := 0.0

	for i := 0; i < 5; i++ {
		antarDuration := (mahaDuration * durations[currentAntarIdx]) / 90.0
		accumulatedAntarMins += antarDuration

		if minutesElapsed <= accumulatedAntarMins {
			return tatwas[currentAntarIdx]
		}

		if !isReverse {
			currentAntarIdx = (currentAntarIdx + 1) % 5
		} else {
			currentAntarIdx = (currentAntarIdx - 1 + 5) % 5
		}
	}

	return tatwas[mahaTatwaIdx]
}

func CalculateBTR(input domain.BTRInput, ctx *domain.CalculationContext) (domain.BTRResult, error) {
	// Parse Local Time
	t, err := time.Parse("15:04:05", input.TimeOfBirth)
	if err != nil {
		t, err = time.Parse("15:04", input.TimeOfBirth)
		if err != nil {
			return domain.BTRResult{}, err
		}
	}
	rawSeconds := t.Hour()*3600 + t.Minute()*60 + t.Second()

	stdMeridian := input.Timezone * 15.0
	// Correct astronomical LMT: (Longitude - StandardMeridian) * 4 minutes (240s)
	// Example: Guntur (80.15) is West of IST (82.5). 80.15 - 82.5 = -2.35 * 240 = -564 seconds.
	lmtCorrectionSeconds := (input.Longitude - stdMeridian) * 240.0
	lmtSeconds := float64(rawSeconds) + lmtCorrectionSeconds

	// Determine the Julian Day for 00:00:00 of the input date
	localStart, _ := time.Parse("2006-01-02 15:04:05", input.DateOfBirth+" 00:00:00")
	utcStart := localStart.UTC()
	startOfDayJD := swisseph.Julday(int32(utcStart.Year()), int32(utcStart.Month()), int32(utcStart.Day()), float64(utcStart.Hour())+float64(utcStart.Minute())/60.0+float64(utcStart.Second())/3600.0, swisseph.GregCal)

	sunriseUTC := getSunrise(startOfDayJD, input.Latitude, input.Longitude)
	sunriseLocalSecs := (sunriseUTC + input.Timezone) * 3600.0
	if sunriseLocalSecs > 86400 {
		sunriseLocalSecs -= 86400
	}
	if input.SunriseOverride != "" {
		if tSun, err := time.Parse("15:04:05", input.SunriseOverride); err == nil {
			sunriseLocalSecs = float64(tSun.Hour()*3600 + tSun.Minute()*60 + tSun.Second())
		}
	}
	sunriseDiffSecs := sunriseLocalSecs - (6.0 * 3600.0)

	finalTimeSecs := lmtSeconds - sunriseDiffSecs

	// Handle wrapping for row calculation
	normalizedSecs := finalTimeSecs
	for normalizedSecs < 0 {
		normalizedSecs += 86400
	}
	for normalizedSecs >= 86400 {
		normalizedSecs -= 86400
	}

	finalMinutes := normalizedSecs / 60.0
	baseRow := int(math.Floor(finalMinutes/3.0)) + 1
	if baseRow > 480 {
		baseRow = 480
	}

	weekday := getWeekday(ctx.JulianDayUT, input.Timezone)

	planetsList, _ := planets.CalculatePlanets(ctx)
	ascendantLon, _, _, _ := houses.CalculateHouses(ctx)

	ascType := getAscendantType(ascendantLon)
	if input.AscendantOverride != "" {
		m := map[string]string{
			"Mesha": "Movable", "Vrishabha": "Fixed", "Mithuna": "Dual",
			"Kataka": "Movable", "Simha": "Fixed", "Kanya": "Dual",
			"Thula": "Movable", "Vrischika": "Fixed", "Dhanus": "Dual",
			"Makara": "Movable", "Kumbha": "Fixed", "Meena": "Dual",
			"Aries": "Movable", "Taurus": "Fixed", "Gemini": "Dual",
			"Cancer": "Movable", "Leo": "Fixed", "Virgo": "Dual",
			"Libra": "Movable", "Scorpio": "Fixed", "Sagittarius": "Dual",
			"Capricorn": "Movable", "Aquarius": "Fixed", "Pisces": "Dual",
		}
		if t, ok := m[input.AscendantOverride]; ok {
			ascType = t
		}
	}
	if input.SignTypeOverride != "" && input.SignTypeOverride != "Auto-detect" {
		ascType = input.SignTypeOverride
	}

	actualStarLord := "Moon"
	for _, p := range planetsList {
		if p.Planet == "Moon" {
			actualStarLord = p.NakshatraLord
			break
		}
	}
	if input.StarLordOverride != "" {
		actualStarLord = input.StarLordOverride
	}

	baseMatch, baseTatwa, baseGender, basePlanet := evaluateRow(baseRow, weekday, ascType, input.Gender, actualStarLord)

	status := "Failed"
	if baseMatch {
		status = "Verified"
	}

	// Helper to format seconds to HH:MM:SS
	formatSecs := func(secs float64) string {
		s := int(secs)
		for s < 0 {
			s += 86400
		}
		for s >= 86400 {
			s -= 86400
		}
		return fmt.Sprintf("%02d:%02d:%02d", s/3600, (s%3600)/60, s%60)
	}

	weekdays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}

	// Calculate base Vinod Planet
	baseVinodPlanet := getNadiPlanet(baseRow, ascType, 0)

	// Calculate Ascendant Name for output
	ascSignIdx := int(math.Floor(ascendantLon / 30.0))
	signNames := []string{"Mesha", "Vrishabha", "Mithuna", "Kataka", "Simha", "Kanya", "Thula", "Vrischika", "Dhanus", "Makara", "Kumbha", "Meena"}
	ascSignName := signNames[ascSignIdx]
	if input.AscendantOverride != "" {
		ascSignName = input.AscendantOverride
	}

	res := domain.BTRResult{
		InputTimeStatus: status,
		InputAnalysis: domain.BTRAnalysis{
			Weekday:               weekdays[weekday],
			Sunrise:               formatSecs(sunriseLocalSecs),
			Lmt:                   formatSecs(lmtSeconds),
			LmtSunrise:            formatSecs(finalTimeSecs),
			StarMatch:             basePlanet == actualStarLord,
			ActualStarLord:        actualStarLord,
			NadiRow:               baseRow,
			CalculatedTatwa:       baseTatwa,
			CalculatedGender:      baseGender,
			GenderMatch:           baseGender == input.Gender,
			CalculatedPlanet:      basePlanet,
			CalculatedPlanetVinod: baseVinodPlanet,
			AscendantSign:         ascSignName,
			AscendantDegree:       math.Round((ascendantLon-float64(ascSignIdx*30))*10) / 10,
		},
		SuggestedRectifications: []domain.BTRCandidate{},
	}

	// Generate candidates strictly within the requested scan window
	scanMinus := input.ScanMinusMinutes
	if scanMinus <= 0 {
		scanMinus = 10
	}
	scanPlus := input.ScanPlusMinutes
	if scanPlus <= 0 {
		scanPlus = 5
	}

	type scoredCandidate struct {
		cand  domain.BTRCandidate
		score int
	}
	var scoredList []scoredCandidate

	// Check a wide enough range of rows to cover the window
	scanRangeRows := (scanMinus+scanPlus)/3 + 2

	for offset := -scanRangeRows; offset <= scanRangeRows; offset++ {
		checkRow := baseRow + offset
		wrappedRow := checkRow
		if wrappedRow < 1 {
			wrappedRow += 480
		} else if wrappedRow > 480 {
			wrappedRow -= 480
		}

		// Calculate exact clock start time of this row
		rowSecsSinceSunrise := float64((wrappedRow - 1) * 180)
		rowLmtSecs := rowSecsSinceSunrise + sunriseDiffSecs
		rowIstSecs := rowLmtSecs - lmtCorrectionSeconds

		for rowIstSecs < 0 {
			rowIstSecs += 86400
		}
		for rowIstSecs >= 86400 {
			rowIstSecs -= 86400
		}

		diffSecs := rowIstSecs - float64(rawSeconds)
		if diffSecs < -43200 {
			diffSecs += 86400
		}
		if diffSecs > 43200 {
			diffSecs -= 86400
		}

		if diffSecs >= -float64(scanMinus*60) && diffSecs <= float64(scanPlus*60) {
			_, checkTatwa, checkGender, checkPlanet := evaluateRow(wrappedRow, weekday, ascType, input.Gender, actualStarLord)

			score := 0
			if checkGender == input.Gender {
				score += 5
			}
			if checkPlanet == actualStarLord {
				score += 5
			}

			sh := int(rowIstSecs) / 3600
			sm := (int(rowIstSecs) % 3600) / 60
			ss := int(rowIstSecs) % 60
			suggStr := fmt.Sprintf("%02d:%02d:%02d", sh, sm, ss)

			diffMins := int(math.Round(diffSecs / 60.0))

			t1Mins := wrappedRow * 3
			t1Str := fmt.Sprintf("%d%02d", t1Mins/60, t1Mins%60)
			if t1Mins/60 == 0 {
				t1Str = fmt.Sprintf("%d", t1Mins)
			}

			genStr := "Male"
			if checkGender == "Female" {
				genStr = "Female"
			}

			_, _, _, vinodPlanet := evaluateRow(wrappedRow, weekday, ascType, input.Gender, actualStarLord)
			vinodPlanet = getNadiPlanet(wrappedRow, ascType, 0) // Override with Vinod planet logic

			scoredList = append(scoredList, scoredCandidate{
				cand: domain.BTRCandidate{
					T1:          t1Str,
					T2:          suggStr,
					Score:       score,
					NadiRow:     wrappedRow,
					Tatwa:       checkTatwa,
					Antar:       getAntarTatwa(wrappedRow, weekday),
					Gender:      genStr,
					Planet90:    checkPlanet,
					PlanetVinod: vinodPlanet,

					SuggestedTime:     suggStr,
					DifferenceMinutes: diffMins,
				},
				score: score,
			})
		}
	}

	// Sort candidates by Score (descending), then by proximity to target
	for i := 0; i < len(scoredList)-1; i++ {
		for j := 0; j < len(scoredList)-i-1; j++ {
			if scoredList[j].score < scoredList[j+1].score {
				scoredList[j], scoredList[j+1] = scoredList[j+1], scoredList[j]
			} else if scoredList[j].score == scoredList[j+1].score {
				if math.Abs(float64(scoredList[j].cand.DifferenceMinutes)) > math.Abs(float64(scoredList[j+1].cand.DifferenceMinutes)) {
					scoredList[j], scoredList[j+1] = scoredList[j+1], scoredList[j]
				}
			}
		}
	}

	for i := 0; i < len(scoredList); i++ {
		scoredList[i].cand.Rank = i + 1
		res.SuggestedRectifications = append(res.SuggestedRectifications, scoredList[i].cand)
	}

	if input.ReturnFullTable {
		for i := 1; i <= 480; i++ {
			t1Mins := i * 3
			t1Str := fmt.Sprintf("%d%02d", t1Mins/60, t1Mins%60)
			if t1Mins/60 == 0 {
				t1Str = fmt.Sprintf("%d", t1Mins)
			}

			normalizedSecs := float64(t1Mins * 60)
			lmtSecs := normalizedSecs + sunriseDiffSecs
			istSecs := lmtSecs - lmtCorrectionSeconds

			for istSecs < 0 {
				istSecs += 86400
			}
			for istSecs >= 86400 {
				istSecs -= 86400
			}

			sh := int(istSecs) / 3600
			sm := (int(istSecs) % 3600) / 60
			t2 := time.Date(2000, 1, 1, sh, sm, 0, 0, time.UTC).Format("15:04")

			tWed, gWed := getTatwaInfo(i, 3)
			tMonFri, gMonFri := getTatwaInfo(i, 1)
			tSunTues, gSunTues := getTatwaInfo(i, 0)
			tSat, gSat := getTatwaInfo(i, 6)
			tThur, gThur := getTatwaInfo(i, 4)

			pMovable := getNadiPlanet(i, "Movable", getWeekdayLordOffset(weekday))
			pFixed := getNadiPlanet(i, "Fixed", getWeekdayLordOffset(weekday))
			pDual := getNadiPlanet(i, "Dual", getWeekdayLordOffset(weekday))

			vMovable := getNadiPlanet(i, "Movable", 0)
			vFixed := getNadiPlanet(i, "Fixed", 0)
			vDual := getNadiPlanet(i, "Dual", 0)

			res.FullTable = append(res.FullTable, domain.BTRTableRow{
				No:           i,
				T1:           t1Str,
				T2:           t2,
				Wed:          formatTatwa(gWed, tWed),
				MonFri:       formatTatwa(gMonFri, tMonFri),
				SunTues:      formatTatwa(gSunTues, tSunTues),
				Sat:          formatTatwa(gSat, tSat),
				Thur:         formatTatwa(gThur, tThur),
				Movable:      pMovable,
				Fixed:        pFixed,
				Dual:         pDual,
				VinodMovable: vMovable,
				VinodFixed:   vFixed,
				VinodDual:    vDual,
			})
		}
	}

	return res, nil
}
