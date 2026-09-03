package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vamsi/astrology_backend_go/internal/astrology/transit"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func UpcomingTransitsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.UpcomingTransitsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Default to birth date if search dates aren't provided
	startDate := input.SearchStartDate
	if startDate == "" {
		startDate = input.DateOfBirth
	}
	startTime := input.SearchStartTime
	if startTime == "" {
		startTime = input.TimeOfBirth
	}

	utcTime, err := astronomyTime.ParseLocalToUTC(startDate, startTime, input.Timezone)
	if err != nil {
		http.Error(w, "Invalid search date/time format", http.StatusBadRequest)
		return
	}

	jd := astronomyTime.UTCToJulianDay(utcTime)
	ayanamsaMode := ephemeris.GetAyanamsaMode(input.Ayanamsa)

	results := transit.CalculateUpcomingTransits(jd, ayanamsaMode)

	res := domain.UpcomingTransitsResult{
		SearchStartDateUTC: utcTime.Format("2006-01-02T15:04:05Z07:00"), // standard RFC3339
		Transits:           results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
