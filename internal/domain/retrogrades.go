package domain

type RetrogradesInput struct {
	StartDate string  `json:"start_date"` // YYYY-MM-DD
	EndDate   string  `json:"end_date"`   // YYYY-MM-DD
	Timezone  float64 `json:"timezone"`
}

type RetrogradePhase struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type RetrogradesResult struct {
	StartDate   string                       `json:"start_date"`
	EndDate     string                       `json:"end_date"`
	Timezone    float64                      `json:"timezone"`
	Retrogrades map[string][]RetrogradePhase `json:"retrogrades"`
}
