package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func postShadbala(t *testing.T, in domain.BirthInput) *httptest.ResponseRecorder {
	t.Helper()
	body, _ := json.Marshal(in)
	req, err := http.NewRequest("POST", "/api/shadbala", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(ShadbalaHandler).ServeHTTP(rr, req)
	return rr
}

func TestShadbalaHandler(t *testing.T) {
	if err := ephemeris.Init("../../../ephe_data"); err != nil {
		t.Fatalf("ephemeris init failed: %v", err)
	}
	defer ephemeris.Close()

	ok := postShadbala(t, domain.BirthInput{
		Name: "Example", DateOfBirth: "1998-04-28", TimeOfBirth: "14:30:00", Latitude: 17.385, Longitude: 78.4867,
		Timezone: 5.5, Ayanamsa: "Lahiri", HouseSystem: "Placidus",
	})
	if ok.Code != http.StatusOK {
		t.Fatalf("ordinary chart: status %d: %s", ok.Code, ok.Body.String())
	}
	var res domain.ShadbalaResult
	if err := json.NewDecoder(ok.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}
	if len(res.Planets) != 7 {
		t.Errorf("got %d planets", len(res.Planets))
	}
	for _, p := range res.Planets {
		if (p.Planet == "Sun" || p.Planet == "Moon") && p.CheshtaBala != 0 {
			t.Errorf("%s has no Cheshta Bala in the total, got %v", p.Planet, p.CheshtaBala)
		}
	}

	// Midnight sun: no sunrise to count from, so a clear client error and not zeros.
	polar := postShadbala(t, domain.BirthInput{
		Name: "Polar", DateOfBirth: "2023-06-21", TimeOfBirth: "12:00:00", Latitude: 80, Longitude: 25,
		Timezone: 0, Ayanamsa: "Lahiri", HouseSystem: "Placidus",
	})
	if polar.Code != http.StatusBadRequest {
		t.Fatalf("polar chart: status %d, want 400: %s", polar.Code, polar.Body.String())
	}
	if !strings.Contains(polar.Body.String(), "does not rise or set") {
		t.Errorf("unhelpful message: %q", polar.Body.String())
	}
}
