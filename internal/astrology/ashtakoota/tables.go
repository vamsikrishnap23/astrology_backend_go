package ashtakoota

// classical-guna-milan-v1

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

// 2. Vashya Mapping (2 points)
// Aries: Chatuspada, Taurus: Chatuspada, Gemini: Manav (Biped), Cancer: Jalchar
// Leo: Vanchar (Wild), Virgo: Manav, Libra: Manav, Scorpio: Keeta (Insect)
// Sagittarius: 1st half Biped, 2nd half Chatuspada (we'll just use Chatuspada as classical generic, or split)
// Wait, classic Vashya: Sagittarius is Kshatriya/Chatuspada, Capricorn is Jalchar (1st half Chatuspada, 2nd half Jalchar).
// Let's use standard simplified: Sagittarius=Chatuspada, Capricorn=Jalchar, Aquarius=Manav, Pisces=Jalchar
var VashyaMapping = map[string]string{
	"Aries": "Chatuspada", "Taurus": "Chatuspada", "Gemini": "Manav", "Cancer": "Jalchar",
	"Leo": "Vanchar", "Virgo": "Manav", "Libra": "Manav", "Scorpio": "Keeta",
	"Sagittarius": "Chatuspada", "Capricorn": "Jalchar", "Aquarius": "Manav", "Pisces": "Jalchar",
}

// Vashya Scores [Groom][Bride]
// Manav, Chatuspada, Jalchar, Vanchar, Keeta
// Same: 2, Enemy: 0, Friendly/Neutral: 1, 0.5 (We'll use classical table)
var VashyaScores = map[string]map[string]float64{
	"Manav":      {"Manav": 2, "Chatuspada": 1, "Jalchar": 1, "Vanchar": 0, "Keeta": 1},
	"Chatuspada": {"Manav": 1, "Chatuspada": 2, "Jalchar": 1, "Vanchar": 0, "Keeta": 1},
	"Jalchar":    {"Manav": 1, "Chatuspada": 1, "Jalchar": 2, "Vanchar": 1, "Keeta": 1},
	"Vanchar":    {"Manav": 0, "Chatuspada": 0, "Jalchar": 1, "Vanchar": 2, "Keeta": 0},
	"Keeta":      {"Manav": 1, "Chatuspada": 1, "Jalchar": 1, "Vanchar": 0, "Keeta": 2},
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

// Row=Groom, Col=Bride
var GanaScores = map[string]map[string]float64{
	"Deva":     {"Deva": 6, "Manushya": 6, "Rakshasa": 1},
	"Manushya": {"Deva": 5, "Manushya": 6, "Rakshasa": 0},
	"Rakshasa": {"Deva": 1, "Manushya": 0, "Rakshasa": 6},
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
