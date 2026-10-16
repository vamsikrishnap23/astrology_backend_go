package matchmaking

// classical-guna-milan-v2

var Nakshatras = []string{
	"Ashwini", "Bharani", "Krittika", "Rohini", "Mrigashira", "Ardra",
	"Punarvasu", "Pushya", "Ashlesha", "Magha", "Purva Phalguni", "Uttara Phalguni",
	"Hasta", "Chitra", "Swati", "Vishakha", "Anuradha", "Jyeshtha",
	"Mula", "Purva Ashadha", "Uttara Ashadha", "Shravana", "Dhanishta", "Shatabhisha",
	"Purva Bhadrapada", "Uttara Bhadrapada", "Revati",
}

var Signs = []string{
	"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo",
	"Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces",
}

var SignLords = map[string]string{
	"Aries": "Mars", "Scorpio": "Mars",
	"Taurus": "Venus", "Libra": "Venus",
	"Gemini": "Mercury", "Virgo": "Mercury",
	"Cancer":      "Moon",
	"Leo":         "Sun",
	"Sagittarius": "Jupiter", "Pisces": "Jupiter",
	"Capricorn": "Saturn", "Aquarius": "Saturn",
}

// 1. Varna Mapping (1 point)
var VarnaMapping = map[string]string{
	"Cancer": "Brahmin", "Scorpio": "Brahmin", "Pisces": "Brahmin",
	"Aries": "Kshatriya", "Leo": "Kshatriya", "Sagittarius": "Kshatriya",
	"Taurus": "Vaishya", "Virgo": "Vaishya", "Capricorn": "Vaishya",
	"Gemini": "Shudra", "Libra": "Shudra", "Aquarius": "Shudra",
}

var VarnaRank = map[string]int{
	"Brahmin":   4,
	"Kshatriya": 3,
	"Vaishya":   2,
	"Shudra":    1,
}

// 2. Vashya (2 points): the book's table (pp. 69-70). Each rasi lists the rasis that are vashya to it;
// a couple scores 2 if either rasi lists the other. The old 5-animal-category rule was removed.
// NOTE: the book prints Aries as "Leo, Virgo" (common tables: "Leo, Scorpio"); kept as printed.
var VashyaRasis = map[string][]string{
	"Aries":       {"Leo", "Virgo"},
	"Taurus":      {"Cancer", "Libra"},
	"Gemini":      {"Virgo"},
	"Cancer":      {"Scorpio", "Sagittarius"},
	"Leo":         {"Libra"},
	"Virgo":       {"Pisces", "Gemini"},
	"Libra":       {"Capricorn", "Virgo"},
	"Scorpio":     {"Cancer"},
	"Sagittarius": {"Pisces"},
	"Capricorn":   {"Aries", "Aquarius"},
	"Aquarius":    {"Aries"},
	"Pisces":      {"Capricorn"},
}

// 4. Yoni Mapping (4 points)
var NakshatraYoni = map[string]string{
	"Ashwini": "Horse", "Bharani": "Elephant", "Krittika": "Sheep", "Rohini": "Serpent",
	"Mrigashira": "Serpent", "Ardra": "Dog", "Punarvasu": "Cat", "Pushya": "Sheep",
	"Ashlesha": "Cat", "Magha": "Rat", "Purva Phalguni": "Rat", "Uttara Phalguni": "Cow",
	"Hasta": "Buffalo", "Chitra": "Tiger", "Swati": "Buffalo", "Vishakha": "Tiger",
	"Anuradha": "Hare", "Jyeshtha": "Hare", "Mula": "Dog", "Purva Ashadha": "Monkey",
	"Uttara Ashadha": "Mongoose", "Shravana": "Monkey", "Dhanishta": "Lion",
	"Shatabhisha": "Horse", "Purva Bhadrapada": "Lion", "Uttara Bhadrapada": "Cow",
	"Revati": "Elephant",
}

var YoniScores = map[string]map[string]float64{
	"Horse":    {"Horse": 4, "Elephant": 2, "Sheep": 2, "Serpent": 3, "Dog": 2, "Cat": 2, "Rat": 2, "Cow": 1, "Buffalo": 0, "Tiger": 1, "Hare": 3, "Monkey": 3, "Mongoose": 2, "Lion": 1},
	"Elephant": {"Horse": 2, "Elephant": 4, "Sheep": 3, "Serpent": 3, "Dog": 2, "Cat": 2, "Rat": 2, "Cow": 2, "Buffalo": 3, "Tiger": 1, "Hare": 2, "Monkey": 3, "Mongoose": 2, "Lion": 0},
	"Sheep":    {"Horse": 2, "Elephant": 3, "Sheep": 4, "Serpent": 2, "Dog": 1, "Cat": 2, "Rat": 1, "Cow": 3, "Buffalo": 3, "Tiger": 1, "Hare": 2, "Monkey": 0, "Mongoose": 2, "Lion": 1},
	"Serpent":  {"Horse": 3, "Elephant": 3, "Sheep": 2, "Serpent": 4, "Dog": 2, "Cat": 1, "Rat": 1, "Cow": 1, "Buffalo": 1, "Tiger": 2, "Hare": 2, "Monkey": 2, "Mongoose": 0, "Lion": 2},
	"Dog":      {"Horse": 2, "Elephant": 2, "Sheep": 1, "Serpent": 2, "Dog": 4, "Cat": 2, "Rat": 1, "Cow": 2, "Buffalo": 2, "Tiger": 1, "Hare": 0, "Monkey": 2, "Mongoose": 1, "Lion": 1},
	"Cat":      {"Horse": 2, "Elephant": 2, "Sheep": 2, "Serpent": 1, "Dog": 2, "Cat": 4, "Rat": 0, "Cow": 2, "Buffalo": 2, "Tiger": 1, "Hare": 3, "Monkey": 3, "Mongoose": 2, "Lion": 1},
	"Rat":      {"Horse": 2, "Elephant": 2, "Sheep": 1, "Serpent": 1, "Dog": 1, "Cat": 0, "Rat": 4, "Cow": 2, "Buffalo": 2, "Tiger": 2, "Hare": 2, "Monkey": 2, "Mongoose": 1, "Lion": 2},
	"Cow":      {"Horse": 1, "Elephant": 2, "Sheep": 3, "Serpent": 1, "Dog": 2, "Cat": 2, "Rat": 2, "Cow": 4, "Buffalo": 3, "Tiger": 0, "Hare": 3, "Monkey": 2, "Mongoose": 2, "Lion": 1},
	"Buffalo":  {"Horse": 0, "Elephant": 3, "Sheep": 3, "Serpent": 1, "Dog": 2, "Cat": 2, "Rat": 2, "Cow": 3, "Buffalo": 4, "Tiger": 1, "Hare": 2, "Monkey": 2, "Mongoose": 2, "Lion": 1},
	"Tiger":    {"Horse": 1, "Elephant": 1, "Sheep": 1, "Serpent": 2, "Dog": 1, "Cat": 1, "Rat": 2, "Cow": 0, "Buffalo": 1, "Tiger": 4, "Hare": 1, "Monkey": 2, "Mongoose": 2, "Lion": 1},
	"Hare":     {"Horse": 3, "Elephant": 2, "Sheep": 2, "Serpent": 2, "Dog": 0, "Cat": 3, "Rat": 2, "Cow": 3, "Buffalo": 2, "Tiger": 1, "Hare": 4, "Monkey": 2, "Mongoose": 2, "Lion": 1},
	"Monkey":   {"Horse": 3, "Elephant": 3, "Sheep": 0, "Serpent": 2, "Dog": 2, "Cat": 3, "Rat": 2, "Cow": 2, "Buffalo": 2, "Tiger": 2, "Hare": 2, "Monkey": 4, "Mongoose": 3, "Lion": 2},
	"Mongoose": {"Horse": 2, "Elephant": 2, "Sheep": 2, "Serpent": 0, "Dog": 1, "Cat": 2, "Rat": 1, "Cow": 2, "Buffalo": 2, "Tiger": 2, "Hare": 2, "Monkey": 3, "Mongoose": 4, "Lion": 2},
	"Lion":     {"Horse": 1, "Elephant": 0, "Sheep": 1, "Serpent": 2, "Dog": 1, "Cat": 1, "Rat": 2, "Cow": 1, "Buffalo": 1, "Tiger": 1, "Hare": 1, "Monkey": 2, "Mongoose": 2, "Lion": 4},
}

// 5. Graha Maitri (5 points)
var GrahaNaturalMaitri = map[string]map[string]string{
	"Sun":     {"Moon": "Friend", "Mars": "Friend", "Jupiter": "Friend", "Mercury": "Neutral", "Venus": "Enemy", "Saturn": "Enemy"},
	"Moon":    {"Sun": "Friend", "Mercury": "Friend", "Mars": "Neutral", "Jupiter": "Neutral", "Venus": "Neutral", "Saturn": "Neutral"},
	"Mars":    {"Sun": "Friend", "Moon": "Friend", "Jupiter": "Friend", "Venus": "Neutral", "Saturn": "Neutral", "Mercury": "Enemy"},
	"Mercury": {"Sun": "Friend", "Venus": "Friend", "Mars": "Neutral", "Jupiter": "Neutral", "Saturn": "Neutral", "Moon": "Enemy"},
	"Jupiter": {"Sun": "Friend", "Moon": "Friend", "Mars": "Friend", "Saturn": "Neutral", "Mercury": "Enemy", "Venus": "Enemy"},
	"Venus":   {"Mercury": "Friend", "Saturn": "Friend", "Mars": "Neutral", "Jupiter": "Neutral", "Sun": "Enemy", "Moon": "Enemy"},
	"Saturn":  {"Mercury": "Friend", "Venus": "Friend", "Jupiter": "Neutral", "Sun": "Enemy", "Moon": "Enemy", "Mars": "Enemy"},
}

// 6. Gana Koota (6 points)
var NakshatraGana = map[string]string{
	"Ashwini": "Deva", "Mrigashira": "Deva", "Punarvasu": "Deva", "Pushya": "Deva",
	"Hasta": "Deva", "Swati": "Deva", "Anuradha": "Deva", "Shravana": "Deva", "Revati": "Deva",

	"Bharani": "Manushya", "Rohini": "Manushya", "Ardra": "Manushya", "Purva Phalguni": "Manushya",
	"Uttara Phalguni": "Manushya", "Purva Ashadha": "Manushya", "Uttara Ashadha": "Manushya",
	"Purva Bhadrapada": "Manushya", "Uttara Bhadrapada": "Manushya",

	"Krittika": "Rakshasa", "Ashlesha": "Rakshasa", "Magha": "Rakshasa", "Chitra": "Rakshasa",
	"Vishakha": "Rakshasa", "Jyeshtha": "Rakshasa", "Mula": "Rakshasa", "Dhanishta": "Rakshasa",
	"Shatabhisha": "Rakshasa",
}

// Row=Bride, Col=Groom (published table: a Deva bride with a Manushya groom scores 6, a Manushya bride
// with a Deva groom scores 5, a Deva bride with a Rakshasa groom 0, a Rakshasa bride with a Deva groom 1).
var GanaScores = map[string]map[string]float64{
	"Deva":     {"Deva": 6, "Manushya": 6, "Rakshasa": 0},
	"Manushya": {"Deva": 5, "Manushya": 6, "Rakshasa": 0},
	"Rakshasa": {"Deva": 1, "Manushya": 0, "Rakshasa": 6},
}

// Bride nakshatras whose gana score is raised to the full 6 whatever the groom's gana.
var GanaExemptBrideNakshatras = map[string]bool{
	"Ashlesha": true, "Chitra": true, "Vishakha": true, "Jyeshtha": true, "Mula": true, "Shatabhisha": true,
}

// 8. Nadi Koota (8 points)
var NakshatraNadi = map[string]string{
	"Ashwini": "Adi", "Ardra": "Adi", "Punarvasu": "Adi", "Uttara Phalguni": "Adi",
	"Hasta": "Adi", "Jyeshtha": "Adi", "Mula": "Adi", "Shatabhisha": "Adi", "Purva Bhadrapada": "Adi",

	"Bharani": "Madhya", "Mrigashira": "Madhya", "Pushya": "Madhya", "Purva Phalguni": "Madhya",
	"Chitra": "Madhya", "Anuradha": "Madhya", "Purva Ashadha": "Madhya", "Dhanishta": "Madhya", "Uttara Bhadrapada": "Madhya",

	"Krittika": "Antya", "Rohini": "Antya", "Ashlesha": "Antya", "Magha": "Antya",
	"Swati": "Antya", "Vishakha": "Antya", "Uttara Ashadha": "Antya", "Shravana": "Antya", "Revati": "Antya",
}

// Rajju Mapping
var NakshatraRajju = map[string]string{
	"Ashwini": "Pada", "Bharani": "Uru", "Krittika": "Nabhi", "Rohini": "Kantha", "Mrigashira": "Shiro",
	"Ardra": "Kantha", "Punarvasu": "Nabhi", "Pushya": "Uru", "Ashlesha": "Pada",
	"Magha": "Pada", "Purva Phalguni": "Uru", "Uttara Phalguni": "Nabhi", "Hasta": "Kantha", "Chitra": "Shiro",
	"Swati": "Kantha", "Vishakha": "Nabhi", "Anuradha": "Uru", "Jyeshtha": "Pada",
	"Mula": "Pada", "Purva Ashadha": "Uru", "Uttara Ashadha": "Nabhi", "Shravana": "Kantha",
	"Dhanishta": "Shiro", "Shatabhisha": "Kantha", "Purva Bhadrapada": "Nabhi", "Uttara Bhadrapada": "Uru", "Revati": "Pada",
}

// Jaimini Maharshi Nadi Exceptions (21 nakshatras). Not applied to the Nadi score: when both nakshatras
// are in this list the rule is only reported in NadiResult.Suggestion.
var JaiminiNadiExceptions = map[string]bool{
	"Uttara Phalguni": true, "Shatabhisha": true, "Purva Bhadrapada": true, "Punarvasu": true, "Ardra": true, "Mula": true, "Ashwini": true, // Adi
	"Purva Ashadha": true, "Anuradha": true, "Dhanishta": true, "Pushya": true, "Chitra": true, "Purva Phalguni": true, "Mrigashira": true, // Madhya
	"Krittika": true, "Vishakha": true, "Ashlesha": true, "Shravana": true, "Magha": true, "Uttara Ashadha": true, "Rohini": true, // Antya
}

// Bhakoot Exceptions (Auspicious pairs that cancel the Dosha)
// We store them as "Sign1-Sign2" strings (alphabetically sorted internally or just both combinations)
var AuspiciousBhakootPairs = map[string]bool{
	// 2/12 Auspicious
	"Aries-Pisces": true, "Pisces-Aries": true,
	"Cancer-Leo": true, "Leo-Cancer": true,
	"Libra-Virgo": true, "Virgo-Libra": true,
	"Aquarius-Capricorn": true, "Capricorn-Aquarius": true,
	"Gemini-Taurus": true, "Taurus-Gemini": true,
	"Sagittarius-Scorpio": true, "Scorpio-Sagittarius": true,

	// 5/9 Auspicious
	"Aries-Leo": true, "Leo-Aries": true,
	"Gemini-Libra": true, "Libra-Gemini": true,
	"Aquarius-Libra": true, "Libra-Aquarius": true,
	"Aries-Sagittarius": true, "Sagittarius-Aries": true,
	"Taurus-Virgo": true, "Virgo-Taurus": true,
	"Leo-Sagittarius": true, "Sagittarius-Leo": true,
	"Pisces-Scorpio": true, "Scorpio-Pisces": true,
	"Capricorn-Taurus": true, "Taurus-Capricorn": true,

	// 6/8 Auspicious
	"Aries-Scorpio": true, "Scorpio-Aries": true,
	"Leo-Pisces": true, "Pisces-Leo": true,
	"Cancer-Sagittarius": true, "Sagittarius-Cancer": true,
	"Capricorn-Gemini": true, "Gemini-Capricorn": true,
	"Libra-Taurus": true, "Taurus-Libra": true,
	"Aquarius-Virgo": true, "Virgo-Aquarius": true,
}

// Mruga Vairam (Animal Enmity for Yoni)
var MrugaEnmity = map[string]string{
	"Horse": "Buffalo", "Buffalo": "Horse",
	"Cow": "Tiger", "Tiger": "Cow",
	"Monkey": "Sheep", "Sheep": "Monkey",
	"Hare": "Dog", "Dog": "Hare",
	"Elephant": "Lion", "Lion": "Elephant",
	"Serpent": "Mongoose", "Mongoose": "Serpent",
	"Cat": "Rat", "Rat": "Cat",
}

// Eka Nakshatra (Same Nakshatra) Exceptions (Page 3 of PDF)
// If same nakshatra, these specific ones are Ashubham (Inauspicious)
// Eka Nakshatra (Same Nakshatra) Exceptions (Page 3 of PDF)
// Only the explicitly Ashubham (Inauspicious) stars are mapped here.
// Missing stars are either Shubham (Good) or Madhyamam (Medium), both of which bypass the Dosha.
var BadEkaNakshatras = map[string]bool{
	"Bharani": true, "Ashlesha": true, "Swati": true,
	"Jyeshtha": true, "Mula": true, "Dhanishta": true,
	"Shatabhisha": true, "Purva Bhadrapada": true,
}

// Ashubha Kurpu (Inauspicious Combinations - 15 explicitly rejected pairs from Page 3)
// Key: Bride Nakshatra, Value: Groom Nakshatra
var AshubhaKurpu = map[string]string{
	"Krittika":          "Ashlesha",
	"Mrigashira":        "Purva Phalguni",
	"Ashlesha":          "Swati",
	"Magha":             "Vishakha",
	"Uttara Phalguni":   "Jyeshtha",
	"Hasta":             "Mula",
	"Chitra":            "Purva Ashadha",
	"Vishakha":          "Shravana",
	"Anuradha":          "Dhanishta",
	"Jyeshtha":          "Shatabhisha",
	"Shravana":          "Ashwini",
	"Dhanishta":         "Bharani",
	"Shatabhisha":       "Krittika",
	"Uttara Bhadrapada": "Mrigashira",
	"Revati":            "Ardra",
}
