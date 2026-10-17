# Shadbala API Endpoint

**Endpoint:** `POST /api/shadbala`

## Description
Six-fold planetary strength (Shadbala) of the seven planets Sun to Saturn, following B.V. Raman's *Graha and Bhava Balas*
(BPHS chapter 27 where Raman is silent). Every rule, its page reference and the book's misprints are in
`docs/shadbala_spec.md`. All values are in **shashtiamsas** (virupas); 60 shashtiamsas = 1 rupa.

## Request
Standard `BirthInput` (`docs/payload_birth_input.md`): date, time, latitude, longitude, timezone offset, `ayanamsa`,
`house_system`. The ayanamsa moves the planets between signs; the house system only decides the Ascendant and Midheaven
used by Dig Bala.

## Response
```json
{
  "rule_set": "raman-shadbala-v1",
  "unit": "shashtiamsas (virupas); 60 = 1 rupa",
  "calculation_time_utc": "2005-11-23T10:06:00Z",
  "planets": [
    {
      "planet": "Sun",
      "sthana_bala": {"uchcha_bala": 9.09, "saptavargaja_bala": 50.63, "ojayugma_bala": 0, "kendradi_bala": 30, "drekkana_bala": 15, "total": 104.72},
      "dig_bala": 42.96,
      "kala_bala": {"nathonnatha_bala": 41.65, "paksha_bala": 28.18, "tribhaga_bala": 0, "varsha_bala": 15, "masa_bala": 0, "dina_bala": 0, "hora_bala": 0, "ayana_bala": 8.99, "yuddha_bala": 0, "total": 93.82},
      "cheshta_bala": 0,
      "naisargika_bala": 60,
      "drik_bala": -31.7,
      "total_shadbala": 269.8,
      "rupas": 4.497,
      "minimum_required": 5,
      "strength_ratio": 0.899,
      "meets_minimum": false,
      "rank": 7
    },
    ...  // Moon, Mars, Mercury, Jupiter, Venus, Saturn in the same shape
  ],
  "tables": {
    "planets": ["Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"],
    "sthana_bala": {
      "title": "Sthana Bala (positional strength)",
      "rows": [
        {"component": "uchcha_bala", "label": "Uchcha Bala", "values": {"Jupiter": 27.63, "Mars": 33.93, "Mercury": 41.6, "Moon": 30.39, "Saturn": 29.12, "Sun": 9.09, "Venus": 28.53}},
        {"component": "saptavargaja_bala", "label": "Saptavargaja Bala", "values": {"Jupiter": 135, "Mars": 86.25, "Mercury": 101.25, "Moon": 120, "Saturn": 84.38, "Sun": 50.63, "Venus": 129.38}},
        {"component": "ojayugma_bala", "label": "Ojayugmarasyamsa Bala", "values": {"Jupiter": 15, "Mars": 30, "Mercury": 15, "Moon": 0, "Saturn": 15, "Sun": 0, "Venus": 0}},
        {"component": "kendradi_bala", "label": "Kendra Bala", "values": {"Jupiter": 60, "Mars": 60, "Mercury": 30, "Moon": 30, "Saturn": 60, "Sun": 30, "Venus": 15}},
        {"component": "drekkana_bala", "label": "Drekkana Bala", "values": {"Jupiter": 0, "Mars": 0, "Mercury": 15, "Moon": 0, "Saturn": 15, "Sun": 15, "Venus": 15}}
      ],
      "total": {"Jupiter": 237.63, "Mars": 210.18, "Mercury": 202.85, "Moon": 180.39, "Saturn": 203.5, "Sun": 104.72, "Venus": 187.91}
    },
    "kala_bala": {
      "title": "Kala Bala (temporal strength)",
      "rows": [
        {"component": "nathonnatha_bala", "label": "Nathonnatha Bala", "values": {"Jupiter": 41.65, "Mars": 18.35, "Mercury": 60, "Moon": 18.35, "Saturn": 18.35, "Sun": 41.65, "Venus": 41.65}},
        {"component": "paksha_bala", "label": "Paksha Bala", "values": {"Jupiter": 31.82, "Mars": 28.18, "Mercury": 28.18, "Moon": 63.63, "Saturn": 28.18, "Sun": 28.18, "Venus": 31.82}},
        {"component": "tribhaga_bala", "label": "Tribhaga Bala", "values": {"Jupiter": 60, "Mars": 0, "Mercury": 0, "Moon": 0, "Saturn": 60, "Sun": 0, "Venus": 0}},
        {"component": "varsha_bala", "label": "Abda (Varsha) Bala", "values": {"Jupiter": 0, "Mars": 0, "Mercury": 0, "Moon": 0, "Saturn": 0, "Sun": 15, "Venus": 0}},
        {"component": "masa_bala", "label": "Masa Bala", "values": {"Jupiter": 30, "Mars": 0, "Mercury": 0, "Moon": 0, "Saturn": 0, "Sun": 0, "Venus": 0}},
        {"component": "dina_bala", "label": "Vara (Dina) Bala", "values": {"Jupiter": 0, "Mars": 0, "Mercury": 45, "Moon": 0, "Saturn": 0, "Sun": 0, "Venus": 0}},
        {"component": "hora_bala", "label": "Hora Bala", "values": {"Jupiter": 0, "Mars": 0, "Mercury": 0, "Moon": 0, "Saturn": 60, "Sun": 0, "Venus": 0}},
        {"component": "ayana_bala", "label": "Ayana Bala", "values": {"Jupiter": 14.31, "Mars": 49.06, "Mercury": 55.77, "Moon": 9.44, "Saturn": 7.63, "Sun": 8.99, "Venus": 0}},
        {"component": "yuddha_bala", "label": "Yuddha Bala", "values": {"Jupiter": 0, "Mars": 0, "Mercury": 0, "Moon": 0, "Saturn": 0, "Sun": 0, "Venus": 0}}
      ],
      "total": {"Jupiter": 177.78, "Mars": 95.6, "Mercury": 188.96, "Moon": 91.42, "Saturn": 174.16, "Sun": 93.82, "Venus": 73.47}
    },
    "dig_bala": {
      "title": "Dig Bala (directional strength)",
      "rows": [
        {"component": "dig_bala", "label": "Dig Bala", "values": {"Jupiter": 2.16, "Mars": 24.06, "Mercury": 11.52, "Moon": 48.86, "Saturn": 33.91, "Sun": 42.96, "Venus": 1.94}}
      ],
      "total": {"Jupiter": 2.16, "Mars": 24.06, "Mercury": 11.52, "Moon": 48.86, "Saturn": 33.91, "Sun": 42.96, "Venus": 1.94}
    },
    ...  // cheshta_bala, naisargika_bala and drik_bala have the same shape (one row each)
    "ishta_kashta_phala": {"title": "Ishta and Kashta Phala (benefic and malefic tendency)", "rows": [...]},  // uchcha_bala, chesta_bala, ishta_phala, kashta_phala by planet; no total
    "shadbala_pinda": {
      "title": "Shadbala Pinda",
      "rows": [
        {"component": "sthana_bala", "label": "Sthana Bala", "values": {"Jupiter": 237.63, "Mars": 210.18, "Mercury": 202.85, "Moon": 180.39, "Saturn": 203.5, "Sun": 104.72, "Venus": 187.91}},
        {"component": "dig_bala", "label": "Dig Bala", "values": {"Jupiter": 2.16, "Mars": 24.06, "Mercury": 11.52, "Moon": 48.86, "Saturn": 33.91, "Sun": 42.96, "Venus": 1.94}},
        {"component": "kala_bala", "label": "Kala Bala", "values": {"Jupiter": 177.78, "Mars": 95.6, "Mercury": 188.96, "Moon": 91.42, "Saturn": 174.16, "Sun": 93.82, "Venus": 73.47}},
        {"component": "cheshta_bala", "label": "Cheshta Bala", "values": {"Jupiter": 8.82, "Mars": 53.21, "Mercury": 56.08, "Moon": 0, "Saturn": 38.74, "Sun": 0, "Venus": 43.17}},
        {"component": "naisargika_bala", "label": "Naisargika Bala", "values": {"Jupiter": 34.28, "Mars": 17.14, "Mercury": 25.71, "Moon": 51.43, "Saturn": 8.57, "Sun": 60, "Venus": 42.85}},
        {"component": "drik_bala", "label": "Drik Bala", "values": {"Jupiter": -39.16, "Mars": 5.71, "Mercury": -30.6, "Moon": -10.9, "Saturn": -12.7, "Sun": -31.7, "Venus": -8.71}}
      ],
      "total": {"Jupiter": 421.51, "Mars": 405.9, "Mercury": 454.52, "Moon": 361.2, "Saturn": 446.18, "Sun": 269.8, "Venus": 340.63},
      "rupas": {"Jupiter": 7.025, "Mars": 6.765, "Mercury": 7.575, "Moon": 6.02, "Saturn": 7.436, "Sun": 4.497, "Venus": 5.677},
      "minimum_required_rupas": {"Jupiter": 6.5, "Mars": 5, "Mercury": 7, "Moon": 6, "Saturn": 5, "Sun": 5, "Venus": 5.5},
      "strength_ratio": {"Jupiter": 1.081, "Mars": 1.353, "Mercury": 1.082, "Moon": 1.003, "Saturn": 1.487, "Sun": 0.899, "Venus": 1.032},
      "rank": {"Jupiter": 4, "Mars": 2, "Mercury": 3, "Moon": 6, "Saturn": 1, "Sun": 7, "Venus": 5}
    }
  },
  "ishta_kashta_phala": [
    {"planet": "Sun", "uchcha_bala": 9.09, "chesta_bala": 9.6, "ishta_phala": 9.34, "kashta_phala": 50.66, "inclination": "kashta"},
    {"planet": "Moon", "uchcha_bala": 30.39, "chesta_bala": 31.82, "ishta_phala": 31.1, "kashta_phala": 28.89, "inclination": "ishta"},
    {"planet": "Mars", "uchcha_bala": 33.93, "chesta_bala": 53.21, "ishta_phala": 42.49, "kashta_phala": 13.31, "inclination": "ishta"},
    {"planet": "Mercury", "uchcha_bala": 41.6, "chesta_bala": 56.08, "ishta_phala": 48.3, "kashta_phala": 8.49, "inclination": "ishta"},
    {"planet": "Jupiter", "uchcha_bala": 27.63, "chesta_bala": 8.82, "ishta_phala": 15.61, "kashta_phala": 40.7, "inclination": "kashta"},
    {"planet": "Venus", "uchcha_bala": 28.53, "chesta_bala": 43.17, "ishta_phala": 35.09, "kashta_phala": 23.02, "inclination": "ishta"},
    {"planet": "Saturn", "uchcha_bala": 29.12, "chesta_bala": 38.74, "ishta_phala": 33.59, "kashta_phala": 25.62, "inclination": "ishta"}
  ],
  "bhava_bala": {
    "bhavas": [
      {"house": 1, "name": "Thanu", "bhava_madhya": 5.64, "sign": "Aries", "lord": "Mars", "sign_group": "chathushpada", "bhavadhipathi_bala": 405.9, "bhava_digbala": 30, "subha_drishti_bala": 61.22, "asubha_drishti_bala": 12.65, "bhava_drishti_bala": 48.57, "total_bhava_bala": 484.47, "rupas": 8.07, "rank": 6},
      {"house": 2, "name": "Dhana", "bhava_madhya": 33.23, "sign": "Taurus", "lord": "Venus", "sign_group": "chathushpada", "bhavadhipathi_bala": 340.63, "bhava_digbala": 20, "subha_drishti_bala": 100.34, "asubha_drishti_bala": 29.57, "bhava_drishti_bala": 70.77, "total_bhava_bala": 431.39, "rupas": 7.19, "rank": 8},
      ...  // houses 3 to 12 in the same shape
    ],
    "strongest_house": 4,
    "weakest_house": 5
  }
}
```

### Top level
| Field | Meaning |
|---|---|
| `rule_set` | the rules followed, `raman-shadbala-v1` |
| `unit` | shashtiamsas (virupas); 60 = 1 rupa |
| `calculation_time_utc` | the birth moment in UTC |
| `planets` | one object per planet, in the order Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn |
| `tables` | the same numbers laid out as tables, one per bala, as in Raman's worked example |
| `ishta_kashta_phala` | benefic and malefic tendency of each planet (below) |
| `bhava_bala` | strength of the twelve houses (below) |

### `planets[]`
| Field | Meaning |
|---|---|
| `sthana_bala` | positional strength: `uchcha_bala`, `saptavargaja_bala`, `ojayugma_bala`, `kendradi_bala`, `drekkana_bala` and their `total` |
| `dig_bala` | directional strength |
| `kala_bala` | temporal strength: `nathonnatha_bala`, `paksha_bala`, `tribhaga_bala`, `varsha_bala`, `masa_bala`, `dina_bala`, `hora_bala`, `ayana_bala`, `yuddha_bala` and their `total` |
| `cheshta_bala` | motional strength; **0 for the Sun and Moon**, which have none in the total (Raman, example 56) |
| `naisargika_bala` | natural strength: 60, 51.43, 42.85, 34.28, 25.71, 17.14, 8.57 for the Sun, Moon, Venus, Jupiter, Mercury, Mars, Saturn |
| `drik_bala` | aspect strength; **can be negative** (malefic aspects exceed benefic ones) |
| `total_shadbala` | Shadbala Pinda: Sthana + Dig + Kala + Cheshta + Naisargika + Drik, in shashtiamsas |
| `rupas` | `total_shadbala` / 60 |
| `minimum_required` | the rupas at which the planet is strong: Sun 5, Moon 6, Mars 5, Mercury 7, Jupiter 6.5, Venus 5.5, Saturn 5 |
| `strength_ratio` | `rupas` / `minimum_required`; 1 or more means the planet has its required strength |
| `meets_minimum` | `strength_ratio` is at least 1 |
| `rank` | 1 is the strongest planet by `strength_ratio`, 7 the weakest |

### `tables`
`tables.planets` lists the planets in order. Each of `sthana_bala`, `dig_bala`, `kala_bala`, `cheshta_bala`,
`naisargika_bala` and `drik_bala` is a table with
- `title`,
- `rows`: one per part (`component` is the machine name, `label` the book's name for display) with `values`, the value for
  each planet by name,
- `total`: the total of the parts for each planet.

The Sthana table has 5 rows and the Kala table 9; the other four have one row, equal to their total.
`shadbala_pinda` is the final table: its `rows` are the six balas (Sthana, Dig, Kala, Cheshta, Naisargika, Drik), and it
adds `total` (shashtiamsas), `rupas`, `minimum_required_rupas`, `strength_ratio` and `rank`, each by planet. This is the
layout of Raman's examples 56 and 57.

Table values are rounded to two decimals; a table's `total` can differ from the sum of its rounded rows by a few hundredths.
Amounts in `planets` are rounded the same way.

### `ishta_kashta_phala[]`
One entry per planet (Raman chapter X).

| Field | Meaning |
|---|---|
| `uchcha_bala` | the Uchcha Bala of the planet (from Sthana Bala) |
| `chesta_bala` | its Chesta Bala; for the Sun and Moon this is their own value (Raman art. 136-137: the Sun's from its Sayana longitude plus 90 degrees, the Moon's from its arc from the Sun), which is not part of the Shadbala total |
| `ishta_phala` | benefic tendency: the square root of `uchcha_bala` x `chesta_bala` |
| `kashta_phala` | malefic tendency: the square root of (60 - `uchcha_bala`) x (60 - `chesta_bala`) |
| `inclination` | `ishta` if the benefic tendency is greater, otherwise `kashta`: whether the planet is inclined to do good or harm in its dasa and bhukti |

`tables.ishta_kashta_phala` has the same four values as rows, by planet.

### `bhava_bala`
The strength of the twelve houses (Raman chapter IX), in shashtiamsas. The middle of each house (`bhava_madhya`) is found by
Sripati's method whatever house system was chosen: the Ascendant, the lower meridian, the Descendant and the Midheaven are
the middles of the 1st, 4th, 7th and 10th houses, and the arc between two of them is divided in three for the houses between.

| Field | Meaning |
|---|---|
| `house`, `name` | 1 to 12 and Thanu, Dhana, Bhratru, Matru, Putra, Satru, Kalatra, Ayur, Bhagya, Karma, Labha, Vyaya |
| `bhava_madhya`, `sign` | sidereal longitude of the middle of the house and its sign |
| `lord` | the planet ruling that sign |
| `bhavadhipathi_bala` | the lord's Shadbala Pinda (`total_shadbala`) |
| `sign_group` | `nara` (Gemini, Virgo, Libra, Aquarius, first half of Sagittarius), `jalachara` (Cancer, Pisces, second half of Capricorn), `chathushpada` (Aries, Taurus, Leo, second half of Sagittarius, first half of Capricorn) or `keeta` (Scorpio) |
| `bhava_digbala` | 0 to 60 from the sign group: the 1st house is strongest in a Nara sign, the 4th in Jalachara, the 7th in Keeta, the 10th in Chathushpada, and weakest in the opposite house |
| `subha_drishti_bala`, `asubha_drishti_bala`, `bhava_drishti_bala` | benefic and malefic aspects on the house middle and their difference; Jupiter's and Mercury's aspects count whole and the others by a quarter; Mercury counts as a benefic here |
| `total_bhava_bala`, `rupas` | the three parts added, in shashtiamsas and rupas |
| `rank` | 1 is the strongest house |

`strongest_house` and `weakest_house` give the numbers of the strongest and weakest houses. BPHS 27.26-31 adds the planets
occupying a house and the rising of its sign; Raman does not, so they are not included.

### Reading the numbers
- A planet is considered strong when `strength_ratio` is 1 or more; the larger the ratio the stronger. `rank` orders the
  planets by it. In the example the Moon (1.003) only just reaches its requirement and the Sun (0.899) falls short.
- Sthana, Dig, Kala, Cheshta and Naisargika are never negative; Drik Bala is positive when benefic aspects dominate and
  negative when malefic ones do.
- `yuddha_bala` is non-zero only for planets less than one degree apart (the Sun and Moon never fight); the winner (smaller
  longitude) gains what the loser loses.
- Mercury counts as a malefic (for `paksha_bala` and for the aspects it casts) when it is in the same sign as the Sun, Mars,
  Saturn, Rahu or Ketu.
- Ayana Bala uses the true declination of the planet.
- The Sun's and Moon's Chesta Bala in `ishta_kashta_phala` is not added to their Shadbala; `cheshta_bala` in `planets` stays 0.

## Errors
| Status | When |
|---|---|
| `400` | the body or date/time cannot be read |
| `400` | `shadbala cannot be calculated: the Sun does not rise or set at this latitude on this date`: in polar day or night there is no sunrise to count Tribhaga, Dina, Hora, Varsha and Masa Bala from; try another date or place |
| `405` | the method is not POST |
| `500` | the ephemeris failed (for example a date outside 1800-2399) |

## Notes for clients
- Field names, planet names and `component` names are stable English identifiers; show the `label`, or translate the
  `component`.
- The response is for one chart; the `tables` repeat the numbers of `planets`, so a client can render either.
