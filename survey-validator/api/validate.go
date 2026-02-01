package api

import (
	"encoding/json"
	"net/http"

	"github.com/survey-validator/engine"
	"github.com/survey-validator/models"
)

// Handler is the Vercel serverless function handler for /api/v1/validate
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed. Use POST.")
		return
	}

	var surveyData models.SurveyData
	if err := json.NewDecoder(r.Body).Decode(&surveyData); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}
	defer r.Body.Close()

	eng := engine.NewEngine()
	report := eng.Validate(&surveyData)
	respondJSON(w, http.StatusOK, report)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error": message,
	})
}
