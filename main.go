package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	store := NewReceiptStore()
	mux := http.NewServeMux()

	mux.HandleFunc("/receipts/process", processReceiptHandler(store))
	mux.HandleFunc("/receipts/{id}/points", getPointsHandler(store))

	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", mux)
}

func processReceiptHandler(store ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rec, err := ParseReceiptJson(r.Body)
		if err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		id := store.Save(rec)

		response := ReceiptIDResponse{ID: id}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

func getPointsHandler(store ReceiptStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Receipt ID is required", http.StatusBadRequest)
			return
		}

		rec, exits := store.Get(id)
		if rec == nil || !exits {
			http.Error(w, fmt.Sprintf("No receipt found with id: %v", id), http.StatusBadRequest)
			return
		}
		if rec.Points == UninitializedPoints {
			CalculateReceiptPoints(rec)
		}

		response := PointsResponse{Points: rec.Points}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
