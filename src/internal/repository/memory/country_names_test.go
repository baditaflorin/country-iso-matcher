package memory

import (
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/pkg/normalizer"
)

func TestCommaCountryNamesAndCongoDisambiguation(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ query, want string }{
		{"Congo", "CG"},
		{"Congo, Democratic Republic of the", "CD"},
		{"Moldova, Republic of", "MD"},
		{"Palestine, State of", "PS"},
		{"Taiwan, Province of China", "TW"},
		{"Tanzania, United Republic of", "TZ"},
	} {
		got, err := repo.FindByName(tc.query)
		if err != nil {
			t.Errorf("lookup %q: %v, want %s", tc.query, err, tc.want)
			continue
		}
		if got.ISO2 != tc.want {
			t.Errorf("lookup %q = %s, want %s", tc.query, got.ISO2, tc.want)
		}
	}
}

func TestProductionOfficeCountryLabelsMatchAcrossLocales(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}

	// These country-only labels were observed in bounded, read-only production
	// office-location samples. Russian/Ukrainian labels are the discovery set;
	// Chinese/Japanese/Korean labels are the disjoint script holdout.
	for _, tc := range []struct{ label, want string }{
		{"Россия", "RU"},
		{"Україна", "UA"},
		{"中国", "CN"},
		{"日本", "JP"},
		{"한국", "KR"},
		{"Magyarország", "HU"},
		{"Česká republika", "CZ"},
		{"Ελλάδα", "GR"},
	} {
		got, err := repo.FindByName(tc.label)
		if err != nil {
			t.Errorf("lookup %q: %v, want %s", tc.label, err, tc.want)
			continue
		}
		if got.ISO2 != tc.want {
			t.Errorf("lookup %q = %s, want %s", tc.label, got.ISO2, tc.want)
		}
	}

	// Nearby negative checks ensure the two Cyrillic endonyms stay attached to
	// their own countries rather than becoming cross-country aliases.
	for _, tc := range []struct{ label, wrong string }{
		{"Россия", "UA"},
		{"Україна", "RU"},
	} {
		got, err := repo.FindByName(tc.label)
		if err == nil && got.ISO2 == tc.wrong {
			t.Errorf("lookup %q incorrectly resolved to %s", tc.label, tc.wrong)
		}
	}
}

func TestProductionAliasesResolveToDomainScopeCountryCodes(t *testing.T) {
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loader)
	if err != nil {
		t.Fatal(err)
	}
	aliases, err := loader.LoadAliases()
	if err != nil {
		t.Fatal(err)
	}

	for _, code := range []string{"HK", "XK"} {
		country, err := repo.FindByCode(code)
		if err != nil {
			t.Fatalf("country code %s: %v", code, err)
		}
		if country == nil || country.ISO2 != code {
			t.Fatalf("country code %s did not return its country record", code)
		}
		for _, alias := range aliases[code] {
			got, err := repo.FindByName(alias)
			if err != nil {
				t.Errorf("alias for %s did not resolve: %v", code, err)
				continue
			}
			if got == nil || got.ISO2 != code {
				t.Errorf("alias for %s resolved to an unexpected country", code)
			}
		}
	}

	kosovo, err := repo.FindByName("Kosovo")
	if err != nil {
		t.Fatalf("Kosovo country name: %v", err)
	}
	if kosovo == nil || kosovo.ISO2 != "XK" || kosovo.ISO3 != "XKX" {
		t.Fatalf("Kosovo country name did not resolve to the DomainScope XK country code")
	}
}
