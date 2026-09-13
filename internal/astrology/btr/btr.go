package btr

import (
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
	rsflag := int32(swisseph.FlagSwieph | swisseph.BitDiscCenter | swisseph.BitNoRefraction)
	resRise := swisseph.RiseTrans(jd, swisseph.Sun, "", rsflag, int32(swisseph.CalcRise), geopos, 0, 0)
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

func getNadiPlanet(row int, ascType string) string {
	offset := 0
	if ascType == "Fixed" {
		offset = 2
	} else if ascType == "Dual" {
		offset = 4
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

func evaluateRow(row int, weekday int, ascType string, userGender string, actualStarLord string) (bool, string, string, string) {
	calcTatwa, calcGender := getTatwaInfo(row, weekday)
	calcPlanet := getNadiPlanet(row, ascType)
	match := (calcGender == userGender) && (calcPlanet == actualStarLord)
	return match, calcTatwa, calcGender, calcPlanet
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
	lmtCorrectionSeconds := (stdMeridian - input.Longitude) * 240.0
	lmtSeconds := float64(rawSeconds) + lmtCorrectionSeconds

	sunriseUTC := getSunrise(ctx.JulianDayUT, input.Latitude, input.Longitude)
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
	_, ascendantLon, _, _ := houses.CalculateHouses(ctx)

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

	res := domain.BTRResult{
		InputTimeStatus: status,
		InputAnalysis: domain.BTRAnalysis{
			GenderMatch:      baseGender == input.Gender,
			StarMatch:        basePlanet == actualStarLord,
			CalculatedTatwa:  baseTatwa,
			CalculatedGender: baseGender,
			CalculatedPlanet: basePlanet,
			ActualStarLord:   actualStarLord,
		},
		SuggestedRectifications: []domain.BTRCandidate{},
	}

	if !baseMatch {
		// Scan nearby rows for candidates
		var candidates []domain.BTRCandidate
		scanRange := 40 // Check +/- 40 rows (2 hours each way)

		for offset := -scanRange; offset <= scanRange; offset++ {
			if offset == 0 {
				continue
			}

			checkRow := baseRow + offset
			wrappedRow := checkRow
			if wrappedRow < 1 {
				wrappedRow += 480
			} else if wrappedRow > 480 {
				wrappedRow -= 480
			}

			match, checkTatwa, _, _ := evaluateRow(wrappedRow, weekday, ascType, input.Gender, actualStarLord)

			if match {
				diffMinutes := offset * 3

				// Calculate suggested time string based on offset
				suggestedSecs := finalTimeSecs + float64(offset*3*60)

				// Un-normalize back to clock time
				suggestedLmtSecs := suggestedSecs + sunriseDiffSecs
				suggestedClockSecs := suggestedLmtSecs - lmtCorrectionSeconds

				for suggestedClockSecs < 0 {
					suggestedClockSecs += 86400
				}
				for suggestedClockSecs >= 86400 {
					suggestedClockSecs -= 86400
				}

				sh := int(suggestedClockSecs) / 3600
				sm := (int(suggestedClockSecs) % 3600) / 60
				ss := int(suggestedClockSecs) % 60

				candidates = append(candidates, domain.BTRCandidate{
					SuggestedTime:     time.Date(2000, 1, 1, sh, sm, ss, 0, time.UTC).Format("15:04:05"),
					DifferenceMinutes: diffMinutes,
					Tatwa:             checkTatwa,
					NadiRow:           wrappedRow,
				})
			}
		}

		// Sort by absolute difference
		for i := 0; i < len(candidates)-1; i++ {
			for j := 0; j < len(candidates)-i-1; j++ {
				if math.Abs(float64(candidates[j].DifferenceMinutes)) > math.Abs(float64(candidates[j+1].DifferenceMinutes)) {
					candidates[j], candidates[j+1] = candidates[j+1], candidates[j]
				}
			}
		}

		// Take top 5 and assign ranks
		limit := 5
		if len(candidates) < limit {
			limit = len(candidates)
		}

		for i := 0; i < limit; i++ {
			candidates[i].Rank = i + 1
			res.SuggestedRectifications = append(res.SuggestedRectifications, candidates[i])
		}
	}

	if input.ReturnFullTable {
		for i := 1; i <= 480; i++ {
			t1 := i * 3

			normalizedSecs := float64(t1 * 60)
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

			pMovable := getNadiPlanet(i, "Movable")
			pFixed := getNadiPlanet(i, "Fixed")
			pDual := getNadiPlanet(i, "Dual")

			res.FullTable = append(res.FullTable, domain.BTRTableRow{
				No:           i,
				T1:           t1,
				T2:           t2,
				Wed:          formatTatwa(gWed, tWed),
				MonFri:       formatTatwa(gMonFri, tMonFri),
				SunTues:      formatTatwa(gSunTues, tSunTues),
				Sat:          formatTatwa(gSat, tSat),
				Thur:         formatTatwa(gThur, tThur),
				Movable:      pMovable,
				Fixed:        pFixed,
				Dual:         pDual,
				VinodMovable: pMovable,
				VinodFixed:   pFixed,
				VinodDual:    pDual,
			})
		}
	}

	return res, nil
}
