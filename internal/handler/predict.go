package handler

import (
	"encoding/json"
	"net/http"
	"log"
)

type PredictionRequest struct {
	Data []float64 `json:"data"`
}

type PredictionResponse struct {
	Prediction float64 `json:"prediction"`
	Message    string  `json:"message"`
}

// PredictHandler handles incoming prediction requests
func PredictHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req PredictionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Dummy prediction logic for demonstration
	// In a real scenario, this would involve calling the loaded ML model
	prediction := 0.0
	if len(req.Data) > 0 {
		for _, val := range req.Data {
			prediction += val
		}
		prediction /= float64(len(req.Data))
	}

	resp := PredictionResponse{
		Prediction: prediction,
		Message:    "Prediction successful (dummy output)",
	}

	json.NewEncoder(w).Encode(resp)
	log.Printf("Received prediction request with data: %v, responded with prediction: %f", req.Data, prediction)
}
