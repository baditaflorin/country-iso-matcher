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
	if len(countries) != 183 {
		t.Fatalf("bundled country count = %d, want 183", len(countries))
	}
	for _, country := range countries {
		if len(country.ISO3) != 3 {
			t.Fatalf("%s ISO3 = %q, want three-letter alpha-3 code", country.ISO2, country.ISO3)
		}
	}
}
