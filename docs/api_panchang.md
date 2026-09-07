# Panchang API

**Endpoint:** `POST /api/panchang`

## Description
Calculates the five principal elements (limbs) of the traditional Vedic Calendar (*Panchangam*): **Tithi** (Lunar Day), **Vara** (Weekday), **Nakshatra** (Constellation), **Yoga** (Sun/Moon combination), and **Karana** (Half-Tithi). 

It dynamically computes exact localized timestamps for astrological boundaries, rather than relying on static averages, using the precise Swiss Ephemeris bisection-search algorithms. It also computes auspicious/inauspicious windows tied to the specific location's sunrise and sunset.

---

## Request Payload

The endpoint accepts the standard `BirthInput` structure.

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
| `date_of_birth` | string | Target date in `YYYY-MM-DD` format. |
| `time_of_birth` | string | Target time in `HH:MM:SS` format. Used as the anchor point for the Panchang. |
| `latitude` | float | Decimal latitude of the location. |
| `longitude` | float | Decimal longitude of the location. |
| `timezone` | float | Decimal timezone offset from UTC (e.g., `5.5` for IST). |
| `ayanamsa` | string | Ayanamsa mode (e.g., `Lahiri`, `Raman`, `KP`). Applies strictly to Sidereal calculations (Nakshatra, Yoga). |

---

## Response Structure

All `start` and `end` timestamps are strictly returned in **ISO-8601 format** localized to the requested `timezone` (e.g., `+05:30`). If a period crosses over midnight, the ISO timestamp will automatically correctly reflect the next calendar day. 

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
    "number": 3,
    "name": "Wednesday",
    "ruler": "Mercury"
  },
  "tithi": {
    "number": 8,
    "name": "Ashtami",
    "paksha": "Krishna",
    "progress": 45.2,
    "start": "2005-11-22T21:40:00+05:30",
    "end": "2005-11-23T20:15:00+05:30"
  },
  "nakshatra": {
    "number": 10,
    "name": "Magha",
    "pada": 1,
    "progress": 15.6,
    "start": "2005-11-23T05:10:00+05:30",
    "end": "2005-11-24T06:20:00+05:30",
    "ruler": "Ketu"
  },
  "yoga": {
    "number": 5,
    "name": "Shobhana",
    "progress": 30.1,
    "start": "2005-11-22T19:30:00+05:30",
    "end": "2005-11-23T18:45:00+05:30"
  },
  "karana": {
    "number": 15,
    "name": "Balava",
    "type": "Moving",
    "progress": 90.4,
    "start": "2005-11-23T09:00:00+05:30",
    "end": "2005-11-23T20:15:00+05:30"
  },
  "rahu_kalam": {
    "start": "2005-11-23T12:00:00+05:30",
    "end": "2005-11-23T13:30:00+05:30"
  },
  "yamaganda": {
    "start": "2005-11-23T07:30:00+05:30",
    "end": "2005-11-23T09:00:00+05:30"
  },
  "durmuhurtam": [
    {
      "start": "2005-11-23T11:45:00+05:30",
      "end": "2005-11-23T12:33:00+05:30"
    }
  ],
  "varjyam": [
    {
      "start": "2005-11-23T20:25:00+05:30",
      "end": "2005-11-23T22:01:00+05:30"
    }
  ],
  "amrutha_ghadiyalu": [
    {
      "start": "2005-11-23T05:15:00+05:30",
      "end": "2005-11-23T06:51:00+05:30"
    }
  ]
}
```

### Important Attributes
- **Rasi:** The astrological Moon Sign representing where the Moon is currently transiting.
- **Progress:** For Tithi, Nakshatra, Yoga, and Karana, the `progress` field indicates the exact percentage (`0.0` to `100.0`) of elapsed time within that element's total duration for the specific requested `time_of_birth`.
- **Amrutha Ghadiyalu (Amrita Kalam):** Indicates highly auspicious localized timing windows.
- **Varjyam (Tyajyam / Visha Ghatis):** Indicates highly inauspicious localized timing windows.
