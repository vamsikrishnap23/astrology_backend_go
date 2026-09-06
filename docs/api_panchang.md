# Panchang API Endpoint

**Endpoint:** `POST /api/panchang`

## Description
Calculates the 5 limbs of the Vedic calendar (Panchang): Tithi, Vara, Nakshatra, Yoga, and Karana, along with daily sun/moon timings and inauspicious periods (Rahu Kalam, Yamaganda, Durmuhurtam, Varjyam).

## Request Payload

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

| Field | Type | Description |
|---|---|---|
| `name` | string | Name of the person/event |
| `date_of_birth` | string | Date in `YYYY-MM-DD` format |
| `time_of_birth` | string | Time in `HH:MM:SS` format |
| `place_of_birth` | string | Location string |
| `latitude` | float | Decimal latitude |
| `longitude` | float | Decimal longitude |
| `timezone` | float | Decimal offset from UTC |
| `ayanamsa` | string | Ayanamsa type (e.g., `Lahiri`) |
| `house_system` | string | House system code (e.g., `Placidus`) |

## Response Structure

```json
{
  "date": "2005-11-23",
  "local_time": "15:35:00",
  "timezone": 5.5,
  "sunrise": "2005-11-23T06:17:00+05:30",
  "sunset": "2005-11-23T17:33:00+05:30",
  "solar_noon": "2005-11-23T11:55:00+05:30",
  "moonrise": "2005-11-23T23:45:00+05:30",
  "moonset": "2005-11-24T12:15:00+05:30",
  "rasi": "Leo",
  "vara": {
    "number": 4,
    "name": "Wednesday",
    "ruler": "Mercury"
  },
  "tithi": {
    "number": 8,
    "name": "Ashtami",
    "paksha": "Krishna",
    "progress": 45.2,
    "start": "2005-11-22T21:40:00Z",
    "end": "2005-11-23T20:15:00Z"
  },
  "nakshatra": {
    "number": 10,
    "name": "Magha",
    "pada": 1,
    "progress": 15.6,
    "start": "2005-11-23T05:10:00Z",
    "end": "2005-11-24T06:20:00Z",
    "ruler": "Ketu"
  },
  "yoga": {
    "number": 5,
    "name": "Sobhana",
    "progress": 30.1,
    "start": "2005-11-22T19:30:00Z",
    "end": "2005-11-23T18:45:00Z"
  },
  "karana": {
    "number": 15,
    "name": "Balava",
    "type": "Movable",
    "progress": 90.4,
    "start": "2005-11-23T09:00:00Z",
    "end": "2005-11-23T20:15:00Z"
  },
  "rahu_kalam": {
    "start": "2005-11-23T06:40:00Z",
    "end": "2005-11-23T08:10:00Z"
  },
  "yamaganda": {
    "start": "2005-11-23T03:40:00Z",
    "end": "2005-11-23T05:10:00Z"
  },
  "durmuhurtam": [
    {
      "start": "2005-11-23T05:36:00Z",
      "end": "2005-11-23T06:24:00Z"
    }
  ],
  "varjyam": [
    {
      "start": "2005-11-23T15:10:00Z",
      "end": "2005-11-23T16:46:00Z"
    }
  ]
}
```
