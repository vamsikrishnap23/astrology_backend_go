package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func TestChartHandler(t *testing.T) {
	// Initialize ephemeris
	if err := ephemeris.Init("../../../ephe_data"); err != nil {
		t.Fatalf("ephemeris init failed: %v", err)
	}
	defer ephemeris.Close()

	input := domain.BirthInput{
		Name:         "Example",
		DateOfBirth:  "1998-04-28",
		TimeOfBirth:  "14:30:00",
		PlaceOfBirth: "Hyderabad",
		Latitude:     17.3850,
		Longitude:    78.4867,
		Timezone:     5.5,
		Ayanamsa:     "Lahiri",
		HouseSystem:  "Placidus",
	}

	body, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/api/v1/chart", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ChartHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res domain.ChartData
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	// Verify standard fields
	if res.SelectedAyanamsa != "Lahiri" {
		t.Errorf("Expected Lahiri, got %v", res.SelectedAyanamsa)
	}
	if len(res.Planets) == 0 {
		t.Error("Expected planets, got none")
	}
	if len(res.Houses) != 12 {
		t.Errorf("Expected 12 houses, got %d", len(res.Houses))
	}
}
