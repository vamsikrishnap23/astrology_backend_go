package site

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/tejzpr/go-swisseph"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/panchang"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	astrotime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// PanchangRequest is the input of /site/panchang.
// Latitude, longitude and timezone are required; date defaults to today
// in the given timezone.
type PanchangRequest struct {
	Date      string   `json:"date"` // YYYY-MM-DD
	Place     string   `json:"place"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Timezone  *float64 `json:"timezone"` // hours from UTC, e.g. 5.5
}

// PanchangResponse is the daily panchang for one date and place.
// The panchang fields are the same as /api/panchang/daily.
type PanchangResponse struct {
	Place     string  `json:"place"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Ayanamsa  string  `json:"ayanamsa"`
	domain.DailyPanchangResult
}

// PanchangHandler serves /site/panchang. GET reads query parameters
// (date, place, lat|latitude, lon|longitude, tz|timezone); POST reads JSON.
func PanchangHandler(w http.ResponseWriter, r *http.Request) {
	var req PanchangRequest
	switch r.Method {
	case http.MethodGet:
		var err error
		if req, err = panchangRequestFromQuery(r); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	case http.MethodPost:
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
	default:
		writeError(w, http.StatusMethodNotAllowed, "use GET or POST")
		return
	}

	res, status, msg := calculatePanchang(req)
	if msg != "" {
		writeError(w, status, msg)
		return
	}
	if r.Method == http.MethodGet {
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}
	writeJSON(w, http.StatusOK, res)
}

func panchangRequestFromQuery(r *http.Request) (PanchangRequest, error) {
	q := r.URL.Query()
	req := PanchangRequest{Date: q.Get("date"), Place: q.Get("place")}
	num := func(names ...string) (*float64, error) {
		for _, n := range names {
			if s := q.Get(n); s != "" {
				v, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return nil, &paramError{n}
				}
				return &v, nil
			}
		}
		return nil, nil
	}
	var err error
	if req.Latitude, err = num("latitude", "lat"); err != nil {
		return req, err
	}
	if req.Longitude, err = num("longitude", "lon"); err != nil {
		return req, err
	}
	if req.Timezone, err = num("timezone", "tz"); err != nil {
		return req, err
	}
	return req, nil
}

type paramError struct{ name string }

func (e *paramError) Error() string { return e.name + " must be a number" }

// calculatePanchang returns the response, or an HTTP status and message.
func calculatePanchang(req PanchangRequest) (PanchangResponse, int, string) {
	if req.Latitude == nil || req.Longitude == nil || req.Timezone == nil {
		return PanchangResponse{}, http.StatusBadRequest, "latitude, longitude and timezone are required"
	}
	lat, lon, tz := *req.Latitude, *req.Longitude, *req.Timezone
	if msg := validLocation(lat, lon, tz); msg != "" {
		return PanchangResponse{}, http.StatusBadRequest, msg
	}

	loc := time.FixedZone("Local", int(tz*3600))
	if req.Date == "" {
		req.Date = time.Now().In(loc).Format("2006-01-02")
	}
	day, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return PanchangResponse{}, http.StatusBadRequest, "date must be YYYY-MM-DD"
	}
	if msg := validYear(day.Year()); msg != "" {
		return PanchangResponse{}, http.StatusBadRequest, msg
	}

	// Local noon for the first pass; only the Telugu calendar reads this
	// instant, and it is redone at sunrise below.
	utc, _ := astrotime.ParseLocalToUTC(req.Date, "12:00:00", tz)
	lahiri := ephemeris.GetAyanamsaMode("lahiri")
	ctx := &domain.CalculationContext{
		Input: domain.BirthInput{
			DateOfBirth: req.Date, TimeOfBirth: "12:00:00",
			Latitude: lat, Longitude: lon, Timezone: tz,
			Ayanamsa: "Lahiri", PlaceOfBirth: req.Place,
		},
		UTCTime:     utc,
		JulianDayUT: astrotime.UTCToJulianDay(utc),
		Config:      domain.CalculationConfig{AyanamsaMode: lahiri},
	}
	daily, err := panchang.CalculateDailyPanchang(ctx)
	if err != nil {
		return PanchangResponse{}, http.StatusInternalServerError, "could not calculate the panchang"
	}

	// The Telugu calendar (masam, paksha, samvatsara...) is the one in
	// force at local sunrise. Fall back to noon where the Sun does not rise.
	calJD := ctx.JulianDayUT
	if sr, err := time.Parse(time.RFC3339, daily.Sunrise); err == nil && sr.In(loc).Format("2006-01-02") == req.Date {
		calJD = astrotime.UTCToJulianDay(sr.UTC())
	}
	ephemeris.Lock()
	swisseph.SetSidMode(int32(lahiri), 0, 0)
	daily.TeluguCalendar = panchang.CalculateTeluguCalendar(calJD)
	ephemeris.Unlock()

	return PanchangResponse{
		Place:               req.Place,
		Latitude:            lat,
		Longitude:           lon,
		Ayanamsa:            "Lahiri",
		DailyPanchangResult: daily,
	}, http.StatusOK, ""
}
