# Astrology Go Backend

## Architecture
The project follows a standard Go layout:
- `cmd/server/`: The main entrypoint for the HTTP server.
- `internal/api/`: HTTP handlers, request, and response logic.
- `internal/astronomy/`: Core astronomical calculations relying on Swiss Ephemeris.
- `internal/domain/`: Core structs representing inputs, configurations, and outputs.

## Swiss Ephemeris Integration
The project uses `github.com/tejzpr/go-swisseph` to interface with the Swiss Ephemeris C library. 

### Ephemeris Files Configuration (Strictly Required)
For astrological calculations, the system **strictly requires** the high-precision Swiss Ephemeris data files (`.se1`). It will **not** silently fall back to the Moshier analytical model.

You must configure the `EPHE_PATH` environment variable to point to a directory containing at least the following standard ephemeris files (which cover standard planets and the moon for modern dates):
- `sepl_18.se1`
- `semo_18.se1`

If these files are missing, the server will intentionally fail to start or return an internal error, ensuring absolute reproducible accuracy.

```bash
export EPHE_PATH=/path/to/ephe_data
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
