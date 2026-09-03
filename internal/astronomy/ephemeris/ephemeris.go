package ephemeris

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tejzpr/go-swisseph"
)

// Init initializes the Swiss Ephemeris.
// It explicitly requires a valid ephemeris path containing .se1 files.
func Init(path string) error {
	if path == "" {
		return errors.New("EPHE_PATH is not set. Swiss Ephemeris requires a valid path to .se1 ephemeris files")
	}

	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("EPHE_PATH '%s' is not a valid directory", path)
	}

	// Verify required basic ephemeris files are present for the contemporary era (1800-2399).
	// If you need broader historical support, you'd require more files (like sepl_12.se1).
	requiredFiles := []string{"sepl_18.se1", "semo_18.se1"}
	for _, req := range requiredFiles {
		fp := filepath.Join(path, req)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			return fmt.Errorf("required ephemeris file '%s' not found in EPHE_PATH '%s'. Do NOT rely on the Moshier analytical fallback", req, path)
		}
	}

	swisseph.SetEphePath(path)
	return nil
}

// Close closes the ephemeris.
func Close() {
	swisseph.Close()
}

// GetAyanamsaMode maps user-facing string to Swiss Ephemeris constant.
func GetAyanamsaMode(name string) int {
	name = strings.ToLower(name)
	switch name {
	case "lahiri", "true chitrapaksha":
		return swisseph.SidmLahiri
	case "raman":
		return swisseph.SidmRaman
	case "krishnamurti":
		return swisseph.SidmKrishnamurti
	case "fagan-bradley":
		return swisseph.SidmFaganBradley
	case "yukteshwar":
		return swisseph.SidmYukteshwar
	case "j.n. bhasin":
		return swisseph.SidmJNBhasin
	case "true revati":
		return swisseph.SidmTrueRevati
	case "true pushya":
		return swisseph.SidmTruePushya
	default:
		// Default to Lahiri
		return swisseph.SidmLahiri
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
