package benchmarks

import (
	"testing"

	"github.com/baditaflorin/country-iso-matcher/internal/repository/memory"
	"github.com/baditaflorin/country-iso-matcher/internal/service"
	"github.com/baditaflorin/country-iso-matcher/pkg/normalizer"
)

func BenchmarkCountryLookup(b *testing.B) {
	// Setup
	normalizer := normalizer.NewTextNormalizer()
	repo := memory.NewCountryRepository(normalizer)
	service := service.NewCountryService(repo)

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
		_, err := service.LookupCountry(country)
		if err != nil {
			b.Errorf("unexpected error: %v", err)
		}
	}
}

func BenchmarkNormalizer(b *testing.B) {
	normalizer := normalizer.NewTextNormalizer()

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
		_ = normalizer.Normalize(input)
	}
}
