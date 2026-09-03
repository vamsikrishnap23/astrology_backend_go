# Ashtakavarga API Endpoint

**Endpoint:** `POST /api/ashtakavarga`

## Description
Calculates the Bhinnashtakavarga (BAV) for the 7 classical planets and the composite Sarvashtakavarga (SAV) points for all 12 signs, precisely following Brihat Parashara Hora Shastra rulesets. 

## Request Payload
Standard `BirthInput`.

## Response Structure
```json
{
  "calculation_time_utc": "2005-11-23T10:05:00Z",
  "bav": [
    {
      "planet": "Sun",
      "total_bindus": 48,
      "signs": [
        {
          "sign": "Aries",
          "sign_index": 1,
          "total_bindus": 4,
          "contributions": [
            { "source_planet": "Sun", "value": 1 },
            { "source_planet": "Mars", "value": 1 }
          ]
        }
      ]
    }
  ],
  "sav": [
    {
      "sign": "Aries",
      "sign_index": 1,
      "total_bindus": 30
    },
    {
      "sign": "Taurus",
      "sign_index": 2,
      "total_bindus": 30
    }
  ],
  "total_sav_bindus": 337
}
```
