package countries

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	normalizedCountryNames map[string]string
	isoToOfficialName      map[string]string
	accentRemover          = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
)

func init() {
	normalizedCountryNames = make(map[string]string)
	isoToOfficialName = make(map[string]string)

	// Build lookups from official list
	for _, country := range List {
		isoToOfficialName[country.Code] = country.Name

		normalized, _, _ := transform.String(accentRemover, country.Name)
		normalizedCountryNames[strings.ToLower(normalized)] = country.Code
		normalizedCountryNames[strings.ToLower(country.Name)] = country.Code
	}

	// Add aliases
	for isoCode, aliases := range Aliases {
		for _, alias := range aliases {
			normalized, _, _ := transform.String(accentRemover, alias)
			normalizedCountryNames[strings.ToLower(normalized)] = isoCode
		}
	}
}

// NormalizeString removes accents and lowercases
func NormalizeString(s string) string {
	normalized, _, _ := transform.String(accentRemover, s)
	return strings.ToLower(normalized)
}

// LookupISO returns iso code + official name given a query
func LookupISO(query string) (isoCode, officialName string, found bool) {
	normalized := NormalizeString(query)
	code, ok := normalizedCountryNames[normalized]
	if !ok {
		return "", "", false
	}
	return code, isoToOfficialName[code], true
}
