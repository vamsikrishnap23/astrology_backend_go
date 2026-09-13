package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/btr"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func BTRHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.BTRInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if input.DateOfBirth == "" || input.TimeOfBirth == "" || input.Gender == "" {
		http.Error(w, "date, time, and gender are required", http.StatusBadRequest)
		return
	}

	utcTime, err := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	if err != nil {
		http.Error(w, "Invalid date/time format", http.StatusBadRequest)
		return
	}

	jd := astronomyTime.UTCToJulianDay(utcTime)

	ayanamsa := input.Ayanamsa
	if ayanamsa == "" {
		ayanamsa = "Lahiri"
	}

	config := domain.CalculationConfig{
		AyanamsaMode: ephemeris.GetAyanamsaMode(ayanamsa),
		HouseCode:    ephemeris.GetHouseSystemCode("Placidus"),
	}

	ctx := domain.CalculationContext{
		Input:       input.BirthInput,
		Config:      config,
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	res, err := btr.CalculateBTR(input, &ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
