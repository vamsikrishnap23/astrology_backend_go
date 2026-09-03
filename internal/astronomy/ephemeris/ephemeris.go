package ephemeris

import (
	"strings"

	"github.com/mshafiee/swephgo"
)

// Init initializes the Swiss Ephemeris.
func Init(path string) {
	if path != "" {
		swephgo.SetEphePath([]byte(path))
	} else {
		// Use built-in or default
		swephgo.SetEphePath([]byte(""))
	}
}

// Close closes the ephemeris.
func Close() {
	swephgo.Close()
}

// GetAyanamsaMode maps user-facing string to Swiss Ephemeris constant.
func GetAyanamsaMode(name string) int {
	name = strings.ToLower(name)
	switch name {
	case "lahiri", "true chitrapaksha":
		return swephgo.SeSidmLahiri
	case "raman":
		return swephgo.SeSidmRaman
	case "krishnamurti":
		return swephgo.SeSidmKrishnamurti
	case "fagan-bradley":
		return swephgo.SeSidmFaganBradley
	case "yukteshwar":
		return swephgo.SeSidmYukteshwar
	case "j.n. bhasin":
		return swephgo.SeSidmJnBhasin
	case "true revati":
		return swephgo.SeSidmTrueRevati
	case "true pushya":
		return swephgo.SeSidmTruePushya
	default:
		// Default to Lahiri
		return swephgo.SeSidmLahiri
	}
}

// GetHouseSystemCode maps user-facing house string to Swiss Ephemeris byte code.
func GetHouseSystemCode(name string) byte {
	name = strings.ToLower(name)
	switch name {
	case "placidus":
		return 'P'
	case "koch":
		return 'K'
	case "porphyry":
		return 'O'
	case "regiomontanus":
		return 'R'
	case "campanus":
		return 'C'
	case "equal":
		return 'E' // or 'A'
	case "whole sign":
		return 'W'
	case "sripati":
		return 'S'
	default:
		// Default to Placidus
		return 'P'
	}
}
