package main

import (
	"encoding/json"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"log"
	"net/http"
	"strings"
	"unicode"
)

// Response structs for clean JSON marshaling
type SuccessResponse struct {
	Query        string `json:"query"`
	OfficialName string `json:"officialName"`
	ISOCode      string `json:"isoCode"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Query string `json:"query,omitempty"`
}

var (
	normalizedCountryNames = make(map[string]string)
	isoToOfficialName      = make(map[string]string)
	accentRemover          = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
)

// init runs once when the application starts to build our lookup maps.
func init() {
	// 1. Build lookup tables from the official country data
	for _, country := range countries {
		isoToOfficialName[country.Code] = country.Name
		normalized, _, _ := transform.String(accentRemover, country.Name)
		normalizedCountryNames[strings.ToLower(normalized)] = country.Code
		normalizedCountryNames[strings.ToLower(country.Name)] = country.Code
	}

	// 2. Add all custom aliases to the lookup table
	for isoCode, aliases := range countryAliases {
		for _, alias := range aliases {
			normalized, _, _ := transform.String(accentRemover, alias)
			normalizedCountryNames[strings.ToLower(normalized)] = isoCode
		}
	}
}

// normalizeString removes accents and converts to lowercase.
func normalizeString(s string) string {
	normalized, _, _ := transform.String(accentRemover, s)
	return strings.ToLower(normalized)
}

// countryHandler processes the API requests.
func countryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	countryName := r.URL.Query().Get("country")

	if countryName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Country query parameter is required."})
		return
	}

	normalizedInput := normalizeString(countryName)
	isoCode, found := normalizedCountryNames[normalizedInput]

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Country not found or name is invalid.", Query: countryName})
		return
	}

	officialName := isoToOfficialName[isoCode]
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
