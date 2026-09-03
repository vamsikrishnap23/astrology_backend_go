# Astrology Go Backend — Requirements & Architecture Specification

## 1. Project Goal

Build a fresh, production-quality astrology calculation engine and REST API in Go for a website.

The system should prioritize:

1. calculation accuracy
2. clear, testable domain logic
3. modular architecture
4. structured API responses
5. reproducibility
6. maintainability
7. good performance

The frontend should consume calculated data from the API. It should not implement astrology calculations itself.

The backend should be designed as a calculation engine first and an HTTP API second, so the core calculation packages can be tested independently of HTTP.

---

# 2. Core Design Principles

## 2.1 Accuracy over convenience

Astrological calculations must not be approximated merely to make implementation easier.

Use Swiss Ephemeris for astronomical calculations wherever appropriate.

For derived astrology calculations, implement the relevant documented/classical formulas explicitly.

Do not use unexplained magic numbers.

When multiple conventions exist, document the selected convention.

## 2.2 Structured data over presentation strings

The calculation layer must return structured domain objects.

Prefer:

```json
{
  "planet": "Moon",
  "sign": "Sagittarius",
  "longitude": 243.123456,
  "degree_in_sign": 3.123456,
  "nakshatra": {
    "name": "Mula",
    "number": 19,
    "pada": 1
  }
}
```

over:

```json
{
  "planet": "Mo",
  "position": "Sagittarius 3:07:24, Mula-1"
}
```

The API may include display-friendly fields in addition to structured fields, but calculation logic must never depend on formatted strings.

## 2.3 One source of astronomical truth

All modules should consume common astronomical primitives.

For example:

```text
Input
  ↓
Time normalization
  ↓
Julian Day / ephemeris context
  ↓
Swiss Ephemeris
  ↓
Tropical planetary positions
  ↓
Sidereal conversion
  ↓
Houses / angles
  ↓
Common astrology primitives
  ↓
KP / Panchang / Vargas / Dasha / etc.
```

Do not independently recalculate planetary positions inside every module.

## 2.4 Test every layer

Calculation modules must have unit tests.

Where practical, tests should expose intermediate calculations so errors can be localized.

The final values will be manually compared against trusted external astrology software/websites.

External comparison is a validation tool, not a replacement for automated tests.

---

# 3. Technology

## Backend

- Go
- REST API
- JSON
- idiomatic Go
- standard library wherever practical
- a lightweight HTTP router only if it provides meaningful value

## Astronomy

Use Swiss Ephemeris through a maintained Go binding/wrapper where appropriate.

Before implementation, verify:

- binding maturity
- supported Swiss Ephemeris version
- supported planets/bodies
- sidereal modes
- house systems
- rise/set calculations
- ephemeris file requirements
- CGO requirements
- build requirements
- licensing

Do not assume that a Go wrapper exposes every Swiss Ephemeris feature.

Swiss Ephemeris documentation is the authority for its astronomical API and house-system/sidereal behavior.

---

# 4. Input Model

The main birth-chart request must support:

```text
name
date_of_birth
time_of_birth
place_of_birth
latitude
longitude
timezone
ayanamsa
house_system
```

## Timezone

The requested timezone representation is a decimal UTC offset.

Examples:

```text
India = 5.5
Nepal = 5.75
UTC = 0
New York EST = -5
```

The backend must convert the local birth datetime plus offset into the correct UT/UTC representation required by Swiss Ephemeris.

Do not silently infer the timezone from the place name.

Do not use the decimal offset as geographic longitude.

Validate:

- latitude: -90 to +90
- longitude: -180 to +180
- timezone offset: reasonable UTC-offset range
- valid date
- valid time

Keep the original local input and normalized calculation timestamp available in the calculation context.

---

# 5. Configuration

## 5.1 Ayanamsa

The UI requirements include options such as:

- KP New
- Lahiri
- Raman
- Krishnamurti
- KP Straight Line
- KP Kullar
- True Revati
- True Pushya
- Fagan-Bradley
- De Luce
- Yukteshwar
- J.N. Bhasin
- True Chitrapaksha

Before exposing these in the API, verify the exact Swiss Ephemeris sidereal modes and map each user-facing name to the correct Swiss Ephemeris constant/mode.

Important:

- Do not invent a mapping.
- Do not assume similar names are identical.
- Distinguish true vs mean/reference-specific ayanamsas where Swiss Ephemeris does.
- Preserve the selected ayanamsa in the calculation context and API response.

Swiss Ephemeris supports many more sidereal modes than the UI requires. Only expose the project-approved set.

## 5.2 House systems

The UI requirements include:

- Placidus
- Koch
- Equal
- Whole Sign
- Sripati
- Porphyry
- Regiomontanus
- Campanus
- Topocentric
- Alcabitius
- Morinus
- Meridian / Axial Rotation
- Azimuthal / Horizontal
- Polich-Page
- Krusinski

Verify exact Swiss Ephemeris support and mappings before implementing.

Important distinction:

Swiss Ephemeris documents geometrically correct house-position methods separately from methods where house positions may be simplified. Preserve that distinction and document any limitation.

Relevant Swiss Ephemeris house codes include, among others:

```text
P = Placidus
K = Koch
O = Porphyry
R = Regiomontanus
C = Campanus
A/E = Equal
W = Whole Sign
B = Alcabitius
X = Meridian / Axial Rotation
H = Horizon / Azimuthal
M = Morinus
T = Polich/Page
U = Krusinski
S = Sripati
```

Verify all mappings against the actual Swiss Ephemeris version/binding used by the project.

---

# 6. Proposed Architecture

```text
cmd/
  server/
    main.go

internal/
  api/
    handlers/
    middleware/
    request/
    response/
    routes/

  config/

  astronomy/
    ephemeris/
    time/
    coordinates/
    planets/
    houses/

  astrology/
    core/
    zodiac/
    nakshatra/
    panchang/
    kp/
    dasha/
    divisional/
    progression/
    transit/
    bhava_chalit/
    ashtakavarga/
    shadbala/

  chart/
    natal/
    rendering_model/

  validation/

  testutil/

pkg/
  (only for intentionally public/reusable packages)

tests/
  integration/
  fixtures/

docs/
  architecture/
  calculations/
  api/

```

The exact package layout can be refined during Phase 0.

---

# 7. Calculation Context

Create a common immutable/read-only calculation context conceptually containing:

```text
Birth/Input
  - name
  - local date
  - local time
  - timezone offset
  - latitude
  - longitude
  - place

Configuration
  - ayanamsa
  - house system
  - node convention if required

Derived astronomical context
  - Julian Day UT
  - Julian Day TT where required
  - ayanamsa value
  - planetary positions
  - planetary speeds
  - houses
  - angles
```

Modules should consume this context rather than independently reconstructing it.

Cache expensive common calculations within a single request.

---

# 8. Common Planet Model

Support at least:

```text
Sun
Moon
Mars
Mercury
Jupiter
Venus
Saturn
Rahu
Ketu
Uranus
Neptune
Pluto
Ascendant
```

For planets/points where meaningful, provide:

```text
longitude
latitude
distance
speed
retrograde
sign
degree_in_sign
degree
minute
second
house
nakshatra
pada
sign_lord
star_lord
sub_lord
sub_sub_lord / deeper KP hierarchy where required
```

Do not assume retrograde for Rahu/Ketu merely from conventional presentation; derive node behavior from the selected node calculation/convention.

---

# 9. Natal / Rashi Chart

The primary natal chart must provide:

- Ascendant
- MC where useful
- all required planets/points
- signs
- exact positions
- houses
- retrograde state
- nakshatra
- pada
- relevant lord hierarchy
- chart tooltip data

The chart response must be structured so the website can render the chart without doing calculations.

---

# 10. Panchang

Return slightly more information than the reference UI while keeping the module focused.

Required:

- Vara
- Tithi
- Tithi number
- Paksha
- Nakshatra
- Nakshatra number
- Nakshatra Pada
- Yoga
- Yoga number
- Karana
- Rashi
- Sunrise
- Sunset
- Moonrise
- Moonset

For Tithi, Nakshatra, Yoga and Karana include transition/end timestamps where applicable.

Example:

```json
{
  "tithi": {
    "name": "Saptami",
    "number": 7,
    "paksha": "Shukla",
    "start": "...",
    "end": "..."
  }
}
```

The Panchang module must clearly document:

- sunrise definition
- local time conversion
- node convention if relevant
- lunar month convention if later added
- any tradition-specific assumptions

Do not add large numbers of extra Panchang features unless explicitly required.

---

# 11. Upcoming Transits

Provide the next sign ingress for relevant planets.

Required output:

```text
planet
destination_sign
transition_datetime
```

Include:

- Sun
- Moon
- Mars
- Mercury
- Jupiter
- Venus
- Saturn
- Rahu
- Ketu

Include outer planets if required by the frontend.

For each planet, search forward from the requested reference datetime and determine the next sign boundary crossing using actual computed longitude, not a fixed lookup table.

Handle retrograde motion correctly.

For Rahu/Ketu, explicitly document the chosen node convention.

---

# 12. Planetary Table

Return a table-ready structured model for:

- Ascendant
- Sun
- Moon
- Mars
- Mercury
- Jupiter
- Venus
- Saturn
- Rahu
- Ketu
- Uranus
- Neptune
- Pluto

Fields:

```text
planet
sign
exact_position
house
nakshatra
pada
sign_lord
star_lord
sub_lord
ss_lord
sss_lord
retrograde
```

The exact KP hierarchy must be calculated, not copied from display strings.

---

# 13. House Table

For houses 1–12 provide:

```text
house_number
cusp_longitude
sign
degree_in_sign
nakshatra
pada
sign_lord
star_lord
sub_lord
ss_lord
sss_lord
```

The house system selected by the user must be respected.

---

# 14. KP Significators

Implement two views.

## 14.1 Planet View

For each relevant planet:

```text
A = Star Lord's Bhava
B = Planet's Bhava
C = Star Lord's Houses
D = Planet's Houses
```

Return structured arrays.

Example:

```json
{
  "planet": "Sun",
  "a": [],
  "b": [10],
  "c": [7, 11, 12],
  "d": [10]
}
```

The exact interpretation of "Bhava" vs "Houses" must be defined in the calculation documentation.

## 14.2 House View

For each house return the corresponding planet lists for A/B/C/D.

Do not return comma-separated display strings as the primary data.

---

# 15. KP Ruling Planets

Calculate:

- Ascendant Star Lord
- Ascendant Sign Lord
- Ascendant Sub Lord
- Moon Star Lord
- Moon Sign Lord
- Moon Sub Lord
- Day Lord
- Nodes as Agents

The exact nodes-as-agents convention must be explicitly documented and tested.

Return structured planet references.

---

# 16. Vimshottari Dasha

Implement four levels:

```text
Mahadasha
  └── Antardasha
      └── Pratyantardasha
          └── Sookshma
```

The initial Mahadasha balance must be derived from the Moon's exact birth Nakshatra position.

Return nested structures:

```json
{
  "planet": "Jupiter",
  "start": "...",
  "end": "...",
  "antardasha": [
    {
      "planet": "Jupiter",
      "start": "...",
      "end": "...",
      "pratyantardasha": [
        {
          "planet": "Jupiter",
          "start": "...",
          "end": "...",
          "sookshma": []
        }
      ]
    }
  ]
}
```

The calculation must account for the exact fractional balance remaining at birth.

Document:

- Vimshottari year convention
- duration calculations
- ordering
- boundary behavior
- timezone handling

---

# 17. KP 4-Step Signification

Implement the four-step signification view shown in the UI.

For each relevant planet expose:

1. Planet
2. Star Lord of Planet
3. Sub Lord of Planet
4. Star Lord of Sub Lord

Each step must include the relevant significated houses.

Example:

```json
{
  "planet": "Sun",
  "star_lord": {
    "planet": "Venus",
    "houses": [7, 11, 12]
  },
  "sub_lord": {
    "planet": "Moon",
    "houses": [6]
  },
  "star_lord_of_sub": {
    "planet": "Sun",
    "houses": [10]
  }
}
```

Keep this calculation separate from the general significator tables.

---

# 18. Divisional Charts / Vargas

Implement all required divisional charts as a reusable Varga engine.

At minimum the architecture must support adding:

```text
D1
D2
D3
D4
D7
D9
D10
D12
D16
D20
D24
D27
D30
D40
D45
D60
```

The final supported list should be confirmed before implementation if the website requires a different set.

Do not write completely independent duplicated algorithms for every Varga.

Create a common mechanism for:

```text
planet longitude
    ↓
Varga rule
    ↓
divisional sign
```

Each Varga may have its own classical mapping rules and must implement those rules explicitly.

## Tooltip requirements

Every divisional chart must expose enough information for tooltips.

Tooltip data should include:

- planet
- exact source longitude
- divisional sign
- relevant degree/position
- nakshatra
- pada
- lord information
- retrograde status
- any chart-specific data required by the UI

The frontend must not recalculate these values.

---

# 19. Progression Chart

Provide a progression chart with a year picker.

The API must accept a selected progression year (and optionally a month/day, or default to exactly N years from birth).

Methodology chosen: **Secondary Progression**
- Rule: 1 day after birth equals 1 tropical year of life.
- The target progressed date = Birth Date + (Age in days).
- Calculate the planetary positions and houses for this precise progressed date and time.
- Standard Sidereal Ayanamsa and House systems must still apply to this progressed date.

The response should include:

- progressed date used
- age (in years)
- planets
- Ascendant
- houses where applicable
- signs
- exact positions
- retrograde state
- tooltip data

The API should make the selected year explicit.

---

# 20. Transit Chart

Provide an arbitrary transit chart using:

```text
date
time
```

The birth chart's location/configuration is the default reference unless a separate transit location is later required.

Return:

- planetary positions
- signs
- houses
- Ascendant
- retrograde state
- nakshatra
- pada
- tooltip metadata

This is distinct from Upcoming Transits.

Upcoming Transits answers:

> When does a planet next enter another sign?

Transit Chart answers:

> Where are the planets at this exact date/time?

---

# 21. Bhava Chalit

Implement a separate Bhava Chalit chart.

Return:

- house cusps
- house numbers
- planetary house placement
- planetary sign placement
- exact longitude
- Ascendant
- tooltip data

Keep zodiac/sign placement and house/bhava placement conceptually separate.

Respect the selected house-system configuration.

Document the exact Bhava Chalit convention used.

---

# 22. Ashtakavarga

Implement:

## Bhinnashtakavarga

For relevant planets, return:

- sign-wise bindu counts
- total
- contributing planet/source data where appropriate

## Sarvashtakavarga

Return:

- sign-wise total bindus
- 12-sign table
- totals/check sums

Use explicit classical rules.

Do not hardcode expected totals without calculating the contributing bindus.

Create tests for contribution logic and aggregate totals.

---

# 23. Shadbala

This is a high-priority accuracy module.

Calculate Shadbala for:

- Sun
- Moon
- Mars
- Mercury
- Jupiter
- Venus
- Saturn

Six major components:

1. Sthana Bala
2. Dig Bala
3. Kala Bala
4. Cheshta Bala
5. Naisargika Bala
6. Drik Bala

Return:

```text
planet
sthana_bala
dig_bala
kala_bala
cheshta_bala
naisargika_bala
drik_bala
total
rupas
required_minimum
strength_ratio
```

Where the classical calculation itself contains subcomponents, expose those too.

For example, Sthana Bala should not be represented as a single opaque number if it is composed of independently calculable parts.

The implementation must document:

- exact formulas
- units
- conversion to Rupas
- minimum strength values
- exaltation calculations
- divisional strength
- positional components
- directional strength
- temporal strength
- motional strength
- natural strength
- aspectual strength
- retrograde handling
- day/night handling
- sign/house dependencies

Do NOT implement Shadbala as an approximation.

Every subcomponent should have independent tests.

This module must be validated against trusted external calculations using several known charts.

---

# 24. Time and Precision

Use floating-point values only where appropriate.

Preserve high precision internally.

Do not repeatedly round intermediate values.

Round only at presentation boundaries.

Store/return enough precision for manual comparison.

Use UTC/UT correctly for ephemeris calls and local time for user-facing timestamps.

Be explicit about:

```text
local civil time
UTC
UT
TT
Julian Day UT
Julian Day TT
```

where relevant.

---

# 25. Error Handling

Use explicit typed/domain errors where useful.

API errors should be consistent.

Examples:

```json
{
  "error": {
    "code": "INVALID_LATITUDE",
    "message": "latitude must be between -90 and 90"
  }
}
```

Do not expose raw CGO/Swiss Ephemeris errors directly to clients.

Log detailed internal errors while returning safe API messages.

---

# 26. API Design

Prefer logically grouped endpoints rather than one huge endpoint.

Possible structure:

```text
POST /api/v1/chart
POST /api/v1/panchang
POST /api/v1/transits/upcoming
POST /api/v1/transits/chart
POST /api/v1/progression
POST /api/v1/divisional/{varga}
POST /api/v1/bhava-chalit
POST /api/v1/kp/significators
POST /api/v1/kp/ruling-planets
POST /api/v1/kp/four-step
POST /api/v1/dasha/vimshottari
POST /api/v1/ashtakavarga
POST /api/v1/shadbala
```

Exact endpoint design may be refined after the domain model is finalized.

All responses should be JSON.

Version the public API.

---

# 27. DTO vs Domain Models

Keep API DTOs separate from calculation/domain structures where doing so improves maintainability.

Conceptually:

```text
HTTP JSON
   ↓
Request DTO
   ↓
Validation
   ↓
Domain Input
   ↓
Calculation Engine
   ↓
Domain Result
   ↓
Response DTO
   ↓
JSON
```

Do not let JSON formatting decisions contaminate calculation algorithms.

---

# 28. Caching

Within a single chart request, cache common expensive calculations:

- Julian dates
- ayanamsa
- planetary positions
- houses
- nakshatras
- KP lord hierarchy

Do not introduce a distributed cache until profiling demonstrates a need.

Correctness comes before caching.

---

# 29. Observability

Add structured logging around:

- request ID
- calculation duration
- selected configuration
- ephemeris errors
- module failures

Never log sensitive user information unnecessarily.

Provide enough instrumentation to benchmark calculation modules.

---

# 30. Testing Strategy

## Unit tests

Each package gets tests.

Examples:

```text
astronomy/time
astronomy/planets
astronomy/houses

astrology/nakshatra
astrology/panchang
astrology/kp
astrology/dasha
astrology/divisional
astrology/ashtakavarga
astrology/shadbala
```

## Golden/reference cases

Create a set of fixed birth inputs with expected intermediate values.

For example:

```text
input
  date
  time
  timezone
  latitude
  longitude
  ayanamsa
  house system

expected
  Julian Day
  ayanamsa
  Sun longitude
  Moon longitude
  Ascendant
  house cusps
  nakshatra
  etc.
```

## Cross-validation

The developer will manually compare results against a trusted astrology website.

The project should make it easy to compare:

- exact planetary longitude
- sign
- degree
- nakshatra
- house
- KP lords
- dasha boundaries
- Varga placements
- Panchang transitions
- Shadbala components

---

# 31. Development Phases

Do not implement everything in one pass.

## Phase 0 — Architecture and dependency verification

Deliver:

- project structure
- Go module
- Swiss Ephemeris binding decision
- dependency/licensing notes
- ephemeris file strategy
- ayanamsa mapping table
- house-system mapping table
- core domain model
- API plan
- test strategy

No large astrology implementation yet.

## Phase 1 — Astronomical foundation

Implement:

- time normalization
- Julian Day
- Swiss Ephemeris initialization
- planetary positions
- node calculation
- ayanamsa
- signs
- degree/minute/second conversion
- retrograde/speed
- houses
- Ascendant
- MC

Build tests before moving on.

## Phase 2 — Core astrology primitives

Implement:

- zodiac
- nakshatras
- padas
- sign lords
- star lords
- KP sub divisions
- reusable planet/house models

## Phase 3 — Natal chart + Panchang

Implement:

- natal chart
- planetary table
- house table
- Panchang
- sunrise
- sunset
- moonrise
- moonset

## Phase 4 — KP

Implement:

- significators planet view
- significators house view
- ruling planets
- 4-step signification

## Phase 5 — Dasha

Implement:

- Vimshottari Mahadasha
- Antardasha
- Pratyantardasha
- Sookshma
- exact birth balance

## Phase 6 — Vargas and chart variants

Implement:

- Varga engine
- all required divisional charts
- tooltip metadata
- Bhava Chalit
- progression
- transit chart

## Phase 7 — Predictive/calculation modules

Implement:

- upcoming transits
- Ashtakavarga
- Shadbala

Shadbala should receive additional validation before being considered complete.

## Phase 8 — API hardening

Implement:

- request validation
- error model
- OpenAPI documentation
- integration tests
- benchmarks
- logging
- performance profiling
- Docker/build documentation

---

# 32. Definition of Done

A module is complete only when:

- implementation is modular
- formulas/conventions are documented
- unit tests exist
- edge cases have been considered
- intermediate calculations are inspectable
- API output is structured
- no calculation is duplicated unnecessarily
- manual comparison has been performed
- discrepancies are investigated rather than hidden by rounding

A matching result is not considered proof of correctness if the underlying formula is unknown.

---

# 33. Important Constraints

Do not:

- blindly translate an existing implementation
- copy calculation constants without documenting them
- put calculations inside HTTP handlers
- use formatted strings as internal data
- silently choose between conflicting traditions
- round intermediate calculations
- invent unsupported Swiss Ephemeris features
- approximate Shadbala
- hardcode transit dates
- make the frontend calculate astrology values
- implement all modules in one giant commit

Do:

- build incrementally
- verify dependencies first
- document formulas
- write tests before/alongside calculations
- keep astronomy and astrology layers separate
- use reusable primitives
- expose intermediate data for validation
- keep the API stable and structured

---

# 34. First Deliverable

Start with Phase 0 only.

Before writing the main calculation engine, produce:

1. proposed repository structure
2. dependency choices
3. Swiss Ephemeris Go binding evaluation
4. ephemeris data/file strategy
5. exact ayanamsa mapping table
6. exact house-system mapping table
7. domain model proposal
8. API endpoint proposal
9. calculation dependency graph
10. testing strategy
11. known ambiguities
12. Phase 1 implementation plan

Then wait for approval before implementing Phase 1.
