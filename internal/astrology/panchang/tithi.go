package panchang

import (
	"github.com/tejzpr/go-swisseph"
	"github.com/vamsi/astrology_backend_go/internal/domain"
	"math"
)

var tithiNames = []string{
	"Pratipada", "Dwitiya", "Tritiya", "Chaturthi", "Panchami",
	"Shashthi", "Saptami", "Ashtami", "Navami", "Dashami",
	"Ekadashi", "Dwadashi", "Trayodashi", "Chaturdashi", "Purnima",
	"Pratipada", "Dwitiya", "Tritiya", "Chaturthi", "Panchami",
	"Shashthi", "Saptami", "Ashtami", "Navami", "Dashami",
	"Ekadashi", "Dwadashi", "Trayodashi", "Chaturdashi", "Amavasya",
}

func calcTithiAngle(jd float64) float64 {
	tflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed)
	sun := swisseph.CalcUT(jd, swisseph.Sun, tflag).Data[0]
	moon := swisseph.CalcUT(jd, swisseph.Moon, tflag).Data[0]
	return math.Mod(moon-sun+360.0, 360.0)
}

func calculateTithi(jd, sunTrop, moonTrop float64) elementData {
	angle := math.Mod(moonTrop-sunTrop+360.0, 360.0)
	interval := 12.0
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

func formatTithi(data elementData, formatter func(float64) string) domain.Tithi {
	paksha := "Shukla"
	if data.Number > 15 {
		paksha = "Krishna"
	}
	return domain.Tithi{
		Number:   data.Number,
		Name:     tithiNames[data.Number-1],
		Paksha:   paksha,
		Progress: data.Progress,
		Start:    formatter(data.StartJD),
		End:      formatter(data.EndJD),
	}
}
