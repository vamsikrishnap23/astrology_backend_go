package ephemeris

import (
	"github.com/mshafiee/swephgo"
	"testing"
)

func TestGetAyanamsaMode(t *testing.T) {
	if GetAyanamsaMode("Lahiri") != swephgo.SeSidmLahiri {
		t.Error("Expected Lahiri mode")
	}
	if GetAyanamsaMode("Raman") != swephgo.SeSidmRaman {
		t.Error("Expected Raman mode")
	}
	if GetAyanamsaMode("Unknown") != swephgo.SeSidmLahiri { // Default
		t.Error("Expected fallback to Lahiri")
	}
}

func TestGetHouseSystemCode(t *testing.T) {
	if GetHouseSystemCode("Placidus") != 'P' {
		t.Error("Expected P")
	}
	if GetHouseSystemCode("Whole Sign") != 'W' {
		t.Error("Expected W")
	}
}
