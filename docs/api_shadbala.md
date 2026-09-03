# Shadbala API Endpoint

**Endpoint:** `POST /api/shadbala`

## Description
Calculates the 6-fold planetary strength (Shadbala) for the 7 classical planets exactly following the classical Brihat Parashara Hora Shastra mathematical formulas. Determines the ultimate strength in Rupas, and compares it against the minimum required thresholds.

## Request Payload
Standard `BirthInput`.

## Response Structure
```json
{
  "calculation_time_utc": "2005-11-23T10:05:00Z",
  "planets": [
    {
      "planet": "Sun",
      "sthana_bala": {
        "uchcha_bala": 14.5,
        "saptavargaja_bala": 112.5,
        "ojayugma_bala": 15.0,
        "kendradi_bala": 60.0,
        "drekkana_bala": 15.0,
        "total": 217.0
      },
      "dig_bala": 22.4,
      "kala_bala": {
        "nathonnatha_bala": 45.2,
        "paksha_bala": 30.1,
        "tribhaga_bala": 60.0,
        "varsha_bala": 15.0,
        "masa_bala": 30.0,
        "dina_bala": 45.0,
        "hora_bala": 60.0,
        "ayana_bala": 24.5,
        "yuddha_bala": 0.0,
        "total": 309.8
      },
      "cheshta_bala": 45.2,
      "naisargika_bala": 60.0,
      "drik_bala": 12.5,
      "total_shadbala": 666.9,
      "rupas": 11.115,
      "minimum_required": 5.0,
      "strength_ratio": 2.223,
      "meets_minimum": true
    }
  ]
}
```
