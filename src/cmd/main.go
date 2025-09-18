package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/baditaflorin/country-iso-matcher/src/pkg/countries"
)

type SuccessResponse struct {
	Query        string `json:"query"`
	OfficialName string `json:"officialName"`
	ISOCode      string `json:"isoCode"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Query string `json:"query,omitempty"`
}

func countryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	countryName := r.URL.Query().Get("country")

	if countryName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Country query parameter is required."})
		return
	}

	isoCode, officialName, found := countries.LookupISO(countryName)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Country not found or name is invalid.", Query: countryName})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SuccessResponse{
		Query:        countryName,
		OfficialName: officialName,
		ISOCode:      isoCode,
	})
}

func main() {
	http.HandleFunc("/api/convert", countryHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go Country to ISO Code Converter is running. Use /api/convert?country=YourCountryName"))
	})

	log.Println("Server starting on port 3030...")
	if err := http.ListenAndServe(":3030", nil); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
