# KP Horary ascendant: "phantom JD" fix

Date: 2026-10-15. Code: `internal/astronomy/houses/houses.go` (`findPhantomJD`, called from `CalculateHouses`).
Tests: `houses_test.go` (`TestHoraryAscendantMatchesTable`, `TestHoraryPolarNeverSilentlyWrong`).

## What it does
For `horary_number` 1..249 the ascendant is forced to the START degree of that number's
sub-division (`HoraryTable[n-1].Start`, built in `table_249.go`). Swiss Ephemeris cannot be asked for
a given ascendant, so we find the time ("phantom JD") at which the real ascendant equals the target
and compute the houses at that time. Planets are still computed at the real birth time; only the
houses use the phantom time. Ayanamsa is taken at the real JD and used both when converting the
target to tropical and when converting the result back, so the sidereal ascendant equals the table
value exactly.

## The bug that was fixed
The old code guessed the phantom JD assuming the ascendant moves at a uniform 360 deg per 0.997 d,
then bisected only +-0.1 d around that guess, without checking convergence. The ascendant's rate
varies by sign and latitude, so the root fell outside the window and the loop silently returned the
unconverged guess (wrong cusps, no error).

Measured with the old algorithm (every number 1..249, 24 times of day):

| Latitude | Wrong | Worst error |
|---|---|---|
| 8.5N - 34N (India) | 0% | none |
| 36N | 1.8% | 36 deg |
| 45N | 18% | 52 deg |
| 51.5N | 30% | 67 deg |
| 55N | 37% | 77 deg |
| 38S | 5.7% | 40 deg |

So results for Indian latitudes are unchanged by the fix; only |lat| above ~36 is affected.

## New algorithm
1. `offset(jd)` = wrapped (asc(jd) - target) in (-180, 180].
2. Scan jd from start-0.5 to start+0.5 in 0.02 d steps; a crossing is `prev < 0 <= cur` with `cur - prev < 180`.
3. Take the crossing nearest the real time, bisect it 60 times.
4. Verify `|offset(root)| <= 1e-4 deg`; otherwise return an error.

## Verified
- Table: 249 rows, contiguous, ends at 360. Spot checks match the standard KP table:
  #1 = 0 deg 00'00" - 0 deg 46'40" (Ketu/Ketu), #2 ends at 3 deg 00'00" (Ketu/Venus), #125 starts at 180 deg 00'00".
- All 249 numbers x 3 locations x 2 times, plus lat 0, 23.4, 60, 66, and Placidus/Whole Sign/Equal/Koch,
  1900/1998/2100 epochs: asc matches the table to < 1e-5 deg, no errors.
- All other test packages still pass.

## Known limits
- **Beyond the polar circle (|lat| > ~66.5)** the ascendant is discontinuous and some numbers are not
  attainable. The function returns an error there (before the final check was added it could return an
  ascendant up to 179 deg off, without error). Handlers surface the error as HTTP 500 "Error calculating houses".
- **Convention not verified against reference software**: this places the ascendant at the START of the
  sub-division. Some KP software uses another point of it (e.g. its middle). If a reference program
  disagrees at an Indian latitude, that is the convention in `table_249.go`/this choice, not the root finder.
- **Nearest crossing**: the ascendant hits a given degree once per sidereal day; we choose the instant
  closest to real time. Only the ascendant/cusps are affected (planets stay at real time), so the
  choice matters only for cusps other than the ascendant, which vary slightly by day.
- Ayanamsa for the cusp conversion is taken at the real JD, not the phantom JD (difference is negligible).
