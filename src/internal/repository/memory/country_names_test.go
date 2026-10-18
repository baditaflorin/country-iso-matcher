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

func TestBareKoreaIsAmbiguousButQualifiedNamesResolve(t *testing.T) {
	for _, tc := range []struct {
		name   string
		loader data.Loader
	}{
		{name: "bundled CSV", loader: data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")},
		{name: "memory", loader: data.NewMemoryLoader()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), tc.loader)
			if err != nil {
				t.Fatal(err)
			}
			if country, err := repo.FindByName("Korea"); err == nil {
				t.Fatalf("ambiguous bare Korea resolved to %s; want not found", country.ISO2)
			}
			for _, want := range []struct{ query, code string }{
				{"South Korea", "KR"},
				{"North Korea", "KP"},
			} {
				country, err := repo.FindByName(want.query)
				if err != nil || country.ISO2 != want.code {
					t.Errorf("lookup %q = %v, %v; want %s", want.query, country, err, want.code)
				}
			}
		})
	}
}

func TestSampledLocalizedCountryAliasesResolve(t *testing.T) {
	cases := []struct{ query, want string }{
		{"CIN", "CN"},
		{"SIP", "CY"},
		{"República Dominicana", "DO"},
		{"日本", "JP"},
		{"Кыргызстан", "KG"},
		{"MAROC", "MA"},
		{"Украина", "UA"},
	}
	for _, loaderCase := range []struct {
		name   string
		loader data.Loader
	}{
		{name: "bundled CSV", loader: data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")},
		{name: "memory", loader: data.NewMemoryLoader()},
	} {
		t.Run(loaderCase.name, func(t *testing.T) {
			repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loaderCase.loader)
			if err != nil {
				t.Fatal(err)
			}
			for i, tc := range cases {
				got, err := repo.FindByName(tc.query)
				if err != nil || got.ISO2 != tc.want {
					t.Errorf("localized alias case %d resolved to %v, err %v; want %s", i, got, err, tc.want)
				}
			}
		})
	}
}
