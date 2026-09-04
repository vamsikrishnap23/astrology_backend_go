package domain

// CharaKaraka represents a single planet's assigned Jaimini Karaka.
type CharaKaraka struct {
	Planet          string  `json:"planet"`
	Karaka          string  `json:"karaka"`
	SourceLongitude float64 `json:"source_longitude"`
	DivisionalSign  string  `json:"divisional_sign"` // Was Sign
	Degree          int     `json:"degree"`
	Minute          int     `json:"minute"`
	Second          float64 `json:"second"`
	Nakshatra       string  `json:"nakshatra"`
	NakshatraPada   int     `json:"nakshatra_pada"`
	NakshatraLord   string  `json:"nakshatra_lord"`
	DegreeInSign    float64 `json:"degree_in_sign"`
	Retrograde      bool    `json:"retrograde"`
}

// JaiminiKarakasResult contains the fully calculated Chara Karakas.
type JaiminiKarakasResult struct {
	CalculationTimeUTC string        `json:"calculation_time_utc"`
	Planets            []CharaKaraka `json:"planets"` // Was Karakas
}
