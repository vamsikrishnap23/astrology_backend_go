package domain

import "time"

// Sookshma represents the 4th level of Vimshottari Dasha.
type Sookshma struct {
	Lord      string    `json:"lord"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// Pratyantardasha represents the 3rd level of Vimshottari Dasha.
type Pratyantardasha struct {
	Lord      string     `json:"lord"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
	Sookshma  []Sookshma `json:"sookshma"`
}

// Antardasha represents the 2nd level of Vimshottari Dasha (Bhukti).
type Antardasha struct {
	Lord            string            `json:"lord"`
	StartDate       time.Time         `json:"start_date"`
	EndDate         time.Time         `json:"end_date"`
	Pratyantardasha []Pratyantardasha `json:"pratyantardasha"`
}

// Mahadasha represents the 1st level of Vimshottari Dasha.
type Mahadasha struct {
	Lord       string       `json:"lord"`
	StartDate  time.Time    `json:"start_date"`
	EndDate    time.Time    `json:"end_date"`
	Antardasha []Antardasha `json:"antardasha"`
}

// VimshottariDashaResult encapsulates the complete 120-year Dasha sequence from the time of birth.
type VimshottariDashaResult struct {
	BalanceYears float64     `json:"balance_years"`
	Mahadasha    []Mahadasha `json:"mahadasha"`
}
