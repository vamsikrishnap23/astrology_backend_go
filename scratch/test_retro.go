package main

import (
	"encoding/json"
	"fmt"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/retrogrades"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/ephemeris"
)

func main() {
	ephemeris.Init("ephe_data")
	defer ephemeris.Close()
	
	input := domain.RetrogradesInput{
		StartDate: "2026-01-01",
		EndDate:   "2026-12-31",
		Timezone:  5.5,
	}

	res, err := retrogrades.CalculateRetrogrades(input)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	
	b, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(b))
}
