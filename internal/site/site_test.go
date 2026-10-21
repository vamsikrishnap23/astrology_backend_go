package site

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
)

func TestMain(m *testing.M) {
	if err := ephemeris.Init("../../ephe_data"); err != nil {
		fmt.Println("ephemeris init failed:", err)
		os.Exit(1)
	}
	code := m.Run()
	ephemeris.Close()
	os.Exit(code)
}

type valuesFile struct {
	EnumPaths map[string]map[string]string         `json:"enum_paths"`
	Enums     map[string]map[string]map[string]any `json:"enums"`
	Labels    map[string]map[string]string         `json:"labels"`
}

func loadValues(t *testing.T) valuesFile {
	t.Helper()
	var v valuesFile
	if err := json.Unmarshal(valuesJSON, &v); err != nil {
		t.Fatalf("values.json is not valid JSON: %v", err)
	}
	for name, enum := range v.Enums {
		for key, tr := range enum {
			if tr["te"] == nil || tr["te"] == "" || tr["en"] == nil {
				t.Errorf("enums.%s.%s needs en and te", name, key)
			}
		}
	}
	return v
}

// collect walks a decoded JSON response and gathers the value at every path
// listed in enum_paths ("a.b", "a[].b", "a[].b[]").
func collect(doc any, path string) []string {
	var out []string
	var walk func(v any, parts []string)
	walk = func(v any, parts []string) {
		if len(parts) == 0 {
			if s, ok := v.(string); ok {
				out = append(out, s)
			}
			return
		}
		p := parts[0]
		list := strings.HasSuffix(p, "[]")
		p = strings.TrimSuffix(p, "[]")
		m, ok := v.(map[string]any)
		if !ok {
			return
		}
		child := m[p]
		if list {
			arr, _ := child.([]any)
			for _, e := range arr {
				walk(e, parts[1:])
			}
		} else {
			walk(child, parts[1:])
		}
	}
	walk(doc, strings.Split(path, "."))
	return out
}

// checkCovered fails if a response value at an enum path is missing from
// values.json, and records which values were seen.
func checkCovered(t *testing.T, v valuesFile, endpoint string, body []byte, seen map[string]map[string]bool) {
	t.Helper()
	var doc any
	if err := json.Unmarshal(body, &doc); err != nil {
		t.Fatalf("bad JSON: %v", err)
	}
	for path, enum := range v.EnumPaths[endpoint] {
		vals := collect(doc, path)
		for _, s := range vals {
			if _, ok := v.Enums[enum][s]; !ok {
				t.Errorf("%s %s = %q is not in enums.%s", endpoint, path, s, enum)
			}
			if seen[enum] == nil {
				seen[enum] = map[string]bool{}
			}
			seen[enum][s] = true
		}
	}
}

func post(h http.HandlerFunc, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	rr := httptest.NewRecorder()
	h(rr, httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b)))
	return rr
}

func TestPanchangValuesCovered(t *testing.T) {
	v := loadValues(t)
	seen := map[string]map[string]bool{}
	lat, lon, tz := 17.385, 78.4867, 5.5
	day := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 400; i++ {
		d := day.AddDate(0, 0, i).Format("2006-01-02")
		rr := post(PanchangHandler, PanchangRequest{Date: d, Latitude: &lat, Longitude: &lon, Timezone: &tz})
		if rr.Code != http.StatusOK {
			t.Fatalf("%s: status %d: %s", d, rr.Code, rr.Body)
		}
		checkCovered(t, v, "panchang", rr.Body.Bytes(), seen)
	}
	// A year sees every weekday, tithi, nakshatra, yoga and karana.
	for enum, want := range map[string]int{"weekdays": 7, "tithis": 16, "nakshatras": 27, "yogas": 27, "karanas": 11, "pakshas": 2, "ayanas": 2, "rutus": 6} {
		if got := len(seen[enum]); got != want {
			t.Errorf("saw %d %s in a year, want %d", got, enum, want)
		}
	}
	if n := len(v.Enums["samvatsaras"]); n != 60 {
		t.Errorf("samvatsaras: %d, want 60", n)
	}
	if n := len(v.Enums["masams"]); n != 24 {
		t.Errorf("masams: %d, want 24 (12 + Adhika)", n)
	}
}

func TestSamvatsaraValuesCovered(t *testing.T) {
	v := loadValues(t)
	seen := map[string]map[string]bool{}
	lat, lon, tz := 17.385, 78.4867, 5.5
	for y := 1990; y < 2050; y++ {
		d := fmt.Sprintf("%d-06-15", y)
		rr := post(PanchangHandler, PanchangRequest{Date: d, Latitude: &lat, Longitude: &lon, Timezone: &tz})
		if rr.Code != http.StatusOK {
			t.Fatalf("%s: status %d", d, rr.Code)
		}
		checkCovered(t, v, "panchang", rr.Body.Bytes(), seen)
	}
	if got := len(seen["samvatsaras"]); got != 60 {
		t.Errorf("saw %d samvatsaras over 60 years, want 60", got)
	}
}

func TestPanchangGetAndErrors(t *testing.T) {
	rr := httptest.NewRecorder()
	PanchangHandler(rr, httptest.NewRequest(http.MethodGet, "/site/panchang?date=2026-10-21&lat=17.385&lon=78.4867&tz=5.5&place=Hyderabad", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("GET status %d: %s", rr.Code, rr.Body)
	}
	var res PanchangResponse
	json.Unmarshal(rr.Body.Bytes(), &res)
	if res.Place != "Hyderabad" || res.Ayanamsa != "Lahiri" || res.Date != "2026-10-21" || res.Vara.Name != "Wednesday" {
		t.Errorf("unexpected response: place=%q ayanamsa=%q date=%q vara=%q", res.Place, res.Ayanamsa, res.Date, res.Vara.Name)
	}
	if !strings.HasPrefix(res.Sunrise, "2026-10-21T06:") {
		t.Errorf("sunrise %q, want about 06:00 in Hyderabad", res.Sunrise)
	}

	for _, url := range []string{
		"/site/panchang?lat=17&lon=78",                        // no tz
		"/site/panchang?date=2026-13-01&lat=17&lon=78&tz=5.5", // bad date
		"/site/panchang?date=1700-01-01&lat=17&lon=78&tz=5.5", // outside ephemeris
		"/site/panchang?lat=95&lon=78&tz=5.5",                 // bad latitude
		"/site/panchang?lat=abc&lon=78&tz=5.5",
	} {
		rr := httptest.NewRecorder()
		PanchangHandler(rr, httptest.NewRequest(http.MethodGet, url, nil))
		if rr.Code != http.StatusBadRequest {
			t.Errorf("%s: status %d, want 400", url, rr.Code)
		}
	}
}

func TestChartValuesCovered(t *testing.T) {
	v := loadValues(t)
	seen := map[string]map[string]bool{}
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 300; i++ {
		ts := time.Date(1950+rng.Intn(80), time.Month(1+rng.Intn(12)), 1+rng.Intn(28), rng.Intn(24), rng.Intn(60), 0, 0, time.UTC)
		lat, lon, tz := -40+rng.Float64()*95, -150+rng.Float64()*300, 5.5
		req := ChartRequest{DateOfBirth: ts.Format("2006-01-02"), TimeOfBirth: ts.Format("15:04"), Latitude: &lat, Longitude: &lon, Timezone: &tz}
		rr := post(ChartHandler, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("%+v: status %d: %s", req, rr.Code, rr.Body)
		}
		checkCovered(t, v, "chart", rr.Body.Bytes(), seen)

		var res ChartResponse
		json.Unmarshal(rr.Body.Bytes(), &res)
		if len(res.Planets) != 9 || len(res.Houses) != 12 || len(res.Signs) != 12 {
			t.Fatalf("counts: %d planets, %d houses, %d signs", len(res.Planets), len(res.Houses), len(res.Signs))
		}
		if res.Houses[0].Sign != res.Ascendant.Sign {
			t.Errorf("house 1 sign %s != ascendant sign %s", res.Houses[0].Sign, res.Ascendant.Sign)
		}
		n := 0
		for _, s := range res.Signs {
			n += len(s.Planets)
		}
		if n != 9 {
			t.Errorf("signs hold %d planets, want 9", n)
		}
	}
	for enum, want := range map[string]int{"signs": 12, "nakshatras": 27, "planets": 9} {
		if got := len(seen[enum]); got != want {
			t.Errorf("saw %d %s, want %d", got, enum, want)
		}
	}
}

func TestChartMatchesAPIChart(t *testing.T) {
	// Hyderabad 1998-04-28 14:30 IST (same chart as handlers/chart_test.go).
	lat, lon, tz := 17.385, 78.4867, 5.5
	rr := post(ChartHandler, ChartRequest{DateOfBirth: "1998-04-28", TimeOfBirth: "14:30:00", Latitude: &lat, Longitude: &lon, Timezone: &tz})
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rr.Code, rr.Body)
	}
	var res ChartResponse
	json.Unmarshal(rr.Body.Bytes(), &res)
	sun := res.Planets[0]
	if sun.Planet != "Sun" || sun.Sign != "Aries" {
		t.Errorf("Sun in %s, want Aries", sun.Sign)
	}
	if res.Ascendant.Sign != "Leo" {
		t.Errorf("ascendant %s, want Leo", res.Ascendant.Sign)
	}
}

func TestPositionBoundaries(t *testing.T) {
	p := position(29.99999)
	if p.Sign != "Aries" || p.Degree != 29 || p.Minute != 59 || p.Second != 59 {
		t.Errorf("29.99999 -> %s %d %d %d", p.Sign, p.Degree, p.Minute, p.Second)
	}
	p = position(359.999999)
	if p.Sign != "Pisces" || p.Nakshatra != "Revati" || p.NakshatraPada != 4 {
		t.Errorf("359.999999 -> %+v", p)
	}
	p = position(0)
	if p.Sign != "Aries" || p.Nakshatra != "Ashwini" || p.NakshatraPada != 1 || p.NakshatraLord != "Ketu" {
		t.Errorf("0 -> %+v", p)
	}
}

func TestCORS(t *testing.T) {
	mux := http.NewServeMux()
	t.Setenv("SITE_ALLOWED_ORIGINS", "https://example.com")
	Register(mux)
	req := httptest.NewRequest(http.MethodOptions, "/site/chart", nil)
	req.Header.Set("Origin", "https://example.com")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent || rr.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("preflight: %d %q", rr.Code, rr.Header().Get("Access-Control-Allow-Origin"))
	}
	req = httptest.NewRequest(http.MethodGet, "/site/values", nil)
	req.Header.Set("Origin", "https://evil.example")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("other origin allowed")
	}
}
