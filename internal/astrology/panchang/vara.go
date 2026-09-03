package panchang

import (
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"time"
)

var weekdayNames = []string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
}

var weekdayRulers = []string{
	"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn",
}

func calculateVara(dateStr string, timezone float64) domain.Vara {
	t, _ := time.Parse("2006-01-02", dateStr)
	wd := int(t.Weekday())
	return domain.Vara{
		Number: wd,
		Name:   weekdayNames[wd],
		Ruler:  weekdayRulers[wd],
	}
}
