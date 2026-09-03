package domain

// CharaKaraka represents a single planet's assigned Jaimini Karaka.
type CharaKaraka struct {
	Planet       string  `json:"planet"`
	Karaka       string  `json:"karaka"`
	Sign         string  `json:"sign"`
	Degree       int     `json:"degree"`
	Minute       int     `json:"minute"`
	Second       float64 `json:"second"`
	DegreeInSign float64 `json:"degree_in_sign"`
	Retrograde   bool    `json:"retrograde"`
}

// JaiminiKarakasResult contains the fully calculated Chara Karakas.
type JaiminiKarakasResult struct {
	CalculationTimeUTC string        `json:"calculation_time_utc"`
	Karakas            []CharaKaraka `json:"karakas"`
}
