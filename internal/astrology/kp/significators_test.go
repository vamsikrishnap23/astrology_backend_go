package kp

import (
	"reflect"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func TestCalculateSignificators(t *testing.T) {
	// Mock tables
	tables := domain.TablesResult{
		PlanetaryTable: []domain.TablePlanet{
			{PlanetName: "Sun", HouseNumber: 10, NakshatraLord: "Ketu"},
			{PlanetName: "Moon", HouseNumber: 4, NakshatraLord: "Venus"},
			{PlanetName: "Ketu", HouseNumber: 1, NakshatraLord: "Sun"},
			{PlanetName: "Venus", HouseNumber: 5, NakshatraLord: "Moon"},
		},
		HouseTable: []domain.TableHouse{
			{HouseNumber: 1, SignLord: "Mars"},
			{HouseNumber: 4, SignLord: "Moon"},
			{HouseNumber: 5, SignLord: "Sun"},
			{HouseNumber: 7, SignLord: "Venus"},
			{HouseNumber: 10, SignLord: "Saturn"},
			{HouseNumber: 11, SignLord: "Saturn"},
			{HouseNumber: 12, SignLord: "Jupiter"},
		},
	}

	result := CalculateSignificators(tables)

	// Validate Planet View
	if len(result.PlanetView) != 4 {
		t.Fatalf("Expected 4 planets, got %d", len(result.PlanetView))
	}

	var sunSig domain.PlanetSignificator
	for _, ps := range result.PlanetView {
		if ps.Planet == "Sun" {
			sunSig = ps
		}
	}

	// Sun is in 10, Star Lord is Ketu. Ketu is in 1.
	// A = SL's house = [1]
	// B = Planet's house = [10]
	// C = Houses owned by SL (Ketu) = []
	// D = Houses owned by Planet (Sun) = [5]
	expectedA := []int{1}
	expectedB := []int{10}
	expectedC := []int{}
	expectedD := []int{5}

	if !reflect.DeepEqual(sunSig.A, expectedA) {
		t.Errorf("Sun A: got %v, want %v", sunSig.A, expectedA)
	}
	if !reflect.DeepEqual(sunSig.B, expectedB) {
		t.Errorf("Sun B: got %v, want %v", sunSig.B, expectedB)
	}
	if !reflect.DeepEqual(sunSig.C, expectedC) {
		t.Errorf("Sun C: got %v, want %v", sunSig.C, expectedC)
	}
	if !reflect.DeepEqual(sunSig.D, expectedD) {
		t.Errorf("Sun D: got %v, want %v", sunSig.D, expectedD)
	}

	// Check House View for House 1
	var house1 domain.HouseSignificator
	for _, hs := range result.HouseView {
		if hs.House == 1 {
			house1 = hs
		}
	}

	// A: Planets whose Star Lord is in House 1. (Ketu is in House 1. Sun's Star Lord is Ketu). So Sun should be in A.
	expectedHouse1A := []string{"Sun"}
	if !reflect.DeepEqual(house1.A, expectedHouse1A) {
		t.Errorf("House 1 A: got %v, want %v", house1.A, expectedHouse1A)
	}
}
