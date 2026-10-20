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

func TestAllPinnedCLDRLocaleCountryNamesResolve(t *testing.T) {
	if display.CLDRVersion != "32" {
		t.Fatalf("update the legacy-locale alias corpus for CLDR %s", display.CLDRVersion)
	}
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	countries, err := loader.LoadCountries()
	if err != nil {
		t.Fatalf("load bundled country data: %v", err)
	}
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loader)
	if err != nil {
		t.Fatalf("load country repository: %v", err)
	}

	// The aliases include CLDR 48.2.3 names and retain CLDR 32 display names
	// from the pinned x/text dependency as legacy inputs used by older sources.
	for _, locale := range []string{"ar", "de", "es", "fr", "hi", "it", "ja", "nl", "pt", "ru", "tr", "uk", "vi", "zh"} {
		names := display.Regions(language.MustParse(locale))
		checked := 0
		for _, country := range countries {
			region, err := language.ParseRegion(country.ISO2)
			if err != nil {
				t.Fatalf("parse ISO2 country code %q: %v", country.ISO2, err)
			}
			query := names.Name(region)
			if query == "" || query == country.ISO2 {
				continue
			}
			got, err := repo.FindByName(query)
			if err != nil {
				t.Errorf("lookup CLDR %s country label for ISO2 %s: %v", locale, country.ISO2, err)
				continue
			}
			if got.ISO2 != country.ISO2 {
				t.Errorf("CLDR %s country label for ISO2 %s resolved to %s", locale, country.ISO2, got.ISO2)
				continue
			}
			checked++
		}
		if checked != len(countries) {
			t.Errorf("resolved %d/%d CLDR %s country labels", checked, len(countries), locale)
		}
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
