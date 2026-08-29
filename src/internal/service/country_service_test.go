package service_test

import (
	"testing"

	"country-iso-matcher/src/internal/domain"
	"country-iso-matcher/src/internal/service"
)

type mockRepository struct {
	countries map[string]*domain.Country
}

func (m *mockRepository) FindByName(name string) (*domain.Country, error) {
	country, exists := m.countries[name]
	if !exists {
		return nil, domain.NewNotFoundError(name)
	}
	return country, nil
}

func (m *mockRepository) FindByCode(code string) (*domain.Country, error) {
	for _, country := range m.countries {
		if country.ISO2 == code {
			return country, nil
		}
	}
	return nil, domain.NewNotFoundError(code)
}

func TestCountryService_LookupCountry(t *testing.T) {
	mockRepo := &mockRepository{
		countries: map[string]*domain.Country{
			"romania": {ISO2: "RO", ISO3: "ROU", Names: map[string]string{"en": "Romania"}},
			"germany": {ISO2: "DE", ISO3: "DEU", Names: map[string]string{"en": "Germany"}},
		},
	}

	svc := service.NewCountryService(mockRepo)

	tests := []struct {
		name          string
		query         string
		expectedISO2  string
		expectedISO3  string
		expectedName  string
		expectedError bool
	}{
		{
			name:         "valid country",
			query:        "romania",
			expectedISO2: "RO",
			expectedISO3: "ROU",
			expectedName: "Romania",
		},
		{
			name:          "empty query",
			query:         "",
			expectedError: true,
		},
		{
			name:          "whitespace query",
			query:         "   ",
			expectedError: true,
		},
		{
			name:          "unknown country",
			query:         "unknown",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.LookupCountry(tt.query)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.ISO2Code != tt.expectedISO2 {
				t.Errorf("expected ISO2 code %s, got %s", tt.expectedISO2, result.ISO2Code)
			}

			if result.ISO3Code != tt.expectedISO3 {
				t.Errorf("expected ISO3 code %s, got %s", tt.expectedISO3, result.ISO3Code)
			}

			if result.OfficialName != tt.expectedName {
				t.Errorf("expected name %s, got %s", tt.expectedName, result.OfficialName)
			}
		})
	}
}
