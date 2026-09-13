package retrogrades

import (
	"fmt"
	"math"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astrotime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

var planetMap = map[string]int{
	"Mercury": swisseph.Mercury,
	"Venus":   swisseph.Venus,
	"Mars":    swisseph.Mars,
	"Jupiter": swisseph.Jupiter,
	"Saturn":  swisseph.Saturn,
	"Uranus":  swisseph.Uranus,
	"Neptune": swisseph.Neptune,
	"Pluto":   swisseph.Pluto,
}

func getPlanetSpeed(jd float64, seID int) float64 {
	iflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed)
	res := swisseph.CalcUT(jd, int32(seID), iflag)
	return res.Data[3] // speed
}

func findStation(startJD, endJD float64, seID int) float64 {
	// Bisection method to find where speed crosses 0
	low := startJD
	high := endJD

	// Verify signs are opposite
	sLow := getPlanetSpeed(low, seID)
	sHigh := getPlanetSpeed(high, seID)
	if sLow*sHigh > 0 {
		return (low + high) / 2.0 // Fallback if no crossing
	}

	for i := 0; i < 50; i++ {
		mid := (low + high) / 2.0
		sMid := getPlanetSpeed(mid, seID)

		if math.Abs(sMid) < 0.000001 {
			return mid
		}

		if (sLow < 0 && sMid < 0) || (sLow > 0 && sMid > 0) {
			low = mid
			sLow = sMid
		} else {
			high = mid
			sHigh = sMid
		}
	}
	return (low + high) / 2.0
}

func formatTime(jd float64, tz float64) string {
	res := swisseph.Revjul(jd, swisseph.GregCal)
	y, m, d, h := res.Year, res.Month, res.Day, res.Hour
	hr := int(h)
	min := int((h - float64(hr)) * 60)
	sec := int(math.Round((h - float64(hr) - float64(min)/60) * 3600))
	if sec >= 60 {
		sec -= 60
		min++
	}
	if min >= 60 {
		min -= 60
		hr++
	}
	utc := time.Date(int(y), time.Month(m), int(d), hr, min, sec, 0, time.UTC)
	tzDuration := time.Duration(tz * float64(time.Hour))
	loc := time.FixedZone("Local", int(tzDuration.Seconds()))
	return utc.In(loc).Format("2006-01-02T15:04:05-07:00")
}

func CalculateRetrogrades(input domain.RetrogradesInput) (domain.RetrogradesResult, error) {
	ephemeris.Mu.Lock()
	defer ephemeris.Mu.Unlock()
	swisseph.SetEphePath(ephemeris.EphePath)

	utcStart, err := astrotime.ParseLocalToUTC(input.StartDate, "00:00:00", input.Timezone)
	if err != nil {
		return domain.RetrogradesResult{}, fmt.Errorf("invalid start date: %v", err)
	}
	utcEnd, err := astrotime.ParseLocalToUTC(input.EndDate, "23:59:59", input.Timezone)
	if err != nil {
		return domain.RetrogradesResult{}, fmt.Errorf("invalid end date: %v", err)
	}

	startJD := astrotime.UTCToJulianDay(utcStart)
	endJD := astrotime.UTCToJulianDay(utcEnd)

	result := domain.RetrogradesResult{
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		Timezone:    input.Timezone,
		Retrogrades: make(map[string][]domain.RetrogradePhase),
	}

	planetOrder := []string{"Mercury", "Venus", "Mars", "Jupiter", "Saturn"}

	for _, pName := range planetOrder {
		seID := planetMap[pName]
		var phases []domain.RetrogradePhase

		var currentPhaseStart float64 = 0

		// Check if already retrograde at startJD
		initialSpeed := getPlanetSpeed(startJD, seID)
		if initialSpeed < 0 {
			// Sweep backwards to find when it went retrograde
			testJD := startJD
			for {
				if getPlanetSpeed(testJD, seID) > 0 {
					currentPhaseStart = findStation(testJD, testJD+1.0, seID)
					break
				}
				testJD -= 1.0
			}
		}

		prevSpeed := initialSpeed
		for jd := startJD + 1.0; jd <= endJD+1.0; jd += 1.0 { // +1.0 ensures we check the end boundary
			currSpeed := getPlanetSpeed(jd, seID)

			if prevSpeed > 0 && currSpeed < 0 {
				// Went retrograde
				currentPhaseStart = findStation(jd-1.0, jd, seID)
			} else if prevSpeed < 0 && currSpeed > 0 {
				// Went direct
				phaseEnd := findStation(jd-1.0, jd, seID)
				if currentPhaseStart > 0 {
					phases = append(phases, domain.RetrogradePhase{
						Start: formatTime(currentPhaseStart, input.Timezone),
						End:   formatTime(phaseEnd, input.Timezone),
					})
					currentPhaseStart = 0
				}
			}
			prevSpeed = currSpeed
		}

		// If still retrograde at endJD, sweep forwards to find when it goes direct
		if currentPhaseStart > 0 {
			testJD := endJD
			for {
				if getPlanetSpeed(testJD, seID) > 0 {
					phaseEnd := findStation(testJD-1.0, testJD, seID)
					phases = append(phases, domain.RetrogradePhase{
						Start: formatTime(currentPhaseStart, input.Timezone),
						End:   formatTime(phaseEnd, input.Timezone),
					})
					break
				}
				testJD += 1.0
			}
		}

		result.Retrogrades[pName] = phases
	}

	return result, nil
}
