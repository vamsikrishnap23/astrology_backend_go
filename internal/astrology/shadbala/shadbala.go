package shadbala

import (
	"math"
	"time"

	"github.com/vamsi/astrology_backend_go/internal/astrology/tables"
	"github.com/vamsi/astrology_backend_go/internal/astrology/vargas"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/houses"
	"github.com/vamsi/astrology_backend_go/internal/astronomy/planets"
	"github.com/vamsi/astrology_backend_go/internal/domain"
)

var TargetPlanets = []string{"Sun", "Moon", "Mars", "Mercury", "Jupiter", "Venus", "Saturn"}

var MinimumRequirements = map[string]float64{
	"Sun":     6.5,
	"Moon":    6.0,
	"Mars":    5.0,
	"Mercury": 7.0,
	"Jupiter": 6.5,
	"Venus":   5.5,
	"Saturn":  5.0,
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
	vargasRes := vargas.CalculateVargas(tablesRes)

	// Create lookup maps
	lons := make(map[string]float64)
	signs := make(map[string]int)
	retro := make(map[string]bool)
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
		retro[p.Planet] = p.Retrograde
		decl[p.Planet] = p.Declination
	}

	var res domain.ShadbalaResult
	res.CalculationTimeUTC = ctx.UTCTime.Format(time.RFC3339)

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
		sapta := calculateSaptavargajaBala(pName, vargasRes, signs)
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
		kala := CalculateKalaBala(pName, ctx, lons["Sun"], lons["Moon"], decl[pName])

		// Cheshta Bala
		cheshta := CalculateCheshtaBala(pName, retro[pName], kala.AyanaBala)

		// Naisargika Bala
		naisargika := CalculateNaisargikaBala(pName)

		// Drik Bala
		drik := CalculateDrikBala(pName, lon, lons)

		total := sthana.Total + dig + kala.Total + cheshta + naisargika + drik
		rupas := total / 60.0
		req := MinimumRequirements[pName]

		ps := domain.PlanetShadbala{
			Planet:          pName,
			SthanaBala:      sthana,
			DigBala:         math.Round(dig*100) / 100,
			KalaBala:        kala,
			CheshtaBala:     math.Round(cheshta*100) / 100,
			NaisargikaBala:  math.Round(naisargika*100) / 100,
			DrikBala:        math.Round(drik*100) / 100,
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
