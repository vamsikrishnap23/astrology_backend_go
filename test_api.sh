#!/bin/bash
curl -X POST http://localhost:8080/api/v1/progression \
-H "Content-Type: application/json" \
-d '{
  "date_of_birth": "1975-03-11",
  "time_of_birth": "03:15:00",
  "latitude": 16.42593,
  "longitude": 79.54041,
  "timezone": 5.5,
  "ayanamsa": "KP",
  "house_system": "Placidus",
  "progression_date": "2026-08-01"
}'
