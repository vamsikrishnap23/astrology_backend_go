# Feature: Ashtakoota / Guna Milan

## Overview
This implements the traditional 8-Koota (36-point) compatibility matching algorithm for Jyotish, based directly on exact Swiss Ephemeris astronomical Moon positions. 

## The 8 Kootas & Maximum Points
1. **Varna** (1 point): Spiritual/ego compatibility based on Moon signs.
2. **Vashya** (2 points): Dominance/influence compatibility based on Moon signs.
3. **Tara/Dina** (3 points): Daily destiny/health compatibility based on Nakshatras.
4. **Yoni** (4 points): Intimacy and physical compatibility based on Nakshatra animals.
5. **Graha Maitri** (5 points): Psychological/mental compatibility based on Moon sign Lords.
6. **Gana** (6 points): Temperament compatibility based on Nakshatra Gana (Deva/Manushya/Rakshasa).
7. **Bhakoot** (7 points): Family and welfare compatibility based on Moon sign relative distances (2/12, 5/9, 6/8 are doshas).
8. **Nadi** (8 points): Genetic and health compatibility based on Nakshatra Nadi (Adi, Madhya, Antya).

Total Maximum Points: **36**.

## Rule Set
**Version**: `classical-guna-milan-v1`
The lookup tables for Vashya, Yoni, and Gana strictly follow the canonical texts (Brihat Parashara Hora Shastra / B.V. Raman tables).
- **Vashya Mapping**: Simplified canonical (Aries/Taurus/Sagittarius=Chatuspada, Gemini/Virgo/Libra/Aquarius=Manav, Cancer/Capricorn/Pisces=Jalchar, Leo=Vanchar, Scorpio=Keeta).
- **Yoni Table**: Complete 14-animal matrix.
- **Bhakoot Cancellation**: Triggered automatically if the Rashi lords are identical or friendly (Graha Maitri score >= 4).
- **Nadi Cancellation**: Triggered automatically if the Nakshatras are identical but Padas differ, or if the Rashis are identical but Nakshatras differ.

## Output
The API outputs both a `raw_score` and an `effective_score` (if cancellations apply), preserving full transparency.

## API Endpoint
`POST /api/v1/ashtakoota`

**Payload:**
```json
{
  "groom": {
    "name": "Groom",
    "date_of_birth": "2005-11-23",
    "time_of_birth": "15:35:00",
    "latitude": 16.3938,
    "longitude": 80.1522,
    "timezone": 5.5,
    "ayanamsa": "Lahiri"
  },
  "bride": {
    "name": "Bride",
    "date_of_birth": "2002-05-14",
    "time_of_birth": "10:00:00",
    "latitude": 12.9716,
    "longitude": 77.5946,
    "timezone": 5.5,
    "ayanamsa": "Lahiri"
  }
}
```
