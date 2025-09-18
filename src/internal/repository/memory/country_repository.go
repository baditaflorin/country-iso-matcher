package memory

import (
	"country-iso-matcher/src/internal/domain"
	"country-iso-matcher/src/pkg/normalizer"
)

type countryRepository struct {
	nameToCode    map[string]string
	codeToCountry map[string]*domain.Country
	normalizer    normalizer.TextNormalizer
}

func NewCountryRepository(normalizer normalizer.TextNormalizer) *countryRepository {
	repo := &countryRepository{
		nameToCode:    make(map[string]string),
		codeToCountry: make(map[string]*domain.Country),
		normalizer:    normalizer,
	}
	repo.loadCountries()
	return repo
}

func (r *countryRepository) FindByName(name string) (*domain.Country, error) {
	normalized := r.normalizer.Normalize(name)
	code, exists := r.nameToCode[normalized]
	if !exists {
		return nil, domain.NewNotFoundError(name)
	}
	return r.codeToCountry[code], nil
}

func (r *countryRepository) FindByCode(code string) (*domain.Country, error) {
	country, exists := r.codeToCountry[code]
	if !exists {
		return nil, domain.NewNotFoundError(code)
	}
	return country, nil
}

func (r *countryRepository) loadCountries() {
	countries := r.getCountryData()
	aliases := r.getAliasData()

	// Load official countries
	for _, country := range countries {
		r.codeToCountry[country.Code] = &country
		normalized := r.normalizer.Normalize(country.Name)
		r.nameToCode[normalized] = country.Code
	}

	// Load aliases
	for isoCode, aliasNames := range aliases {
		for _, alias := range aliasNames {
			normalized := r.normalizer.Normalize(alias)
			r.nameToCode[normalized] = isoCode
		}
	}
}

func (r *countryRepository) getCountryData() []domain.Country {
	return []domain.Country{
		{"AF", "Afghanistan"}, {"AL", "Albania"}, {"DZ", "Algeria"}, {"AD", "Andorra"},
		{"AO", "Angola"}, {"AG", "Antigua and Barbuda"}, {"AR", "Argentina"}, {"AM", "Armenia"},
		{"AU", "Australia"}, {"AT", "Austria"}, {"AZ", "Azerbaijan"}, {"BS", "Bahamas"},
		{"BH", "Bahrain"}, {"BD", "Bangladesh"}, {"BB", "Barbados"}, {"BY", "Belarus"},
		{"BE", "Belgium"}, {"BZ", "Belize"}, {"BJ", "Benin"}, {"BT", "Bhutan"},
		{"BO", "Bolivia (Plurinational State of)"}, {"BA", "Bosnia and Herzegovina"}, {"BW", "Botswana"},
		{"BR", "Brazil"}, {"BN", "Brunei Darussalam"}, {"BG", "Bulgaria"}, {"BF", "Burkina Faso"},
		{"BI", "Burundi"}, {"CV", "Cabo Verde"}, {"KH", "Cambodia"}, {"CM", "Cameroon"},
		{"CA", "Canada"}, {"CF", "Central African Republic"}, {"TD", "Chad"}, {"CL", "Chile"},
		{"CN", "China"}, {"CO", "Colombia"}, {"KM", "Comoros"}, {"CG", "Congo"},
		{"CD", "Congo, Democratic Republic of the"}, {"CR", "Costa Rica"}, {"CI", "Côte d'Ivoire"},
		{"HR", "Croatia"}, {"CU", "Cuba"}, {"CY", "Cyprus"}, {"CZ", "Czechia"},
		{"DK", "Denmark"}, {"DJ", "Djibouti"}, {"DM", "Dominica"}, {"DO", "Dominican Republic"},
		{"EC", "Ecuador"}, {"EG", "Egypt"}, {"SV", "El Salvador"}, {"GQ", "Equatorial Guinea"},
		{"ER", "Eritrea"}, {"EE", "Estonia"}, {"SZ", "Eswatini"}, {"ET", "Ethiopia"},
		{"FI", "Finland"}, {"FR", "France"}, {"GA", "Gabon"}, {"GM", "Gambia"},
		{"GE", "Georgia"}, {"DE", "Germany"}, {"GH", "Ghana"}, {"GR", "Greece"},
		{"GD", "Grenada"}, {"GT", "Guatemala"}, {"GN", "Guinea"}, {"GW", "Guinea-Bissau"},
		{"GY", "Guyana"}, {"HT", "Haiti"}, {"HN", "Honduras"}, {"HU", "Hungary"},
		{"IS", "Iceland"}, {"IN", "India"}, {"ID", "Indonesia"}, {"IR", "Iran (Islamic Republic of)"},
		{"IQ", "Iraq"}, {"IE", "Ireland"}, {"IL", "Israel"}, {"IT", "Italy"},
		{"JM", "Jamaica"}, {"JP", "Japan"}, {"JO", "Jordan"}, {"KZ", "Kazakhstan"},
		{"KE", "Kenya"}, {"KW", "Kuwait"}, {"KG", "Kyrgyzstan"}, {"LA", "Lao People's Democratic Republic"},
		{"LV", "Latvia"}, {"LB", "Lebanon"}, {"LS", "Lesotho"}, {"LR", "Liberia"},
		{"LY", "Libya"}, {"LI", "Liechtenstein"}, {"LT", "Lithuania"}, {"LU", "Luxembourg"},
		{"MG", "Madagascar"}, {"MW", "Malawi"}, {"MY", "Malaysia"}, {"MV", "Maldives"},
		{"ML", "Mali"}, {"MT", "Malta"}, {"MR", "Mauritania"}, {"MU", "Mauritius"},
		{"MX", "Mexico"}, {"MD", "Moldova, Republic of"}, {"MC", "Monaco"}, {"MN", "Mongolia"},
		{"ME", "Montenegro"}, {"MA", "Morocco"}, {"MZ", "Mozambique"}, {"MM", "Myanmar"},
		{"NA", "Namibia"}, {"NP", "Nepal"}, {"NL", "Netherlands"}, {"NZ", "New Zealand"},
		{"NI", "Nicaragua"}, {"NE", "Niger"}, {"NG", "Nigeria"}, {"KP", "North Korea"},
		{"MK", "North Macedonia"}, {"NO", "Norway"}, {"OM", "Oman"}, {"PK", "Pakistan"},
		{"PS", "Palestine, State of"}, {"PA", "Panama"}, {"PY", "Paraguay"}, {"PE", "Peru"},
		{"PH", "Philippines"}, {"PL", "Poland"}, {"PT", "Portugal"}, {"PR", "Puerto Rico"},
		{"QA", "Qatar"}, {"RO", "Romania"}, {"RU", "Russian Federation"}, {"RW", "Rwanda"},
		{"KN", "Saint Kitts and Nevis"}, {"LC", "Saint Lucia"}, {"VC", "Saint Vincent and the Grenadines"},
		{"SM", "San Marino"}, {"ST", "Sao Tome and Principe"}, {"SA", "Saudi Arabia"},
		{"SN", "Senegal"}, {"RS", "Serbia"}, {"SC", "Seychelles"}, {"SL", "Sierra Leone"},
		{"SG", "Singapore"}, {"SK", "Slovakia"}, {"SI", "Slovenia"}, {"SO", "Somalia"},
		{"ZA", "South Africa"}, {"KR", "South Korea"}, {"SS", "South Sudan"}, {"ES", "Spain"},
		{"LK", "Sri Lanka"}, {"SD", "Sudan"}, {"SE", "Sweden"}, {"CH", "Switzerland"},
		{"SY", "Syrian Arab Republic"}, {"TW", "Taiwan, Province of China"}, {"TJ", "Tajikistan"},
		{"TZ", "Tanzania, United Republic of"}, {"TH", "Thailand"}, {"TG", "Togo"},
		{"TT", "Trinidad and Tobago"}, {"TN", "Tunisia"}, {"TR", "Turkey"}, {"TM", "Turkmenistan"},
		{"UG", "Uganda"}, {"UA", "Ukraine"}, {"AE", "United Arab Emirates"},
		{"GB", "United Kingdom of Great Britain and Northern Ireland"}, {"US", "United States of America"},
		{"UY", "Uruguay"}, {"UZ", "Uzbekistan"}, {"VE", "Venezuela (Bolivarian Republic of)"},
		{"VN", "Viet Nam"}, {"EH", "Western Sahara"}, {"YE", "Yemen"}, {"ZM", "Zambia"},
		{"ZW", "Zimbabwe"},
	}
}

func (r *countryRepository) getAliasData() map[string][]string {
	return map[string][]string{
		"KR": {"south korea", "republic of korea", "korea south", "corée du sud", "südkorea"},
		"US": {"usa", "united states", "america", "états-unis", "vereinigte staaten"},
		"GB": {"uk", "united kingdom", "britain", "england", "royaume-uni", "vereinigtes königreich"},
		"RU": {"russia", "russie", "russland"},
		"AT": {"austria", "österreich", "autriche"},
		"BE": {"belgium", "belgique", "belgië", "belgien"},
		"BG": {"bulgaria", "българия", "bulgarie", "bulgarien"},
		"HR": {"croatia", "hrvatska", "croatie", "kroatien"},
		"CY": {"cyprus", "κύπρος", "chypre", "zypern"},
		"CZ": {"czech republic", "česká republika", "république tchèque", "tschechische republik"},
		"DK": {"denmark", "danmark", "danemark", "dänemark"},
		"EE": {"estonia", "eesti", "estonie", "estland"},
		"FI": {"finland", "suomi", "finlande", "finnland"},
		"FR": {"france", "frankreich"},
		"DE": {"germany", "deutschland", "allemagne"},
		"GR": {"greece", "ελλάδα", "grèce", "griechenland"},
		"HU": {"hungary", "magyarország", "hongrie", "ungarn"},
		"IE": {"ireland", "éire", "irlande", "irland"},
		"IT": {"italy", "italia", "italie", "italien"},
		"LV": {"latvia", "latvija", "lettonie", "lettland"},
		"LT": {"lithuania", "lietuva", "lituanie", "litauen"},
		"LU": {"luxembourg", "luxemburg"},
		"MT": {"malta", "malte"},
		"NL": {"netherlands", "nederland", "pays-bas", "niederlande"},
		"PL": {"poland", "polska", "pologne", "polen"},
		"PT": {"portugal"},
		"RO": {"romania", "românia", "roumanie", "rumänien"},
		"SK": {"slovakia", "slovensko", "slovaquie", "slowakei"},
		"SI": {"slovenia", "slovenija", "slovénie", "slowenien"},
		"ES": {"spain", "españa", "espagne", "spanien"},
		"SE": {"sweden", "sverige", "suède", "schweden"},
		"TR": {"turkey", "türkiye"},
		"BA": {"bosnia", "bosnia and herzegovina"},
		"CH": {"switzerland", "suisse", "schweiz"},
		"MK": {"north macedonia", "macedonia"},
		"MD": {"moldova", "republic of moldova"},
	}
}
