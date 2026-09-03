# Feature 23: Shadbala (Six-fold Planetary Strength)

## Overview & Authority
This document establishes the classical algorithms used for Feature 23 (Shadbala). 
The formulas are derived from **Brihat Parashara Hora Shastra (BPHS)** and mathematically formulated based on the standard computational text **"Graha and Bhava Balas" by Dr. B.V. Raman**.

All calculations are performed mathematically (e.g., using angular distances) rather than simple lookups.
Shadbala is calculated for the 7 primary planets: Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn.
Output values are represented in **Shashtiamsas** (1 Rupa = 60 Shashtiamsas).

---

## 1. Sthana Bala (Positional Strength)
Comprises 5 sub-components:

### 1.1 Uchcha Bala (Exaltation Strength)
Measures distance from deep debilitation.
- **Exaltation Points:** Sun (Aries 10°), Moon (Taurus 3°), Mars (Capricorn 28°), Mercury (Virgo 15°), Jupiter (Cancer 5°), Venus (Pisces 27°), Saturn (Libra 20°).
- **Debilitation Points:** Exaltation + 180°.
- **Formula:** 
  `Arc = |Planet_Longitude - Debilitation_Longitude|`
  `If Arc > 180°: Arc = 360° - Arc`
  `Uchcha Bala = Arc / 3` (Max 60, Min 0)

### 1.2 Saptavargaja Bala (Divisional Strength)
Based on placement in 7 divisional charts (Rasi, Hora, Drekkana, Saptamsha, Navamsha, Dwadashamsha, Trishamsha).
- **Rule Set (B.V. Raman):** 
  - Moolatrikona: 45
  - Swakshetra (Own House): 30
  - Adhi Mitra (Extreme Friend): 22.5
  - Mitra (Friend): 15
  - Sama (Neutral): 7.5
  - Shatru (Enemy): 3.75
  - Adhi Shatru (Extreme Enemy): 1.875
*(Note: Total Saptavargaja Bala is the sum of these values across all 7 Vargas)*

### 1.3 Ojayugma Bala (Odd/Even Sign Strength)
Evaluated in Rasi and Navamsha charts.
- **Rule:** 
  - Venus and Moon get 15 Shashtiamsas if in an Even (Yugma) sign.
  - Sun, Mars, Jupiter, Mercury, Saturn get 15 Shashtiamsas if in an Odd (Oja) sign.
  - Applied separately in both Rasi and Navamsha (Max 30 total).

### 1.4 Kendradi Bala (Angular Strength)
Based on Rasi house placement (1-indexed from Ascendant).
- **Kendra (Angles):** 1, 4, 7, 10 = 60 Shashtiamsas
- **Panaphara (Succedent):** 2, 5, 8, 11 = 30 Shashtiamsas
- **Apoklima (Cadent):** 3, 6, 9, 12 = 15 Shashtiamsas

### 1.5 Drekkana Bala (Decanate Strength)
Based on which 10° segment (Drekkana) of a sign the planet occupies.
- Male planets (Sun, Mars, Jupiter): 15 pts in 1st Drekkana (0°-10°).
- Female planets (Moon, Venus): 15 pts in 2nd Drekkana (10°-20°).
- Eunuch planets (Mercury, Saturn): 15 pts in 3rd Drekkana (20°-30°).
- Otherwise 0.

---

## 2. Dig Bala (Directional Strength)
Measures angular distance from a planet's weakest direction (zero power point).
- **Max Power Points (100% = 60 Shashtiamsas):**
  - Sun, Mars: 10th House cusp (Midheaven / MC).
  - Jupiter, Mercury: 1st House cusp (Ascendant / East).
  - Saturn: 7th House cusp (Descendant / West).
  - Moon, Venus: 4th House cusp (IC / North).
- **Zero Power Points (0%):** Max Power + 180°.
- **Formula:** 
  `Arc = |Planet_Longitude - Zero_Power_Longitude|`
  `If Arc > 180°: Arc = 360° - Arc`
  `Dig Bala = Arc / 3` (Max 60)

---

## 3. Kala Bala (Temporal Strength)
Comprises multiple sub-components based on time.

### 3.1 Nathonnatha Bala (Day/Night Strength)
- Evaluates birth time relative to midnight/noon.
- **Moon, Mars, Saturn:** Strong at midnight (60), zero at noon (0). 
  `Formula: |Time_Difference_From_Noon_in_Hours| * 5`
- **Sun, Jupiter, Venus:** Strong at noon (60), zero at midnight (0).
  `Formula: |Time_Difference_From_Midnight_in_Hours| * 5`
- **Mercury:** Always 60.

### 3.2 Paksha Bala (Lunar Phase Strength)
- Evaluates the angle between Sun and Moon.
- `Angle = Moon_Longitude - Sun_Longitude` (wrapped 0-360)
- `If Angle > 180: Angle = 360 - Angle`
- `Base = Angle / 3` (Max 60)
- Benefics (Jupiter, Venus, Moon, Mercury if benefic): Bala = Base
- Malefics (Sun, Mars, Saturn, Mercury if malefic): Bala = 60 - Base
- **Special Moon Rule:** Moon's Paksha Bala is always multiplied by 2 (Max 120).

### 3.3 Tribhaga Bala (Day/Night Thirds Strength)
- Day and Night are each divided into 3 equal parts.
- Depending on the birth segment, one specific planet receives 60 points:
  - Day Part 1: Mercury
  - Day Part 2: Sun
  - Day Part 3: Saturn
  - Night Part 1: Moon
  - Night Part 2: Venus
  - Night Part 3: Mars
  - Jupiter always receives 60 points regardless of birth time.

### 3.4 Varsha, Masa, Dina, Hora Bala (Year, Month, Day, Hour Lords)
- Evaluated based on Astrological/Vedic calendar rules.
- **Dina (Weekday) Lord:** 45 points to the ruler of the weekday.
- *(Note: Simplified versions often skip Varsha/Masa/Hora without a full panchang engine, but we will implement Dina Bala fully and assign 0 to others if unresolved, OR implement the full Ahargana-based Hora/Masa/Varsha lords if feasible).* 

### 3.5 Ayana Bala (Equinoctial Strength)
- Based on planetary declination (Krantivritta).
- Maximum Declination = 23°27' (or 24° classically).
- Sun, Mars, Jupiter, Venus gain strength in Northern Declination (Max 60 at +24°, Min 0 at -24°).
- Moon, Saturn gain strength in Southern Declination (Max 60 at -24°, Min 0 at +24°).
- Mercury is always 30 (or scales differently based on tradition; B.V. Raman says proportional to declination, but often just given max at both extremes? We will use: 60 at North, 60 at South, 30 at equator).
- **Sun Special Rule:** Sun's Ayana Bala is always multiplied by 2 (Max 120).

### 3.6 Yuddha Bala (Planetary War)
*(Deferred or simplified to 0 for this phase unless exactly conjunct).*

---

## 4. Cheshta Bala (Motional Strength)
- Based on apparent speed and retrogradation.
- Retrogression (Vakra): 60
- Sun & Moon do not retrograde. Their Cheshta Bala = Ayana Bala (often mapped this way in BPHS).

---

## 5. Naisargika Bala (Natural Strength)
Fixed values based on intrinsic planetary luminosity/power.
- Sun: 60.00 (1.0 Rupa)
- Moon: 51.43 (0.857 Rupa)
- Venus: 42.85 (0.714 Rupa)
- Jupiter: 34.28 (0.571 Rupa)
- Mercury: 25.71 (0.428 Rupa)
- Mars: 17.14 (0.285 Rupa)
- Saturn: 8.57 (0.142 Rupa)

---

## 6. Drik Bala (Aspect Strength)
Calculated mathematically based on angular distance between interacting planets.
- Formula evaluates the longitudinal angle (Distance = Aspecting - Aspected).
- Full aspect (60 points) typically occurs at 180°.
- Special full aspects: Mars (90°, 210°), Jupiter (120°, 240°), Saturn (60°, 270°).
- Benefic aspects add to total; malefic aspects subtract.
- Total Drik Bala is a quarter of the aspect sum.

---
## Summary of Total and Minimums
- Total Shadbala = Sum of the 6 Balas.
- Converted to Rupas (Total / 60).
- Minimum Requirements (in Rupas): Sun (6.5), Moon (6.0), Mars (5.0), Mercury (7.0), Jupiter (6.5), Venus (5.5), Saturn (5.0).
