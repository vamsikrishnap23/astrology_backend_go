package domain

type BTRInput struct {
	Name             string  `json:"name"`
	DateOfBirth      string  `json:"date_of_birth"`
	TimeOfBirth      string  `json:"time_of_birth"`
	PlaceOfBirth     string  `json:"place_of_birth"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Timezone         float64 `json:"timezone"`
	Ayanamsa         string  `json:"ayanamsa"`
	HouseSystem      string  `json:"house_system"`
	Gender           string  `json:"gender"`
	ScanMinusMinutes int     `json:"scan_minus_minutes"`
	ScanPlusMinutes  int     `json:"scan_plus_minutes"`
}

type BTRAstronomicalContext struct {
	Sunrise       string `json:"sunrise"`
	LMT           string `json:"lmt"`
	DayLord       string `json:"day_lord"`
	MoonNakshatra string `json:"moon_nakshatra"`
	StarLord      string `json:"star_lord"`
}

type BTRScanConfig struct {
	MinusMinutes int `json:"minus_minutes"`
	PlusMinutes  int `json:"plus_minutes"`
}

type TatwaDetails struct {
	MainTatwa     string  `json:"main_tatwa"`
	AntarTatwa    string  `json:"antar_tatwa"`
	MainDuration  float64 `json:"main_duration_mins"`
	AntarDuration float64 `json:"antar_duration_mins"`
	IsAroha       bool    `json:"is_aroha"`
	CycleIndex    int     `json:"cycle_index"`
}

type NadiRow struct {
	RowNumber          int    `json:"row_number,omitempty"`
	NinetyMinutePlanet string `json:"ninety_minute_planet"`
	AscendantSignType  string `json:"ascendant_sign_type"`
	StarLord           string `json:"star_lord"`
	Gender             string `json:"gender"`
	ExpectedTatwa      string `json:"expected_tatwa"`
	ExpectedAntarTatwa string `json:"expected_antar_tatwa"`
}

type BTRCandidate struct {
	CandidateTimeUTC   string       `json:"candidate_time_utc"`
	CandidateTimeLocal string       `json:"candidate_time_local"`
	OffsetMinutes      float64      `json:"offset_minutes"`
	Tatwa              TatwaDetails `json:"tatwa"`
	AscendantDegree    float64      `json:"ascendant_degree"`
	AscendantSign      string       `json:"ascendant_sign"`
	AscendantSignType  string       `json:"ascendant_sign_type"`
	NinetyMinutePlanet string       `json:"ninety_minute_planet"`
	MoonNakshatra      string       `json:"moon_nakshatra"`
	StarLord           string       `json:"star_lord"`
	Gender             string       `json:"gender"`
	NadiRow            NadiRow      `json:"nadi_row"`
	MatchStatus        string       `json:"match_status"` // "match", "no_match", "unresolved"
	MatchExplanation   string       `json:"match_explanation"`
}

type BTRResult struct {
	Status              string                 `json:"status"` // e.g., "partially_implemented"
	RecordedTime        string                 `json:"recorded_time"`
	AstronomicalContext BTRAstronomicalContext `json:"astronomical_context"`
	Scan                BTRScanConfig          `json:"scan"`
	Candidates          []BTRCandidate         `json:"candidates"`
	Result              map[string]interface{} `json:"result"`
	Warnings            []string               `json:"warnings"`
}
