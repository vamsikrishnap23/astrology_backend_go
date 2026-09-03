# Astrology Backend Go - Frontend API Reference

**Base URL:** `http://localhost:8080`

All endpoints accept and return `application/json`. 
Virtually all endpoints require a standard `BirthInput` payload to define the astrological chart being queried.

---

## Standard Base Payload (`BirthInput`)

Whenever an endpoint requires standard birth data, use this exact structure:

```json
{
  "name": "User Name",
  "date_of_birth": "2005-11-23",     // Format: YYYY-MM-DD
  "time_of_birth": "15:35:00",       // Format: HH:MM:SS
  "place_of_birth": "Sattenapalle",  // String representation
  "latitude": 16.3938,               // Decimal degrees
  "longitude": 80.1522,              // Decimal degrees
  "timezone": 5.5,                   // Decimal offset from UTC
  "ayanamsa": "Lahiri",              // e.g. "Lahiri", "Raman", "KP"
  "house_system": "Placidus"         // e.g. "Placidus", "WholeSign"
}
```

---

## 1. Natal Chart
**Endpoint:** `POST /api/chart`

**Payload:** `BirthInput`

**Response:**
Returns the complete calculation time, ayanamsa used, ascendant degree, MC, and the arrays of planetary positions and house cusps.

---

## 2. Panchang (Vedic Calendar)
**Endpoint:** `POST /api/panchang`

**Payload:** `BirthInput`

**Response:**
Returns sunrise/set, moonrise/set, Tithi, Vara, Nakshatra, Yoga, Karana, and auspicious/inauspicious daily periods.

---

## 3. Mathematical Tables
**Endpoint:** `POST /api/tables`

**Payload:** `BirthInput`

**Response:**
Returns exhaustive tables for `planetary_table` and `house_table` including retrograde status, exact longitude, Nakshatra padas, and the full hierarchy of Sign Lord -> Nakshatra Lord -> Sub Lord -> Sub Sub Lord -> Sub Sub Sub Lord.

---

## 4. KP Significators
**Endpoint:** `POST /api/significators`

**Payload:** `BirthInput`

**Response:**
Returns `planet_view` and `house_view` resolving the KP A, B, C, D significators.

---

## 5. Four-Step Theory (Sunil Gondhalekar)
**Endpoint:** `POST /api/four-step`

**Payload:** `BirthInput`

**Response:**
Returns `four_step_view` mapping the Planet -> Star Lord -> Sub Lord -> Star Lord of Sub Lord.

---

## 6. KP Ruling Planets
**Endpoint:** `POST /api/ruling-planets`

**Payload:** `BirthInput`

**Response:**
Returns the resolved `ruling_planets` (Ascendant/Moon lords and Day Lord) including node agent assignments (Rahu/Ketu acting as agents).

---

## 7. Vimshottari Dasha
**Endpoint:** `POST /api/dasha`

**Payload:** `BirthInput`

**Response:**
Returns `balance_years` and a deeply nested `mahadasha` array down to the 4th level (`antardasha` -> `pratyantardasha` -> `sookshma`).

---

## 8. Divisional Charts (Vargas)
**Endpoint:** `POST /api/vargas`

**Payload:** `BirthInput`

**Response:**
Returns `vargas` array mapping D1 through D60, including the Ascendant and planetary placements within those harmonic divisions.

---

## 9. Secondary Progressions
**Endpoint:** `POST /api/progression`

**Payload:** `BirthInput` + `"progression_year": 2040`

**Response:**
Returns the exact progressed chart variables (`progressed_planets` and `progressed_houses`) utilizing the day-for-a-year solar system.

---

## 10. Transits (Snapshot Chart)
**Endpoint:** `POST /api/transits/chart`

**Payload:** `BirthInput` + `"transit_date": "2024-05-15"` + `"transit_time": "12:00:00"`

**Response:**
Returns a `transit_data` block identical to the `TablesResult` to allow for easy UI overlay/comparison.

---

## 11. Transits (Upcoming Events)
**Endpoint:** `POST /api/transits/upcoming`

**Payload:** `BirthInput` + `"search_start_date": "2024-05-15"` + `"search_start_time": "12:00:00"`

**Response:**
Returns `transits` array showing the next exact UTC timestamps when planets cross into new signs.

---

## 12. Bhava Chalit
**Endpoint:** `POST /api/bhava-chalit`

**Payload:** `BirthInput`

**Response:**
Returns `houses` array containing the exact mathematical cusps and the `occupants` shifted into those realigned boundaries.

---

## 13. Ashtakavarga
**Endpoint:** `POST /api/ashtakavarga`

**Payload:** `BirthInput`

**Response:**
Returns the Bhinnashtakavarga (`bav` array for the 7 planets) and the composite Sarvashtakavarga (`sav` array for the 12 signs).

---

## 14. Shadbala
**Endpoint:** `POST /api/shadbala`

**Payload:** `BirthInput`

**Response:**
Returns `planets` array breaking down the 6-fold strength (Sthana, Dig, Kala, Cheshta, Naisargika, Drik) resolving to the ultimate strength ratio and `rupas`.

---

## 15. Jaimini Chara Karakas
**Endpoint:** `POST /api/jaimini-karakas`

**Payload:** `BirthInput`

**Response:**
Returns `karakas` array sorting the 7 planets into AK, AmK, BK, MK, PK, GK, and DK based on fractional sign longitude.

---

## 16. Ashtakoota (Guna Milan Compatibility)
**Endpoint:** `POST /api/ashtakoota`

**Payload:**
Requires dual birth inputs:
```json
{
  "groom": { /* BirthInput */ },
  "bride": { /* BirthInput */ }
}
```

**Response:**
Returns `kootas` dictionary resolving Varna, Vashya, Tara, Yoni, Graha Maitri, Gana, Bhakoot, and Nadi with exact scores, dosha evaluations, and a final `summary`.

---

## 17. Birth Time Rectification (BTR)
**Endpoint:** `POST /api/btr`

**Payload:** `BirthInput` + `"scan_minus_minutes": 10` + `"scan_plus_minutes": 5`

**Response:**
Returns `astronomical_context` and a `candidates` array mapping 90-min rulers and Tatwas. Note: The final match evaluation is marked `unresolved` pending a proprietary classical table injection.

