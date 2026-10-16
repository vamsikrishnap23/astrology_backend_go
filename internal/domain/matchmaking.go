package domain

type MatchInput struct {
	Groom BirthInput `json:"groom"`
	Bride BirthInput `json:"bride"`
}

// MatchmakingResult is the /api/vivaha-pontana response. It is organised in sections:
//
//	rule_set    which rule set and ayanamsa were used
//	groom/bride each person: name, Moon details, individual chart risks
//	lagna       ascendant compatibility (informational)
//	ashtakoota  the eight scored kootas (36 points) and the summary
//	checks      the other compatibility checks (not scored): each one satisfied or not
type MatchmakingResult struct {
	RuleSet    RuleSetInfo        `json:"rule_set"`
	Groom      PersonMatchDetails `json:"groom"`
	Bride      PersonMatchDetails `json:"bride"`
	Lagna      LagnaCompatibility `json:"lagna"`
	Ashtakoota Ashtakoota         `json:"ashtakoota"`
	Checks     Checks             `json:"checks"`
}

type RuleSetInfo struct {
	Name     string `json:"name"`
	Ayanamsa string `json:"ayanamsa"`
}

// ---- people ----

type MoonDetails struct {
	Longitude     float64 `json:"longitude"`
	Sign          string  `json:"sign"`
	Degree        float64 `json:"degree"`
	Nakshatra     string  `json:"nakshatra"`
	Pada          int     `json:"pada"`
	NakshatraLord string  `json:"nakshatra_lord"`
	RashiLord     string  `json:"rashi_lord"`
}

type IndividualRiskCheck struct {
	HasRisk     bool     `json:"has_risk"`
	RiskFactors []string `json:"risk_factors"`
}

type IndividualRisks struct {
	VaidhavyaDosham   IndividualRiskCheck `json:"vaidhavya_dosham"`
	DwikalatraDosham  IndividualRiskCheck `json:"dwikalatra_dosham"`
	Napumsakatvam     IndividualRiskCheck `json:"napumsakatvam"`
	DampatyaMarakatva IndividualRiskCheck `json:"dampatya_marakatva"`
}

type PersonMatchDetails struct {
	Name  string          `json:"name"`
	Moon  MoonDetails     `json:"moon"`
	Risks IndividualRisks `json:"risks"`
}

// ---- lagna ----

// LagnaCompatibility compares the two ascendants. Distance is counted from the bride's lagna to the
// groom's. Quality is Auspicious (1, 3, 5, 9, 11), Moderate (2, 7) or Inauspicious (4, 6, 8, 10, 12).
type LagnaCompatibility struct {
	GroomLagna        string `json:"groom_lagna"`
	BrideLagna        string `json:"bride_lagna"`
	GroomLagnaLord    string `json:"groom_lagna_lord"`
	BrideLagnaLord    string `json:"bride_lagna_lord"`
	Distance          int    `json:"distance"`
	Quality           string `json:"quality"`
	IsAuspicious      bool   `json:"is_auspicious"`
	LordsRelationship string `json:"lords_relationship"`
}

// ---- ashtakoota (scored) ----

type Ashtakoota struct {
	Kootas  Kootas       `json:"kootas"`
	Summary MatchSummary `json:"summary"`
}

type Kootas struct {
	Varna       Koota `json:"varna"`
	Vashya      Koota `json:"vashya"`
	Tara        Koota `json:"tara"`
	Yoni        Koota `json:"yoni"`
	GrahaMaitri Koota `json:"graha_maitri"`
	Gana        Koota `json:"gana"`
	Bhakoot     Koota `json:"bhakoot"`
	Nadi        Koota `json:"nadi"`
}

// Koota is one scored koota. Every koota has the same shape.
//
//	raw_score   the score before any cancellation
//	score       the score counted in the total (equals raw_score unless a cancellation applied)
//	groom/bride the value each person was classified as (varna, rasi, nakshatra, animal, lord, gana, nadi)
//	dosha       true only for the kootas that have a named dosha (bhakoot, nadi) when it is present
//	suggestion  a rule that may apply but is deliberately not counted in the score
//	details     extra fields specific to the koota (tara, graha_maitri, bhakoot)
type Koota struct {
	Name                string  `json:"name"`
	Maximum             float64 `json:"maximum"`
	RawScore            float64 `json:"raw_score"`
	Score               float64 `json:"score"`
	Groom               string  `json:"groom"`
	Bride               string  `json:"bride"`
	Dosha               bool    `json:"dosha"`
	CancellationApplied bool    `json:"cancellation_applied"`
	CancellationReason  string  `json:"cancellation_reason,omitempty"`
	Suggestion          string  `json:"suggestion,omitempty"`
	Explanation         string  `json:"explanation"`
	Details             any     `json:"details,omitempty"`
}

// TaraDetails are the two directional counts behind the Tara koota.
type TaraDetails struct {
	GroomToBride TaraDirection `json:"groom_to_bride"`
	BrideToGroom TaraDirection `json:"bride_to_groom"`
}

type TaraDirection struct {
	NakshatraCount int     `json:"nakshatra_count"`
	Tara           string  `json:"tara"`
	Auspicious     bool    `json:"auspicious"`
	Score          float64 `json:"score"`
}

// GrahaMaitriDetails: the koota's groom/bride values are the two rashi lords.
type GrahaMaitriDetails struct {
	GroomRashi                       string `json:"groom_rashi"`
	BrideRashi                       string `json:"bride_rashi"`
	GroomLordRelationshipToBrideLord string `json:"groom_lord_relationship_to_bride_lord"`
	BrideLordRelationshipToGroomLord string `json:"bride_lord_relationship_to_groom_lord"`
}

// BhakootDetails: the koota's groom/bride values are the two Moon signs.
type BhakootDetails struct {
	GroomToBrideDistance int    `json:"groom_to_bride_distance"`
	BrideToGroomDistance int    `json:"bride_to_groom_distance"`
	Relationship         string `json:"relationship"`
}

type MatchSummary struct {
	RawTotal                float64 `json:"raw_total"`
	EffectiveTotal          float64 `json:"effective_total"`
	Maximum                 float64 `json:"maximum"`
	Percentage              float64 `json:"percentage"`
	TraditionalThresholdMet bool    `json:"traditional_threshold_met"`
}

// ---- checks (not scored) ----

// Checks are compatibility checks that do not add to the 36 points. Each is a Check:
// satisfied is true when there is no outstanding problem.
type Checks struct {
	Rajju        Check `json:"rajju"`
	Vedha        Check `json:"vedha"`
	Mahendra     Check `json:"mahendra"`
	StreeDeergha Check `json:"stree_deergha"`
	MrugaVairam  Check `json:"mruga_vairam"`
	SpecialTara  Check `json:"special_tara"`
	KalaSarpa    Check `json:"kala_sarpa"`
}

// Check is one non-scored check. groom/bride describe what each person was classified as
// (rajju group, nakshatra, yoni animal, or Present/Absent for kala sarpa).
type Check struct {
	Name                string `json:"name"`
	Satisfied           bool   `json:"satisfied"`
	Groom               string `json:"groom,omitempty"`
	Bride               string `json:"bride,omitempty"`
	CancellationApplied bool   `json:"cancellation_applied,omitempty"`
	CancellationReason  string `json:"cancellation_reason,omitempty"`
	Explanation         string `json:"explanation"`
	Details             any    `json:"details,omitempty"`
}

// CountDetails: the count from the bride's nakshatra to the groom's (inclusive).
type CountDetails struct {
	BrideToGroomCount int `json:"bride_to_groom_count"`
}

type StreeDeerghaDetails struct {
	BrideToGroomCount           int    `json:"bride_to_groom_count"`
	Navaka                      int    `json:"navaka"`
	Quality                     string `json:"quality"`
	SatisfiedByAlternateOpinion bool   `json:"satisfied_by_alternate_opinion"`
	ExceptionApplies            bool   `json:"exception_applies"`
}
