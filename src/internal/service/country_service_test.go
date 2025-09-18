package service_test

import (
	"errors"
	"testing"

	"github.com/baditaflorin/country-iso-matcher/internal/domain"
	"github.com/baditaflorin/country-iso-matcher/internal/service"
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
		if country.Code == code {
			return country, nil
		}
	}
	return nil, domain.NewNotFoundError(code)
}

func TestCountryService_LookupCountry(t *testing.T) {
	// Setup
	mockRepo := &mockRepository{
		countries: map[string]*domain.Country{
			"romania": {Code: "RO", Name: "Romania"},
			"germany": {Code: "DE", Name: "Germany"},
		},
	}

	service := service.NewCountryService(mockRepo)

	tests := []struct {
		name          string
		query         string
		expectedCode  string
		expectedName  string
		expectedError bool
	}{
		{
			name:         "valid country",
			query:        "romania",
			expectedCode: "RO",
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
			result, err := service.LookupCountry(tt.query)

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

			if result.ISOCode != tt.expectedCode {
				t.Errorf("expected ISO code %s, got %s", tt.expectedCode, result.ISOCode)
			}

			if result.OfficialName != tt.expectedName {
				t.Errorf("expected name %s, got %s", tt.expectedName, result.OfficialName)
			}
		})
	}
}
