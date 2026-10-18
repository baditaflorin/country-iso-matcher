package memory

import (
	"encoding/csv"
	"os"
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/pkg/normalizer"
)

func TestDomainScopeLocalizedCountryAliases(t *testing.T) {
	fixture, err := os.Open("testdata/domainscope_country_aliases.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer fixture.Close()

	rows, err := csv.NewReader(fixture).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatalf("load bundled country data: %v", err)
	}
	for _, row := range rows[1:] {
		if len(row) != 2 {
			t.Fatalf("invalid fixture row with %d columns", len(row))
		}
		got, err := repo.FindByName(row[0])
		if err != nil {
			t.Errorf("country alias lookup failed: %v", err)
			continue
		}
		if got.ISO2 != row[1] {
			t.Errorf("country alias resolved to %s, want %s", got.ISO2, row[1])
		}
	}
}
