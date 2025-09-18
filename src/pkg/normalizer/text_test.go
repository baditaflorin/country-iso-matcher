package normalizer_test

import (
	"testing"

	"country-iso-matcher/src/pkg/normalizer"
)

func TestTextNormalizer_Normalize(t *testing.T) {
	textNormalizer := normalizer.NewTextNormalizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic lowercase",
			input:    "Hello",
			expected: "hello",
		},
		{
			name:     "remove accents",
			input:    "Côte d'Ivoire",
			expected: "cote d'ivoire",
		},
		{
			name:     "trim whitespace",
			input:    "  Germany  ",
			expected: "germany",
		},
		{
			name:     "complex with accents and case",
			input:    "  FRANÇAIS  ",
			expected: "francais",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := textNormalizer.Normalize(tt.input)
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
