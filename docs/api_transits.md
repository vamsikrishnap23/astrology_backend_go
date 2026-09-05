# Transits API Endpoints

## 1. Transit Chart
**Endpoint:** `POST /api/transits/chart`

### Description
Calculates the snapshot of current/historical planetary positions for a specified target date/time, localized to the original birth location constraints. The transit house cusps are dynamically built as a **Whole Sign** system starting from the **Natal Moon's Sign (Chandra Lagna)**, which is the most standard Vedic astrological method for reading Gochara (transits). Reuses the exhaustive `TablesResult` structure to provide identical formatting for direct overlay comparisons in UIs.

### Request Payload
Standard `BirthInput`, plus the `transit_date` and `transit_time`.

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
  "transit_date": "2024-05-15",
  "transit_time": "12:00:00"
}
```

### Response Structure
```json
{
  "transit_date_utc": "2024-05-15T06:30:00Z",
  "julian_day": 2460445.7708,
  "ayanamsa": 24.19,
  "ascendant": 115.3,
  "transit_data": {
    "planetary_table": [ ... ],
    "house_table": [ ... ]
  }
}
```

---

## 2. Upcoming Transits
**Endpoint:** `POST /api/transits/upcoming`

### Description
Calculates the next upcoming sign ingress/transition events for the 9 classical planets (Sun - Ketu). Determines the exact mathematical UTC timestamp where the sidereal planetary longitude crosses a 30-degree boundary into a new sign.

### Request Payload
Standard `BirthInput`, plus an optional search start window.

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
  "search_start_date": "2024-05-15",
  "search_start_time": "12:00:00"
}
```

### Response Structure
```json
{
  "search_start_date_utc": "2024-05-15T06:30:00Z",
  "transits": [
    {
      "planet": "Sun",
      "destination_sign": "Taurus",
      "transition_datetime": "2024-05-14T11:45:00Z"
    },
    {
      "planet": "Jupiter",
      "destination_sign": "Taurus",
      "transition_datetime": "2024-05-01T15:20:00Z"
    }
  ]
}
```
