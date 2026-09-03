package domain

// BhavaChalitPlanet represents a planet's placement in the Bhava Chalit chart.
type BhavaChalitPlanet struct {
	PlanetName     string  `json:"planet_name"`
	HouseNumber    int     `json:"house_number"`
	Sign           string  `json:"sign"`
	Degree         int     `json:"degree"`
	Minute         int     `json:"minute"`
	Second         float64 `json:"second"`
	ExactLongitude float64 `json:"exact_longitude"`
}

// BhavaChalitHouse represents a house in the Bhava Chalit chart.
type BhavaChalitHouse struct {
	HouseNumber   int                 `json:"house_number"`
	CuspLongitude float64             `json:"cusp_longitude"`
	Sign          string              `json:"sign"`
	Degree        int                 `json:"degree"`
	Minute        int                 `json:"minute"`
	Second        float64             `json:"second"`
	Occupants     []BhavaChalitPlanet `json:"occupants"`
}

// BhavaChalitResult represents the full Bhava Chalit response.
type BhavaChalitResult struct {
	Ascendant float64            `json:"ascendant"`
	Houses    []BhavaChalitHouse `json:"houses"`
}
