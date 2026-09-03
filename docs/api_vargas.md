# Divisional Charts (Vargas) API Endpoint

**Endpoint:** `POST /api/vargas`

## Description
Generates a complete suite of Divisional Charts (Vargas) ranging from D1 (Rasi) to D60 (Shashtiamsha). Calculates exact divisional sign placements, ascending points, and lords for every harmonic fraction.

## Request Payload
Standard `BirthInput`.

## Response Structure
```json
{
  "vargas": [
    {
      "division": 9,
      "name": "Navamsha",
      "ascendant": {
        "planet": "Ascendant",
        "source_longitude": 5.3453,
        "divisional_sign": "Taurus",
        "degree": 18,
        "minute": 10,
        "second": 43.1,
        "nakshatra": "Rohini",
        "nakshatra_pada": 3,
        "sign_lord": "Venus",
        "retrograde": false
      },
      "planets": [
        {
          "planet": "Sun",
          "source_longitude": 217.259,
          "divisional_sign": "Cancer",
          "degree": 5,
          "minute": 31,
          "second": 12.0,
          "nakshatra": "Pushya",
          "nakshatra_pada": 1,
          "sign_lord": "Moon",
          "retrograde": false
        }
      ]
    }
  ]
}
```
