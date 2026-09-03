package kp

import (
	"reflect"
	"testing"

	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func TestCalculateFourStepSignificators(t *testing.T) {
	tables := domain.TablesResult{
		PlanetaryTable: []domain.TablePlanet{
			{PlanetName: "Sun", NakshatraLord: "Venus", SubLord: "Moon"},
			{PlanetName: "Moon", NakshatraLord: "Jupiter"},
			{PlanetName: "Venus", NakshatraLord: "Mars"},
		},
	}

	sigs := domain.KPSignificatorsResult{
		PlanetView: []domain.PlanetSignificator{
			{Planet: "Sun", B: []int{10}, D: []int{5}},
			{Planet: "Venus", B: []int{7}, D: []int{11, 12}},
			{Planet: "Moon", B: []int{6}, D: []int{4}},
			{Planet: "Jupiter", B: []int{2}, D: []int{9, 12}},
		},
	}

	res := CalculateFourStepSignificators(tables, sigs)

	if len(res.FourStepView) != 3 {
		t.Fatalf("Expected 3 planets in 4-step view, got %d", len(res.FourStepView))
	}

	// Verify Sun
	var sun domain.FourStepSignificator
	for _, fs := range res.FourStepView {
		if fs.Planet == "Sun" {
			sun = fs
			break
		}
	}

	// Step 1: Planet = Sun (B=10, D=5 => 5, 10)
	if sun.PlanetDetails.Planet != "Sun" {
		t.Errorf("Expected Sun PlanetDetails.Planet to be Sun, got %s", sun.PlanetDetails.Planet)
	}
	expectedSunHouses := []int{5, 10}
	if !reflect.DeepEqual(sun.PlanetDetails.Houses, expectedSunHouses) {
		t.Errorf("Expected Sun houses %v, got %v", expectedSunHouses, sun.PlanetDetails.Houses)
	}

	// Step 2: Star Lord = Venus (B=7, D=11, 12 => 7, 11, 12)
	if sun.StarLord.Planet != "Venus" {
		t.Errorf("Expected Sun StarLord to be Venus, got %s", sun.StarLord.Planet)
	}
	expectedVenusHouses := []int{7, 11, 12}
	if !reflect.DeepEqual(sun.StarLord.Houses, expectedVenusHouses) {
		t.Errorf("Expected Venus houses %v, got %v", expectedVenusHouses, sun.StarLord.Houses)
	}

	// Step 3: Sub Lord = Moon (B=6, D=4 => 4, 6)
	if sun.SubLord.Planet != "Moon" {
		t.Errorf("Expected Sun SubLord to be Moon, got %s", sun.SubLord.Planet)
	}
	expectedMoonHouses := []int{4, 6}
	if !reflect.DeepEqual(sun.SubLord.Houses, expectedMoonHouses) {
		t.Errorf("Expected Moon houses %v, got %v", expectedMoonHouses, sun.SubLord.Houses)
	}

	// Step 4: Star Lord of Sub Lord = Jupiter (B=2, D=9, 12 => 2, 9, 12)
	if sun.StarLordOfSub.Planet != "Jupiter" {
		t.Errorf("Expected Sun StarLordOfSub to be Jupiter, got %s", sun.StarLordOfSub.Planet)
	}
	expectedJupiterHouses := []int{2, 9, 12}
	if !reflect.DeepEqual(sun.StarLordOfSub.Houses, expectedJupiterHouses) {
		t.Errorf("Expected Jupiter houses %v, got %v", expectedJupiterHouses, sun.StarLordOfSub.Houses)
	}
}
