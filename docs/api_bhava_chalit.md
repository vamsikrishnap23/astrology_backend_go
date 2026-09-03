# Bhava Chalit API Endpoint

**Endpoint:** `POST /api/bhava-chalit`

## Description
Generates a traditional Bhava Chalit chart, which realigns the planetary placements based on exact mathematical house cusps rather than fixed 30-degree Rasi sign bounds. 

## Request Payload
Standard `BirthInput`.

## Response Structure
```json
{
  "ascendant": 5.3453,
  "houses": [
    {
      "house_number": 1,
      "cusp_longitude": 5.3453,
      "sign": "Aries",
      "degree": 5,
      "minute": 20,
      "second": 43.1,
      "nakshatra": "Ashwini",
      "nakshatra_pada": 2,
      "nakshatra_lord": "Ketu",
      "occupants": [
        {
          "planet_name": "Rahu",
          "house_number": 1,
          "sign": "Pisces",
          "degree": 17,
          "minute": 5,
          "second": 12.0,
          "exact_longitude": 347.086,
          "nakshatra": "Revati",
          "nakshatra_pada": 1,
          "nakshatra_lord": "Mercury"
        }
      ]
    },
    {
      "house_number": 2,
      "cusp_longitude": 36.4,
      "sign": "Taurus",
      "degree": 6,
      "minute": 24,
      "second": 0.0,
      "occupants": []
    }
  ]
}
```
