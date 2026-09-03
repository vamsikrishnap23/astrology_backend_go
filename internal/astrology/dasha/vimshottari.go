package dasha

import (
	"math"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

var vimshottariLords = []string{"Ketu", "Venus", "Sun", "Moon", "Mars", "Rahu", "Jupiter", "Saturn", "Mercury"}
var vimshottariYears = []float64{7, 20, 6, 10, 7, 18, 16, 19, 17}

// VimshottariYearDays defines the standard solar year length used in modern Vimshottari calculations.
const VimshottariYearDays = 365.2425
const nakshatraLen = 13.0 + 1.0/3.0

func indexOfLord(lord string) int {
	for i, l := range vimshottariLords {
		if l == lord {
			return i
		}
	}
	return 0
}

func addYears(t time.Time, years float64) time.Time {
	// 1 year = 365.2425 days
	duration := time.Duration(years * VimshottariYearDays * 24 * float64(time.Hour))
	return t.Add(duration)
}

// CalculateVimshottari derives the 4-level Vimshottari Dasha periods for a full 120-year cycle.
func CalculateVimshottari(moonLongitude float64, birthTimeUTC time.Time) domain.VimshottariDashaResult {
	moonLongitude = math.Mod(moonLongitude, 360.0)
	if moonLongitude < 0 {
		moonLongitude += 360.0
	}

	nakIdx := int(math.Floor(moonLongitude/nakshatraLen + 1e-12))
	rem := moonLongitude - float64(nakIdx)*nakshatraLen
	if rem < 0 {
		rem = 0
	}

	fractionElapsed := rem / nakshatraLen
	fractionRemaining := 1.0 - fractionElapsed

	startLordIdx := nakIdx % 9
	startLordTotalYears := vimshottariYears[startLordIdx]
	balanceYears := fractionRemaining * startLordTotalYears

	passedYears := fractionElapsed * startLordTotalYears
	theoreticalStartUTC := addYears(birthTimeUTC, -passedYears)

	res := domain.VimshottariDashaResult{
		BalanceYears: balanceYears,
	}

	currentTime := theoreticalStartUTC
	currentLordIdx := startLordIdx

	for i := 0; i < 9; i++ {
		mdLord := vimshottariLords[currentLordIdx]
		mdYears := vimshottariYears[currentLordIdx]
		mdEnd := addYears(currentTime, mdYears)

		if mdEnd.After(birthTimeUTC) {
			md := domain.Mahadasha{
				Lord:      mdLord,
				StartDate: currentTime,
				EndDate:   mdEnd,
			}

			// Generate Antardashas
			adTime := currentTime
			adLordIdx := currentLordIdx
			for j := 0; j < 9; j++ {
				adLord := vimshottariLords[adLordIdx]
				adYears := (mdYears * vimshottariYears[adLordIdx]) / 120.0
				adEnd := addYears(adTime, adYears)

				if adEnd.After(birthTimeUTC) {
					ad := domain.Antardasha{
						Lord:      adLord,
						StartDate: adTime,
						EndDate:   adEnd,
					}

					// Generate Pratyantardashas
					pdTime := adTime
					pdLordIdx := adLordIdx
					for k := 0; k < 9; k++ {
						pdLord := vimshottariLords[pdLordIdx]
						pdYears := (adYears * vimshottariYears[pdLordIdx]) / 120.0
						pdEnd := addYears(pdTime, pdYears)

						if pdEnd.After(birthTimeUTC) {
							pd := domain.Pratyantardasha{
								Lord:      pdLord,
								StartDate: pdTime,
								EndDate:   pdEnd,
							}

							// Generate Sookshma
							sdTime := pdTime
							sdLordIdx := pdLordIdx
							for m := 0; m < 9; m++ {
								sdLord := vimshottariLords[sdLordIdx]
								sdYears := (pdYears * vimshottariYears[sdLordIdx]) / 120.0
								sdEnd := addYears(sdTime, sdYears)

								if sdEnd.After(birthTimeUTC) {
									sd := domain.Sookshma{
										Lord:      sdLord,
										StartDate: sdTime,
										EndDate:   sdEnd,
									}
									pd.Sookshma = append(pd.Sookshma, sd)
								}
								sdTime = sdEnd
								sdLordIdx = (sdLordIdx + 1) % 9
							}
							ad.Pratyantardasha = append(ad.Pratyantardasha, pd)
						}
						pdTime = pdEnd
						pdLordIdx = (pdLordIdx + 1) % 9
					}
					md.Antardasha = append(md.Antardasha, ad)
				}
				adTime = adEnd
				adLordIdx = (adLordIdx + 1) % 9
			}
			res.Mahadasha = append(res.Mahadasha, md)
		}
		currentTime = mdEnd
		currentLordIdx = (currentLordIdx + 1) % 9
	}

	return res
}
