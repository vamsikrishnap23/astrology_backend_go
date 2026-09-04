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
  "planets": [
    {
      "planet": "Rahu",
      "source_longitude": 347.086,
      "divisional_sign": "Pisces",
      "degree": 17,
      "minute": 5,
      "second": 12.0,
      "nakshatra": "Revati",
      "nakshatra_pada": 1,
      "nakshatra_lord": "Mercury",
      "sign_lord": "Jupiter",
      "retrograde": true,
      "house_number": 1
    }
  ],
  "houses": [
    {
      "house_number": 1,
      "longitude": 5.3453,
      "sign": "Aries",
      "degree": 5,
      "minute": 20,
      "second": 43.1,
      "nakshatra": "Ashwini",
      "nakshatra_pada": 2,
      "nakshatra_lord": "Ketu"
    },
    {
      "house_number": 2,
      "longitude": 36.4,
      "sign": "Taurus",
      "degree": 6,
      "minute": 24,
      "second": 0.0,
      "nakshatra": "Krittika",
      "nakshatra_pada": 3,
      "nakshatra_lord": "Sun"
    }
  ]
}
```
