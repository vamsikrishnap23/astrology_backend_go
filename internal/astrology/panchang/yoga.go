package panchang

import (
	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"math"
)

var yogaNames = []string{
	"Vishkambha", "Priti", "Ayushman", "Saubhagya", "Shobhana", "Atiganda",
	"Sukarma", "Dhriti", "Shula", "Ganda", "Vriddhi", "Dhruva", "Vyaghata",
	"Harshana", "Vajra", "Siddhi", "Vyatipata", "Variyan", "Parigha",
	"Shiva", "Siddha", "Sadhya", "Shubha", "Shukla", "Brahma", "Indra",
	"Vaidhriti",
}

func calcYogaAngle(jd float64) float64 {
	sflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed | swisseph.FlagSidereal)
	sun := swisseph.CalcUT(jd, swisseph.Sun, sflag).Data[0]
	moon := swisseph.CalcUT(jd, swisseph.Moon, sflag).Data[0]
	return math.Mod(sun+moon+360.0, 360.0)
}

func calculateYoga(jd, sunSid, moonSid float64) elementData {
	angle := math.Mod(sunSid+moonSid+360.0, 360.0)
	interval := 13.0 + 1.0/3.0
	idx := int(math.Floor(angle / interval))
	progress := math.Mod(angle, interval) / interval * 100.0
	start, end := findElementBoundaries(jd, angle, interval, calcYogaAngle)
	return elementData{
		Number:   idx + 1,
		Progress: progress,
		StartJD:  start,
		EndJD:    end,
	}
}

func formatYoga(data elementData, formatter func(float64) string) domain.Yoga {
	return domain.Yoga{
		Number:   data.Number,
		Name:     yogaNames[data.Number-1],
		Progress: data.Progress,
		Start:    formatter(data.StartJD),
		End:      formatter(data.EndJD),
	}
}
