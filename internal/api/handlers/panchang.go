package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vamsi/astrology_backend_go/internal/astrology/panchang"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

// PanchangHandler handles panchang calculations.
func PanchangHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.BirthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Parse local time and convert to UTC
	utcT, err := time.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	if err != nil {
		http.Error(w, "Invalid date/time format", http.StatusBadRequest)
		return
	}

	jdUT := time.UTCToJulianDay(utcT)

	// Context mapping
	ctx := &domain.CalculationContext{
		Input:       input,
		JulianDayUT: jdUT,
		Config: domain.CalculationConfig{
			AyanamsaMode: ephemeris.GetAyanamsaMode(input.Ayanamsa),
			HouseCode:    ephemeris.GetHouseSystemCode(input.HouseSystem),
		},
	}

	res, err := panchang.CalculatePanchang(ctx)
	if err != nil {
		http.Error(w, "Error calculating panchang", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
