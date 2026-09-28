# Frontend Integration Guide: Advanced Astrological Features

This document outlines the four new advanced astrological properties added to the Go backend. All these properties have been seamlessly injected into the existing `PlanetPosition` JSON objects returned by the API.

## 1. Updated JSON Payload
You do not need to change any API endpoints. Every planet object in the natal chart, transits, or progressions will now include these 5 new fields:

```json
{
  "planet": "Jupiter",
  "sign": "Sagittarius",
  "degree_in_sign": 21.5,
  "is_pushkaramsa": true,
  "is_pushkara_bhaga": true,
  "is_mrityu_bhaga": false,
  "avastha_age": "Vriddha",
  "avastha_mood": "Swastha"
}
```

---

## 2. Feature Implementation & UI Display Strategies

### A. Pushkaramsa & Pushkara Bhaga (Auspicious Degrees)
*   **Keys:** `is_pushkaramsa` (boolean), `is_pushkara_bhaga` (boolean)
*   **Meaning:** These are highly auspicious and nourishing placements. A Pushkara Bhaga is a highly concentrated "golden point" that always falls inside a Pushkaramsa window.
*   **UI Suggestion:** 
    *   Inside the D1 chart tooltip, if `is_pushkara_bhaga` is true, display a brilliant glowing star or text like **"Pushkara Bhaga (Highly Auspicious)"**.
    *   If only `is_pushkaramsa` is true, display a standard star or text like **"Pushkaramsa"**.

### B. Mrithyu Bhaga (Fatal Degrees)
*   **Key:** `is_mrityu_bhaga` (boolean)
*   **Meaning:** The planet is sitting at its exact mathematical fatal/death-inflicting degree for that specific sign. It becomes severely afflicted.
*   **UI Suggestion:** 
    *   Display a warning icon (⚠️) or red text in the tooltip: **"Mrithyu Bhaga (Fatal Degree)"**.

### C. Graha Avasthas (Planetary States)
*   **Keys:** `avastha_age` (string), `avastha_mood` (string)
*   **Meaning:** `avastha_age` calculates the physical maturity (Baaladi) based on Odd/Even signs and 6-degree chunks. `avastha_mood` calculates the psychological dignity (Deeptadi) based on friendships, exaltation, and combustion.
*   **UI Suggestion:** 
    *   Add a line in the tooltip: `State: {Translated_Age} / {Translated_Mood}`
    *   *Example:* `స్థితి: యౌవన / ముదిత`

---

## 3. Telugu Translation Dictionary

Map the exact English strings returned by the API to your frontend i18n Telugu dictionary using the following tables.

### Avastha Age (Baaladi Avastha)
These values represent the physical maturity of the planet.

| API String (Key) | Telugu Translation | English Meaning |
| :--- | :--- | :--- |
| `Bala` | బాల్య | Infant / Child |
| `Kumara` | కౌమార | Youth / Teen |
| `Yuva` | యౌవన | Adult / Prime |
| `Vriddha` | వృద్ధ | Old Age |
| `Mrita` | మృత | Dead / Inactive |

### Avastha Mood (Deeptadi Avastha)
These values represent the psychological dignity and comfort of the planet.

| API String (Key) | Telugu Translation | English Meaning (Astrological Reason) |
| :--- | :--- | :--- |
| `Deepta` | దీప్త | Radiant (Exalted / ఉచ్చ) |
| `Swastha` | స్వస్థ | Confident (Own Sign / స్వక్షేత్రం) |
| `Mudita` | ముదిత | Happy (Friend's Sign / మిత్ర రాశి) |
| `Shanta` | శాంత | Peaceful (Neutral Sign / సమ రాశి) |
| `Dukhita` | దుఃఖిత | Sad / Distressed (Enemy's Sign / శత్రు రాశి) |
| `Deena` | దీన | Depressed (Debilitated / నీచ) |
| `Kopa` | కోప | Furious (Combust by Sun / అస్తంగతం) |

### Feature Flags
For the boolean flags, if you are displaying raw text instead of icons, you can use these translations:

| JSON Key Trigger | Telugu Translation |
| :--- | :--- |
| `is_pushkaramsa: true` | పుష్కర నవాంశ |
| `is_pushkara_bhaga: true` | పుష్కర భాగ |
| `is_mrityu_bhaga: true` | మృత్యు భాగ |
