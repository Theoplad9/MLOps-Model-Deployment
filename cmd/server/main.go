package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Theoplad9/MLOps-Model-Deployment/internal/handler"
	"github.com/Theoplad9/MLOps-Model-Deployment/internal/model"
)

func main() {
	// Load the ML model (dummy for now)
	model.LoadModel()

	r := mux.NewRouter()
	r.HandleFunc("/predict", handler.PredictHandler).Methods("POST")

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
