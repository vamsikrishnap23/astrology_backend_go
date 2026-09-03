# KP Ruling Planets API Endpoint

**Endpoint:** `POST /api/ruling-planets`

## Description
Extracts the KP Ruling Planets for a given chart (Ascendant Star/Sign/Sub Lords, Moon Star/Sign/Sub Lords, Day Lord) and dynamically resolves Node Agents (Rahu/Ketu acting on behalf of other planets).

## Request Payload
Standard `BirthInput`.

## Response Structure
```json
{
  "ruling_planets": [
    {
      "planet": "Ketu",
      "source": "Ascendant Star Lord"
    },
    {
      "planet": "Mars",
      "source": "Ascendant Sign Lord"
    },
    {
      "planet": "Ketu",
      "source": "Moon Star Lord"
    },
    {
      "planet": "Sun",
      "source": "Moon Sign Lord"
    },
    {
      "planet": "Mercury",
      "source": "Day Lord"
    },
    {
      "planet": "Rahu",
      "source": "Node Agent",
      "agent_for": "Mercury"
    }
  ]
}
```
