package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/transit"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func TransitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.TransitInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if input.TransitDate == "" || input.TransitTime == "" {
		http.Error(w, "transit_date and transit_time are required", http.StatusBadRequest)
		return
	}

	// Calculate Transit time in UTC based on natal timezone (assuming transit is happening locally there)
	utcTime, err := astronomyTime.ParseLocalToUTC(input.TransitDate, input.TransitTime, input.Timezone)
	if err != nil {
		http.Error(w, "Invalid transit date/time format", http.StatusBadRequest)
		return
	}

	jd := astronomyTime.UTCToJulianDay(utcTime)

	config := domain.CalculationConfig{
		AyanamsaMode: ephemeris.GetAyanamsaMode(input.Ayanamsa),
		HouseCode:    ephemeris.GetHouseSystemCode(input.HouseSystem),
	}

	transitCtx := domain.CalculationContext{
		Input:       input.BirthInput,
		Config:      config,
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	// Run Transit
	res, err := transit.CalculateTransitChart(&transitCtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
