# Progression API Endpoint

**Endpoint:** `POST /api/progression`

## Description
Calculates a Secondary Progressed Chart ("A day for a year") based on the exact solar orbital mechanics. The system maps the native's age to days past the birth date, progressing all planetary and house cusp coordinates accordingly.

## Request Payload
Standard `BirthInput`, plus the `progression_year` target.

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
  "house_system": "Placidus",
  "progression_year": 2040
}
```

## Response Structure
```json
{
  "natal_date_utc": "2005-11-23T10:05:00Z",
  "target_progression_year": 2040,
  "age_in_years": 34.1,
  "progressed_date_utc": "2005-12-27T12:29:00Z",
  "progressed_julian_day": 2453732.02,
  "progressed_ayanamsa": 23.95,
  "ascendant": 12.35,
  "mc": 284.15,
  "progressed_planets": [
    {
      "planet": "Sun",
      "tropical_longitude": 275.4,
      "sidereal_longitude": 251.45,
      "speed": 1.019,
      "retrograde": false,
      "sign": "Sagittarius",
      "degree_in_sign": 11.45,
      "degree": 11,
      "minute": 27,
      "second": 0.0
    }
  ],
  "progressed_houses": [
    {
      "house_number": 1,
      "longitude": 12.35,
      "sign": "Aries"
    }
  ]
}
```
