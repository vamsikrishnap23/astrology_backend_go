# Astrology Go Backend

## Architecture
The project follows a standard Go layout:
- `cmd/server/`: The main entrypoint for the HTTP server.
- `internal/api/`: HTTP handlers, request, and response logic.
- `internal/astronomy/`: Core astronomical calculations relying on Swiss Ephemeris.
- `internal/domain/`: Core structs representing inputs, configurations, and outputs.

## Swiss Ephemeris Integration
The project uses `github.com/mshafiee/swephgo` to interface with the Swiss Ephemeris C library. 

### Ephemeris Files Configuration
By default, Swiss Ephemeris calculates positions using internal analytical models if ephemeris files (`.se1`) are missing. For high precision, you must provide `.se1` files.
- You can set the `EPHE_PATH` environment variable to the directory containing your `.se1` files before running the server.

```bash
export EPHE_PATH=/path/to/ephe
```

## Running the Server
```bash
go run cmd/server/main.go
```

## Running Tests
```bash
go test ./...
```

## API Example
### Request
```bash
curl -X POST http://localhost:8080/api/v1/chart \
-H "Content-Type: application/json" \
-d '{
  "name": "Example",
  "date_of_birth": "1998-04-28",
  "time_of_birth": "14:30:00",
  "place_of_birth": "Hyderabad",
  "latitude": 17.3850,
  "longitude": 78.4867,
  "timezone": 5.5,
  "ayanamsa": "Lahiri",
  "house_system": "Placidus"
}'
```

### Known Limitations
- Node calculations use Mean Node by default. True Node toggle isn't exposed in the API yet.
- Precision depends on the availability of `.se1` ephemeris files. Without them, it falls back to Moshier's analytical model (which is still very accurate but not the maximum precision provided by JPL data).
