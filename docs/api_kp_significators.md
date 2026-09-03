# KP Significators & Four-Step API Endpoints

## 1. Standard KP Significators
**Endpoint:** `POST /api/significators`

### Description
Calculates Krishnamurti Paddhati (KP) planetary and house significators (A, B, C, D levels) for a birth chart.

- **A**: House occupied by the star lord of the planet
- **B**: House occupied by the planet itself
- **C**: Houses owned by the star lord of the planet
- **D**: Houses owned by the planet itself

### Request Payload
Standard `BirthInput` (see `/api/panchang`).

### Response Structure
```json
{
  "planet_view": [
    {
      "planet": "Sun",
      "a": [8],
      "b": [8],
      "c": [10, 11],
      "d": [5]
    }
  ],
  "house_view": [
    {
      "house": 8,
      "a": ["Sun", "Moon"],
      "b": ["Mars", "Venus"],
      "c": ["Jupiter"],
      "d": ["Saturn"]
    }
  ]
}
```

---

## 2. Four-Step Significators
**Endpoint:** `POST /api/four-step`

### Description
Implements Sunil Gondhalekar's Four-Step Theory of KP Astrology, evaluating the Planet, Star Lord, Sub Lord, and Star Lord of the Sub Lord.

### Request Payload
Standard `BirthInput`.

### Response Structure
```json
{
  "four_step_view": [
    {
      "planet": "Sun",
      "planet_details": {
        "planet": "Sun",
        "houses": [8, 5]
      },
      "star_lord": {
        "planet": "Saturn",
        "houses": [4, 10, 11]
      },
      "sub_lord": {
        "planet": "Mercury",
        "houses": [8, 3, 6]
      },
      "star_lord_of_sub": {
        "planet": "Jupiter",
        "houses": [7, 9, 12]
      }
    }
  ]
}
```
