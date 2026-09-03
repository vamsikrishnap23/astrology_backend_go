package domain

// RulingPlanetSource defines the origin of the ruling planet in the KP hierarchy.
type RulingPlanetSource string

const (
	SourceAscendantStarLord RulingPlanetSource = "Ascendant Star Lord"
	SourceAscendantSignLord RulingPlanetSource = "Ascendant Sign Lord"
	SourceAscendantSubLord  RulingPlanetSource = "Ascendant Sub Lord"
	SourceMoonStarLord      RulingPlanetSource = "Moon Star Lord"
	SourceMoonSignLord      RulingPlanetSource = "Moon Sign Lord"
	SourceMoonSubLord       RulingPlanetSource = "Moon Sub Lord"
	SourceDayLord           RulingPlanetSource = "Day Lord"
	SourceNodeAgent         RulingPlanetSource = "Node Agent"
)

// RulingPlanet represents a single ruling planet and its source.
type RulingPlanet struct {
	Planet string             `json:"planet"`
	Source RulingPlanetSource `json:"source"`
	// If this is a Node Agent, AgentFor contains the planet it is representing.
	AgentFor string `json:"agent_for,omitempty"`
}

// KPRulingPlanetsResult represents the set of ruling planets.
type KPRulingPlanetsResult struct {
	RulingPlanets []RulingPlanet `json:"ruling_planets"`
}
