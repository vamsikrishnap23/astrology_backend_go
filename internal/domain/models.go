package domain

import "time"

// BirthInput represents the birth details input by the user.
type BirthInput struct {
	Name         string  `json:"name"`
	DateOfBirth  string  `json:"date_of_birth"` // YYYY-MM-DD
	TimeOfBirth  string  `json:"time_of_birth"` // HH:MM:SS
	PlaceOfBirth string  `json:"place_of_birth"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Timezone     float64 `json:"timezone"` // Decimal offset from UTC
	Ayanamsa     string  `json:"ayanamsa"`
	HouseSystem  string  `json:"house_system"`
}

// CalculationConfig represents the resolved configuration.
type CalculationConfig struct {
	AyanamsaMode int
	HouseCode    byte
	NodeMode     int
}

// CalculationContext contains all derived base variables.
type CalculationContext struct {
	Input       BirthInput
	Config      CalculationConfig
	LocalTime   time.Time
	UTCTime     time.Time
	JulianDayUT float64
	JulianDayTT float64
	Ayanamsa    float64
}

// PlanetPosition represents the astronomical position of a body.
type PlanetPosition struct {
	Planet            string  `json:"planet"`
	TropicalLongitude float64 `json:"tropical_longitude"`
	SiderealLongitude float64 `json:"sidereal_longitude"`
	Latitude          float64 `json:"latitude,omitempty"`
	Distance          float64 `json:"distance,omitempty"`
	Speed             float64 `json:"speed"`
	Retrograde        bool    `json:"retrograde"`
	Sign              string  `json:"sign"`
	DegreeInSign      float64 `json:"degree_in_sign"`
	Degree            int     `json:"degree"`
	Minute            int     `json:"minute"`
	Second            float64 `json:"second"`
}

// HouseCusp represents a single house cusp.
type HouseCusp struct {
	HouseNumber int     `json:"house_number"`
	Longitude   float64 `json:"longitude"`
	Sign        string  `json:"sign"`
}

// ChartData is the overall chart response.
type ChartData struct {
	CalculationTimeUTC  string           `json:"calculation_time_utc"`
	JulianDay           float64          `json:"julian_day"`
	SelectedAyanamsa    string           `json:"selected_ayanamsa"`
	AyanamsaValue       float64          `json:"ayanamsa_value"`
	SelectedHouseSystem string           `json:"selected_house_system"`
	Ascendant           float64          `json:"ascendant"`
	MC                  float64          `json:"mc"`
	Planets             []PlanetPosition `json:"planets"`
	Houses              []HouseCusp      `json:"houses"`
}
