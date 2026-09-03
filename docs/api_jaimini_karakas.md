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
  "karakas": [
    {
      "planet": "Venus",
      "karaka": "AK",
      "sign": "Sagittarius",
      "degree": 28,
      "minute": 15,
      "second": 42.1,
      "degree_in_sign": 28.2616,
      "retrograde": false
    },
    {
      "planet": "Mercury",
      "karaka": "AmK",
      "sign": "Scorpio",
      "degree": 24,
      "minute": 10,
      "second": 12.0,
      "degree_in_sign": 24.1700,
      "retrograde": false
    }
    // ... all 7 Karakas
  ]
}
```
