package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vamsi/astrology_backend_go/internal/astrology/progression"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func ProgressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.ProgressionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if input.ProgressionYear == 0 {
		http.Error(w, "progression_year is required", http.StatusBadRequest)
		return
	}

	// Normal time conversion
	utcTime, err := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	if err != nil {
		http.Error(w, "Invalid date/time format", http.StatusBadRequest)
		return
	}

	jd := astronomyTime.UTCToJulianDay(utcTime)

	config := domain.CalculationConfig{
		AyanamsaMode: ephemeris.GetAyanamsaMode(input.Ayanamsa),
		HouseCode:    ephemeris.GetHouseSystemCode(input.HouseSystem),
	}

	natalCtx := domain.CalculationContext{
		Input:       input.BirthInput,
		Config:      config,
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	// Run Progression
	res, err := progression.CalculateSecondaryProgression(&natalCtx, input.ProgressionYear)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
