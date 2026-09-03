package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astrology/btr"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func BTRHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.BTRInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Set defaults if needed
	if input.ScanMinusMinutes == 0 {
		input.ScanMinusMinutes = 10
	}
	if input.ScanPlusMinutes == 0 {
		input.ScanPlusMinutes = 5
	}
	if input.Ayanamsa == "" {
		input.Ayanamsa = "Lahiri"
	}

	res, err := btr.ProcessBTR(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
