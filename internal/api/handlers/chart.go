package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func ChartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.BirthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// 1. Time Normalization
	utcTime, err := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	if err != nil {
		http.Error(w, "Invalid date/time format", http.StatusBadRequest)
		return
	}

	jd := astronomyTime.UTCToJulianDay(utcTime)

	// 2. Setup Context
	config := domain.CalculationConfig{
		AyanamsaMode: ephemeris.GetAyanamsaMode(input.Ayanamsa),
		HouseCode:    ephemeris.GetHouseSystemCode(input.HouseSystem),
	}

	ctx := domain.CalculationContext{
		Input:       input,
		Config:      config,
		UTCTime:     utcTime,
		JulianDayUT: jd,
	}

	// 3. Calculation
	planetPositions, err := planets.CalculatePlanets(&ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ascendant, mc, houseCusps, err := houses.CalculateHouses(&ctx)
	if err != nil {
		http.Error(w, "Error calculating houses", http.StatusInternalServerError)
		return
	}

	// 4. Build Response
	res := domain.ChartData{
		CalculationTimeUTC:  utcTime.Format(time.RFC3339),
		JulianDay:           jd,
		SelectedAyanamsa:    input.Ayanamsa,
		AyanamsaValue:       ctx.Ayanamsa,
		SelectedHouseSystem: input.HouseSystem,
		Ascendant:           ascendant,
		MC:                  mc,
		Planets:             planetPositions,
		Houses:              houseCusps,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
