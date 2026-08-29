package service_test

import (
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/internal/repository/memory"
	"country-iso-matcher/src/internal/service"
	"country-iso-matcher/src/pkg/normalizer"
)

// BenchmarkCountryLookup exercises the real CSV-backed in-memory repository.
// Lives here rather than in benchmarks/ because Go forbids importing
// src/internal/... from outside the src/ tree.
func BenchmarkCountryLookup(b *testing.B) {
	textNormalizer := normalizer.NewTextNormalizer()
	loader := data.NewCSVLoader("../../../data/countries.csv", "../../../data/aliases.csv")

	repo, err := memory.NewCountryRepository(textNormalizer, loader)
	if err != nil {
		b.Fatalf("failed to build repository: %v", err)
	}
	svc := service.NewCountryService(repo)

	countries := []string{
		"Romania",
		"Germany",
		"United States",
		"France",
		"italy",
		"SPAIN",
		"united kingdom",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		country := countries[i%len(countries)]
		if _, err := svc.LookupCountry(country); err != nil {
			b.Errorf("unexpected error for %q: %v", country, err)
		}
	}
}
