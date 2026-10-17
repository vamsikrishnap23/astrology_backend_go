# Shadbala: what the code follows and where it comes from

Rule set: **Raman first, BPHS where Raman is silent** (owner's decision, 2026-10-17). Primary source: B.V. Raman,
*Graha and Bhava Balas* (13th ed., 1996; page numbers below are the printed ones). Secondary: *Brihat Parashara Hora
Shastra*, ch. 27 "Evaluation of Strengths" and ch. 3 (the English edition in `books/`, not committed). Unit: shashtiamsas
(virupas); 60 = 1 rupa. Planets: Sun to Saturn only. Test chart: Raman's "standard horoscope" (16 Oct 1918, 2:06:16 pm LMT,
13 N 77 35' E, Raman ayanamsa); its numbers are the regression targets (`shadbala_raman_test.go` and the per-component
tests).

## Components

| Component | Rule | Source | Code |
|---|---|---|---|
| Uchcha | arc from the deep-debilitation point (reflected to 0-180) / 3 | BPHS 27.1; Raman art. 20, p. 11-13 | `sthana.go` |
| Saptavargaja | value of the planet's sign in D1, D2, D3, D7, D9, D12, D30: Moolatrikona 45, own 30, adhimitra 22.5, mitra 15, sama 7.5, shatru 3.75, adhishatru 1.875; relations are panchadha (natural + temporary in the Rasi). **Moolatrikona (45) counts only in the Rasi chart, for the whole Moolatrikona sign** (Sun Leo, Moon Taurus, Mars Aries, Mercury Virgo, Jupiter Sagittarius, Venus Libra, Saturn Aquarius); in the other vargas the same sign is the planet's own (30). BPHS 3.51-54 limits the Moolatrikona to part of the sign; those degrees are kept as a comment in `relationships.go`, not used || Raman art. 30, p. 22-24 | `sthana.go`, `relationships.go` |
| Ojayugma | 15 in the Rasi and 15 in the Navamsa for the Moon and Venus in an even sign, the others in an odd sign | Raman art. 31, p. 24-26; BPHS 27.4 | `sthana.go` |
| Kendradi | 60 / 30 / 15 for the kendra / panaphara / apoklima *sign* from the Ascendant sign | Raman art. 32-35 | `sthana.go` |
| Drekkana | 15 for the Sun, Mars, Jupiter in the 1st decanate, Mercury and Saturn in the 2nd, Moon and Venus in the 3rd | Raman art. 36-39, p. 27; summary p. 105 | `sthana.go` |
| Dig | arc from the powerless point / 3 (Sun, Mars: 4th cusp; Jupiter, Mercury: 7th; Moon, Venus: 10th; Saturn: Ascendant) | Raman art. 44-45, p. 30-33 | `dig.go` |
| Nathonnatha | local *apparent* time, angle from midnight reflected to 0-180: Sun, Jupiter, Venus = angle/3; Moon, Mars, Saturn = 60 minus that; Mercury 60 | Raman art. 47-51, p. 34-37; BPHS 27.8-9 | `nathonnatha.go` |
| Paksha | arc Moon-Sun (reflected) / 3 for benefics, 60 minus it for malefics, the Moon's doubled; Mercury is malefic when it shares a sign with the Sun, Mars, Saturn, Rahu or Ketu | Raman art. 52-55, p. 37-39; BPHS 3.11 | `paksha.go` |
| Tribhaga | day or night split in three: Mercury, Sun, Saturn by day; Moon, Venus, Mars by night; Jupiter always | Raman art. 56-57 | `tribhaga.go` |
| Varsha, Masa | lord of the 360-day year (15) and 30-day month (30) of the Hindu day, from the Ahargana (days since 2 May 1827, day 1) | Raman art. 58-67, p. 41-49 and Tables I-III, p. 110-112; BPHS 27.13 | `ahargana.go` |
| Dina | lord of the weekday of the latest sunrise: 45 | Raman art. 67 | `hora.go`, `kala.go` |
| Hora | 60-minute horas counted from the latest sunrise, the first ruled by the weekday lord, then Sun, Venus, Mercury, Moon, Saturn, Jupiter, Mars | Raman art. 68-70, p. 49-53 | `hora.go` |
| Ayana | (24 deg + kranti)/48 x 60: north kranti adds for the Sun, Mars, Jupiter, Venus, south adds for the Moon and Saturn, Mercury always adds; Sun doubled; limited to 0-60 (**true declination from the ephemeris**, owner's decision) | Raman art. 71-75, p. 53-60 | `ayana.go` |
| Yuddha | planets (not Sun or Moon) less than 1 deg apart: difference of (Sthana + Dig + Kala up to Hora) divided by the difference of the disc diameters, added to the winner (smaller longitude) and taken from the loser | Raman art. 76-77, p. 60-61 | `yuddha.go` |
| Cheshta | Chesta Kendra = Seeghrochcha - (mean + true)/2, reduced to 0-180, / 3, for Mars, Mercury, Jupiter, Venus, Saturn; **none for the Sun and Moon** in the total | Raman art. 79-107, p. 64-79, Tables IV-IX p. 112-117, example 56 p. 89 | `cheshta.go` |
| Naisargika | 60, 51.43, 42.86, 34.29, 25.71, 17.14, 8.57 (Sun, Moon, Venus, Jupiter, Mercury, Mars, Saturn) | Raman art. 108, p. 80 | `naisargika.go` |
| Drik | quarter of (benefic - malefic) aspect values; aspect values by angle and the special aspects of Mars (+15 on 4th, 8th), Jupiter (+30 on 5th, 9th), Saturn (+45 on 3rd, 10th); Moon benefic from 84 to 264 deg elongation; Mercury as in Paksha | Raman art. 109-120, p. 82-88 (formulas p. 85) | `drik.go` |
| Total | Sthana + Dig + Kala + Cheshta + Naisargika +/- Drik; rupas = total/60; minimum 300, 360, 300, 420, 390, 330, 300 shashtiamsas | Raman art. 121-122, p. 89-90 and p. 109 | `shadbala.go` |
| Bhava Bala | per house: the lord's Shadbala Pinda + Bhava Digbala (from the sign group of the Sripati house middle: Nara, Jalachara, Chathushpada, Keeta) + Bhava Drishti Bala (benefic minus malefic aspects on the middle, Jupiter's and Mercury's whole, the others by a quarter, Mercury a full benefic) | Raman ch. IX, art. 124-134, pp. 91-100 | `bhava.go` |
| Ishta and Kashta Phala | Ishta = square root of Uchcha x Chesta Bala; Kashta = square root of (60 - Uchcha)(60 - Chesta); the Sun's Chesta from its Sayana longitude + 90 deg, the Moon's from its arc from the Sun, both reduced to 0-180 and divided by 3 | Raman ch. X, art. 135-140, pp. 101-104 | `ishta.go` |

## Where the code goes beyond the book
- **Bhava Bala** always uses Sripati's house middles, whatever house system the request names.
- **Cheshta, zodiac:** Raman's mean longitudes are in his own zodiac (the Raman ayanamsa); they are shifted into the
  chart's ayanamsa (`ramanFrameShift`) so the result does not depend on that choice. With the Raman ayanamsa the
  shift is 0 and his numbers are reproduced.
- **Cheshta, wrap:** the mean and true longitude are averaged across the 0/360 line.
- **Yuddha, "Kala up to Hora":** art. 76 excludes Ayana Bala from the aggregate; the summary (p. 107) says just Kala Bala.
  Art. 76 is followed.

## Where Raman and BPHS differ (Raman followed)
Bhava Bala (BPHS 27.26-31 adds +/- 1 rupa for the occupying planets and 15 virupas for the rising of the sign; Raman has only the three parts) and Ishta/Kashta Phala (BPHS ch. 28 uses the Uchcha and Chesta Rasmis; Kashta = 60 - Ishta); Saptavargaja values (BPHS 45/30/20/15/10/4/2) and the Moolatrikona (BPHS limits it to degrees of the sign, Raman uses the sign); Drekkana (BPHS translation: male 1st, female 2nd, hermaphrodite 3rd);
Ayana constant (BPHS 23 deg 27'); the Moon's Paksha Bala (not doubled in BPHS); Cheshta of the Sun and Moon (BPHS 27.18:
Ayana and Paksha Bala); Drik Bala (BPHS adds the whole aspect of Mercury and Jupiter); Yuddha (BPHS 27.20: the
difference of the two Shadbalas); the Sun's minimum (BPHS 390, Raman 300).

## Misprints in the book (the tests use the corrected values)
| Page | Printed | Correct |
|---|---|---|
| 89 | Sun total 424.24 | 424.94 (its components and the printed 7.08 rupas) |
| 89 | Mars Kala Bala 33.06 | 30.06 (p. 62) |
| 90 | Saturn Naisargika 8.37 | 8.57 (p. 80) |
| 87 | Moon's Drishti Pinda -86.23 | -86.93 (p. 88) |
| 83 | Jupiter on Mars 145 20' | 145 30' (from the longitudes) |
| 24 | Moon 49.750, Mercury 13.500, Venus 11.625 | 48.75, 135.0, 116.25 |
| 60 | Mars Ayana 1.84, Venus 23.75 | 1.94, 23.80 by his formula |
| 98 | header of example 59: bhava madhya of the 2nd house 311 10' | 331 10' (example 8 and the trisection; the printed Digbala 40 needs it) |
| 98 | Saturn's aspect on the 1st house 13.00 | 12.03 (its Drishti Kendra 174 04'; his malefic total 25.81 agrees) |
| 98 | Subhadrishtibala of the 4th house 62.89 | 61.89 (its parts 42.50 + 9.31 + 10.08), net 27.10 |
| 99 | Total Bhava Bala of the 12th house 426.76, 7.09 rupas | 526.76, 8.78 (389.21 + 40 + 97.55) |
| 99 | Total of the 7th house 7.00 rupas; Mars in the 10th house 298.24 | 7.09; 298.14 |
| 102 | Sun's Chesta Bala 22.66 | 22.61 (67 50' / 3) |
| 103 | Mars's Chesta 22.03 in the Ishta line; Kashta of Mercury 49.16 and Jupiter 13.19; Ishta of Venus 3.49 and Saturn 27.00; Kashta of Saturn 31.50 | 22.23; 17.81 and 9.53; 3.41 and 27.12; 31.29 (his own formula) |
| 111 | Ahargana at 31 Dec 2000 = 63,431 | 63,432 (a leap year) |
| 112, 114-117 | Table IV 146.0265; Table VI hundreds row 8; Table VII rows 8-9; Table VIII 133.32; Table IX correction 0.0001 t; Example 3 Jupiter 240.48 | 136.0265; inconsistent with its column; inconsistent with its column; 132.33; 0.001 t (Example 48); 239.08 |

## Known differences from the book's printed totals
- Ayana Bala for the whole chart (true declination against his table): Moon about -6, Sun +0.3, Mars -0.8, Mercury -1.1,
  Jupiter -1.3, Venus +1.7, Saturn -1.2. All other differences are rounding (<= 0.2).
- The Moon's benefic/malefic window for aspects (84 to 264 deg) is a reading of art. 53; the book's example Moon is benefic
  under any reading.
- A planet in more than one war accumulates; the book treats only pairs.
- Where the Sun does not rise or set (polar day or night), Tribhaga, Dina, Hora, Varsha and Masa Bala have nothing to count
  from, so the calculation returns `ErrNoSunriseSunset` and the API answers 400 ("the Sun does not rise or set at this
  latitude on this date") instead of returning zeros. All five use one lookup, `hinduDayAt` (latest sunrise, the sunset
  after it, the next sunrise).
- Whether Mercury "joins a malefic" is judged by sidereal sign (same sign as the Sun, Mars, Saturn, Rahu or Ketu), so when
  Mercury and such a planet are close to a sign boundary two ayanamsas can disagree about it; Paksha and Drik Bala of
  Mercury (and the aspects Mercury casts) then differ between them. Every other component is independent of the ayanamsa.
- The response also carries `rule_set` (`raman-shadbala-v1`), a `rank` per planet (by strength ratio, Raman example 57) and
  `tables` laid out as his examples 13, 35, 56 and 57 (`tables.go`); all parts are rounded to two decimals for display and the
  totals are formed from the unrounded values.
