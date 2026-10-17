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
	Rank            int     `json:"rank"` // 1 = strongest by strength_ratio
}

// ShadbalaTableRow is one line of a table: a component and its value for each planet.
type ShadbalaTableRow struct {
	Component string             `json:"component"`
	Label     string             `json:"label"`
	Values    map[string]float64 `json:"values"` // by planet name
}

// ShadbalaTable is a component of Shadbala laid out as in Raman's worked tables: one row per part and
// one value per planet, with the total of the parts for each planet.
type ShadbalaTable struct {
	Title string             `json:"title"`
	Rows  []ShadbalaTableRow `json:"rows"`
	Total map[string]float64 `json:"total,omitempty"` // per planet; absent where parts are not added
}

// ShadbalaPindaTable is the final table: the six balas, their total in shashtiamsas and rupas, the
// requirement, the ratio and the rank of each planet (Raman examples 56-57).
type ShadbalaPindaTable struct {
	Title                string             `json:"title"`
	Rows                 []ShadbalaTableRow `json:"rows"`
	Total                map[string]float64 `json:"total"` // shashtiamsas
	Rupas                map[string]float64 `json:"rupas"`
	MinimumRequiredRupas map[string]float64 `json:"minimum_required_rupas"`
	StrengthRatio        map[string]float64 `json:"strength_ratio"`
	Rank                 map[string]int     `json:"rank"`
}

// ShadbalaTables are the tables of the response: Planets gives the planets in order.
type ShadbalaTables struct {
	Planets        []string           `json:"planets"`
	SthanaBala     ShadbalaTable      `json:"sthana_bala"`
	DigBala        ShadbalaTable      `json:"dig_bala"`
	KalaBala       ShadbalaTable      `json:"kala_bala"`
	CheshtaBala    ShadbalaTable      `json:"cheshta_bala"`
	NaisargikaBala ShadbalaTable      `json:"naisargika_bala"`
	DrikBala       ShadbalaTable      `json:"drik_bala"`
	ShadbalaPinda  ShadbalaPindaTable `json:"shadbala_pinda"`

	IshtaKashtaPhala ShadbalaTable `json:"ishta_kashta_phala"`
}

// BhavaBalaHouse is the strength of one house (bhava), Raman chapter IX.
type BhavaBalaHouse struct {
	House             int     `json:"house"`        // 1..12
	Name              string  `json:"name"`         // Thanu, Dhana, Bhratru, Matru, Putra, Satru, Kalatra, Ayur, Bhagya, Karma, Labha, Vyaya
	BhavaMadhya       float64 `json:"bhava_madhya"` // sidereal longitude of the middle of the house (Sripati)
	Sign              string  `json:"sign"`
	Lord              string  `json:"lord"`               // lord of the sign of the bhava madhya
	SignGroup         string  `json:"sign_group"`         // nara, jalachara, chathushpada or keeta
	BhavadhipathiBala float64 `json:"bhavadhipathi_bala"` // Shadbala Pinda of the lord, in shashtiamsas
	BhavaDigbala      float64 `json:"bhava_digbala"`
	SubhaDrishtiBala  float64 `json:"subha_drishti_bala"`  // benefic aspects on the bhava madhya
	AsubhaDrishtiBala float64 `json:"asubha_drishti_bala"` // malefic aspects
	BhavaDrishtiBala  float64 `json:"bhava_drishti_bala"`  // subha - asubha
	TotalBhavaBala    float64 `json:"total_bhava_bala"`    // shashtiamsas
	Rupas             float64 `json:"rupas"`
	Rank              int     `json:"rank"` // 1 = strongest house
}

// BhavaBalaResult is the strength of the twelve houses.
type BhavaBalaResult struct {
	Bhavas         []BhavaBalaHouse `json:"bhavas"`
	StrongestHouse int              `json:"strongest_house"`
	WeakestHouse   int              `json:"weakest_house"`
}

// IshtaKashtaPlanet is the benefic (ishta) and malefic (kashta) tendency of a planet, Raman chapter X.
type IshtaKashtaPlanet struct {
	Planet      string  `json:"planet"`
	UchchaBala  float64 `json:"uchcha_bala"`
	ChestaBala  float64 `json:"chesta_bala"` // the Sun's and Moon's come from Raman art. 136-137 and are not part of the Shadbala total
	IshtaPhala  float64 `json:"ishta_phala"`
	KashtaPhala float64 `json:"kashta_phala"`
	Inclination string  `json:"inclination"` // ishta if the benefic tendency is greater, otherwise kashta
}

type ShadbalaResult struct {
	RuleSet            string           `json:"rule_set"` // which book's rules were followed
	Unit               string           `json:"unit"`
	CalculationTimeUTC string           `json:"calculation_time_utc"`
	Planets            []PlanetShadbala `json:"planets"`
	Tables             ShadbalaTables   `json:"tables"`

	IshtaKashtaPhala []IshtaKashtaPlanet `json:"ishta_kashta_phala"`
	BhavaBala        BhavaBalaResult     `json:"bhava_bala"`
}
