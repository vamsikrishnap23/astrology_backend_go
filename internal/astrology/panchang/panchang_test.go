package panchang

import (
	"math"
	"testing"
)

func TestBisectionSearch(t *testing.T) {
	// A simple dummy function y = x - 150
	calcFunc := func(x float64) float64 {
		return math.Mod(x, 360)
	}

	res := bisectionSearch(100.0, 200.0, 150.0, calcFunc)
	if math.Abs(res-150.0) > 0.001 {
		t.Errorf("Expected 150.0, got %f", res)
	}
}

func TestVara(t *testing.T) {
	res := calculateVara("2005-11-23", 5.5)
	if res.Name != "Wednesday" {
		t.Errorf("Expected Wednesday, got %s", res.Name)
	}
	if res.Ruler != "Mercury" {
		t.Errorf("Expected Mercury, got %s", res.Ruler)
	}
}

func TestKaranaNames(t *testing.T) {
	name, ktype := getKaranaDetails(0)
	if name != "Kinstughna" || ktype != "Fixed" {
		t.Errorf("Expected Kinstughna Fixed, got %s %s", name, ktype)
	}

	name, ktype = getKaranaDetails(57)
	if name != "Shakuni" || ktype != "Fixed" {
		t.Errorf("Expected Shakuni Fixed, got %s %s", name, ktype)
	}

	name, ktype = getKaranaDetails(1)
	if name != "Bava" || ktype != "Moving" {
		t.Errorf("Expected Bava Moving, got %s %s", name, ktype)
	}
}
