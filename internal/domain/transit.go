package domain

type TransitInput struct {
	BirthInput
	TransitDate string `json:"transit_date"` // YYYY-MM-DD
	TransitTime string `json:"transit_time"` // HH:MM:SS
}

// TransitResult wraps the transit calculation output.
// We reuse the verified TablesResult because it contains everything requested:
// planets, signs, houses, Ascendant, retrograde, nakshatra, pada, tooltip metadata (which are all in the tables).
type TransitResult struct {
	TransitDateUTC string       `json:"transit_date_utc"`
	JulianDay      float64      `json:"julian_day"`
	Ayanamsa       float64      `json:"ayanamsa"`
	Ascendant      float64      `json:"ascendant"`
	TransitData    TablesResult `json:"transit_data"`
}
