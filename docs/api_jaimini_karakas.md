# Jaimini Chara Karakas API Endpoint

**Endpoint:** `POST /api/jaimini-karakas`

## Description
Calculates the 7-planet Jaimini Chara Karakas scheme. It ranks the classical planets (excluding Rahu/Ketu) strictly based on their fractional longitudinal degree *within their current sign* (from highest to lowest). 

The ranked karakas are: Atmakaraka (AK), Amatyakaraka (AmK), Bhratrikaraka (BK), Matrikaraka (MK), Pitrikaraka (PK), Gnatikaraka (GK), and Darakaraka (DK).

## Request Payload
Standard `BirthInput`.

## Response Structure
```json
{
  "calculation_time_utc": "2005-11-23T10:05:00Z",
  "planets": [
    {
      "planet": "Venus",
      "karaka": "AK",
      "source_longitude": 268.2616,
      "divisional_sign": "Sagittarius",
      "degree": 28,
      "minute": 15,
      "second": 42.1,
      "nakshatra": "Uttara Ashadha",
      "nakshatra_pada": 1,
      "nakshatra_lord": "Sun",
      "degree_in_sign": 28.2616,
      "retrograde": false
    },
    {
      "planet": "Mercury",
      "karaka": "AmK",
      "source_longitude": 234.1700,
      "divisional_sign": "Scorpio",
      "degree": 24,
      "minute": 10,
      "second": 12.0,
      "nakshatra": "Jyeshtha",
      "nakshatra_pada": 3,
      "nakshatra_lord": "Mercury",
      "degree_in_sign": 24.1700,
      "retrograde": false
    }
    // ... all 7 Karakas
  ]
}
```
