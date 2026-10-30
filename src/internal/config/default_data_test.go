package config_test

import (
	"testing"

	"country-iso-matcher/src/internal/config"
	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/internal/repository/memory"
	"country-iso-matcher/src/pkg/normalizer"
)

// The shipped default must serve the same full country set as the container.
func TestDefaultDataIncludesCountriesOutsideSampleJSON(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Data.CountriesFile = "../../../data/countries.csv"
	cfg.Data.AliasesFile = "../../../data/aliases.csv"
	cfg.Data.CountriesDir = "../../../data/countries"
	loader, err := data.NewLoader(&cfg.Data)
	if err != nil {
		t.Fatal(err)
	}
	repo, err := memory.NewCountryRepository(normalizer.NewTextNormalizer(), loader)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, code string }{{"Canada", "CA"}, {"Finland", "FI"}, {"Mexico", "MX"}, {"Hong Kong", "HK"}} {
		country, err := repo.FindByName(tc.name)
		if err != nil {
			t.Errorf("lookup %q: %v", tc.name, err)
			continue
		}
		if country.ISO2 != tc.code {
			t.Errorf("lookup %q = %s, want %s", tc.name, country.ISO2, tc.code)
		}
	}
}

func TestExplicitJSONSourceStillSupported(t *testing.T) {
	t.Setenv("DATA_SOURCE", "json")
	cfg, err := config.Load("")
	if err != nil {
		t.Fatal(err)
	}
	cfg.Data.CountriesDir = "../../../data/countries"
	loader, err := data.NewLoader(&cfg.Data)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := loader.LoadCountries(); err != nil {
		t.Fatalf("explicit JSON: %v", err)
	}
}
