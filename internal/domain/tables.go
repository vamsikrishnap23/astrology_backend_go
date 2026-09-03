package domain

type TablePlanet struct {
	PlanetName     string  `json:"planet_name"`
	Sign           string  `json:"sign"`
	Degree         int     `json:"degree"`
	Minute         int     `json:"minute"`
	Second         float64 `json:"second"`
	ExactLongitude float64 `json:"exact_longitude"`
	Retrograde     bool    `json:"retrograde"`
	Speed          float64 `json:"speed"`
	HouseNumber    int     `json:"house_number"`
	Nakshatra      string  `json:"nakshatra"`
	NakshatraPada  int     `json:"nakshatra_pada"`
	SignLord       string  `json:"sign_lord"`
	NakshatraLord  string  `json:"nakshatra_lord"`
	SubLord        string  `json:"sub_lord"`
	SubSubLord     string  `json:"sub_sub_lord"`
	SubSubSubLord  string  `json:"sub_sub_sub_lord"`
}

type TableHouse struct {
	HouseNumber   int      `json:"house_number"`
	CuspLongitude float64  `json:"cusp_longitude"`
	Sign          string   `json:"sign"`
	Degree        int      `json:"degree"`
	Minute        int      `json:"minute"`
	Second        float64  `json:"second"`
	SignLord      string   `json:"sign_lord"`
	NakshatraLord string   `json:"nakshatra_lord"`
	SubLord       string   `json:"sub_lord"`
	SubSubLord    string   `json:"sub_sub_lord"`
	SubSubSubLord string   `json:"sub_sub_sub_lord"`
	Occupants     []string `json:"occupants"`
}

type TablesResult struct {
	PlanetaryTable []TablePlanet `json:"planetary_table"`
	HouseTable     []TableHouse  `json:"house_table"`
}
