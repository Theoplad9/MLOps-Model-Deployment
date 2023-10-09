package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type MLModel interface {
	Predict([]float64) (float64, error)
	Load(string) error
}

type SimpleLinearModel struct {
	weights []float64
	bias    float64
}

func (slm *SimpleLinearModel) Predict(features []float64) (float64, error) {
	if len(features) != len(slm.weights) {
		return 0, fmt.Errorf("feature count mismatch: expected %d, got %d", len(slm.weights), len(features))
	}
	prediction := slm.bias
	for i, f := range features {
		prediction += f * slm.weights[i]
	}
	return prediction, nil
}

func (slm *SimpleLinearModel) Load(modelPath string) error {
	log.Printf("Loading model from %s...", modelPath)
	time.Sleep(1 * time.Second)
	slm.weights = []float64{0.5, 1.2, -0.3}
	slm.bias = 0.1
	log.Println("SimpleLinearModel loaded successfully.")
	return nil
}

var globalModel MLModel

func init() {
	globalModel = &SimpleLinearModel{}
	err := globalModel.Load("/models/dummy_model.json")
	if err != nil {
		log.Fatalf("Failed to load initial model: %v", err)
	}
}

type PredictionRequest struct {
	Features []float64 `json:"features"`
}

type PredictionResponse struct {
	Prediction float64 `json:"prediction"`
	Message    string  `json:"message"`
	Error      string  `json:"error,omitempty"`
}

func PredictHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req PredictionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request payload: %v", err), http.StatusBadRequest)
		return
	}

	if len(req.Features) == 0 {
		http.Error(w, "Features array cannot be empty", http.StatusBadRequest)
		return
	}

	prediction, err := globalModel.Predict(req.Features)
	if err != nil {
		resp := PredictionResponse{Error: err.Error(), Message: "Prediction failed"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := PredictionResponse{
		Prediction: prediction,
		Message:    "Prediction successful",
	}
	json.NewEncoder(w).Encode(resp)
	log.Printf("Received prediction request with features: %v, responded with prediction: %f", req.Features, prediction)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/predict", PredictHandler).Methods("POST")
	r.HandleFunc("/health", healthCheckHandler).Methods("GET")

	port := ":8080"
	fmt.Printf("Server starting on %s
", port)
	log.Fatal(http.ListenAndServe(port, r))
}
