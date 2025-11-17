package benchmarks

import (
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/internal/repository/memory"
	"country-iso-matcher/src/internal/service"
	"country-iso-matcher/src/pkg/normalizer"
)

func BenchmarkCountryLookup(b *testing.B) {
	// Setup
	norm := normalizer.NewTextNormalizer()
	loader := data.NewMemoryLoader()
	repo, err := memory.NewCountryRepository(norm, loader)
	if err != nil {
		b.Fatalf("failed to create repository: %v", err)
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
		_, err := svc.LookupCountry(country)
		if err != nil {
			b.Errorf("unexpected error: %v", err)
		}
	}
}

func BenchmarkNormalizer(b *testing.B) {
	norm := normalizer.NewTextNormalizer()

	inputs := []string{
		"Côte d'Ivoire",
		"DEUTSCHLAND",
		"  United States of America  ",
		"République française",
		"中国",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)]
		_ = norm.Normalize(input)
	}
}
