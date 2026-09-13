package domain

type BTRInput struct {
	BirthInput
	Gender            string `json:"gender"` // "Male" or "Female"
	ReturnFullTable   bool   `json:"return_full_table"`
	ScanMinusMinutes  int    `json:"scan_minus_minutes,omitempty"`
	ScanPlusMinutes   int    `json:"scan_plus_minutes,omitempty"`
	SignTypeOverride  string `json:"sign_type_override,omitempty"`
	SunriseOverride   string `json:"sunrise_override,omitempty"`
	AscendantOverride string `json:"ascendant_override,omitempty"`
	StarLordOverride  string `json:"star_lord_override,omitempty"`
}

type BTRTableRow struct {
	No           int    `json:"no"`
	T1           int    `json:"t1"`
	T2           string `json:"t2"`
	Wed          string `json:"wed"`
	MonFri       string `json:"mon_fri"`
	SunTues      string `json:"sun_tues"`
	Sat          string `json:"sat"`
	Thur         string `json:"thur"`
	Movable      string `json:"movable"`
	Fixed        string `json:"fixed"`
	Dual         string `json:"dual"`
	VinodMovable string `json:"vinod_movable"`
	VinodFixed   string `json:"vinod_fixed"`
	VinodDual    string `json:"vinod_dual"`
}

type BTRCandidate struct {
	Rank              int    `json:"rank"`
	SuggestedTime     string `json:"suggested_time"`
	DifferenceMinutes int    `json:"difference_minutes"`
	Tatwa             string `json:"tatwa"`
	NadiRow           int    `json:"nadi_row"`
}

type BTRAnalysis struct {
	GenderMatch      bool   `json:"gender_match"`
	StarMatch        bool   `json:"star_match"`
	CalculatedTatwa  string `json:"calculated_tatwa"`
	CalculatedGender string `json:"calculated_gender"`
	CalculatedPlanet string `json:"calculated_planet"`
	ActualStarLord   string `json:"actual_star_lord"`
}

type BTRResult struct {
	InputTimeStatus         string         `json:"input_time_status"`
	InputAnalysis           BTRAnalysis    `json:"input_analysis"`
	SuggestedRectifications []BTRCandidate `json:"suggested_rectifications"`
	FullTable               []BTRTableRow  `json:"full_table,omitempty"`
}
