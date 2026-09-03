package domain

type SthanaBala struct {
	UchchaBala       float64 `json:"uchcha_bala"`
	SaptavargajaBala float64 `json:"saptavargaja_bala"`
	OjayugmaBala     float64 `json:"ojayugma_bala"`
	KendradiBala     float64 `json:"kendradi_bala"`
	DrekkanaBala     float64 `json:"drekkana_bala"`
	Total            float64 `json:"total"`
}

type KalaBala struct {
	NathonnathaBala float64 `json:"nathonnatha_bala"`
	PakshaBala      float64 `json:"paksha_bala"`
	TribhagaBala    float64 `json:"tribhaga_bala"`
	VarshaBala      float64 `json:"varsha_bala"`
	MasaBala        float64 `json:"masa_bala"`
	DinaBala        float64 `json:"dina_bala"`
	HoraBala        float64 `json:"hora_bala"`
	AyanaBala       float64 `json:"ayana_bala"`
	YuddhaBala      float64 `json:"yuddha_bala"`
	Total           float64 `json:"total"`
}

type PlanetShadbala struct {
	Planet         string     `json:"planet"`
	SthanaBala     SthanaBala `json:"sthana_bala"`
	DigBala        float64    `json:"dig_bala"`
	KalaBala       KalaBala   `json:"kala_bala"`
	CheshtaBala    float64    `json:"cheshta_bala"`
	NaisargikaBala float64    `json:"naisargika_bala"`
	DrikBala       float64    `json:"drik_bala"`

	TotalShadbala   float64 `json:"total_shadbala"`
	Rupas           float64 `json:"rupas"`
	MinimumRequired float64 `json:"minimum_required"`
	StrengthRatio   float64 `json:"strength_ratio"`
	MeetsMinimum    bool    `json:"meets_minimum"`
}

type ShadbalaResult struct {
	CalculationTimeUTC string           `json:"calculation_time_utc"`
	Planets            []PlanetShadbala `json:"planets"`
}
