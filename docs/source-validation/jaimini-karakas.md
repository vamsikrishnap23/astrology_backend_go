# Feature 24: Jaimini Chara Karakas

## Overview
This document outlines the Jaimini Chara Karakas calculation logic, providing dynamic 7-karaka planetary assignments based on longitudinal degrees.

## 7-Karaka System
Chara Karakas (variable significators) represent different souls/figures in a person's life. The 7-karaka scheme implemented does **not** include Rahu or Ketu. The available significators are:

1. **AK (Atmakaraka):** Planet with the highest degree in its sign.
2. **AmK (Amatyakaraka):** Second highest degree.
3. **BK (Bhratrikaraka):** Third highest degree.
4. **MK (Matrikaraka):** Fourth highest degree.
5. **PK (Pitrikaraka):** Fifth highest degree.
6. **GK (Gnatikaraka):** Sixth highest degree.
7. **DK (Darakaraka):** Lowest degree.

## Included Planets
- Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn.
*(Rahu, Ketu, Uranus, Neptune, Pluto are strictly ignored).*

## Ranking Method
Ranking is based exclusively on a planet's longitude **within its current sign** (0°00'00" to 29°59'59"). The absolute zodiacal longitude is modulated by 30° to achieve the sign-agnostic fractional degree.

All calculations rely on exact Swiss Ephemeris sidereal longitudes without pre-rounding. 

### Tie-Handling
If two planets hold the exact same raw floating-point degree-within-sign (an astronomically near-impossible event but computationally possible in synthetic tests), the engine falls back to a deterministic resolution using the traditional week-day planetary sequence: Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn. The planet appearing earlier in the sequence wins the tie.

## Validation & Regression
The calculation is continuously validated against the reference chart for Vamsi (2005-11-23 15:35:00 IST at Sattenapalle, Lahiri Ayanamsa). 

Expected degree-in-sign and order:
1. **Venus** (22°34'28") -> AK
2. **Saturn** (17°22'18") -> AmK
3. **Mars** (16°12'46") -> BK
4. **Jupiter** (12°07'15") -> MK
5. **Mercury** (10°12'15") -> PK
6. **Sun** (07°16'32") -> GK
7. **Moon** (01°49'11") -> DK

## API Endpoint
`POST /api/v1/jaimini-karakas`

Returns:
```json
{
  "calculation_time_utc": "...",
  "karakas": [
    {
      "planet": "Venus",
      "karaka": "AK",
      "sign": "Sagittarius",
      "degree": 22,
      "minute": 34,
      "second": 28.169,
      "degree_in_sign": 22.574491
    }
  ]
}
```
