package domain

type ManglikAnalysis struct {
	IsPresent bool `json:"is_present"`
	House     int  `json:"house"`
}

type ManglikCancellation struct {
	Rule string `json:"rule"`
}

type ManglikResult struct {
	IsManglik     bool                       `json:"is_manglik"`
	Status        string                     `json:"status"`
	BaseAnalysis  map[string]ManglikAnalysis `json:"base_analysis"`
	Cancellations []ManglikCancellation      `json:"cancellations"`
	Remedies      []string                   `json:"remedies"`
}
