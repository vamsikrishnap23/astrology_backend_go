# Tables API Endpoint

**Endpoint:** `POST /api/tables`

## Description
Generates exhaustive mathematical planetary and house tables based on Krishnamurti Paddhati (KP) sub-lord logic and Placidus/selected house system.

## Request Payload

Standard `BirthInput` struct.

```json
{
  "name": "Vamsi",
  "date_of_birth": "2005-11-23",
  "time_of_birth": "15:35:00",
  "place_of_birth": "Sattenapalle",
  "latitude": 16.3938,
  "longitude": 80.1522,
  "timezone": 5.5,
  "ayanamsa": "Lahiri",
  "house_system": "Placidus"
}
```

## Response Structure

```json
{
  "planetary_table": [
    {
      "planet_name": "Sun",
      "sign": "Scorpio",
      "degree": 7,
      "minute": 15,
      "second": 32.5,
      "exact_longitude": 217.25902,
      "retrograde": false,
      "speed": 1.002,
      "house_number": 8,
      "nakshatra": "Anuradha",
      "nakshatra_pada": 2,
      "sign_lord": "Mars",
      "nakshatra_lord": "Saturn",
      "sub_lord": "Mercury",
      "sub_sub_lord": "Venus",
      "sub_sub_sub_lord": "Sun"
    }
    // ... 9 planetary objects (Sun to Ketu)
  ],
  "house_table": [
    {
      "house_number": 1,
      "cusp_longitude": 5.3453,
      "sign": "Aries",
      "degree": 5,
      "minute": 20,
      "second": 43.1,
      "nakshatra": "Ashwini",
      "nakshatra_pada": 2,
      "sign_lord": "Mars",
      "nakshatra_lord": "Ketu",
      "sub_lord": "Mars",
      "sub_sub_lord": "Jupiter",
      "sub_sub_sub_lord": "Venus",
      "occupants": ["Moon", "Rahu"]
    }
    // ... 12 house objects
  ]
}
```
