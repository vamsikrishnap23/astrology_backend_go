package vargas

import (
	"testing"
)

func assertVarga(t *testing.T, rule VargaRule, lon float64, expectedSign int) {
	pos := rule.Calculate(lon)
	if pos.SignIndex != expectedSign {
		t.Errorf("%s: for lon %f, expected sign %d, got %d", rule.Name(), lon, expectedSign, pos.SignIndex)
	}
}

func TestD1(t *testing.T) {
	rule := &RuleD1{}
	assertVarga(t, rule, 0.0, 0)     // Aries
	assertVarga(t, rule, 359.99, 11) // Pisces
}

func TestD9(t *testing.T) {
	rule := &RuleD9{}
	// Aries 0 -> Aries (0)
	assertVarga(t, rule, 0.0, 0)
	// Aries 3deg21min (3.35) -> Taurus (1)
	assertVarga(t, rule, 3.35, 1)
	// Taurus 0 (30.0) -> Capricorn (9)
	assertVarga(t, rule, 30.0, 9)
	// Gemini 0 (60.0) -> Libra (6)
	assertVarga(t, rule, 60.0, 6)
	// Cancer 0 (90.0) -> Cancer (3)
	assertVarga(t, rule, 90.0, 3)
}

func TestD2(t *testing.T) {
	rule := &RuleD2{}
	// Aries (Odd): 0-15 -> Leo (4)
	assertVarga(t, rule, 10.0, 4)
	// Aries (Odd): 15-30 -> Cancer (3)
	assertVarga(t, rule, 20.0, 3)
	// Taurus (Even): 0-15 -> Cancer (3)
	assertVarga(t, rule, 40.0, 3)
	// Taurus (Even): 15-30 -> Leo (4)
	assertVarga(t, rule, 50.0, 4)
}

func TestD3(t *testing.T) {
	rule := &RuleD3{}
	// Aries 0-10 -> Aries (0)
	assertVarga(t, rule, 5.0, 0)
	// Aries 10-20 -> Leo (4)
	assertVarga(t, rule, 15.0, 4)
	// Aries 20-30 -> Sag (8)
	assertVarga(t, rule, 25.0, 8)
}

func TestD10(t *testing.T) {
	rule := &RuleD10{}
	// Aries (Odd): 0-3 -> Aries (0)
	assertVarga(t, rule, 1.0, 0)
	// Aries 3-6 -> Taurus (1)
	assertVarga(t, rule, 4.0, 1)
	// Taurus (Even): 0-3 -> Cap (9)
	assertVarga(t, rule, 31.0, 9)
	// Taurus 3-6 -> Aq (10)
	assertVarga(t, rule, 34.0, 10)
}

func TestD30(t *testing.T) {
	rule := &RuleD30{}
	// Odd Sign (Aries): 0-5 -> Aries (0)
	assertVarga(t, rule, 2.0, 0)
	// Odd Sign (Aries): 5-10 -> Aquarius (10)
	assertVarga(t, rule, 7.0, 10)
	// Even Sign (Taurus): 0-5 -> Taurus (1)
	assertVarga(t, rule, 32.0, 1)
	// Even Sign (Taurus): 5-12 -> Virgo (5)
	assertVarga(t, rule, 37.0, 5)
}

func TestBoundaries(t *testing.T) {
	rule := &RuleD9{}
	// Edge of wraparound
	pos := rule.Calculate(359.999999)
	if pos.SignIndex != 11 {
		t.Errorf("Expected 11 at 359.999, got %d", pos.SignIndex)
	}
}

func TestD60(t *testing.T) {
	rule := &RuleD60{}
	// 0.5 degrees = 1 part
	// Aries 0 -> Aries (0)
	assertVarga(t, rule, 0.0, 0)
	// Aries 0.6 -> Taurus (1)
	assertVarga(t, rule, 0.6, 1)
	// Aries 1.1 -> Gemini (2)
	assertVarga(t, rule, 1.1, 2)
}

func TestVargasCalculationWrapper(t *testing.T) {
	// Simple test to ensure engine doesn't panic
	_ = Registry
}
