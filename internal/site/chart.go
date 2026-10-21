package site

import (
	"encoding/json"
	"math"
	"net/http"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	astrotime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// ChartRequest is the input of /site/chart. Field names match the birth
// input of /api/*; ayanamsa and house system are fixed (Lahiri, Placidus).
type ChartRequest struct {
	Name         string   `json:"name"`
	DateOfBirth  string   `json:"date_of_birth"` // YYYY-MM-DD
	TimeOfBirth  string   `json:"time_of_birth"` // HH:MM or HH:MM:SS
	PlaceOfBirth string   `json:"place_of_birth"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Timezone     *float64 `json:"timezone"`
}

// Position is a point of the zodiac, shared by planets, ascendant and cusps.
type Position struct {
	Longitude     float64 `json:"longitude"`   // sidereal, 0-360
	Sign          string  `json:"sign"`        // English sign name
	SignNumber    int     `json:"sign_number"` // 1 = Aries ... 12 = Pisces
	SignLord      string  `json:"sign_lord"`
	Degree        int     `json:"degree"` // inside the sign
	Minute        int     `json:"minute"`
	Second        int     `json:"second"`
	Nakshatra     string  `json:"nakshatra"`
	NakshatraPada int     `json:"nakshatra_pada"`
	NakshatraLord string  `json:"nakshatra_lord"`
}

// ChartPlanet is one graha in the D1 chart.
type ChartPlanet struct {
	Planet string `json:"planet"`
	Position
	House      int  `json:"house"` // Placidus bhava (cusp to next cusp)
	Retrograde bool `json:"retrograde"`
	Combust    bool `json:"combust"`
}

// ChartHouse is one Placidus cusp (start of the house).
type ChartHouse struct {
	House int `json:"house"`
	Position
}

// ChartSign is one box of the Rasi chart: who sits in that sign.
type ChartSign struct {
	Sign        string   `json:"sign"`
	SignNumber  int      `json:"sign_number"`
	IsAscendant bool     `json:"is_ascendant"`
	Planets     []string `json:"planets"`
}

// ChartResponse is the D1 chart.
type ChartResponse struct {
	Name               string        `json:"name"`
	DateOfBirth        string        `json:"date_of_birth"`
	TimeOfBirth        string        `json:"time_of_birth"`
	PlaceOfBirth       string        `json:"place_of_birth"`
	Latitude           float64       `json:"latitude"`
	Longitude          float64       `json:"longitude"`
	Timezone           float64       `json:"timezone"`
	CalculationTimeUTC string        `json:"calculation_time_utc"`
	Ayanamsa           string        `json:"ayanamsa"`
	AyanamsaValue      float64       `json:"ayanamsa_value"`
	HouseSystem        string        `json:"house_system"`
	Ascendant          Position      `json:"ascendant"`
	Planets            []ChartPlanet `json:"planets"`
	Houses             []ChartHouse  `json:"houses"`
	Signs              []ChartSign   `json:"signs"`
}

var signNames = []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}

var signLords = []string{"Mars", "Venus", "Mercury", "Moon", "Sun", "Mercury", "Venus", "Mars", "Jupiter", "Saturn", "Saturn", "Jupiter"}

var nakshatraNames = []string{
	"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
	"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni",
	"Uttara Phalguni", "Hasta", "Chitra", "Swati", "Vishakha", "Anuradha",
	"Jyeshtha", "Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana",
	"Dhanishta", "Shatabhisha", "Purva Bhadrapada", "Uttara Bhadrapada",
	"Revati",
}

var nakshatraLords = []string{"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury"}

// grahas are the nine planets shown in the chart (outer planets left out).
var grahas = []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn", "Rahu", "Ketu"}

// ChartHandler serves POST /site/chart.
func ChartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "use POST")
		return
	}
	var req ChartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	res, status, msg := calculateChart(req)
	if msg != "" {
		writeError(w, status, msg)
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func calculateChart(req ChartRequest) (ChartResponse, int, string) {
	if req.Latitude == nil || req.Longitude == nil || req.Timezone == nil {
		return ChartResponse{}, http.StatusBadRequest, "latitude, longitude and timezone are required"
	}
	lat, lon, tz := *req.Latitude, *req.Longitude, *req.Timezone
	if msg := validLocation(lat, lon, tz); msg != "" {
		return ChartResponse{}, http.StatusBadRequest, msg
	}
	if len(req.TimeOfBirth) == 5 { // HH:MM
		req.TimeOfBirth += ":00"
	}
	utc, err := astrotime.ParseLocalToUTC(req.DateOfBirth, req.TimeOfBirth, tz)
	if err != nil {
		return ChartResponse{}, http.StatusBadRequest, "date_of_birth must be YYYY-MM-DD and time_of_birth HH:MM or HH:MM:SS"
	}
	if msg := validYear(utc.Year()); msg != "" {
		return ChartResponse{}, http.StatusBadRequest, msg
	}

	ctx := &domain.CalculationContext{
		Input: domain.BirthInput{
			Name: req.Name, DateOfBirth: req.DateOfBirth, TimeOfBirth: req.TimeOfBirth,
			PlaceOfBirth: req.PlaceOfBirth, Latitude: lat, Longitude: lon, Timezone: tz,
			Ayanamsa: "Lahiri", HouseSystem: "Placidus",
		},
		Config: domain.CalculationConfig{
			AyanamsaMode: ephemeris.GetAyanamsaMode("lahiri"),
			HouseCode:    ephemeris.GetHouseSystemCode("placidus"),
		},
		UTCTime:     utc,
		JulianDayUT: astrotime.UTCToJulianDay(utc),
	}
	pls, err := planets.CalculatePlanets(ctx)
	if err != nil {
		return ChartResponse{}, http.StatusInternalServerError, "could not calculate planet positions"
	}
	asc, _, cusps, err := houses.CalculateHouses(ctx)
	if err != nil {
		return ChartResponse{}, http.StatusInternalServerError, "could not calculate houses"
	}

	res := ChartResponse{
		Name: req.Name, DateOfBirth: req.DateOfBirth, TimeOfBirth: req.TimeOfBirth,
		PlaceOfBirth: req.PlaceOfBirth, Latitude: lat, Longitude: lon, Timezone: tz,
		CalculationTimeUTC: utc.Format(time.RFC3339),
		Ayanamsa:           "Lahiri",
		AyanamsaValue:      ctx.Ayanamsa,
		HouseSystem:        "Placidus",
		Ascendant:          position(asc),
	}

	cuspLons := make([]float64, len(cusps))
	for i, c := range cusps {
		cuspLons[i] = c.Longitude
		res.Houses = append(res.Houses, ChartHouse{House: c.HouseNumber, Position: position(c.Longitude)})
	}

	byName := map[string]domain.PlanetPosition{}
	for _, p := range pls {
		byName[p.Planet] = p
	}
	for i := range signNames {
		res.Signs = append(res.Signs, ChartSign{Sign: signNames[i], SignNumber: i + 1, Planets: []string{}})
	}
	res.Signs[res.Ascendant.SignNumber-1].IsAscendant = true
	for _, name := range grahas {
		p, ok := byName[name]
		if !ok {
			return ChartResponse{}, http.StatusInternalServerError, "missing planet " + name
		}
		pos := position(p.SiderealLongitude)
		res.Planets = append(res.Planets, ChartPlanet{
			Planet:     name,
			Position:   pos,
			House:      houseOf(p.SiderealLongitude, cuspLons),
			Retrograde: p.Retrograde,
			Combust:    p.Combust,
		})
		res.Signs[pos.SignNumber-1].Planets = append(res.Signs[pos.SignNumber-1].Planets, name)
	}
	return res, http.StatusOK, ""
}

// position splits a sidereal longitude into sign, D/M/S and nakshatra.
func position(lon float64) Position {
	lon = math.Mod(lon, 360)
	if lon < 0 {
		lon += 360
	}
	// Truncate (not round) to whole arc-seconds, so a planet at 29°59'59.7"
	// stays in its sign, like the D/M/S of /api/chart.
	totalSec := int(lon * 3600)
	if totalSec >= 360*3600 {
		totalSec = 360*3600 - 1
	}
	sign := totalSec / (30 * 3600)
	inSign := totalSec % (30 * 3600)

	nakSpan := 360.0 / 27
	nak := int(lon / nakSpan)
	if nak > 26 {
		nak = 26
	}
	pada := int(math.Mod(lon, nakSpan)/(nakSpan/4)) + 1
	if pada > 4 {
		pada = 4
	}
	return Position{
		Longitude:     lon,
		Sign:          signNames[sign],
		SignNumber:    sign + 1,
		SignLord:      signLords[sign],
		Degree:        inSign / 3600,
		Minute:        inSign % 3600 / 60,
		Second:        inSign % 60,
		Nakshatra:     nakshatraNames[nak],
		NakshatraPada: pada,
		NakshatraLord: nakshatraLords[nak%9],
	}
}

// houseOf returns the Placidus house whose cusp interval holds lon
// (same rule as /api/tables: cusp i up to, not including, cusp i+1).
func houseOf(lon float64, cusps []float64) int {
	n := len(cusps)
	for i := 0; i < n; i++ {
		a, b := cusps[i], cusps[(i+1)%n]
		if a < b {
			if lon >= a && lon < b {
				return i + 1
			}
		} else if lon >= a || lon < b {
			return i + 1
		}
	}
	return 1
}
