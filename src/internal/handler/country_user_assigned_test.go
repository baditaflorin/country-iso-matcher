package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/internal/repository/memory"
	"country-iso-matcher/src/internal/service"
	"country-iso-matcher/src/pkg/normalizer"
)

func TestUserAssignedKosovoCodesRequireExplicitOptIn(t *testing.T) {
	repo, err := memory.NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../data/countries.csv", "../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}
	h := NewCountryHandler(service.NewCountryService(repo), slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, tc := range []struct {
		name, optIn string
		status      int
		code        string
		assigned    bool
	}{
		{"Kosovo", "", 404, "", false},
		{"XK", "false", 404, "", false},
		{"Republic of Kosovo", "1", 404, "", false},
		{"Kosovo", "true", 200, "XK", true},
		{"republic of kosovo", "true", 200, "XK", true},
		{"xk", "true", 200, "XK", true},
		{"XKX", "true", 200, "XK", true},
		{"Germany", "true", 200, "DE", false},
		{"unknown", "true", 404, "", false},
	} {
		t.Run(tc.name+"/"+tc.optIn, func(t *testing.T) {
			values := url.Values{"country": {tc.name}}
			if tc.optIn != "" {
				values.Set("includeUserAssigned", tc.optIn)
			}
			w := httptest.NewRecorder()
			h.ConvertCountry(w, httptest.NewRequest(http.MethodGet, "/api/convert?"+values.Encode(), nil))
			if w.Code != tc.status {
				t.Fatalf("status = %d, want %d; response = %s", w.Code, tc.status, w.Body.String())
			}
			if tc.status != 200 {
				return
			}
			var response struct {
				ISO2     string `json:"iso2Code"`
				ISO3     string `json:"iso3Code"`
				Assigned bool   `json:"userAssigned"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if response.ISO2 != tc.code || response.Assigned != tc.assigned {
				t.Errorf("response = %+v", response)
			}
			if tc.assigned && response.ISO3 != "XKX" {
				t.Errorf("user-assigned ISO3-like code = %q, want XKX", response.ISO3)
			}
		})
	}
}
