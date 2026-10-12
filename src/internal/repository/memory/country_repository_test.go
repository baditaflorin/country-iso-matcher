package memory

import (
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/pkg/normalizer"
)

func TestBundledIraqAliasIsNotShadowedByIran(t *testing.T) {
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	for attempt := 0; attempt < 20; attempt++ {
		repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loader)
		if err != nil {
			t.Fatalf("load bundled country data: %v", err)
		}
		for _, query := range []string{"irak", "irák"} {
			country, err := repo.FindByName(query)
			if err != nil {
				t.Fatalf("lookup %q: %v", query, err)
			}
			if country.ISO2 != "IQ" {
				t.Fatalf("lookup %q = %s, want IQ (Iraq); attempt %d", query, country.ISO2, attempt+1)
			}
		}
		country, err := repo.FindByName("Persia")
		if err != nil {
			t.Fatalf("lookup Persia: %v", err)
		}
		if country.ISO2 != "IR" {
			t.Fatalf("lookup Persia = %s, want IR (Iran)", country.ISO2)
		}
	}
}
