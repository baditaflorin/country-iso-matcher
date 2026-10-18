package memory

import (
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/pkg/normalizer"
)

func TestDefunctMultiCountryStatesDoNotResolveToOneSuccessor(t *testing.T) {
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loader)
	if err != nil {
		t.Fatalf("load bundled country data: %v", err)
	}

	for _, name := range []string{
		"Czechoslovakia",
		"Soviet Union",
		"USSR",
		"Union of Soviet Socialist Republics",
	} {
		t.Run(name, func(t *testing.T) {
			if got, err := repo.FindByName(name); err == nil {
				t.Fatalf("historical state %q resolved to %s; it has no single current ISO 3166-1 code", name, got.ISO2)
			}
		})
	}
}

func TestCurrentCountryAliasMatrixAfterRemovingHistoricalAliases(t *testing.T) {
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loader)
	if err != nil {
		t.Fatalf("load bundled country data: %v", err)
	}

	for _, tc := range []struct{ name, want string }{
		{"Czechia", "CZ"},
		{"Slovakia", "SK"},
		{"Russia", "RU"},
		{"ÖSTERREICH", "AT"},
		{"o\u0308sterreich", "AT"}, // decomposed Unicode still matches the composed endonym.
		{"Shqipëria", "AL"},
		{"Magyarország", "HU"},
		{"Nippon", "JP"},
		{"Suomi", "FI"},
		{"Côte d'Ivoire", "CI"},
	} {
		country, err := repo.FindByName(tc.name)
		if err != nil {
			t.Errorf("lookup %q: %v", tc.name, err)
			continue
		}
		if country.ISO2 != tc.want {
			t.Errorf("lookup %q = %s, want %s", tc.name, country.ISO2, tc.want)
		}
	}
}
