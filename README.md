# Astrology Backend Engine

A high-precision, enterprise-grade Vedic Astrology (Jyotish) REST API backend written in Go. 

While many astrological APIs rely on fixed mathematical approximations or generic time-steps, this engine is powered by a Go-native implementation of the Swiss Ephemeris (`go-swisseph`). It prioritizes strict mathematical root-finding, arcsecond accuracy, and uncompromised adherence to classical astrological texts such as *Jataka Parijata*, *Brihat Parashara Hora Shastra*, and *Phaldeepika*.

## Comprehensive Feature Set

### 1. Astronomical Foundation & Ephemeris
* **Sub-Second Precision:** Utilizes the Swiss Ephemeris to calculate continuous planetary longitudes, retrograde velocities, and exact combustion statuses.
* **Ayanamsa Configuration:** Supports multiple Ayanamsas, including Lahiri (Chitra Paksha), Krishnamurti (KP), Raman, and Fagan/Bradley. The Ayanamsa is dynamically configurable per API request.
* **House System Mathematical Models:** Generates mathematical house cusps using standard systems (Placidus, Equal, Porphyry, Campanus).
* **Deep Planetary Metrics:** For every planet (Sun through Ketu, plus Ascendant and Midheaven), the engine returns:
  * Absolute 360-degree Sidereal and Tropical Longitudes.
  * Localized Sign placement, Degree-in-Sign, Minutes, and Seconds.
  * Orbital Velocity and Declination.
  * Exact Nakshatra (Constellation), Nakshatra Pada (Quarter), Nakshatra Lord, Sub-Lord, and Sub-Sub-Lord calculations.

### 2. Divisional Charts (Vargas) Engine
* **Complete 16-Varga Generation:** Dynamically calculates mathematical subdivisional charts including D1 (Rasi), D2 (Hora), D3 (Drekkana), D4 (Chaturthamsha), D7 (Saptamsha), D9 (Navamsha), D10 (Dashamsha), D12 (Dwadashamsha), D16 (Shodashamsha), D20 (Vimshamsha), D24 (Chaturvimshamsha), D27 (Saptavimshamsha), D30 (Trishamsha), D40 (Khavedamsha), D45 (Akshavedamsha), and D60 (Shashtiamsha).
* **Dynamic Varga Dignities:** A planet's dignity is evaluated locally within its Varga. For example, a planet may be debilitated in the D1 physical chart but evaluated mathematically as exalted in the D9 chart.
* **Payload Optimization:** Advanced planetary states (Pushkaramsa, Avasthas, etc.) are implemented using Go pointers and `omitempty` JSON tags. This ensures these compute-heavy, specific fields are only populated for the relevant charts (D1 and D9), automatically stripping `nil` pointers from higher Varga payloads to maintain a lightweight API response.

### 3. Advanced Astrological Dignities & States
* **Graha Avasthas (Planetary States):**
  * **Baaladi Avasthas (Age/Physical State):** Calculated based on odd/even sign placements and strict 6-degree increments (Bala, Kumara, Yuva, Vriddha, Mrita).
  * **Deeptadi Avasthas (Mood/Psychological State):** Calculated dynamically across charts based on temporary/natural relationships, exaltation, debilitation, own house, friends' houses, enemies' houses, and combustion (Deepta, Swastha, Mudita, Shanta, Dukhita, Deena, Kopa).
* **Pushkaramsa & Pushkara Bhaga:** Accurately calculates nourishing degrees. The engine evaluates exact elemental sign boundaries (Fire, Earth, Air, Water) against planetary degrees, incorporating a strict 1-degree tolerance orb for Pushkara Bhaga.
* **Mrityu Bhaga (Fatal Degrees):** Implements a comprehensive 108-degree fatal matrix based strictly on *Jataka Parijata*. Every primary planet's exact degree is checked against its specific fatal degree for the sign it occupies.
* **Shadbala (Six-Fold Strength):** Full mathematical evaluation of planetary strengths, encompassing Sthana Bala (Positional), Dik Bala (Directional), Kala Bala (Temporal), Cheshta Bala (Motional), Naisargika Bala (Natural), and Drik Bala (Aspectual).

### 4. Precision Panchang Engine
* **Mathematical Root-Finding (Bisection Algorithm):** Rather than approximating Tithi or Nakshatra durations using fixed time additions, the engine utilizes a continuous Universal Time (UT) bisection algorithm. It calculates exact boundary crossing timestamps down to the minute by finding the precise mathematical root where the longitudinal distance between the Sun and Moon reaches exact increments (e.g., multiples of 12 degrees for Tithi).
* **Daily Elements:** Calculates exact start and end boundaries for Tithi, Nakshatra, Yoga, and Karana.
* **Muhurta & Daily Timings:** Computes highly localized Sunrise and Sunset utilizing exact latitude/longitude coordinates, subsequently calculating dynamic daily periods like Rahu Kalam, Yamaganda, and Gulika Kalam.

### 5. Dosha & Compatibility Modules
* **Kuja Dosha (Manglik) Engine:** Determines the presence of Manglik Dosha using Mars' placement from the Ascendant, Moon, and Venus.
* **Classical Cancellations:** Implements deep exception logic, including:
  * Mars in own signs (Aries, Scorpio) or exalted (Capricorn).
  * Sign-specific house exceptions (e.g., Mars in the 8th house in Sagittarius/Pisces).
  * Nakshatra-based cancellations evaluating the Moon's exact constellation.
  * Conjunct mitigating planets (Jupiter, Moon).

## System Architecture

### Directory Structure
```text
astrology_backend_go/
├── cmd/
│   ├── server/           # Main application entry point
│   └── verify_panchang/  # CLI tool for isolated algorithm verification
├── internal/
│   ├── api/              # HTTP Handlers, Routers, and Request/Response definitions
│   ├── astrology/        # Core Vedic logic (Avasthas, Manglik, Panchang, Pushkara, Vargas)
│   ├── astronomy/        # Integration with go-swisseph (Planets, Houses, Time conversions)
│   └── domain/           # Shared Structs, Models, and JSON schema definitions
├── go.mod                # Go module dependencies
└── README.md             # Project documentation
```

## API Integration Examples

The backend is designed for seamless frontend consumption, abstracting complex mathematical structs into strictly typed JSON arrays.

### Example Request (Birth Chart / Natal)
`POST /api/v1/chart/natal`
```json
{
  "name": "User",
  "date_of_birth": "1990-05-15",
  "time_of_birth": "14:30:00",
  "place_of_birth": "New Delhi",
  "latitude": 28.6139,
  "longitude": 77.2090,
  "timezone": 5.5,
  "ayanamsa": "Krishnamurti",
  "house_system": "Placidus"
}
```

### Example Payload Structure (Vargas Endpoint)
`POST /api/v1/vargas` (Returns an array of dynamically calculated charts)
```json
{
  "vargas": [
    {
      "division": 9,
      "name": "D9",
      "ascendant": {
        "planet": "Ascendant",
        "divisional_sign": "Taurus"
      },
      "planets": [
        {
          "planet": "Jupiter",
          "divisional_sign": "Capricorn",
          "degree": 14,
          "minute": 20,
          "second": 15.5,
          "nakshatra": "Shravana",
          "avastha_mood": "Deena",
          "is_pushkaramsa": false
        }
      ]
    }
  ]
}
```
*(Note: Because of the optimization engine, fields like `avastha_mood` and `is_pushkaramsa` are automatically omitted in higher subdivisional charts like D10 or D60).*

## Getting Started

### Prerequisites
* Go 1.20 or higher.
* `go-swisseph` dependencies installed.

### Installation & Execution

1. Clone the repository and navigate into the directory:
   ```bash
   git clone <repository_url>
   cd astrology_backend_go
   ```
2. Install module dependencies:
   ```bash
   go mod tidy
   go mod download
   ```
3. Run the development server:
   ```bash
   go run cmd/server/main.go
   ```
4. The REST API will be accessible via standard HTTP (default port 8080 or based on environment variables).
