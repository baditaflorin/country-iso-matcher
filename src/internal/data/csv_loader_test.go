package data

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCSVLoaderUsesISO3ColumnAndKeepsLegacyFallback(t *testing.T) {
	dir := t.TempDir()
	countries := filepath.Join(dir, "countries.csv")
	aliases := filepath.Join(dir, "aliases.csv")
	if err := os.WriteFile(countries, []byte("code,iso3,name\nRO,ROU,Romania\nUS,United States\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(aliases, []byte("code,alias\nRO,romania\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := NewCSVLoader(countries, aliases).LoadCountries()
	if err != nil {
		t.Fatal(err)
	}
	if got[0].ISO2 != "RO" || got[0].ISO3 != "ROU" {
		t.Fatalf("three-column CSV country = %+v, want RO/ROU", got[0])
	}
	if got[1].ISO2 != "US" || got[1].ISO3 != "US" {
		t.Fatalf("legacy CSV country = %+v, want US fallback", got[1])
	}
}

func TestBundledCSVProvidesISO3ForEveryCountry(t *testing.T) {
	countries, err := NewCSVLoader("../../../data/countries.csv", "").LoadCountries()
	if err != nil {
		t.Fatal(err)
	}
	if len(countries) != 249 {
		t.Fatalf("bundled country count = %d, want all 249 assigned ISO 3166-1 entries", len(countries))
	}
	for _, country := range countries {
		if len(country.ISO3) != 3 {
			t.Fatalf("%s ISO3 = %q, want three-letter alpha-3 code", country.ISO2, country.ISO3)
		}
	}
	byCode := make(map[string]bool, len(countries))
	for _, country := range countries {
		byCode[country.ISO2] = true
	}
	for _, code := range []string{"AI", "AX", "BQ", "CW", "SJ", "VA", "YT"} {
		if !byCode[code] {
			t.Errorf("missing long-tail ISO2 code %s", code)
		}
	}
}

func TestBundledAliasesCoverEveryCountry(t *testing.T) {
	loader := NewCSVLoader("../../../data/countries.csv", "../../../data/aliases.csv")
	countries, err := loader.LoadCountries()
	if err != nil {
		t.Fatal(err)
	}
	aliases, err := loader.LoadAliases()
	if err != nil {
		t.Fatal(err)
	}
	for _, country := range countries {
		if len(aliases[country.ISO2]) == 0 {
			t.Errorf("%s has no multilingual/endonym aliases", country.ISO2)
		}
	}
}
