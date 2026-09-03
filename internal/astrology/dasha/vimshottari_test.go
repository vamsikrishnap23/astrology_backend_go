package dasha

import (
	"math"
	"testing"
	"time"
)

func TestCalculateVimshottari(t *testing.T) {
	// DOB: 2005-11-23 15:35:00 IST
	// UTC: 2005-11-23 10:05:00 UTC
	birthTimeUTC, _ := time.Parse("2006-01-02 15:04:05", "2005-11-23 10:05:00")

	// Reference Moon Sidereal Longitude: 121.819768 (Magha, Ketu)
	moonLon := 121.819768

	res := CalculateVimshottari(moonLon, birthTimeUTC)

	// Ketu is 7 years total.
	// 121.819768 / 13.333333 = 9.1364826 (Nakshatra 9 = Magha = Ketu lord)
	// Remainder = 1.819768 degrees
	// fraction elapsed = 1.819768 / 13.333333 = 0.1364826
	// fraction remaining = 1.0 - 0.1364826 = 0.8635174
	// Balance = 0.8635174 * 7 = 6.0446218
	expectedBalance := 6.044622
	if math.Abs(res.BalanceYears-expectedBalance) > 1e-5 {
		t.Errorf("Expected balance %f, got %f", expectedBalance, res.BalanceYears)
	}

	if len(res.Mahadasha) != 9 {
		t.Fatalf("Expected 9 Mahadashas, got %d", len(res.Mahadasha))
	}

	if res.Mahadasha[0].Lord != "Ketu" {
		t.Errorf("Expected first MD to be Ketu, got %s", res.Mahadasha[0].Lord)
	}
	if res.Mahadasha[1].Lord != "Venus" {
		t.Errorf("Expected second MD to be Venus, got %s", res.Mahadasha[1].Lord)
	}

	passedYears := 0.1364826 * 7.0
	theoreticalStart := addYears(birthTimeUTC, -passedYears)

	// First MD Start Date should be theoretical start date now
	if math.Abs(res.Mahadasha[0].StartDate.Sub(theoreticalStart).Seconds()) > 1.0 {
		t.Errorf("First MD start date should be theoretical start time %v, got %v", theoreticalStart, res.Mahadasha[0].StartDate)
	}

	// First MD End Date = birthTimeUTC + BalanceYears
	expectedKetuEnd := addYears(birthTimeUTC, res.BalanceYears)
	if math.Abs(res.Mahadasha[0].EndDate.Sub(expectedKetuEnd).Seconds()) > 1.0 {
		t.Errorf("Expected Ketu MD end %v, got %v", expectedKetuEnd, res.Mahadasha[0].EndDate)
	}

	// Second MD End Date = First MD End Date + 20 years
	expectedVenusEnd := addYears(expectedKetuEnd, 20.0)
	if math.Abs(res.Mahadasha[1].EndDate.Sub(expectedVenusEnd).Seconds()) > 1.0 {
		t.Errorf("Expected Venus MD end %v, got %v", expectedVenusEnd, res.Mahadasha[1].EndDate)
	}

	// Let's verify the first AD in the result is Venus!
	firstAD := res.Mahadasha[0].Antardasha[0]
	if firstAD.Lord != "Venus" {
		t.Errorf("Expected first AD to be Venus (since Ketu AD passed), got %s", firstAD.Lord)
	}

	expectedKetuVenusStart := addYears(theoreticalStart, (7.0*7.0)/120.0)
	if math.Abs(firstAD.StartDate.Sub(expectedKetuVenusStart).Seconds()) > 1.0 {
		t.Errorf("First valid AD start date should be theoretical %v, got %v", expectedKetuVenusStart, firstAD.StartDate)
	}

	// Ketu-Venus ends at: Theoretical Ketu Start + Ketu-Ketu duration + Ketu-Venus duration
	expectedKetuVenusEnd := addYears(theoreticalStart, (7.0*7.0)/120.0+(7.0*20.0)/120.0)

	// We need to compare ignoring tiny nanosecond differences due to floating point math
	diff := firstAD.EndDate.Sub(expectedKetuVenusEnd).Seconds()
	if math.Abs(diff) > 1.0 { // Allow 1 second variance due to floating point accumulation
		t.Errorf("Expected Ketu-Venus AD end %v, got %v (diff: %f seconds)", expectedKetuVenusEnd, firstAD.EndDate, diff)
	}

	// Check Sookshma levels are populated
	if len(firstAD.Pratyantardasha) == 0 {
		t.Fatalf("Expected PDs to be generated")
	}
	if len(firstAD.Pratyantardasha[0].Sookshma) == 0 {
		t.Fatalf("Expected Sookshmas to be generated")
	}
}

// Test exact nakshatra boundary
func TestBoundaryWraparound(t *testing.T) {
	birthTimeUTC := time.Now().UTC()
	// Just past Revati (Mercury) -> Ashwini (Ketu)
	moonLon := 360.000001
	res := CalculateVimshottari(moonLon, birthTimeUTC)
	if res.Mahadasha[0].Lord != "Ketu" {
		t.Errorf("Expected Ketu at 0 degrees, got %s", res.Mahadasha[0].Lord)
	}

	// Just before 0 degrees
	moonLon = 359.999999
	res = CalculateVimshottari(moonLon, birthTimeUTC)
	if res.Mahadasha[0].Lord != "Mercury" {
		t.Errorf("Expected Mercury just before 360, got %s", res.Mahadasha[0].Lord)
	}
}
