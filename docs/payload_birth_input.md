# Standard BirthInput Payload

Across the API documentation, you will see references to the **Standard `BirthInput`**. Almost all endpoints require this base structure to calculate the astrological state for a given person or event.

## JSON Structure

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

## Field Definitions

| Field | Type | Description | Example |
|---|---|---|---|
| `name` | string | Name of the person or event. Used purely for metadata. | `"Vamsi"` |
| `date_of_birth` | string | The local date of birth in strict `YYYY-MM-DD` format. | `"2005-11-23"` |
| `time_of_birth` | string | The local time of birth in strict `HH:MM:SS` (24-hour) format. | `"15:35:00"` |
| `place_of_birth`| string | Name of the location. Used purely for metadata. | `"Sattenapalle"` |
| `latitude` | float | Decimal latitude. Positive for North, Negative for South. | `16.3938` |
| `longitude` | float | Decimal longitude. Positive for East, Negative for West. | `80.1522` |
| `timezone` | float | Decimal offset from UTC. Example: IST (+5:30) is `5.5`. EST (-5:00) is `-5.0`. | `5.5` |
| `ayanamsa` | string | The sidereal zodiac correction factor. **Supported values (case-insensitive):**<br>- `"Lahiri"` or `"True ChitraPaksha"` (Default)<br>- `"Raman"`<br>- `"Krishnamurti"`<br>- `"Fagan-Bradley"`<br>- `"Yukteshwar"`<br>- `"J.N. Bhasin"`<br>- `"True Revati"`<br>- `"True Pushya"` | `"Lahiri"` |
| `house_system` | string | The algorithm used to calculate house cusps. **Supported values (case-insensitive):**<br>- `"Placidus"` (Default)<br>- `"Koch"`<br>- `"Porphyry"`<br>- `"Regiomontanus"`<br>- `"Campanus"`<br>- `"Equal"`<br>- `"Whole Sign"`<br>- `"Sripati"` | `"Placidus"` |
