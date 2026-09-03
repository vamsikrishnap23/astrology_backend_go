# Ashtakoota (Guna Milan) API Endpoint

**Endpoint:** `POST /api/ashtakoota`

## Description
Performs the traditional 8-Koota (36-guna) chart compatibility match between a Groom and a Bride based strictly on precise lunar sidereal mechanics. It includes dosha cancellations for Bhakoot and Nadi.

## Request Payload
Requires two standard `BirthInput` objects grouped under `groom` and `bride`.

```json
{
  "groom": {
    "name": "Groom",
    "date_of_birth": "2005-11-23",
    "time_of_birth": "15:35:00",
    "place_of_birth": "Sattenapalle",
    "latitude": 16.3938,
    "longitude": 80.1522,
    "timezone": 5.5,
    "ayanamsa": "Lahiri"
  },
  "bride": {
    "name": "Bride",
    "date_of_birth": "2002-05-14",
    "time_of_birth": "10:00:00",
    "place_of_birth": "Bangalore",
    "latitude": 12.9716,
    "longitude": 77.5946,
    "timezone": 5.5,
    "ayanamsa": "Lahiri"
  }
}
```

## Response Structure
```json
{
  "rule_set": {
    "name": "classical-guna-milan-v1",
    "ayanamsa": "Lahiri"
  },
  "groom": {
    "name": "Groom",
    "moon": {
      "longitude": 121.81,
      "sign": "Leo",
      "degree": 1.81,
      "nakshatra": "Magha",
      "pada": 1,
      "nakshatra_lord": "Ketu",
      "rashi_lord": "Sun"
    }
  },
  "bride": {
    "name": "Bride",
    "moon": { ... }
  },
  "kootas": {
    "varna": {
      "score": 1,
      "maximum": 1,
      "groom_value": "Kshatriya",
      "bride_value": "Vaishya",
      "explanation": "Groom varna Kshatriya (Rank 3) vs Bride varna Vaishya (Rank 2)"
    },
    "vashya": {
      "score": 0,
      "maximum": 2,
      "groom_value": "Vanchar",
      "bride_value": "Chatuspada",
      "explanation": "Score mapped from standard vashya table: Vanchar vs Chatuspada = 0.0"
    },
    "tara": { ... },
    "yoni": { ... },
    "graha_maitri": { ... },
    "gana": { ... },
    "bhakoot": {
      "raw_score": 0,
      "effective_score": 7,
      "maximum": 7,
      "groom_rashi": "Leo",
      "bride_rashi": "Taurus",
      "groom_to_bride_distance": 10,
      "bride_to_groom_distance": 4,
      "relationship": "4/10",
      "dosha": false,
      "cancellation_applied": true,
      "cancellation_reason": "Graha Maitri is high",
      "explanation": "..."
    },
    "nadi": { ... }
  },
  "summary": {
    "raw_total": 10.5,
    "effective_total": 17.5,
    "maximum": 36,
    "percentage": 48.61,
    "traditional_threshold_met": true
  }
}
```
