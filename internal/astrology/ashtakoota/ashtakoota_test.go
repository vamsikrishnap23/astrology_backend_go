package ashtakoota

import (
	"github.com/tejzpr/go-swisseph"
	"path/filepath"
	"testing"

	"github.com/vamsi/astrology_backend_go/internal/astronomy/ephemeris"
	astronomyTime "github.com/vamsi/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

func TestAshtakoota(t *testing.T) {
	absPath, _ := filepath.Abs("../../../ephe_data")
	ephemeris.Init(absPath)

	utcTime1, _ := astronomyTime.ParseLocalToUTC("2005-11-23", "15:35:00", 5.5)
	jd1 := astronomyTime.UTCToJulianDay(utcTime1)
	swisseph.SetEphePath(absPath)
	swisseph.SetSidMode(swisseph.SidmLahiri, 0, 0)

	ctx1 := domain.CalculationContext{
		Input:       domain.BirthInput{Name: "Groom", Ayanamsa: "Lahiri"},
		UTCTime:     utcTime1,
		JulianDayUT: jd1,
		Ayanamsa:    swisseph.GetAyanamsaUT(jd1),
	}

	utcTime2, _ := astronomyTime.ParseLocalToUTC("2002-05-14", "10:00:00", 5.5)
	jd2 := astronomyTime.UTCToJulianDay(utcTime2)
	ctx2 := domain.CalculationContext{
		Input:       domain.BirthInput{Name: "Bride", Ayanamsa: "Lahiri"},
		UTCTime:     utcTime2,
		JulianDayUT: jd2,
		Ayanamsa:    swisseph.GetAyanamsaUT(jd2),
	}

	res, err := CalculateMatch(&ctx1, &ctx2)
	if err != nil {
		t.Fatalf("Failed to calculate ashtakoota: %v", err)
	}

	if res.Summary.Maximum != 36 {
		t.Errorf("Expected max 36, got %f", res.Summary.Maximum)
	}

	if res.Kootas.Nadi.Maximum != 8 {
		t.Errorf("Expected Nadi max 8")
	}
}
