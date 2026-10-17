package shadbala

import (
	"math"
	"time"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/tables"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/vargas"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

var TargetPlanets = []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"}

// MinimumRequirements are the Shadbala Pindas, in rupas, at which a planet is strong (Raman art. 122 and
// the table on p. 109: 300, 360, 300, 420, 390, 330 and 300 shashtiamsas).
var MinimumRequirements = map[string]float64{
	"Sun":     5.0,
	"Moon":    6.0,
	"Mars":    5.0,
	"Mercury": 7.0,
	"Jupiter": 6.5,
	"Venus":   5.5,
	"Saturn":  5.0,
}

// planetParts holds the components of one planet's Shadbala before the totals are formed.
type planetParts struct {
	sthana                         domain.SthanaBala
	kala                           domain.KalaBala
	dig, cheshta, naisargika, drik float64
}

func CalculateShadbala(ctx *domain.CalculationContext) (domain.ShadbalaResult, error) {
	calcPlanets, err := planets.CalculatePlanets(ctx)
	if err != nil {
		return domain.ShadbalaResult{}, err
	}

	ascLon, mcLon, _, err := houses.CalculateHouses(ctx)
	if err != nil {
		return domain.ShadbalaResult{}, err
	}

	ascSign := int(math.Floor(math.Mod(ascLon, 360.0) / 30.0))

	// Get Vargas
	hCusps := []domain.HouseCusp{{HouseNumber: 1, Longitude: ascLon}} // Minimal houses for tables if full isn't required. Or calculate full.
	tablesRes := tables.GenerateTables(calcPlanets, hCusps)
	vargasRes := vargas.CalculateVargas(tablesRes, hCusps)

	// Create lookup maps
	lons := make(map[string]float64)
	signs := make(map[string]int)
	decl := make(map[string]float64)

	var navamshaVarga domain.VargaChart
	for _, v := range vargasRes.Vargas {
		if v.Division == 9 {
			navamshaVarga = v
		}
	}

	for _, p := range calcPlanets {
		lons[p.Planet] = p.SiderealLongitude
		signs[p.Planet] = int(math.Floor(math.Mod(p.SiderealLongitude, 360.0) / 30.0))
		decl[p.Planet] = p.Declination
	}

	var res domain.ShadbalaResult
	res.CalculationTimeUTC = ctx.UTCTime.Format(time.RFC3339)

	mercuryMalefic := mercuryIsMalefic(signs)
	benefic := beneficMap(lons, signs)

	parts := map[string]*planetParts{}
	for _, pName := range TargetPlanets {
		lon := lons[pName]

		var navSign int
		for _, vp := range navamshaVarga.Planets {
			if vp.Planet == pName {
				navSign = getSignIndexByName(vp.DivisionalSign)
				break
			}
		}

		// Sthana Bala
		uchcha := calculateUchchaBala(pName, lon)
		sapta := calculateSaptavargajaBala(pName, vargasRes, signs, lon)
		oja := calculateOjayugmaBala(pName, signs[pName], navSign)
		kendra := calculateKendradiBala(pName, signs[pName], ascSign)
		drekkana := calculateDrekkanaBala(pName, lon)

		sthana := domain.SthanaBala{
			UchchaBala:       math.Round(uchcha*100) / 100,
			SaptavargajaBala: math.Round(sapta*100) / 100,
			OjayugmaBala:     math.Round(oja*100) / 100,
			KendradiBala:     math.Round(kendra*100) / 100,
			DrekkanaBala:     math.Round(drekkana*100) / 100,
		}
		sthana.Total = sthana.UchchaBala + sthana.SaptavargajaBala + sthana.OjayugmaBala + sthana.KendradiBala + sthana.DrekkanaBala

		// Dig Bala
		dig := CalculateDigBala(pName, lon, ascLon, mcLon)

		// Kala Bala
		kala, err := CalculateKalaBala(pName, ctx, lons["Sun"], lons["Moon"], decl[pName], mercuryMalefic)
		if err != nil {
			return domain.ShadbalaResult{}, err
		}

		// Cheshta Bala
		cheshta := CalculateCheshtaBala(pName, ctx, lons)

		// Naisargika Bala
		naisargika := CalculateNaisargikaBala(pName)

		// Drik Bala
		drik := CalculateDrikBala(pName, lons, benefic)

		parts[pName] = &planetParts{sthana: sthana, dig: dig, kala: kala, cheshta: cheshta, naisargika: naisargika, drik: drik}
	}

	// Yuddha Bala: planets less than a degree apart fight, and the winner's Kala Bala gains what the
	// loser's loses (Raman art. 76). The aggregates use Sthana, Dig and Kala up to Hora Bala.
	aggregate := map[string]float64{}
	for _, pName := range TargetPlanets {
		pp := parts[pName]
		aggregate[pName] = pp.sthana.Total + pp.dig + pp.kala.Total - pp.kala.AyanaBala
	}
	for pName, y := range yuddhaBala(lons, aggregate) {
		pp := parts[pName]
		pp.kala.YuddhaBala = y
		pp.kala.Total += y
	}

	for _, pName := range TargetPlanets {
		pp := parts[pName]
		total := pp.sthana.Total + pp.dig + pp.kala.Total + pp.cheshta + pp.naisargika + pp.drik
		rupas := total / 60.0
		req := MinimumRequirements[pName]

		ps := domain.PlanetShadbala{
			Planet:          pName,
			SthanaBala:      pp.sthana,
			DigBala:         math.Round(pp.dig*100) / 100,
			KalaBala:        pp.kala,
			CheshtaBala:     math.Round(pp.cheshta*100) / 100,
			NaisargikaBala:  math.Round(pp.naisargika*100) / 100,
			DrikBala:        math.Round(pp.drik*100) / 100,
			TotalShadbala:   math.Round(total*100) / 100,
			Rupas:           math.Round(rupas*1000) / 1000,
			MinimumRequired: req,
			StrengthRatio:   math.Round((rupas/req)*1000) / 1000,
			MeetsMinimum:    rupas >= req,
		}

		res.Planets = append(res.Planets, ps)
	}

	return res, nil
}
