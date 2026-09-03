package domain

// BinduContribution represents a single bindu contributed by a source planet to a specific sign
type BinduContribution struct {
	SourcePlanet string `json:"source_planet"`
	Value        int    `json:"value"` // Always 1 in traditional rules
}

// SignBindu represents the total bindus for a specific sign
type SignBindu struct {
	Sign          string              `json:"sign"`
	SignIndex     int                 `json:"sign_index"` // 1 for Aries, 12 for Pisces
	TotalBindus   int                 `json:"total_bindus"`
	Contributions []BinduContribution `json:"contributions,omitempty"` // Used only in BAV
}

// Bhinnashtakavarga represents the BAV for a specific planet
type Bhinnashtakavarga struct {
	Planet      string      `json:"planet"`
	TotalBindus int         `json:"total_bindus"`
	Signs       []SignBindu `json:"signs"` // exactly 12 elements
}

// AshtakavargaResult is the complete payload containing all BAVs and the SAV
type AshtakavargaResult struct {
	CalculationTimeUTC string              `json:"calculation_time_utc"`
	BAV                []Bhinnashtakavarga `json:"bav"`
	SAV                []SignBindu         `json:"sav"` // 12 elements
	TotalSAVBindus     int                 `json:"total_sav_bindus"`
}
