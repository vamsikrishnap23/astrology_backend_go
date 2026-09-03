package domain

type MatchInput struct {
	Groom BirthInput `json:"groom"`
	Bride BirthInput `json:"bride"`
}

type MoonDetails struct {
	Longitude     float64 `json:"longitude"`
	Sign          string  `json:"sign"`
	Degree        float64 `json:"degree"`
	Nakshatra     string  `json:"nakshatra"`
	Pada          int     `json:"pada"`
	NakshatraLord string  `json:"nakshatra_lord"`
	RashiLord     string  `json:"rashi_lord"`
}

type PersonMatchDetails struct {
	Name string      `json:"name"`
	Moon MoonDetails `json:"moon"`
}

type KootaResult struct {
	Score       float64 `json:"score"`
	Maximum     float64 `json:"maximum"`
	GroomValue  string  `json:"groom_value,omitempty"`
	BrideValue  string  `json:"bride_value,omitempty"`
	Explanation string  `json:"explanation"`
}

type TaraDirection struct {
	NakshatraCount int     `json:"nakshatra_count"`
	Tara           string  `json:"tara"`
	Auspicious     bool    `json:"auspicious"`
	Score          float64 `json:"score"`
}

type TaraKootaResult struct {
	Score        float64       `json:"score"`
	Maximum      float64       `json:"maximum"`
	GroomToBride TaraDirection `json:"groom_to_bride"`
	BrideToGroom TaraDirection `json:"bride_to_groom"`
	Explanation  string        `json:"explanation"`
}

type GrahaMaitriResult struct {
	Score                            float64 `json:"score"`
	Maximum                          float64 `json:"maximum"`
	GroomRashi                       string  `json:"groom_rashi"`
	BrideRashi                       string  `json:"bride_rashi"`
	GroomRashiLord                   string  `json:"groom_rashi_lord"`
	BrideRashiLord                   string  `json:"bride_rashi_lord"`
	GroomLordRelationshipToBrideLord string  `json:"groom_lord_relationship_to_bride_lord"`
	BrideLordRelationshipToGroomLord string  `json:"bride_lord_relationship_to_groom_lord"`
	Explanation                      string  `json:"explanation"`
}

type BhakootResult struct {
	RawScore             float64 `json:"raw_score"`
	EffectiveScore       float64 `json:"effective_score"`
	Maximum              float64 `json:"maximum"`
	GroomRashi           string  `json:"groom_rashi"`
	BrideRashi           string  `json:"bride_rashi"`
	GroomToBrideDistance int     `json:"groom_to_bride_distance"`
	BrideToGroomDistance int     `json:"bride_to_groom_distance"`
	Relationship         string  `json:"relationship"`
	Dosha                bool    `json:"dosha"`
	CancellationApplied  bool    `json:"cancellation_applied"`
	CancellationReason   string  `json:"cancellation_reason,omitempty"`
	Explanation          string  `json:"explanation"`
}

type NadiResult struct {
	RawScore            float64 `json:"raw_score"`
	EffectiveScore      float64 `json:"effective_score"`
	Maximum             float64 `json:"maximum"`
	GroomNadi           string  `json:"groom_nadi"`
	BrideNadi           string  `json:"bride_nadi"`
	SameNadi            bool    `json:"same_nadi"`
	Dosha               bool    `json:"dosha"`
	CancellationApplied bool    `json:"cancellation_applied"`
	CancellationReason  string  `json:"cancellation_reason,omitempty"`
	Explanation         string  `json:"explanation"`
}

type Kootas struct {
	Varna       KootaResult       `json:"varna"`
	Vashya      KootaResult       `json:"vashya"`
	Tara        TaraKootaResult   `json:"tara"`
	Yoni        KootaResult       `json:"yoni"`
	GrahaMaitri GrahaMaitriResult `json:"graha_maitri"`
	Gana        KootaResult       `json:"gana"`
	Bhakoot     BhakootResult     `json:"bhakoot"`
	Nadi        NadiResult        `json:"nadi"`
}

type MatchSummary struct {
	RawTotal                float64 `json:"raw_total"`
	EffectiveTotal          float64 `json:"effective_total"`
	Maximum                 float64 `json:"maximum"`
	Percentage              float64 `json:"percentage"`
	TraditionalThresholdMet bool    `json:"traditional_threshold_met"`
}

type RuleSetInfo struct {
	Name     string `json:"name"`
	Ayanamsa string `json:"ayanamsa"`
}

type AshtakootaResult struct {
	RuleSet RuleSetInfo        `json:"rule_set"`
	Groom   PersonMatchDetails `json:"groom"`
	Bride   PersonMatchDetails `json:"bride"`
	Kootas  Kootas             `json:"kootas"`
	Summary MatchSummary       `json:"summary"`
}
