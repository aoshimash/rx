package main

import (
	"log"
	"net/http"
)

func main() {
	// TODO: Initialize chi router and OpenAPI handlers
	// This will be implemented after code generation works

	http.HandleFunc("/api/v1/workouts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]")) // Return empty array for now
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
