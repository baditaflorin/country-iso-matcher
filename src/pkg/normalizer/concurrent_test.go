package normalizer_test

import (
	"country-iso-matcher/src/pkg/normalizer"
	"sync"
	"testing"
)

func TestNormalizeConcurrentMultilingualInputs(t *testing.T) {
	n := normalizer.NewTextNormalizer()
	in := []string{"Côte d’Ivoire", "España", "Türkiye", "ｱﾒﾘｶ合衆国", "Việt Nam", "Ελλάδα"}
	expected := make([]string, len(in))
	for i, s := range in {
		expected[i] = n.Normalize(s)
	}
	var wg sync.WaitGroup
	for k := 0; k < 32; k++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				j := i % len(in)
				if got := n.Normalize(in[j]); got != expected[j] {
					t.Errorf("inconsistent normalization")
					return
				}
			}
		}()
	}
	wg.Wait()
}
