package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// Matchmaking always uses the Lahiri ayanamsa, whatever the request asks for.
func TestMatchmakingAlwaysUsesLahiri(t *testing.T) {
	if err := ephemeris.Init("../../../ephe_data"); err != nil {
		t.Fatalf("ephemeris init failed: %v", err)
	}
	defer ephemeris.Close()

	person := func(ayanamsa string) domain.BirthInput {
		return domain.BirthInput{
			Name: "P", DateOfBirth: "2005-11-23", TimeOfBirth: "15:35:00",
			Latitude: 16.3938, Longitude: 80.1522, Timezone: 5.5, Ayanamsa: ayanamsa, HouseSystem: "Placidus",
		}
	}
	run := func(groomAyanamsa string) domain.MatchmakingResult {
		body, _ := json.Marshal(domain.MatchInput{Groom: person(groomAyanamsa), Bride: person(groomAyanamsa)})
		req := httptest.NewRequest(http.MethodPost, "/api/vivaha-pontana", bytes.NewReader(body))
		rr := httptest.NewRecorder()
		MatchmakingHandler(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status %d: %s", rr.Code, rr.Body.String())
		}
		var res domain.MatchmakingResult
		if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
			t.Fatal(err)
		}
		return res
	}

	lahiri := run("Lahiri")
	other := run("Raman") // would shift the Moon by about 1.4 degrees if it were honoured
	if other.RuleSet.Ayanamsa != "Lahiri" {
		t.Errorf("rule_set ayanamsa = %q, want Lahiri", other.RuleSet.Ayanamsa)
	}
	if other.Groom.Moon.Longitude != lahiri.Groom.Moon.Longitude {
		t.Errorf("moon differs: %v vs %v", other.Groom.Moon.Longitude, lahiri.Groom.Moon.Longitude)
	}
}
