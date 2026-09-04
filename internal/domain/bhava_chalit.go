package domain

// BhavaChalitPlanet represents a planet's placement in the Bhava Chalit chart.
type BhavaChalitPlanet struct {
	Planet          string  `json:"planet"`
	SourceLongitude float64 `json:"source_longitude"`
	DivisionalSign  string  `json:"divisional_sign"` // Equivalent to Sign
	Degree          int     `json:"degree"`
	Minute          int     `json:"minute"`
	Second          float64 `json:"second"`
	Nakshatra       string  `json:"nakshatra"`
	NakshatraPada   int     `json:"nakshatra_pada"`
	NakshatraLord   string  `json:"nakshatra_lord"`
	SignLord        string  `json:"sign_lord,omitempty"`
	Retrograde      bool    `json:"retrograde"`
	HouseNumber     int     `json:"house_number"` // Keep this as extra info for Bhava Chalit
}

// BhavaChalitResult represents the full Bhava Chalit response.
type BhavaChalitResult struct {
	Ascendant float64             `json:"ascendant"`
	Planets   []BhavaChalitPlanet `json:"planets"`
	Houses    []HouseCusp         `json:"houses"` // Use standard HouseCusp model
}
