package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/astronomy/retrogrades"
	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

func RetrogradesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input domain.RetrogradesInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if input.StartDate == "" || input.EndDate == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	result, err := retrogrades.CalculateRetrogrades(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
