package panchang

import (
	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"math"
)

var nakshatraNames = []string{
	"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
	"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni",
	"Uttara Phalguni", "Hasta", "Chitra", "Swati", "Vishakha", "Anuradha",
	"Jyeshtha", "Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana",
	"Dhanishta", "Shatabhisha", "Purva Bhadrapada", "Uttara Bhadrapada",
	"Revati",
}

var nakshatraRulers = []string{
	"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
	"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
	"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury",
}

func calcNakshatraAngle(jd float64) float64 {
	sflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed | swisseph.FlagSidereal)
	moon := swisseph.CalcUT(jd, swisseph.Moon, sflag).Data[0]
	return math.Mod(moon+360.0, 360.0)
}

func calculateNakshatra(jd, moonSid float64) elementData {
	interval := 13.0 + 1.0/3.0
	idx := int(math.Floor(moonSid / interval))
	progress := math.Mod(moonSid, interval) / interval * 100.0
	start, end := findElementBoundaries(jd, moonSid, interval, calcNakshatraAngle)
	return elementData{
		Number:   idx + 1,
		Progress: progress,
		StartJD:  start,
		EndJD:    end,
	}
}

func formatNakshatra(data elementData, formatter func(float64) string) domain.Nakshatra {
	padaProgress := data.Progress / 25.0
	pada := int(math.Floor(padaProgress)) + 1
	idx := data.Number - 1
	return domain.Nakshatra{
		Number:   data.Number,
		Name:     nakshatraNames[idx],
		Pada:     pada,
		Progress: data.Progress,
		Start:    formatter(data.StartJD),
		End:      formatter(data.EndJD),
		Ruler:    nakshatraRulers[idx],
	}
}
