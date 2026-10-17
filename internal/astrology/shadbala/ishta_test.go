package shadbala

import "testing"

// Raman art. 136, example 60: the Sun's Sayana longitude 202 10' plus 90 deg is 292 10', reduced to 67 50',
// divided by 3. He prints 22.66; 67 50' / 3 is 22.61.
func TestSunChestaBala(t *testing.T) {
	near(t, "Sun", sunChestaBala(dms(202, 10, 0)), 22.61, 0.006)
	near(t, "Chesta Kendra 180", sunChestaBala(90), 60, 1e-9) // 90 + 90 = 180 deg: full
	near(t, "Chesta Kendra 0", sunChestaBala(270), 0, 1e-9)   // 270 + 90 = 360
	near(t, "reduced", sunChestaBala(10), 100.0/3.0, 1e-9)    // 100 deg
}

// Art. 137, example 61: the Moon's Chesta Bala is its arc from the Sun, reduced, divided by 3.
func TestMoonChestaBala(t *testing.T) {
	near(t, "Moon", moonChestaBala(dms(180, 54, 0), dms(311, 17, 0)), 43.46, 0.006)
	near(t, "waning side", moonChestaBala(dms(311, 17, 0), dms(180, 54, 0)), 43.46, 0.006)
	near(t, "full moon", moonChestaBala(0, 180), 60, 1e-9)
	near(t, "new moon", moonChestaBala(100, 100), 0, 1e-9)
}

// Examples 62 and 63 (p. 103): Ishta and Kashta Phala from the Uchcha and Chesta Balas of the standard
// horoscope. The Sun, Moon, Mars and Venus (Kashta) are reproduced as printed. Mars's Chesta Bala is printed
// 22.03 in the Ishta line and 22.23 elsewhere; 22.23 gives his 28.70. The other printed values do not follow
// his own formula, so the formula's values are used: Mercury's Kashta (printed 49.16, formula 17.81),
// Jupiter's Kashta (13.19, 9.53), Venus's Ishta (3.49, 3.41), Saturn's Ishta (27.00, 27.12) and Kashta (31.50, 31.29).
func TestIshtaKashtaRamanExamples(t *testing.T) {
	cases := []struct {
		planet                        string
		uchcha, chesta, ishta, kashta float64
	}{
		{"Sun", 3.0, 22.66, 8.25, 46.13}, {"Moon", 32.75, 43.46, 37.73, 21.23}, {"Mars", 37.06, 22.23, 28.70, 29.44},
		{"Mercury", 54.50, 2.30, 11.20, 17.81}, {"Jupiter", 56.33, 35.26, 44.57, 9.53}, {"Venus", 1.95, 5.95, 3.41, 56.00},
		{"Saturn", 34.80, 21.14, 27.12, 31.29},
	}
	for _, c := range cases {
		i, k := ishtaKashta(c.uchcha, c.chesta)
		near(t, c.planet+" Ishta", i, c.ishta, 0.06)
		near(t, c.planet+" Kashta", k, c.kashta, 0.06)
	}
	// the limits: nothing in either gives no Ishta and the full Kashta, full strength in both the reverse
	i, k := ishtaKashta(0, 0)
	near(t, "0,0 ishta", i, 0, 1e-9)
	near(t, "0,0 kashta", k, 60, 1e-9)
	i, k = ishtaKashta(60, 60)
	near(t, "60,60 ishta", i, 60, 1e-9)
	near(t, "60,60 kashta", k, 0, 1e-9)
}

// The engine on his standard horoscope: his remarks of art. 140 that Venus's Kashta predominates over her
// Ishta (3.49 against 56.00) and that Jupiter and the Moon are mainly beneficial, the Sun, Mercury and
// Saturn mainly harmful. (Mars is 28.70 against 29.44, too close to fix.)
func TestIshtaKashtaRamanChart(t *testing.T) {
	ctx := ramanChartContext(t, "raman")
	res, err := CalculateShadbala(&ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.IshtaKashtaPhala) != 7 {
		t.Fatalf("%d planets", len(res.IshtaKashtaPhala))
	}
	want := map[string]string{"Sun": "kashta", "Moon": "ishta", "Mercury": "kashta", "Jupiter": "ishta", "Venus": "kashta", "Saturn": "kashta"}
	ishta := map[string]float64{"Sun": 8.25, "Moon": 37.73, "Mercury": 11.20, "Jupiter": 44.57, "Venus": 3.49, "Saturn": 27.00}
	for _, p := range res.IshtaKashtaPhala {
		if w, ok := want[p.Planet]; ok {
			if p.Inclination != w {
				t.Errorf("%s: inclination %s, want %s", p.Planet, p.Inclination, w)
			}
			near(t, p.Planet+" Ishta", p.IshtaPhala, ishta[p.Planet], 1.0)
		}
		if p.IshtaPhala < 0 || p.IshtaPhala > 60 || p.KashtaPhala < 0 || p.KashtaPhala > 60 {
			t.Errorf("%s out of range: %+v", p.Planet, p)
		}
		if (p.Inclination == "ishta") != (p.IshtaPhala > p.KashtaPhala) {
			t.Errorf("%s: inclination %s with ishta %.2f kashta %.2f", p.Planet, p.Inclination, p.IshtaPhala, p.KashtaPhala)
		}
	}
	// the table in the response repeats the list
	tb := res.Tables.IshtaKashtaPhala
	if len(tb.Rows) != 4 || len(tb.Total) != 0 {
		t.Fatalf("table: %d rows, total %v", len(tb.Rows), tb.Total)
	}
	for _, p := range res.IshtaKashtaPhala {
		near(t, p.Planet+" table Ishta", tb.Rows[2].Values[p.Planet], p.IshtaPhala, 0)
		near(t, p.Planet+" table Kashta", tb.Rows[3].Values[p.Planet], p.KashtaPhala, 0)
	}
}
