package shadbala

import (
	"math"
	"sort"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// RuleSet names the rules the module follows: Raman's Graha and Bhava Balas, with BPHS chapter 27 where
// Raman is silent. See docs/shadbala_spec.md.
const RuleSet = "raman-shadbala-v1"

func round2(v float64) float64 { return math.Round(v*100) / 100 }

// buildTables lays the results out as the tables of Raman's worked example: one table per bala with a row
// per part and a value per planet, and the final Shadbala Pinda table with rupas, requirement and rank.
func buildTables(ps []domain.PlanetShadbala) domain.ShadbalaTables {
	names := make([]string, len(ps))
	for i, p := range ps {
		names[i] = p.Planet
	}
	row := func(component, label string, get func(domain.PlanetShadbala) float64) domain.ShadbalaTableRow {
		r := domain.ShadbalaTableRow{Component: component, Label: label, Values: map[string]float64{}}
		for _, p := range ps {
			r.Values[p.Planet] = round2(get(p))
		}
		return r
	}
	totals := func(get func(domain.PlanetShadbala) float64) map[string]float64 {
		m := map[string]float64{}
		for _, p := range ps {
			m[p.Planet] = round2(get(p))
		}
		return m
	}
	table := func(title string, total func(domain.PlanetShadbala) float64, rows ...domain.ShadbalaTableRow) domain.ShadbalaTable {
		return domain.ShadbalaTable{Title: title, Rows: rows, Total: totals(total)}
	}

	var t domain.ShadbalaTables
	t.Planets = names
	t.SthanaBala = table("Sthana Bala (positional strength)", func(p domain.PlanetShadbala) float64 { return p.SthanaBala.Total },
		row("uchcha_bala", "Uchcha Bala", func(p domain.PlanetShadbala) float64 { return p.SthanaBala.UchchaBala }),
		row("saptavargaja_bala", "Saptavargaja Bala", func(p domain.PlanetShadbala) float64 { return p.SthanaBala.SaptavargajaBala }),
		row("ojayugma_bala", "Ojayugmarasyamsa Bala", func(p domain.PlanetShadbala) float64 { return p.SthanaBala.OjayugmaBala }),
		row("kendradi_bala", "Kendra Bala", func(p domain.PlanetShadbala) float64 { return p.SthanaBala.KendradiBala }),
		row("drekkana_bala", "Drekkana Bala", func(p domain.PlanetShadbala) float64 { return p.SthanaBala.DrekkanaBala }))
	t.DigBala = table("Dig Bala (directional strength)", func(p domain.PlanetShadbala) float64 { return p.DigBala },
		row("dig_bala", "Dig Bala", func(p domain.PlanetShadbala) float64 { return p.DigBala }))
	t.KalaBala = table("Kala Bala (temporal strength)", func(p domain.PlanetShadbala) float64 { return p.KalaBala.Total },
		row("nathonnatha_bala", "Nathonnatha Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.NathonnathaBala }),
		row("paksha_bala", "Paksha Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.PakshaBala }),
		row("tribhaga_bala", "Tribhaga Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.TribhagaBala }),
		row("varsha_bala", "Abda (Varsha) Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.VarshaBala }),
		row("masa_bala", "Masa Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.MasaBala }),
		row("dina_bala", "Vara (Dina) Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.DinaBala }),
		row("hora_bala", "Hora Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.HoraBala }),
		row("ayana_bala", "Ayana Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.AyanaBala }),
		row("yuddha_bala", "Yuddha Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.YuddhaBala }))
	t.CheshtaBala = table("Cheshta Bala (motional strength)", func(p domain.PlanetShadbala) float64 { return p.CheshtaBala },
		row("cheshta_bala", "Cheshta Bala", func(p domain.PlanetShadbala) float64 { return p.CheshtaBala }))
	t.NaisargikaBala = table("Naisargika Bala (natural strength)", func(p domain.PlanetShadbala) float64 { return p.NaisargikaBala },
		row("naisargika_bala", "Naisargika Bala", func(p domain.PlanetShadbala) float64 { return p.NaisargikaBala }))
	t.DrikBala = table("Drik Bala (aspect strength)", func(p domain.PlanetShadbala) float64 { return p.DrikBala },
		row("drik_bala", "Drik Bala", func(p domain.PlanetShadbala) float64 { return p.DrikBala }))

	pinda := domain.ShadbalaPindaTable{
		Title: "Shadbala Pinda",
		Rows: []domain.ShadbalaTableRow{
			row("sthana_bala", "Sthana Bala", func(p domain.PlanetShadbala) float64 { return p.SthanaBala.Total }),
			row("dig_bala", "Dig Bala", func(p domain.PlanetShadbala) float64 { return p.DigBala }),
			row("kala_bala", "Kala Bala", func(p domain.PlanetShadbala) float64 { return p.KalaBala.Total }),
			row("cheshta_bala", "Cheshta Bala", func(p domain.PlanetShadbala) float64 { return p.CheshtaBala }),
			row("naisargika_bala", "Naisargika Bala", func(p domain.PlanetShadbala) float64 { return p.NaisargikaBala }),
			row("drik_bala", "Drik Bala", func(p domain.PlanetShadbala) float64 { return p.DrikBala }),
		},
		Total:                totals(func(p domain.PlanetShadbala) float64 { return p.TotalShadbala }),
		Rupas:                map[string]float64{},
		MinimumRequiredRupas: map[string]float64{},
		StrengthRatio:        map[string]float64{},
		Rank:                 map[string]int{},
	}
	for _, p := range ps {
		pinda.Rupas[p.Planet] = p.Rupas
		pinda.MinimumRequiredRupas[p.Planet] = p.MinimumRequired
		pinda.StrengthRatio[p.Planet] = p.StrengthRatio
		pinda.Rank[p.Planet] = p.Rank
	}
	t.ShadbalaPinda = pinda
	return t
}

// assignRanks sets Rank (1 = strongest) from the strength ratio, as in Raman's example 57; planets with the
// same ratio keep their usual order (Sun to Saturn).
func assignRanks(ps []domain.PlanetShadbala) {
	idx := make([]int, len(ps))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(a, b int) bool { return ps[idx[a]].StrengthRatio > ps[idx[b]].StrengthRatio })
	for rank, i := range idx {
		ps[i].Rank = rank + 1
	}
}

// ishtaKashtaTable lays the Ishta and Kashta Phala out like the other tables: the parts they are made from
// and the two phalas, each by planet (there is no total; the parts are not added).
func ishtaKashtaTable(list []domain.IshtaKashtaPlanet) domain.ShadbalaTable {
	row := func(component, label string, get func(domain.IshtaKashtaPlanet) float64) domain.ShadbalaTableRow {
		r := domain.ShadbalaTableRow{Component: component, Label: label, Values: map[string]float64{}}
		for _, p := range list {
			r.Values[p.Planet] = get(p)
		}
		return r
	}
	return domain.ShadbalaTable{
		Title: "Ishta and Kashta Phala (benefic and malefic tendency)",
		Rows: []domain.ShadbalaTableRow{
			row("uchcha_bala", "Uchcha Bala", func(p domain.IshtaKashtaPlanet) float64 { return p.UchchaBala }),
			row("chesta_bala", "Chesta Bala", func(p domain.IshtaKashtaPlanet) float64 { return p.ChestaBala }),
			row("ishta_phala", "Ishta Phala", func(p domain.IshtaKashtaPlanet) float64 { return p.IshtaPhala }),
			row("kashta_phala", "Kashta Phala", func(p domain.IshtaKashtaPlanet) float64 { return p.KashtaPhala }),
		},
	}
}
