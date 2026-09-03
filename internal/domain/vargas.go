package domain

// VargaPlanet represents a planet's calculated position in a specific Varga.
type VargaPlanet struct {
	Planet          string  `json:"planet"`
	SourceLongitude float64 `json:"source_longitude"`
	DivisionalSign  string  `json:"divisional_sign"`
	Degree          int     `json:"degree"`
	Minute          int     `json:"minute"`
	Second          float64 `json:"second"`
	Nakshatra       string  `json:"nakshatra"`
	NakshatraPada   int     `json:"nakshatra_pada"`
	NakshatraLord   string  `json:"nakshatra_lord"`
	SignLord        string  `json:"sign_lord"`
	Retrograde      bool    `json:"retrograde"`
}

// VargaChart represents a complete divisional chart (like D1, D9).
type VargaChart struct {
	Division  int           `json:"division"`
	Name      string        `json:"name"`
	Ascendant VargaPlanet   `json:"ascendant"`
	Planets   []VargaPlanet `json:"planets"`
	Houses    []HouseCusp   `json:"houses,omitempty"`
}

// VargasResult contains the complete set of requested Varga charts.
type VargasResult struct {
	Vargas []VargaChart `json:"vargas"`
}
