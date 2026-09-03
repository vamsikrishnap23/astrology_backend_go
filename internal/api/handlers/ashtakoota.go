package handlers

import (
	"encoding/json"
	"github.com/tejzpr/go-swisseph"
	"net/http"

	"github.com/vamsi/astrology_backend_go/internal/astrology/ashtakoota"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func buildContext(input domain.BirthInput) (*domain.CalculationContext, error) {
	utcTime, err := astronomyTime.ParseLocalToUTC(input.DateOfBirth, input.TimeOfBirth, input.Timezone)
	if err != nil {
		return nil, err
	}

	jd := astronomyTime.UTCToJulianDay(utcTime)

	config := domain.CalculationConfig{
		AyanamsaMode: ephemeris.GetAyanamsaMode(input.Ayanamsa),
		HouseCode:    ephemeris.GetHouseSystemCode(input.HouseSystem),
	}

	ephemeris.Mu.Lock()
	swisseph.SetEphePath(ephemeris.EphePath)
	swisseph.SetSidMode(int32(config.AyanamsaMode), 0, 0)
	ayanamsa := swisseph.GetAyanamsaUT(jd)
	ephemeris.Mu.Unlock()

	return &domain.CalculationContext{
		Input:       input,
		Config:      config,
		UTCTime:     utcTime,
		JulianDayUT: jd,
		Ayanamsa:    ayanamsa,
	}, nil
}

func AshtakootaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.MatchInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	groomCtx, err := buildContext(input.Groom)
	if err != nil {
		http.Error(w, "Invalid groom date/time format", http.StatusBadRequest)
		return
	}

	brideCtx, err := buildContext(input.Bride)
	if err != nil {
		http.Error(w, "Invalid bride date/time format", http.StatusBadRequest)
		return
	}

	res, err := ashtakoota.CalculateMatch(groomCtx, brideCtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
