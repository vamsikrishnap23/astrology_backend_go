package domain

type UpcomingTransitsInput struct {
	BirthInput
	SearchStartDate string `json:"search_start_date,omitempty"` // YYYY-MM-DD
	SearchStartTime string `json:"search_start_time,omitempty"` // HH:MM:SS
}

type UpcomingTransitPlanet struct {
	Planet             string `json:"planet"`
	DestinationSign    string `json:"destination_sign"`
	TransitionDateTime string `json:"transition_datetime"` // ISO-8601 UTC string
}

type UpcomingTransitsResult struct {
	SearchStartDateUTC string                  `json:"search_start_date_utc"`
	Transits           []UpcomingTransitPlanet `json:"transits"`
}
