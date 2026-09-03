package domain

// PlanetSignificator represents the KP significators (A, B, C, D) for a single planet.
type PlanetSignificator struct {
	Planet string `json:"planet"`
	A      []int  `json:"a"`
	B      []int  `json:"b"`
	C      []int  `json:"c"`
	D      []int  `json:"d"`
}

// HouseSignificator represents the planets acting as significators (A, B, C, D) for a single house.
type HouseSignificator struct {
	House int      `json:"house"`
	A     []string `json:"a"`
	B     []string `json:"b"`
	C     []string `json:"c"`
	D     []string `json:"d"`
}

// KPSignificatorsResult contains both Planet and House views for KP Significators.
type KPSignificatorsResult struct {
	PlanetView []PlanetSignificator `json:"planet_view"`
	HouseView  []HouseSignificator  `json:"house_view"`
}

// StepDetails represents the planet and its significated houses for a single step in the 4-step theory.
type StepDetails struct {
	Planet string `json:"planet"`
	Houses []int  `json:"houses"`
}

// FourStepSignificator represents the 4-step signification for a single planet.
type FourStepSignificator struct {
	Planet        string      `json:"planet"`
	PlanetDetails StepDetails `json:"planet_details"`
	StarLord      StepDetails `json:"star_lord"`
	SubLord       StepDetails `json:"sub_lord"`
	StarLordOfSub StepDetails `json:"star_lord_of_sub"`
}

// FourStepSignificatorsResult contains the 4-step view for all planets.
type FourStepSignificatorsResult struct {
	FourStepView []FourStepSignificator `json:"four_step_view"`
}
