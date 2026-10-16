package ephemeris_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	astronomyTime "github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/time"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

type calcJob struct {
	ayanamsa string
	utc      time.Time
}

// calc returns every planet's sidereal longitude, the ascendant and the ayanamsa for one job.
func (j calcJob) calc() ([]float64, error) {
	ctx := domain.CalculationContext{
		Input:       domain.BirthInput{Latitude: 17.385, Longitude: 78.4867},
		JulianDayUT: astronomyTime.UTCToJulianDay(j.utc),
		Config: domain.CalculationConfig{
			AyanamsaMode: ephemeris.GetAyanamsaMode(j.ayanamsa),
			HouseCode:    ephemeris.GetHouseSystemCode("Placidus"),
		},
	}
	pl, err := planets.CalculatePlanets(&ctx)
	if err != nil {
		return nil, err
	}
	asc, _, _, err := houses.CalculateHouses(&ctx)
	if err != nil {
		return nil, err
	}
	out := []float64{asc, ctx.Ayanamsa}
	for _, p := range pl {
		out = append(out, p.SiderealLongitude)
	}
	return out, nil
}

// The Swiss Ephemeris C library keeps its state (ephemeris path, sidereal mode) per OS thread. Go may move
// a goroutine to another thread between cgo calls, so many requests at once used to fail with a Moshier
// fallback error or read the wrong ayanamsa. Every concurrent result must equal the single-thread result.
func TestConcurrentCalculationsAreStable(t *testing.T) {
	// Work out the expected values one job at a time on a single pinned thread.
	runtime.LockOSThread()
	if err := ephemeris.Init("../../../ephe_data"); err != nil {
		runtime.UnlockOSThread()
		t.Fatalf("ephemeris init failed: %v", err)
	}
	defer ephemeris.Close()

	ayanamsas := []string{"Lahiri", "Krishnamurti", "Raman", "Fagan-Bradley"}
	times := []time.Time{
		time.Date(1987, 3, 14, 4, 20, 0, 0, time.UTC),
		time.Date(2005, 11, 23, 10, 5, 0, 0, time.UTC),
		time.Date(2001, 7, 2, 22, 45, 0, 0, time.UTC),
		time.Date(1963, 12, 9, 13, 10, 0, 0, time.UTC),
	}
	var jobs []calcJob
	var want [][]float64
	for _, a := range ayanamsas {
		for _, tm := range times {
			j := calcJob{a, tm}
			v, err := j.calc()
			if err != nil {
				runtime.UnlockOSThread()
				t.Fatalf("baseline %v failed: %v", j, err)
			}
			jobs = append(jobs, j)
			want = append(want, v)
		}
	}
	runtime.UnlockOSThread()

	const workers, iterations = 16, 120
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		failures []string
	)
	record := func(msg string) {
		mu.Lock()
		defer mu.Unlock()
		if len(failures) < 5 {
			failures = append(failures, msg)
		}
	}
	bad := 0
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				idx := (w*7 + i) % len(jobs)
				got, err := jobs[idx].calc()
				if err != nil {
					mu.Lock()
					bad++
					mu.Unlock()
					record(fmt.Sprintf("job %d (%s): error: %v", idx, jobs[idx].ayanamsa, err))
					continue
				}
				for k := range got {
					if got[k] != want[idx][k] {
						mu.Lock()
						bad++
						mu.Unlock()
						record(fmt.Sprintf("job %d (%s): value %d = %v, want %v", idx, jobs[idx].ayanamsa, k, got[k], want[idx][k]))
						break
					}
				}
			}
		}(w)
	}
	wg.Wait()

	if bad > 0 {
		t.Errorf("%d of %d concurrent calculations were wrong or failed; first few:\n%v", bad, workers*iterations, failures)
	}
}
