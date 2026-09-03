package ephemeris

import (
	"github.com/tejzpr/go-swisseph"
	"testing"
)

func TestGetAyanamsaMode(t *testing.T) {
	if GetAyanamsaMode("Lahiri") != swisseph.SidmLahiri {
		t.Error("Expected Lahiri mode")
	}
	if GetAyanamsaMode("Raman") != swisseph.SidmRaman {
		t.Error("Expected Raman mode")
	}
	if GetAyanamsaMode("Unknown") != swisseph.SidmLahiri { // Default
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
