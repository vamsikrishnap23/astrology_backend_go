# Vimshottari Dasha API Endpoint

**Endpoint:** `POST /api/dasha`

## Description
Generates a complete 120-year Vimshottari Dasha cycle calculated precisely down to the 4th level (Sookshma) based on the exact sidereal longitudinal progression of the Moon through its birth Nakshatra.

## Request Payload
Standard `BirthInput`.

## Response Structure
```json
{
  "balance_years": 4.56,
  "mahadasha": [
    {
      "lord": "Ketu",
      "start_date": "2005-11-23T15:35:00Z",
      "end_date": "2010-06-15T00:00:00Z",
      "antardasha": [
        {
          "lord": "Venus",
          "start_date": "2005-11-23T15:35:00Z",
          "end_date": "2006-08-15T00:00:00Z",
          "pratyantardasha": [
            {
              "lord": "Sun",
              "start_date": "2005-11-23T15:35:00Z",
              "end_date": "2005-12-15T00:00:00Z",
              "sookshma": [
                {
                  "lord": "Moon",
                  "start_date": "2005-11-23T15:35:00Z",
                  "end_date": "2005-11-25T00:00:00Z"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```
