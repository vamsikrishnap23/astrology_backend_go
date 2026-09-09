package houses

import "fmt"

type HoraryDivision struct {
	Number int
	Start  float64
	End    float64
	Sign   int
	Star   int
	Sub    int
}

var HoraryTable []HoraryDivision

func init() {
	dashaYears := []float64{7, 20, 6, 10, 7, 18, 16, 19, 17} // Ketu to Mercury

	HoraryTable = make([]HoraryDivision, 0, 249)
	num := 1
	currentDegree := 0.0

	for nak := 0; nak < 27; nak++ {
		starLordIdx := nak % 9

		subLordIdx := starLordIdx
		for s := 0; s < 9; s++ {
			durationMins := (dashaYears[subLordIdx] / 120.0) * 800.0
			durationDegs := durationMins / 60.0

			endDegree := currentDegree + durationDegs

			signIndex := int(currentDegree / 30.0)
			if currentDegree >= 360.0 {
				signIndex = 11
			}
			nextSignBoundary := float64(signIndex+1) * 30.0

			// A true split only happens if the current degree is strictly LESS than the boundary
			// AND the end degree is strictly GREATER than the boundary.
			if currentDegree < nextSignBoundary-0.0001 && endDegree > nextSignBoundary+0.0001 {
				HoraryTable = append(HoraryTable, HoraryDivision{
					Number: num,
					Start:  currentDegree,
					End:    nextSignBoundary,
					Sign:   signIndex,
					Star:   starLordIdx,
					Sub:    subLordIdx,
				})
				num++

				HoraryTable = append(HoraryTable, HoraryDivision{
					Number: num,
					Start:  nextSignBoundary,
					End:    endDegree,
					Sign:   signIndex + 1,
					Star:   starLordIdx,
					Sub:    subLordIdx,
				})
				num++
			} else {
				HoraryTable = append(HoraryTable, HoraryDivision{
					Number: num,
					Start:  currentDegree,
					End:    endDegree,
					Sign:   signIndex,
					Star:   starLordIdx,
					Sub:    subLordIdx,
				})
				num++
			}

			currentDegree = endDegree
			subLordIdx = (subLordIdx + 1) % 9
		}
	}
}

// GetHoraryAscendant returns the starting degree (Ascendant) for a given KP Horary number (1-249).
func GetHoraryAscendant(number int) (float64, error) {
	if number < 1 || number > 249 {
		return 0, fmt.Errorf("horary number must be between 1 and 249")
	}
	return HoraryTable[number-1].Start, nil
}
