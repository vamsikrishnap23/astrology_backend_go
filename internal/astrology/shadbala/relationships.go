package shadbala

type Relationship int

const (
	AdhiShatru Relationship = iota
	Shatru
	Sama
	Mitra
	AdhiMitra
	OwnHouse
	Moolatrikona
)

var NaisargikaMaitri = map[string]map[string]Relationship{
	"Sun": {
		"Moon": Mitra, "Mars": Mitra, "Jupiter": Mitra,
		"Venus": Shatru, "Saturn": Shatru,
		"Mercury": Sama,
	},
	"Moon": {
		"Sun": Mitra, "Mercury": Mitra,
		"Mars": Sama, "Jupiter": Sama, "Venus": Sama, "Saturn": Sama,
	},
	"Mars": {
		"Sun": Mitra, "Moon": Mitra, "Jupiter": Mitra,
		"Mercury": Shatru,
		"Venus":   Sama, "Saturn": Sama,
	},
	"Mercury": {
		"Sun": Mitra, "Venus": Mitra,
		"Moon": Shatru,
		"Mars": Sama, "Jupiter": Sama, "Saturn": Sama,
	},
	"Jupiter": {
		"Sun": Mitra, "Moon": Mitra, "Mars": Mitra,
		"Mercury": Shatru, "Venus": Shatru,
		"Saturn": Sama,
	},
	"Venus": {
		"Mercury": Mitra, "Saturn": Mitra,
		"Sun": Shatru, "Moon": Shatru,
		"Mars": Sama, "Jupiter": Sama,
	},
	"Saturn": {
		"Mercury": Mitra, "Venus": Mitra,
		"Sun": Shatru, "Moon": Shatru, "Mars": Shatru,
		"Jupiter": Sama,
	},
}

// Temporary relationship in the natal Rasi chart
func getTatkalikaMaitri(planet1Sign, planet2Sign int) Relationship {
	if planet1Sign == planet2Sign {
		return Shatru
	}
	distance := (planet2Sign - planet1Sign + 12) % 12
	// 2, 3, 4, 10, 11, 12 from planet1 (0-indexed: 1, 2, 3, 9, 10, 11)
	if distance == 1 || distance == 2 || distance == 3 || distance == 9 || distance == 10 || distance == 11 {
		return Mitra
	}
	// 1, 5, 6, 7, 8, 9 (0-indexed: 0, 4, 5, 6, 7, 8)
	return Shatru
}

func GetPanchadhaMaitri(planet1, planet2 string, p1Sign, p2Sign int) Relationship {
	natural := NaisargikaMaitri[planet1][planet2]
	temp := getTatkalikaMaitri(p1Sign, p2Sign)

	if natural == Mitra && temp == Mitra {
		return AdhiMitra
	} else if natural == Mitra && temp == Shatru {
		return Sama
	} else if natural == Sama && temp == Mitra {
		return Mitra
	} else if natural == Sama && temp == Shatru {
		return Shatru
	} else if natural == Shatru && temp == Mitra {
		return Sama
	} else if natural == Shatru && temp == Shatru {
		return AdhiShatru
	}
	return Sama
}

var OwnSigns = map[string][]string{
	"Sun":     {"Leo"},
	"Moon":    {"Cancer"},
	"Mars":    {"Aries", "Scorpio"},
	"Mercury": {"Gemini", "Virgo"},
	"Jupiter": {"Sagittarius", "Pisces"},
	"Venus":   {"Taurus", "Libra"},
	"Saturn":  {"Capricorn", "Aquarius"},
}

// moolatrikonaPortion is the part of a sign that is a planet's Moolatrikona, in degrees from the start
// of the sign (BPHS 3.51-54): the Sun's first 20 deg of Leo, the Moon's Taurus beyond 3 deg (the first 3
// are its exaltation), Mars's first 12 deg of Aries, Mercury's 15 to 20 deg of Virgo (the first 15 are
// its exaltation), Jupiter's first 10 deg of Sagittarius, Venus's first 15 deg of Libra and Saturn's
// first 20 deg of Aquarius. The rest of such a sign is the planet's own sign (the Moon's Taurus is not).
type moolatrikonaPortion struct {
	sign     string
	from, to float64
}

var moolatrikonaPortions = map[string]moolatrikonaPortion{
	"Sun":     {"Leo", 0, 20},
	"Moon":    {"Taurus", 3, 30},
	"Mars":    {"Aries", 0, 12},
	"Mercury": {"Virgo", 15, 20},
	"Jupiter": {"Sagittarius", 0, 10},
	"Venus":   {"Libra", 0, 15},
	"Saturn":  {"Aquarius", 0, 20},
}

// isMoolatrikona reports whether a planet standing at degreeInSign (0..30) of the named sign is in its
// Moolatrikona portion.
func isMoolatrikona(planet, sign string, degreeInSign float64) bool {
	m, ok := moolatrikonaPortions[planet]
	return ok && m.sign == sign && degreeInSign >= m.from && degreeInSign < m.to
}

func getSignIndexByName(name string) int {
	signs := []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"}
	for i, s := range signs {
		if s == name {
			return i
		}
	}
	return 0
}

// GetPlacementStrength is the relationship of planet to the sign it occupies in a varga. moolatrikona
// says that the planet is in its Moolatrikona portion of the Rasi chart: Moolatrikona counts only
// there, never in the other six vargas (Raman art. 30), where the same sign is just the planet's own.
func GetPlacementStrength(planet string, vargaSign string, natalSigns map[string]int, moolatrikona bool) Relationship {
	if moolatrikona {
		return Moolatrikona
	}
	// Then Own Sign
	for _, os := range OwnSigns[planet] {
		if os == vargaSign {
			return OwnHouse
		}
	}

	// Find Lord of the Varga Sign
	lord := ""
	for p, signs := range OwnSigns {
		for _, s := range signs {
			if s == vargaSign {
				lord = p
				break
			}
		}
		if lord != "" {
			break
		}
	}

	// If the planet itself is the lord (handled above, but just in case)
	if lord == planet {
		return OwnHouse
	}

	return GetPanchadhaMaitri(planet, lord, natalSigns[planet], natalSigns[lord])
}
