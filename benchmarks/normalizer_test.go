package benchmarks

import (
	"testing"

	"country-iso-matcher/src/pkg/normalizer"
)

func BenchmarkNormalizer(b *testing.B) {
	textNormalizer := normalizer.NewTextNormalizer()

	inputs := []string{
		"Côte d'Ivoire",
		"DEUTSCHLAND",
		"  United States of America  ",
		"République française",
		"中国",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		input := inputs[i%len(inputs)]
		_ = textNormalizer.Normalize(input)
	}
}
