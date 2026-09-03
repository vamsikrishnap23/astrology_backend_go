package panchang

import (
	"github.com/vamsi/astrology_backend_go/internal/domain"
	"math"
)

var movingKaranas = []string{
	"Bava", "Balava", "Kaulava", "Taitila", "Gara", "Vanija", "Vishti",
}

func getKaranaDetails(karanaIndex int) (string, string) {
	if karanaIndex == 0 {
		return "Kinstughna", "Fixed"
	}
	if karanaIndex >= 57 {
		fixedNames := []string{"Shakuni", "Chatushpada", "Naga"}
		return fixedNames[karanaIndex-57], "Fixed"
	}
	movingIdx := (karanaIndex - 1) % 7
	return movingKaranas[movingIdx], "Moving"
}

func calculateKarana(jd, sunTrop, moonTrop float64) elementData {
	angle := math.Mod(moonTrop-sunTrop+360.0, 360.0)
	interval := 6.0
	idx := int(math.Floor(angle / interval))
	progress := math.Mod(angle, interval) / interval * 100.0
	start, end := findElementBoundaries(jd, angle, interval, calcTithiAngle)
	return elementData{
		Number:   idx + 1,
		Progress: progress,
		StartJD:  start,
		EndJD:    end,
	}
}

func formatKarana(data elementData, formatter func(float64) string) domain.Karana {
	idx := data.Number - 1
	name, ktype := getKaranaDetails(idx)
	return domain.Karana{
		Number:   data.Number,
		Name:     name,
		Type:     ktype,
		Progress: data.Progress,
		Start:    formatter(data.StartJD),
		End:      formatter(data.EndJD),
	}
}
