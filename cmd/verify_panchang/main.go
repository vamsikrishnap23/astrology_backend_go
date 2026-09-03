package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"time"

	"github.com/vamsi/astrology_backend_go/internal/api/handlers"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func main() {
	_ = ephemeris.Init("ephe_data")
	defer ephemeris.Close()

	input := domain.BirthInput{
		Name:         "Verify",
		DateOfBirth:  "2005-11-23",
		TimeOfBirth:  "15:35:00",
		PlaceOfBirth: "Unknown",
		Latitude:     16.3938,
		Longitude:    80.1522,
		Timezone:     5.5,
		Ayanamsa:     "Lahiri",
		HouseSystem:  "Placidus",
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/api/v1/panchang", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handlers.PanchangHandler(rr, req)

	respBody, _ := io.ReadAll(rr.Body)
	var res domain.PanchangResult
	json.Unmarshal(respBody, &res)

	fmt.Println("--- PANCHANG INDEPENDENT VERIFICATION ---")
	fmt.Printf("Date: %s %s\n", res.Date, res.LocalTime)

	fmt.Printf("\nSunrise: %s\n", res.Sunrise)
	fmt.Printf("Sunset: %s\n", res.Sunset)
	fmt.Printf("Solar Noon: %s\n", res.SolarNoon)
	fmt.Printf("\nMoonrise: %s\n", res.Moonrise)
	fmt.Printf("Moonset: %s\n", res.Moonset)

	fmt.Printf("\nVara (Weekday): %s (Ruler: %s)\n", res.Vara.Name, res.Vara.Ruler)

	fmt.Printf("\nTithi: %d - %s (%s Paksha)\n", res.Tithi.Number, res.Tithi.Name, res.Tithi.Paksha)
	fmt.Printf("Tithi Progress: %.2f%%\n", res.Tithi.Progress)
	fmt.Printf("Tithi Start: %s\n", res.Tithi.Start)
	fmt.Printf("Tithi End: %s\n", res.Tithi.End)

	fmt.Printf("\nNakshatra: %d - %s\n", res.Nakshatra.Number, res.Nakshatra.Name)
	fmt.Printf("Nakshatra Pada: %d\n", res.Nakshatra.Pada)
	fmt.Printf("Nakshatra Progress: %.2f%%\n", res.Nakshatra.Progress)
	fmt.Printf("Nakshatra Start: %s\n", res.Nakshatra.Start)
	fmt.Printf("Nakshatra End: %s\n", res.Nakshatra.End)

	fmt.Printf("\nYoga: %d - %s\n", res.Yoga.Number, res.Yoga.Name)
	fmt.Printf("Yoga Progress: %.2f%%\n", res.Yoga.Progress)
	fmt.Printf("Yoga Start: %s\n", res.Yoga.Start)
	fmt.Printf("Yoga End: %s\n", res.Yoga.End)

	fmt.Printf("\nKarana: %d - %s (%s)\n", res.Karana.Number, res.Karana.Name, res.Karana.Type)
	fmt.Printf("Karana Progress: %.2f%%\n", res.Karana.Progress)
	fmt.Printf("Karana Start: %s\n", res.Karana.Start)
	fmt.Printf("Karana End: %s\n", res.Karana.End)

	// Add delays to prevent execution issues
	time.Sleep(50 * time.Millisecond)
}
