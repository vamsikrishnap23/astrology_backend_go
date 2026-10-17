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
        {"component": "saptavargaja_bala", "label": "Saptavargaja Bala", "values": {"Jupiter": 135, "Mars": 71.25, "Mercury": 101.25, "Moon": 120, "Saturn": 84.38, "Sun": 50.63, "Venus": 129.38}},
        {"component": "ojayugma_bala", "label": "Ojayugmarasyamsa Bala", "values": {"Jupiter": 15, "Mars": 30, "Mercury": 15, "Moon": 0, "Saturn": 15, "Sun": 0, "Venus": 0}},
        {"component": "kendradi_bala", "label": "Kendra Bala", "values": {"Jupiter": 60, "Mars": 60, "Mercury": 30, "Moon": 30, "Saturn": 60, "Sun": 30, "Venus": 15}},
        {"component": "drekkana_bala", "label": "Drekkana Bala", "values": {"Jupiter": 0, "Mars": 0, "Mercury": 15, "Moon": 0, "Saturn": 15, "Sun": 15, "Venus": 15}}
      ],
      "total": {"Jupiter": 237.63, "Mars": 195.18, "Mercury": 202.85, "Moon": 180.39, "Saturn": 203.5, "Sun": 104.72, "Venus": 187.91}
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
    "shadbala_pinda": {
      "title": "Shadbala Pinda",
      "rows": [
        {"component": "sthana_bala", "label": "Sthana Bala", "values": {"Jupiter": 237.63, "Mars": 195.18, "Mercury": 202.85, "Moon": 180.39, "Saturn": 203.5, "Sun": 104.72, "Venus": 187.91}},
        {"component": "dig_bala", "label": "Dig Bala", "values": {"Jupiter": 2.16, "Mars": 24.06, "Mercury": 11.52, "Moon": 48.86, "Saturn": 33.91, "Sun": 42.96, "Venus": 1.94}},
        {"component": "kala_bala", "label": "Kala Bala", "values": {"Jupiter": 177.78, "Mars": 95.6, "Mercury": 188.96, "Moon": 91.42, "Saturn": 174.16, "Sun": 93.82, "Venus": 73.47}},
        {"component": "cheshta_bala", "label": "Cheshta Bala", "values": {"Jupiter": 8.82, "Mars": 53.21, "Mercury": 56.08, "Moon": 0, "Saturn": 38.74, "Sun": 0, "Venus": 43.17}},
        {"component": "naisargika_bala", "label": "Naisargika Bala", "values": {"Jupiter": 34.28, "Mars": 17.14, "Mercury": 25.71, "Moon": 51.43, "Saturn": 8.57, "Sun": 60, "Venus": 42.85}},
        {"component": "drik_bala", "label": "Drik Bala", "values": {"Jupiter": -39.16, "Mars": 5.71, "Mercury": -30.6, "Moon": -10.9, "Saturn": -12.7, "Sun": -31.7, "Venus": -8.71}}
      ],
      "total": {"Jupiter": 421.51, "Mars": 390.9, "Mercury": 454.52, "Moon": 361.2, "Saturn": 446.18, "Sun": 269.8, "Venus": 340.63},
      "rupas": {"Jupiter": 7.025, "Mars": 6.515, "Mercury": 7.575, "Moon": 6.02, "Saturn": 7.436, "Sun": 4.497, "Venus": 5.677},
      "minimum_required_rupas": {"Jupiter": 6.5, "Mars": 5, "Mercury": 7, "Moon": 6, "Saturn": 5, "Sun": 5, "Venus": 5.5},
      "strength_ratio": {"Jupiter": 1.081, "Mars": 1.303, "Mercury": 1.082, "Moon": 1.003, "Saturn": 1.487, "Sun": 0.899, "Venus": 1.032},
      "rank": {"Jupiter": 4, "Mars": 2, "Mercury": 3, "Moon": 6, "Saturn": 1, "Sun": 7, "Venus": 5}
    }
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
