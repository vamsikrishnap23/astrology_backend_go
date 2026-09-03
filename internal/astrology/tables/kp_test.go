package tables

import (
	"testing"
)

func TestGetKPLords(t *testing.T) {
	tests := []struct {
		name      string
		longitude float64
		wantSL    string
		wantNL    string
		wantSSL   string
		wantSSSL  string
		wantSSSSL string
	}{
		{
			name:      "0 degrees (Aries start)",
			longitude: 0.0,
			wantSL:    "Mars",
			wantNL:    "Ketu",
			wantSSL:   "Ketu",
			wantSSSL:  "Ketu",
			wantSSSSL: "Ketu",
		},
		{
			name:      "13 deg 20 min (Taurus boundary)",
			longitude: 13.33333333333333,
			wantSL:    "Mars", // technically Aries 13:20
			wantNL:    "Venus",
			wantSSL:   "Venus",
			wantSSSL:  "Venus",
			wantSSSSL: "Venus",
		},
		{
			name:      "Reference House VI Cusp (148°46'48\")",
			longitude: 148.780000,
			wantSL:    "Sun",
			wantNL:    "Sun",
			wantSSL:   "Mars",
			wantSSSL:  "Saturn",
			wantSSSSL: "Sun",
		},
		{
			name:      "Reference House VII Cusp (185°20'43\")",
			longitude: 185.34527777777778, // 185 + 20/60 + 43/3600
			wantSL:    "Venus",
			wantNL:    "Mars",
			wantSSL:   "Sun",
			wantSSSL:  "Mercury",
			wantSSSSL: "Venus",
		},
		{
			name:      "Reference Mercury (220°10'12.11\")",
			longitude: 220.17003055555555, // 220 + 10/60 + 12.11/3600
			wantSL:    "Mars",
			wantNL:    "Saturn",
			wantSSL:   "Venus",
			wantSSSL:  "Mercury",
			wantSSSSL: "Saturn",
		},
		{
			name:      "Reference Venus (262°34'24\")",
			longitude: 262.5733333333333, // 262 + 34/60 + 24/3600
			wantSL:    "Jupiter",
			wantNL:    "Venus",
			wantSSL:   "Saturn",
			wantSSSL:  "Ketu",
			wantSSSSL: "Mars",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl, nl, ssl, sssl, ssssl := GetKPLords(tt.longitude)
			if sl != tt.wantSL || nl != tt.wantNL || ssl != tt.wantSSL || sssl != tt.wantSSSL || ssssl != tt.wantSSSSL {
				t.Errorf("GetKPLords(%f) = %s, %s, %s, %s, %s; want %s, %s, %s, %s, %s",
					tt.longitude, sl, nl, ssl, sssl, ssssl, tt.wantSL, tt.wantNL, tt.wantSSL, tt.wantSSSL, tt.wantSSSSL)
			}
		})
	}
}
