package time

import (
	"math"
	"testing"
	"time"
)

func TestParseLocalToUTC(t *testing.T) {
	tests := []struct {
		name        string
		date        string
		time        string
		tzOffset    float64
		expectedUTC time.Time
	}{
		{
			name:        "UTC 0",
			date:        "1998-04-28",
			time:        "14:30:00",
			tzOffset:    0,
			expectedUTC: time.Date(1998, 4, 28, 14, 30, 0, 0, time.UTC),
		},
		{
			name:        "India +5.5",
			date:        "1998-04-28",
			time:        "14:30:00",
			tzOffset:    5.5,
			expectedUTC: time.Date(1998, 4, 28, 9, 0, 0, 0, time.UTC),
		},
		{
			name:        "Nepal +5.75",
			date:        "1998-04-28",
			time:        "14:30:00",
			tzOffset:    5.75, // 5 hours 45 minutes
			expectedUTC: time.Date(1998, 4, 28, 8, 45, 0, 0, time.UTC),
		},
		{
			name:        "Negative -5.0",
			date:        "1998-04-28",
			time:        "14:30:00",
			tzOffset:    -5.0,
			expectedUTC: time.Date(1998, 4, 28, 19, 30, 0, 0, time.UTC),
		},
		{
			name:        "Cross Midnight",
			date:        "1998-04-28",
			time:        "02:30:00",
			tzOffset:    5.5,
			expectedUTC: time.Date(1998, 4, 27, 21, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utc, err := ParseLocalToUTC(tt.date, tt.time, tt.tzOffset)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !utc.Equal(tt.expectedUTC) {
				t.Errorf("expected %v, got %v", tt.expectedUTC, utc)
			}
		})
	}
}

func TestDecimalToDMS(t *testing.T) {
	sign, deg, min, sec := DecimalToDMS(0.0)
	if sign != "Aries" || deg != 0 || min != 0 || sec != 0.0 {
		t.Errorf("expected Aries 0 0 0.0, got %s %d %d %f", sign, deg, min, sec)
	}

	sign, deg, min, sec = DecimalToDMS(359.5)
	if sign != "Pisces" || deg != 29 || min != 30 || sec != 0.0 {
		t.Errorf("expected Pisces 29 30 0.0, got %s %d %d %f", sign, deg, min, sec)
	}
}

func TestUTCToJulianDay(t *testing.T) {
	utc := time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC)
	jd := UTCToJulianDay(utc)
	if math.Abs(jd-2451545.0) > 0.00001 {
		t.Errorf("expected ~2451545.0, got %f", jd)
	}
}
