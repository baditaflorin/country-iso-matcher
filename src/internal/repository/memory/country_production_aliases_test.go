package memory

import (
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/pkg/normalizer"
)

func TestDomainScopeOfficialCountryNameMissesResolve(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}

	// These Unicode CLDR English territory names were absent from the alias corpus,
	// although DomainScope country references assign them to the same ISO alpha-2 code.
	for _, tc := range []struct{ name, want string }{
		{"Brunei", "BN"},
		{"Bolivia", "BO"},
		{"Micronesia", "FM"},
		{"Laos", "LA"},
		{"Pitcairn Islands", "PN"},
		{"Sint Maarten", "SX"},
		{"Tanzania", "TZ"},
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

func TestUnknownCountryNameStaysUnresolvedAfterAliasExpansion(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.FindByName("fictional-country-name-20261018"); err == nil {
		t.Fatal("unknown country unexpectedly resolved")
	}
}
