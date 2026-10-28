package normalizer

import (
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

type TextNormalizer interface {
	Normalize(text string) string
}

type textNormalizer struct {
	transformers sync.Pool
}

func NewTextNormalizer() TextNormalizer {
	return &textNormalizer{
		transformers: sync.Pool{New: func() any {
			return transform.Chain(width.Fold, norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
		}},
	}
}

func (n *textNormalizer) Normalize(text string) string {
	// Transformers have mutable internal buffers and are not safe to share
	// between concurrent HTTP requests. Each call owns a chain until it returns.
	transformer := n.transformers.Get().(transform.Transformer)
	defer n.transformers.Put(transformer)
	normalized, _, _ := transform.String(transformer, text)
	return strings.ToLower(strings.TrimSpace(normalized))
}
