package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http/httptest"
	"time"

	"github.com/tejzpr/go-swisseph"
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
	req := httptest.NewRequest("POST", "/api/v1/chart", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handlers.ChartHandler(rr, req)

	respBody, _ := io.ReadAll(rr.Body)
	var res domain.ChartData
	json.Unmarshal(respBody, &res)

	// Independently calculate using raw Swiss Ephemeris API
	utcTime := time.Date(2005, 11, 23, 10, 5, 0, 0, time.UTC)
	hour := float64(utcTime.Hour()) + float64(utcTime.Minute())/60.0 + float64(utcTime.Second())/3600.0
	jdRaw := swisseph.Julday(int32(utcTime.Year()), int32(utcTime.Month()), int32(utcTime.Day()), hour, swisseph.GregCal)

	swisseph.SetSidMode(swisseph.SidmLahiri, 0, 0)
	ayaRaw := swisseph.GetAyanamsaUT(jdRaw)

	fmt.Println("--- INDEPENDENT VERIFICATION TABLE ---")
	fmt.Printf("%-20s | %-15s | %-15s | %-15s\n", "Category", "Go App", "Raw Swisseph", "Diff")
	fmt.Println("-----------------------------------------------------------------------------")
	fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", "Julian Day", res.JulianDay, jdRaw, math.Abs(res.JulianDay-jdRaw))
	fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", "Ayanamsa (Lahiri)", res.AyanamsaValue, ayaRaw, math.Abs(res.AyanamsaValue-ayaRaw))

	planets := map[string]int32{
		"Sun": swisseph.Sun, "Moon": swisseph.Moon, "Mars": swisseph.Mars,
		"Mercury": swisseph.Mercury, "Jupiter": swisseph.Jupiter, "Venus": swisseph.Venus,
		"Saturn": swisseph.Saturn, "Uranus": swisseph.Uranus, "Neptune": swisseph.Neptune,
		"Pluto": swisseph.Pluto, "Rahu": swisseph.MeanNode,
	}

	var rawRahuTrop, rawRahuSid, rawRahuSpeed float64
	for _, p := range res.Planets {
		if p.Planet == "Ketu" {
			// Direct check for Ketu = Rahu + 180
			kTrop := math.Mod(rawRahuTrop+180.0, 360.0)
			kSid := math.Mod(rawRahuSid+180.0, 360.0)

			fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", "Ketu (Trop)", p.TropicalLongitude, kTrop, math.Abs(p.TropicalLongitude-kTrop))
			fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", "Ketu (Sid)", p.SiderealLongitude, kSid, math.Abs(p.SiderealLongitude-kSid))
			fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", "Ketu (Speed)", p.Speed, rawRahuSpeed, math.Abs(p.Speed-rawRahuSpeed))
			continue
		}

		seID := planets[p.Planet]

		tflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed)
		tres := swisseph.CalcUT(jdRaw, seID, tflag)
		sflag := int32(swisseph.FlagSwieph | swisseph.FlagSpeed | swisseph.FlagSidereal)
		sres := swisseph.CalcUT(jdRaw, seID, sflag)

		if p.Planet == "Rahu" {
			rawRahuTrop = tres.Data[0]
			rawRahuSid = sres.Data[0]
			rawRahuSpeed = tres.Data[3]
		}

		fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", p.Planet+" (Trop)", p.TropicalLongitude, tres.Data[0], math.Abs(p.TropicalLongitude-tres.Data[0]))
		fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", p.Planet+" (Sid)", p.SiderealLongitude, sres.Data[0], math.Abs(p.SiderealLongitude-sres.Data[0]))
		fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", p.Planet+" (Speed)", p.Speed, tres.Data[3], math.Abs(p.Speed-tres.Data[3]))

		// Retrograde status
		isRetRaw := tres.Data[3] < 0
		if p.Retrograde != isRetRaw {
			fmt.Printf("!!! Retrograde mismatch for %s\n", p.Planet)
		}
	}

	resHouse := swisseph.HousesEx(jdRaw, int32(swisseph.FlagSidereal), input.Latitude, input.Longitude, 'P')
	fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", "Ascendant", res.Ascendant, resHouse.Points[0], math.Abs(res.Ascendant-resHouse.Points[0]))
	fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", "MC", res.MC, resHouse.Points[1], math.Abs(res.MC-resHouse.Points[1]))

	for i, h := range res.Houses {
		fmt.Printf("%-20s | %-15.6f | %-15.6f | %-15.9f\n", fmt.Sprintf("House %d", i+1), h.Longitude, resHouse.Houses[i], math.Abs(h.Longitude-resHouse.Houses[i]))
	}
}
