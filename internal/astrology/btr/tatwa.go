package btr

type Element struct {
	Name     string
	Duration float64
}

var Elements = []Element{
	{"Prithvi", 6.0},
	{"Jala", 12.0},
	{"Tejo", 18.0},
	{"Vayu", 24.0},
	{"Akash", 30.0},
}

var WeekdayLords = []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"}

// Day start index into Elements array (0=Prithvi, 1=Jala, 2=Tejo, 3=Vayu, 4=Akash)
var DayStartIndex = map[string]int{
	"Sun":     2, // Tejo
	"Moon":    1, // Jala
	"Mars":    2, // Tejo
	"Mercury": 0, // Prithvi
	"Jupiter": 4, // Akash
	"Venus":   1, // Jala
	"Saturn":  3, // Vayu
}

func getTatwaSequence(dayLord string, isAroha bool) []Element {
	startIdx := DayStartIndex[dayLord]

	directSeq := make([]Element, 5)
	for i := 0; i < 5; i++ {
		directSeq[i] = Elements[(startIdx+i)%5]
	}

	if isAroha {
		return directSeq
	}

	// Avaroha: Reverse of direct sequence
	revSeq := make([]Element, 5)
	for i := 0; i < 5; i++ {
		revSeq[i] = directSeq[4-i]
	}
	return revSeq
}

func getAntarTatwaSequence(mainTatwa string, isAroha bool) []Element {
	// Find index of mainTatwa in standard Elements
	startIdx := 0
	for i, e := range Elements {
		if e.Name == mainTatwa {
			startIdx = i
			break
		}
	}

	directSeq := make([]Element, 5)
	for i := 0; i < 5; i++ {
		directSeq[i] = Elements[(startIdx+i)%5]
	}

	if isAroha {
		return directSeq
	}

	revSeq := make([]Element, 5)
	for i := 0; i < 5; i++ {
		revSeq[i] = directSeq[4-i]
	}
	return revSeq
}

// Calculate ruling Tatwa and Antar-Tatwa for a given elapsed time in minutes
func CalculateTatwa(dayLord string, elapsedMins float64) (main Element, antar Element, isAroha bool, cycleIdx int) {
	if elapsedMins < 0 {
		elapsedMins = 0 // Cap at sunrise if searching backwards crosses it
	}

	cycleIdx = int(elapsedMins / 90.0)
	minsInCycle := elapsedMins - float64(cycleIdx)*90.0

	isAroha = (cycleIdx%2 == 0) // Even cycles are Aroha, Odd are Avaroha

	mainSeq := getTatwaSequence(dayLord, isAroha)

	accum := 0.0
	for _, t := range mainSeq {
		if minsInCycle < accum+t.Duration {
			main = t
			break
		}
		accum += t.Duration
	}

	// Now find Antar-Tatwa
	minsInMain := minsInCycle - accum
	antarSeq := getAntarTatwaSequence(main.Name, isAroha)

	antarAccum := 0.0
	for _, at := range antarSeq {
		antarDur := (main.Duration * at.Duration) / 90.0
		if minsInMain < antarAccum+antarDur {
			antar = Element{Name: at.Name, Duration: antarDur}
			break
		}
		antarAccum += antarDur
	}

	return main, antar, isAroha, cycleIdx
}
