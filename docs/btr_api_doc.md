# Birth Time Rectification (BTR) API Documentation

## Endpoint
`POST /api/btr`

## Description
This endpoint performs Birth Time Rectification using the classical Nadi Tatwa (Elemental) theory. It mathematically normalizes the birth time (adjusting for LMT and local sunrise), maps it to a 3-minute Nadi Row, and verifies if the calculated Tatwa Gender and Nadi Planet align with the user's actual Gender and Nakshatra Lord. 

If the inputted time fails validation, a time-scanning engine iterates through a ±2 hour window to find the closest mathematically verified times and returns them as ranked candidates.

## Request

**Content-Type:** `application/json`

### Body Parameters

| Field | Type | Description | Example |
|---|---|---|---|
| `name` | string | Name of the native | `"John Doe"` |
| `date_of_birth` | string | Format: YYYY-MM-DD | `"2026-09-13"` |
| `time_of_birth` | string | Format: HH:MM:SS (24-hour) | `"01:44:00"` |
| `gender` | string | Biological gender (`"Male"` or `"Female"`) | `"Male"` |
| `place_of_birth` | string | City/Location name | `"Hyderabad"` |
| `latitude` | float64 | Decimal latitude | `17.38405` |
| `longitude` | float64 | Decimal longitude | `78.45636` |
| `timezone` | float64 | Offset from UTC in hours | `5.5` |
| `ayanamsa` | string | Ayanamsa to use (defaults to Lahiri) | `"Lahiri"` |
| `return_full_table` | boolean | Set to true to receive the entire 480-row day matrix | `true` |
| `scan_minus_minutes` | integer | Number of minutes to scan backward if time fails (default 120) | `10` |
| `scan_plus_minutes` | integer | Number of minutes to scan forward if time fails (default 120) | `5` |
| `sign_type_override` | string | Override auto-detected sign type ("Movable", "Fixed", "Dual") | `"Dual"` |
| `sunrise_override` | string | Override calculated sunrise time (e.g. "06:07:54") | `"06:07:54"` |
| `ascendant_override` | string | Override calculated Ascendant sign to derive modality | `"Mithuna"` |
| `star_lord_override` | string | Override calculated Nakshatra lord (e.g. "Moon", "Mars") | `"Mars"` |

### Example Request Payload
```json
{
  "name": "Test User",
  "date_of_birth": "2026-09-13",
  "time_of_birth": "01:44:00",
  "gender": "Male",
  "place_of_birth": "Hyderabad, Telangana, India",
  "latitude": 17.38405,
  "longitude": 78.45636,
  "timezone": 5.5,
  "ayanamsa": "Lahiri",
  "return_full_table": true,
  "scan_minus_minutes": 10,
  "scan_plus_minutes": 5,
  "sign_type_override": "",
  "sunrise_override": "06:07:54",
  "ascendant_override": "Mithuna",
  "star_lord_override": "Moon"
}
```

---

## Response

**Content-Type:** `application/json`

### Output Structure

| Field | Type | Description |
|---|---|---|
| `input_time_status` | string | Will be `"Verified"` if the exact input time matches, otherwise `"Failed"`. |
| `input_analysis` | object | Contains the deep analysis comparing the calculated Nadi values against the actual astrological values for the exact inputted time. |
| `suggested_rectifications` | array | If the input fails, this array lists the closest mathematically valid birth times (ranked by proximity). |
| `full_table` | array | If requested, contains the full list of all 480 Nadi rows and their calculated attributes for the entire day. |

### `input_analysis` Object Fields
| Field | Type | Description |
|---|---|---|
| `gender_match` | boolean | `true` if the calculated Tatwa gender matches the inputted gender. |
| `star_match` | boolean | `true` if the calculated Nadi planet matches the actual Nakshatra Lord. |
| `calculated_tatwa` | string | The active element for that minute (Prithvi, Jala, Tejo, Vayu, Akash). |
| `calculated_gender` | string | The gender governing that Tatwa (Male/Female). |
| `calculated_planet` | string | The ruling planet for that 3-minute Nadi row. |
| `actual_star_lord` | string | The actual ruling planet of the native's natal Moon Nakshatra. |

### `suggested_rectifications` Array Fields
| Field | Type | Description |
|---|---|---|
| `rank` | integer | Ranked from 1 (closest) to 5. |
| `suggested_time` | string | The mathematically verified alternative clock time (HH:MM:SS). |
| `difference_minutes` | integer | How many minutes the suggested time is offset from the originally inputted time. |
| `tatwa` | string | The governing Tatwa of the suggested time. |
| `nadi_row` | integer | The discrete Nadi Row (1-480) for the suggested time. |


### `full_table` Array Fields (Screenshot Parity)
| Field | Type | Description |
|---|---|---|
| `no` | integer | The discrete Nadi Row Number (1-480). |
| `t1` | integer | Elapsed minutes in the cycle (Row * 3). |
| `t2` | string | Exact IST/Standard Clock Time the row ends (HH:MM). |
| `wed` | string | Formatted Tatwa string (e.g. "MalePrithvi") for Wednesday. |
| `mon_fri` | string | Formatted Tatwa string for Monday/Friday. |
| `sun_tues` | string | Formatted Tatwa string for Sunday/Tuesday. |
| `sat` | string | Formatted Tatwa string for Saturday. |
| `thur` | string | Formatted Tatwa string for Thursday. |
| `movable` | string | 90-Min Planet for Movable Ascendant. |
| `fixed` | string | 90-Min Planet for Fixed Ascendant. |
| `dual` | string | 90-Min Planet for Dual Ascendant. |
| `vinod_movable` | string | Vinod Planet for Movable Ascendant. |
| `vinod_fixed` | string | Vinod Planet for Fixed Ascendant. |
| `vinod_dual` | string | Vinod Planet for Dual Ascendant. |

### Example Response (Failed Input -> Candidate Generation)
```json
{
  "input_time_status": "Failed",
  "input_analysis": {
    "gender_match": true,
    "star_match": false,
    "calculated_tatwa": "Akash",
    "calculated_gender": "Male",
    "calculated_planet": "Mercury",
    "actual_star_lord": "Moon"
  },
  "suggested_rectifications": [
    {
      "rank": 1,
      "suggested_time": "01:56:00",
      "difference_minutes": 12,
      "tatwa": "Akash",
      "nadi_row": 43
    },
    {
      "rank": 2,
      "suggested_time": "01:02:00",
      "difference_minutes": -42,
      "tatwa": "Prithvi",
      "nadi_row": 25
    },
    {
      "rank": 3,
      "suggested_time": "02:50:00",
      "difference_minutes": 66,
      "tatwa": "Tejo",
      "nadi_row": 61
    }
  ]
}
```
