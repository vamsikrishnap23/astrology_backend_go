# Birth Time Rectification (BTR) - Nadi Tatwa Theory

This document outlines the detailed calculation logic for the Nadi Astrology-based Birth Time Rectification (BTR) system, utilizing Tatwa Shodhana (Elemental Rectification) and 3-minute Nadi Row offsets.

## 1. Time Normalization
Convert the raw clock time (IST/Standard Time) into an idealized astrological time that accounts for local geography and seasonal sunrise variations.

### Step 1.1: Local Mean Time (LMT) Correction
* **Formula:** `Correction = (Standard Meridian - Birth Longitude) × 240 seconds`
* **Rule:** If the birth city is *West* of the standard meridian, ADD the correction. If *East*, SUBTRACT it.
* **LMT:** `LMT = Birth Time + LMT Correction`

### Step 1.2: Sunrise Adjustment
The astrological day begins at exactly 06:00:00 AM in this specific Nadi system.
* **Formula:** `Sunrise Difference = Actual Local Sunrise - 06:00:00`
* **Final Adjusted Time:** `Final Time = LMT - Sunrise Difference`

## 2. Nadi Row Calculation
The 24-hour day (1440 minutes) is divided into 480 discrete "rows" or intervals of exactly 3 minutes each.
* **Minutes:** Convert Final Adjusted Time into total minutes from midnight.
* **Row Number:** `Row = floor(Minutes / 3) + 1`
* **T1 Time (End of Interval):** `Row × 3 minutes`

## 3. Tatwa Determination (The 90-Minute Cycle)
Every 3-minute interval is governed by a Tatwa (Element) and a Gender. The cycle of Tatwas operates in exact 90-minute waves.

### The 5 Tatwas & Their Exact Durations
Within a 90-minute wave, elements do not have equal time. They are strictly distributed:
1. **Prithvi (Earth) / Male:** 6 minutes (2 Rows)
2. **Jala (Water) / Female:** 12 minutes (4 Rows)
3. **Tejo (Fire) / Male:** 18 minutes (6 Rows)
4. **Vayu (Air) / Female:** 24 minutes (8 Rows)
5. **Akash (Ether) / Male:** 30 minutes (10 Rows)
*(Total = 90 minutes = 30 Rows).*

### Starting Tatwa by Weekday
The 90-minute cycle begins exactly at 06:00 AM (adjusted). The *first* Tatwa in the sequence depends on the weekday:
* **Wednesday:** Prithvi ➔ Jala ➔ Tejo ➔ Vayu ➔ Akash
* **Monday & Friday:** Jala ➔ Tejo ➔ Vayu ➔ Akash ➔ Prithvi
* **Sunday & Tuesday:** Tejo ➔ Vayu ➔ Akash ➔ Prithvi ➔ Jala
* **Saturday:** Vayu ➔ Akash ➔ Prithvi ➔ Jala ➔ Tejo
* **Thursday:** Akash ➔ Prithvi ➔ Jala ➔ Tejo ➔ Vayu

### Aroha vs Avaroha (The 90-Minute Flip)
* **Aroha (Forward):** The sequence runs forward for the first 90 minutes.
* **Avaroha (Reverse):** For the next 90 minutes, the sequence runs perfectly backward (e.g., if Wednesday's Aroha ends on Akash, the Avaroha immediately starts on Akash and runs backward: Akash ➔ Vayu ➔ Tejo ➔ Jala ➔ Prithvi).
* A full 24-hour day consists of 16 half-cycles (8 Aroha, 8 Avaroha).

## 4. Planet 90-Min and Vinod Lookups
The engine cross-references the 3-minute Nadi Row with the Ascendant's sign type (Movable, Fixed, or Dual) to find the governing planet.

### The 9-Planet Sequence
`0: Sun, 1: Moon, 2: Mars, 3: Rahu, 4: Jupiter, 5: Saturn, 6: Mercury, 7: Ketu, 8: Venus`

### The Modulo-9 Calculation
The system calculates the ruling planet using Modulo-9 arithmetic and a specific offset based on the Ascendant type:
* **Movable (Aries, Cancer, Libra, Capricorn):** Offset = 0 (Starts at Sun)
* **Fixed (Taurus, Leo, Scorpio, Aquarius):** Offset = 2 (Starts at Mars)
* **Dual (Gemini, Virgo, Sagittarius, Pisces):** Offset = 4 (Starts at Jupiter)

**Formula:**
`Planet Index = ( (Row - 1) + Offset ) % 9`
*(The resulting index is mapped to the 9-Planet Sequence above).*

## 5. BTR Validation
The final step of the calculator compares the generated data against the known facts of the native:
1. **Gender Match:** Does the calculated Tatwa Gender match the Native's Gender?
2. **Star Lord Match:** Does the 90-Min/Vinod planet match the Lord of the Native's Birth Nakshatra?

If both match, the birth time is considered mathematically "Rectified". If not, the system indicates that the time should be nudged backward or forward to find the closest valid 3-minute Nadi row.
