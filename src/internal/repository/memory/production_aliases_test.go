package memory

import (
	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/pkg/normalizer"
	"testing"
)

// Labels observed in independent DomainScope registration and office records.
// The native names are corroborated by Unicode CLDR territory data for fa and th.
func TestHistoricalRegistrationAndOfficeCountryLabels(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ label, want string }{
		{"Netherlands (the)", "NL"},
		{"United States of America (the)", "US"},
		{"United Kingdom of Great Britain and Northern Ireland (the)", "GB"},
		{"Virgin Islands (British)", "VG"},
		{"Hong Kong", "HK"},
		{"Hong Kong SAR", "HK"},
		{"ایران", "IR"},
		{"ไทย", "TH"},
		{"UNITED STATES (US)", "US"},
		{"UNITED KINGDOM (UK)", "GB"},
	} {
		got, err := repo.FindByName(tc.label)
		if err != nil {
			t.Errorf("lookup %q: %v", tc.label, err)
			continue
		}
		if got.ISO2 != tc.want {
			t.Errorf("lookup %q = %s, want %s", tc.label, got.ISO2, tc.want)
		}
	}
	// A region or company suffix is not itself a country name.
	for _, label := range []string{"New York", "LLC", "Hà Nội", "Virgin Islands"} {
		if got, err := repo.FindByName(label); err == nil {
			t.Errorf("ambiguous/non-country %q unexpectedly resolved to %s", label, got.ISO2)
		}
	}
}
