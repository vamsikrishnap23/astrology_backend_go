# Feature: Ashtakoota / Guna Milan

## Overview
The traditional 8-koota (36-point) compatibility matching for Jyotish, based on the sidereal Moon positions from the Swiss Ephemeris (Lahiri ayanamsa, always). Besides the eight scored kootas, the response carries a set of non-scored compatibility checks (`checks`) and an ascendant comparison (`lagna`). Endpoint: `POST /api/vivaha-pontana`. Field reference: `docs/api_ashtakoota.md`.

## The 8 Kootas & Maximum Points
1. **Varna** (1): groom's varna equal to or above the bride's, by Moon sign.
2. **Vashya** (2): the Telugu book's rasi-pair table (pp. 69-70); 2 if either rasi lists the other.
3. **Tara** (3): nakshatras counted both ways; 1.5 per favourable direction.
4. **Yoni** (4): animal compatibility of the birth nakshatras.
5. **Graha Maitri** (5): natural friendship of the two Moon-sign lords.
6. **Gana** (6): Deva / Manushya / Rakshasa; the table is read as [bride gana][groom gana].
7. **Bhakoot** (7): Moon-sign distances; 2/12, 5/9, 6/8 are doshas.
8. **Nadi** (8): Adi / Madhya / Antya; same nadi is a dosha.

Total maximum: **36**; `traditional_threshold_met` when the effective total is 18 or more.

## Rule Set
**Version**: `classical-guna-milan-v2`.
- **Vashya**: the book's rasi-pair table. The earlier 5-animal-category rule was removed.
- **Gana**: table read as [bride][groom] (a Deva bride with a Manushya groom scores 6, a Manushya bride with a Deva groom 5, a Deva bride with a Rakshasa groom 0, a Rakshasa bride with a Deva groom 1). A Rakshasa bride whose nakshatra is Ashlesha, Chitra, Vishakha, Jyeshtha, Mula or Shatabhisha scores the full 6. When the gana is mismatched but the bride-to-groom nakshatra count is more than 14, a `suggestion` says the harm is greatly reduced (book p. 65); the score is not changed.
- **Bhakoot cancellation**: an auspicious sign pair, or Graha Maitri of 4 or more.
- **Nadi cancellation**: same nakshatra with a different pada, or same rasi with a different nakshatra. If both nakshatras are in the Jaimini Maharshi list of 21, a `suggestion` states the rule; the score is **not** changed.
- **Checks (not scored)**: Rajju, Vedha, Mahendra, Stree Deergha, Mruga Vairam, Special Tara, Kala Sarpa. Mahendra, Stree Deergha and Vedha come from the Telugu muhurta book pages kept in `ref/`.

## Output
Every koota reports `raw_score` (before cancellation) and `score` (counted in the total), so a cancellation is always visible. Every check reports `satisfied`.

## Sources and validation
- Telugu muhurta book pages (`ref/`): Dina, Gana distance, Mahendra, Stree Deergha, Rasi koota, Vashya, Vedha.
- Review of every calculation against the book and online sources: `docs/matchmaking_report_te.pdf` (Telugu), with the owner's decisions recorded in `CLAUDE.md`.
- Open points still to verify against the printed book are listed in the TODO table in `CLAUDE.md` (rows starting "Matchmaking").
