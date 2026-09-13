package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/transit"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
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

	// Calculate Natal Context
	natalUTC, err := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	if err != nil {
		http.Error(w, "Invalid natal date/time format", http.StatusBadRequest)
		return
	}
	natalJD := astronomyTime.UTCToJulianDay(natalUTC)

	config := domain.CalculationConfig{
		AyanamsaMode: ephemeris.GetAyanamsaMode(input.Ayanamsa),
		HouseCode:    ephemeris.GetHouseSystemCode(input.HouseSystem),
	}

	natalCtx := domain.CalculationContext{
		Input:       input.BirthInput,
		Config:      config,
		UTCTime:     natalUTC,
		JulianDayUT: natalJD,
	}

	natalPlanets, err := planets.CalculatePlanets(&natalCtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var natalMoonLon float64
	for _, p := range natalPlanets {
		if p.Planet == "Moon" {
			natalMoonLon = p.SiderealLongitude
			break
		}
	}

	// Calculate Transit time in UTC
	transitUTC, err := astronomyTime.ParseLocalToUTC(input.TransitDate, input.TransitTime, input.Timezone)
	if err != nil {
		http.Error(w, "Invalid transit date/time format", http.StatusBadRequest)
		return
	}

	transitJD := astronomyTime.UTCToJulianDay(transitUTC)

	transitCtx := domain.CalculationContext{
		Input:       input.BirthInput,
		Config:      config,
		UTCTime:     transitUTC,
		JulianDayUT: transitJD,
	}

	// Run Transit
	res, err := transit.CalculateTransitChart(&transitCtx, natalMoonLon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func RasiTransitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.RasiTransitInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if input.TransitDate == "" || input.Rasi == "" {
		http.Error(w, "rasi and transit_date are required", http.StatusBadRequest)
		return
	}

	// Determine faux Moon Longitude from Rasi
	rasiNames := []string{
		"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
		"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}

	rasiIdx := -1
	for i, name := range rasiNames {
		if name == input.Rasi {
			rasiIdx = i
			break
		}
	}
	if rasiIdx == -1 {
		http.Error(w, "Invalid rasi name", http.StatusBadRequest)
		return
	}
	fauxMoonLon := float64(rasiIdx)*30.0 + 15.0

	// Calculate Transit time in UTC
	transitTime := input.TransitTime
	if transitTime == "" {
		transitTime = "12:00:00" // Default to noon if not provided
	}

	transitUTC, err := astronomyTime.ParseLocalToUTC(input.TransitDate, transitTime, input.Timezone)
	if err != nil {
		http.Error(w, "Invalid transit date/time format", http.StatusBadRequest)
		return
	}

	transitJD := astronomyTime.UTCToJulianDay(transitUTC)

	ayanamsa := input.Ayanamsa
	if ayanamsa == "" {
		ayanamsa = "Lahiri"
	}

	config := domain.CalculationConfig{
		AyanamsaMode: ephemeris.GetAyanamsaMode(ayanamsa),
		HouseCode:    ephemeris.GetHouseSystemCode("Placidus"), // standard fallback
	}

	// Create a dummy BirthInput for context
	dummyInput := domain.BirthInput{
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
		Timezone:  input.Timezone,
		Ayanamsa:  ayanamsa,
	}

	transitCtx := domain.CalculationContext{
		Input:       dummyInput,
		Config:      config,
		UTCTime:     transitUTC,
		JulianDayUT: transitJD,
	}

	// Run Transit
	res, err := transit.CalculateTransitChart(&transitCtx, fauxMoonLon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
