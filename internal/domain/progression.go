package domain

// ProgressionInput adds progression-specific parameters to the standard BirthInput.
type ProgressionInput struct {
	BirthInput
	ProgressionDate string `json:"progression_date"` // The target date (e.g. "2040-05-15") to progress to
}

// ProgressedPlanet represents a planet's progressed position.
type ProgressedPlanet struct {
	PlanetPosition
	// Can add more fields if needed, but PlanetPosition holds exact pos, sign, speed, retro.
}

// ProgressionResult represents the final progressed chart.
type ProgressionResult struct {
	NatalDateUTC          string           `json:"natal_date_utc"`
	TargetProgressionDate string           `json:"target_progression_date"`
	AgeInYears            float64          `json:"age_in_years"`
	ProgressedDateUTC     string           `json:"progressed_date_utc"`
	ProgressedJulianDay   float64          `json:"progressed_julian_day"`
	ProgressedAyanamsa    float64          `json:"progressed_ayanamsa"`
	Ascendant             float64          `json:"ascendant"`
	MC                    float64          `json:"mc"`
	ProgressedPlanets     []PlanetPosition `json:"progressed_planets"`
	ProgressedHouses      []HouseCusp      `json:"progressed_houses"`
}
