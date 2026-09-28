package avastha

// GetBaaladiAvastha calculates the Age-based state of a planet.
func GetBaaladiAvastha(sign string, degree float64) string {
	isOdd := false
	switch sign {
	case "Aries", "Gemini", "Leo", "Libra", "Sagittarius", "Aquarius":
		isOdd = true
	}

	chunk := int(degree / 6.0)
	if chunk > 4 {
		chunk = 4 // safety to prevent out of bounds
	}

	if isOdd {
		switch chunk {
		case 0:
			return "Bala"
		case 1:
			return "Kumara"
		case 2:
			return "Yuva"
		case 3:
			return "Vriddha"
		case 4:
			return "Mrita"
		}
	} else {
		switch chunk {
		case 0:
			return "Mrita"
		case 1:
			return "Vriddha"
		case 2:
			return "Yuva"
		case 3:
			return "Kumara"
		case 4:
			return "Bala"
		}
	}
	return ""
}

// GetDeeptadiAvastha calculates the Dignity/Mood-based state of a planet using basic natural friendship.
func GetDeeptadiAvastha(planet, sign string, isCombust bool) string {
	if isCombust {
		return "Kopa" // Furious / Combust
	}

	type Dignity struct {
		Own     []string
		Exalt   string
		Debil   string
		Friends []string
		Enemies []string
	}

	dignities := map[string]Dignity{
		"Sun": {
			Own:     []string{"Leo"},
			Exalt:   "Aries",
			Debil:   "Libra",
			Friends: []string{"Cancer", "Scorpio", "Sagittarius", "Pisces"},
			Enemies: []string{"Taurus", "Capricorn", "Aquarius"},
		},
		"Moon": {
			Own:     []string{"Cancer"},
			Exalt:   "Taurus",
			Debil:   "Scorpio",
			Friends: []string{"Leo", "Gemini", "Virgo"},
			Enemies: []string{},
		},
		"Mars": {
			Own:     []string{"Aries", "Scorpio"},
			Exalt:   "Capricorn",
			Debil:   "Cancer",
			Friends: []string{"Leo", "Sagittarius", "Pisces"}, // Moon is cancer, but it's debil there
			Enemies: []string{"Gemini", "Virgo"},
		},
		"Mercury": {
			Own:     []string{"Gemini", "Virgo"},
			Exalt:   "Virgo",
			Debil:   "Pisces",
			Friends: []string{"Leo", "Taurus", "Libra"},
			Enemies: []string{"Cancer"},
		},
		"Jupiter": {
			Own:     []string{"Sagittarius", "Pisces"},
			Exalt:   "Cancer",
			Debil:   "Capricorn",
			Friends: []string{"Leo", "Aries", "Scorpio"},
			Enemies: []string{"Gemini", "Virgo", "Taurus", "Libra"},
		},
		"Venus": {
			Own:     []string{"Taurus", "Libra"},
			Exalt:   "Pisces",
			Debil:   "Virgo",
			Friends: []string{"Gemini", "Capricorn", "Aquarius"}, // Mercury is Virgo (debil)
			Enemies: []string{"Leo", "Cancer"},
		},
		"Saturn": {
			Own:     []string{"Capricorn", "Aquarius"},
			Exalt:   "Libra",
			Debil:   "Aries",
			Friends: []string{"Gemini", "Virgo", "Taurus"}, // Libra is exalt
			Enemies: []string{"Leo", "Cancer", "Scorpio"}, // Aries is debil
		},
	}

	dig, ok := dignities[planet]
	if !ok {
		return "" // Return empty for Nodes (Rahu/Ketu)
	}

	if sign == dig.Exalt {
		return "Deepta"
	}
	if sign == dig.Debil {
		return "Deena"
	}
	for _, s := range dig.Own {
		if sign == s {
			return "Swastha"
		}
	}
	for _, s := range dig.Friends {
		if sign == s {
			return "Mudita"
		}
	}
	for _, s := range dig.Enemies {
		if sign == s {
			return "Dukhita"
		}
	}

	return "Shanta" // Neutral
}
