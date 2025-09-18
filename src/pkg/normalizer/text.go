package normalizer

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type TextNormalizer interface {
	Normalize(text string) string
}

type textNormalizer struct {
	transformer transform.Transformer
}

func NewTextNormalizer() TextNormalizer {
	return &textNormalizer{
		transformer: transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC),
	}
}

func (n *textNormalizer) Normalize(text string) string {
	normalized, _, _ := transform.String(n.transformer, text)
	return strings.ToLower(strings.TrimSpace(normalized))
}
