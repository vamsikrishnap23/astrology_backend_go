# Website endpoints (`/site/*`)

Public endpoints for the company website. Code: `internal/site/`. Registered in `cmd/server/main.go` via `site.Register`.

- **No login.** Unlike `/api/*`, these need no Supabase token.
- **CORS** is on. `SITE_ALLOWED_ORIGINS` (comma separated, e.g. `https://example.com,https://www.example.com`) limits which sites may call them from a browser; unset = any site.
- **Fixed settings**: Lahiri ayanamsa (both), Placidus houses (chart). Request values for these are ignored.
- Errors are JSON: `{"error": "..."}` with status 400 (bad input), 405 (wrong method) or 500.
- Dates must be within 1800-2399 (bundled ephemeris).

## `GET /site/values`

Every text value the two endpoints can return, with English and Telugu, plus field labels. Source file: `internal/site/values.json` (embedded in the binary). Cached for a day.

```js
const values = await (await fetch("/site/values")).json();
values.enums.tithis["Dashami"].te           // "దశమి"
values.enums.planets["Sun"].short_te        // "ర" (for chart boxes)
values.enum_paths.panchang["tithis[].name"] // "tithis"  (which list a field uses)
values.labels.rahu_kalam.te                 // "రాహుకాలం"
```

Lists: `weekdays` (7), `planets` (9, with `short_en`/`short_te`), `signs` (12), `nakshatras` (27), `tithis` (16 names; Shukla and Krishna share names, the paksha is a separate field), `pakshas` (2), `yogas` (27), `karanas` (11), `karana_types` (2), `samvatsaras` (60), `ayanas` (2), `rutus` (6), `masams` (24 = 12 + "Adhika ..."), `ayanamsas`, `house_systems`.

Numbers and times (degree, pada, start, end, progress) are not translated. `internal/site/site_test.go` fails if an endpoint returns a value missing from the file, so keep it in step with any rename.

## `GET|POST /site/panchang`

Daily panchang for one date and place.

GET: `/site/panchang?date=2026-10-21&lat=17.385&lon=78.4867&tz=5.5&place=Hyderabad` (`latitude`/`longitude`/`timezone` also accepted). Cached for an hour.

POST body:
```json
{"date": "2026-10-21", "place": "Hyderabad", "latitude": 17.385, "longitude": 78.4867, "timezone": 5.5}
```
`latitude`, `longitude`, `timezone` are required. `date` defaults to today in that timezone. `timezone` is a fixed offset in hours (no DST handling).

Response: `place, latitude, longitude, ayanamsa` plus the same fields as `/api/panchang/daily`:
`date, timezone, sunrise, sunset, solar_noon, moonrise, moonset, vara{number,name,ruler}, tithis[], nakshatras[], yogas[], karanas[], telugu_calendar{samvatsara,ayana,rutu,masam,paksha}, rahu_kalam[], yamaganda[], durmuhurtam[], varjyam[], amrutha_ghadiyalu[]`.

- `tithis`, `nakshatras`, `yogas`, `karanas`: every element that runs during the local day (midnight to midnight), with `start`/`end` (may lie on the previous/next day) and `progress` (% done at the start of the day).
- Times are ISO 8601 in the requested offset, e.g. `2026-10-21T06:10:30+05:30`.
- `telugu_calendar` is the one in force at local sunrise (differs from `/api/panchang/daily`, which uses the request time).
- `varjyam` / `amrutha_ghadiyalu` are listed for each nakshatra touching the day, so a period can fall outside the day.

## `POST /site/chart`

D1 (Rasi) chart.

```json
{"name": "Test", "date_of_birth": "1998-04-28", "time_of_birth": "14:30", "place_of_birth": "Hyderabad",
 "latitude": 17.385, "longitude": 78.4867, "timezone": 5.5}
```
`time_of_birth` is `HH:MM` or `HH:MM:SS`. `latitude`, `longitude`, `timezone` are required.

Response:

| Field | Meaning |
|---|---|
| `name, date_of_birth, time_of_birth, place_of_birth, latitude, longitude, timezone` | echoed input |
| `calculation_time_utc` | birth moment in UTC |
| `ayanamsa` (`"Lahiri"`), `ayanamsa_value` (degrees), `house_system` (`"Placidus"`) | settings |
| `ascendant` | a position (below) |
| `planets[]` | 9 grahas in order Sun, Moon, Mars, Mercury, Jupiter, Venus, Saturn, Rahu, Ketu: `planet` + position + `house` (Placidus bhava 1-12), `retrograde`, `combust` |
| `houses[]` | 12 Placidus cusps: `house` + position |
| `signs[]` | 12 boxes Aries..Pisces for drawing the Rasi chart: `sign, sign_number, is_ascendant, planets[]` (planet names) |

A **position** is: `longitude` (sidereal 0-360), `sign`, `sign_number` (1 = Aries), `sign_lord`, `degree`, `minute`, `second` (inside the sign, seconds truncated), `nakshatra`, `nakshatra_pada`, `nakshatra_lord`.

Notes:
- A planet's `sign` (for the Rasi chart) and its `house` (Placidus bhava, cusp to next cusp) can disagree; that is expected.
- Rahu/Ketu are the mean node; they normally show `retrograde: true`.
- Uranus, Neptune and Pluto are left out.
- At polar latitudes Swiss Ephemeris replaces Placidus with Porphyry without saying so.
