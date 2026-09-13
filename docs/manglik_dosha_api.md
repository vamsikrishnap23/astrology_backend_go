# Manglik Dosha API Documentation

## Endpoint
`POST /api/manglik-dosha`

## Description
Calculates Manglik Dosha (Kuja Dosha) based on the traditional principles of Vedic Astrology. This endpoint analyzes the position of Mars relative to the Ascendant, Moon, and Venus. It also rigorously applies classical cancellation rules (Kuja Dosha Bhanga) to determine if the dosha is nullified, including the specific Nakshatra exceptions.

## Request

**Content-Type:** `application/json`

### Body Parameters
The endpoint accepts a standard birth details payload.

| Field | Type | Description | Example |
|---|---|---|---|
| `name` | string | Name of the native | `"John Doe"` |
| `date_of_birth` | string | Format: YYYY-MM-DD | `"1995-10-05"` |
| `time_of_birth` | string | Format: HH:MM:SS (24-hour) | `"12:30:00"` |
| `place_of_birth` | string | City/Location name | `"New Delhi"` |
| `latitude` | float64 | Decimal latitude | `28.6139` |
| `longitude` | float64 | Decimal longitude | `77.2090` |
| `timezone` | float64 | Offset from UTC in hours | `5.5` |
| `ayanamsa` | string | Ayanamsa to use (e.g. Lahiri, Raman) | `"Lahiri"` |
| `house_system` | string | House calculation system | `"Placidus"` |

### Example Request
```json
{
  "name": "Test User",
  "date_of_birth": "1995-10-05",
  "time_of_birth": "12:00:00",
  "place_of_birth": "Guntur",
  "latitude": 16.3067,
  "longitude": 80.4365,
  "timezone": 5.5,
  "ayanamsa": "Lahiri",
  "house_system": "Placidus"
}
```

## Response

**Content-Type:** `application/json`

### Output Structure

| Field | Type | Description |
|---|---|---|
| `is_manglik` | boolean | Final boolean flag indicating if the dosha is active. |
| `status` | string | Classification of the dosha (e.g., `"High Manglik"`, `"Partial Manglik"`, `"Cancelled Manglik"`, `"Non-Manglik"`). |
| `base_analysis` | object | Contains the house and presence of Mars relative to `from_ascendant`, `from_moon`, and `from_venus`. |
| `cancellations` | array | List of specific classical rules that were triggered to nullify the dosha (if applicable). |
| `remedies` | array | List of suggested Vedic remedies (populated if `is_manglik` is true). |

### Example Response (Cancelled Manglik)
```json
{
  "is_manglik": false,
  "status": "Cancelled Manglik",
  "base_analysis": {
    "from_ascendant": {
      "is_present": false,
      "house": 3
    },
    "from_moon": {
      "is_present": false,
      "house": 9
    },
    "from_venus": {
      "is_present": true,
      "house": 4
    }
  },
  "cancellations": [
    {
      "rule": "Mars is in an exempt Nakshatra (Swati), completely nullifying the dosha."
    }
  ],
  "remedies": []
}
```

## Exceptions / Cancellations Implemented
This endpoint checks the following classical exceptions to cancel out the dosha:
1. **Benefic Ascendant:** Ascendant is Leo or Cancer.
2. **Dignity:** Mars is in its Own Sign (Aries, Scorpio), Exalted (Capricorn), or Debilitated (Cancer).
3. **Conjunction:** Mars is conjunct Jupiter or the Moon (Chandra-Mangala Yoga).
4. **House-Specific:**
    * 2nd House in Gemini or Virgo
    * 4th House in Aries or Scorpio
    * 7th House in Cancer or Capricorn
    * 8th House in Sagittarius or Pisces
    * 12th House in Taurus or Libra
5. **Nakshatra Bhanga:** Mars is positioned in one of the following 13 exempt Nakshatras: Ashwini, Mrigashira, Punarvasu, Pushya, Ashlesha, Uttara Phalguni, Swati, Anuradha, Purva Ashadha, Uttara Ashadha, Shravana, Uttara Bhadrapada, Revati.
