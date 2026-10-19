package memory

import (
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/pkg/normalizer"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

func TestAllCLDRKoreanCountryNamesResolve(t *testing.T) {
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	countries, err := loader.LoadCountries()
	if err != nil {
		t.Fatalf("load bundled country data: %v", err)
	}
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loader)
	if err != nil {
		t.Fatalf("load country repository: %v", err)
	}

	koreanNames := display.Regions(language.Korean)
	matched := 0
	for _, country := range countries {
		region, err := language.ParseRegion(country.ISO2)
		if err != nil {
			t.Fatalf("parse ISO2 country code %q: %v", country.ISO2, err)
		}
		query := koreanNames.Name(region)
		got, err := repo.FindByName(query)
		if err != nil {
			t.Errorf("lookup Korean CLDR name for ISO2 %s: %v", country.ISO2, err)
			continue
		}
		if got.ISO2 != country.ISO2 {
			t.Errorf("Korean CLDR name for ISO2 %s resolved to %s", country.ISO2, got.ISO2)
			continue
		}
		matched++
	}
	if matched != len(countries) {
		t.Fatalf("resolved %d/%d Korean CLDR country names", matched, len(countries))
	}
}

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

func TestCanonicalISOCodesTakePrecedenceOverConflictingAliases(t *testing.T) {
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loader)
	if err != nil {
		t.Fatalf("load bundled country data: %v", err)
	}

	for _, tc := range []struct {
		query string
		want  string
	}{
		{query: "SA", want: "SA"}, // South Africa's short alias must not replace Saudi Arabia's ISO2 code.
		{query: "SK", want: "SK"}, // South Korea's short alias must not replace Slovakia's ISO2 code.
		{query: "South Africa", want: "ZA"},
		{query: "South Korea", want: "KR"},
	} {
		country, err := repo.FindByName(tc.query)
		if err != nil {
			t.Fatalf("lookup %q: %v", tc.query, err)
		}
		if country.ISO2 != tc.want {
			t.Errorf("lookup %q = %s, want %s", tc.query, country.ISO2, tc.want)
		}
	}
}

func TestCommonShortNamesResolveForSyriaAndVenezuela(t *testing.T) {
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loader)
	if err != nil {
		t.Fatalf("load bundled country data: %v", err)
	}

	for _, tc := range []struct {
		query string
		want  string
	}{
		{query: "Syria", want: "SY"},
		{query: "Venezuela", want: "VE"},
	} {
		country, err := repo.FindByName(tc.query)
		if err != nil {
			t.Errorf("lookup %q: %v", tc.query, err)
			continue
		}
		if country.ISO2 != tc.want {
			t.Errorf("lookup %q = %s, want %s", tc.query, country.ISO2, tc.want)
		}
	}
}
