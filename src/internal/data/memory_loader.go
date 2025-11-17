package data

import "country-iso-matcher/src/internal/domain"

// MemoryLoader loads country data from memory (hardcoded data)
// This is useful for testing or when external files are not available
type MemoryLoader struct{}

// NewMemoryLoader creates a new memory loader
func NewMemoryLoader() *MemoryLoader {
	return &MemoryLoader{}
}

// LoadCountries returns the hardcoded list of countries
func (l *MemoryLoader) LoadCountries() ([]domain.Country, error) {
	return getCountryData(), nil
}

// LoadAliases returns the hardcoded list of country aliases
func (l *MemoryLoader) LoadAliases() (map[string][]string, error) {
	return getAliasData(), nil
}

// getCountryData returns hardcoded country data
func getCountryData() []domain.Country {
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

// getAliasData returns hardcoded alias data
// This function contains extensive aliases for better country name matching
func getAliasData() map[string][]string {
	return map[string][]string{
		"US": {
			"usa", "united states", "america", "états-unis", "vereinigte staaten",
			"'merica", "murica", "united statas", "united stetes", "united staes",
			"united stets", "united staates", "untied states", "estados unidos",
			"amérique", "us", "u.s.a.", "u.s.a",
		},
		"GB": {
			"uk", "united kingdom", "britain", "england", "royaume-uni",
			"vereinigtes königreich", "great britain", "u.k.", "u.k",
			"northern ireland", "scotland", "wales", "écosse", "angleterre",
		},
		"BR": {
			"brasil", "brazil", "brasil", "brasilz", "braszil", "brazyl",
		},
		"CN": {
			"china", "chine", "chaina", "chyna", "chinia", "chinna", "chinah",
			"mainland china", "prc", "people's republic of china",
		},
		"KR": {
			"south korea", "republic of korea", "korea south", "corée du sud",
			"südkorea", "soth korea", "south koria", "hanguk", "rok", "sk", "korea",
		},
		"RU": {
			"russia", "russie", "russland", "russian federation", "rossiya",
			"rossia", "rusia", "rusija", "russa", "russha", "soviet union",
			"ussr", "union of soviet socialist republics", "sovjet",
		},
		"DE": {
			"germany", "deutschland", "allemagne", "germania", "alemania",
			"deutchland", "deutchlnd", "deutcheland",
		},
		"FR": {
			"france", "frankreich", "francia", "francais", "franse", "franc", "francz",
		},
		"IT": {
			"italy", "italia", "italie", "italien", "it", "itly", "itali", "italya",
		},
		"ES": {
			"spain", "españa", "espagne", "spanien", "spagna", "espanya",
			"espania", "spane", "spian",
		},
		"PL": {
			"poland", "polska", "pologne", "polen", "polland", "polan",
			"p0land", "polend", "polad",
		},
		"NL": {
			"netherlands", "nederland", "pays-bas", "niederlande", "holland",
			"the netherlands", "países bajos", "netherland", "neterlands",
		},
		"CH": {
			"switzerland", "suisse", "schweiz", "svizzera", "helvetia",
			"swizerland", "switserland", "ch",
		},
		"NO": {
			"norway", "norge", "noreg", "norwey", "norvay",
		},
		"SE": {
			"sweden", "sverige", "suède", "schweden", "swden", "sweeden",
		},
		"GR": {
			"greece", "ελλάδα", "grèce", "griechenland", "grecia",
			"hellas", "ellada", "grece", "greese",
		},
		"RO": {
			"romania", "românia", "roumanie", "roumania", "roumanía",
			"roumaniya", "rouminia", "rumänien", "rumænien", "rumunia",
			"rumunija", "rumänia", "rumānija", "rumānīyā", "rumunska",
			"rumunsko", "rumunjska", "rumyniya", "rumuniya", "rumuniia",
			"rumenia", "rumenía", "rúmenía", "rumeenia", "rumínia",
			"rumanio", "rumania", "rumanía", "romanya", "roménia",
			"romênia", "roménya", "roemenië",
		},
		"HU": {
			"hungary", "hunggary", "hungry", "magyarország", "magyar-kok",
			"hongrie", "hongaria", "hongarije", "hongarye", "hongria",
			"hungaria", "hungario", "hungari", "hungariá", "hungarya",
			"hungariya", "hungriya", "hungyria", "hungria", "hungría",
			"ungarn", "ungern", "ungari", "ungaría", "ungaria", "ungārija",
			"ungarija", "ungariya", "ungheria", "węgry",
		},
		"CZ": {
			"czech republic", "česká republika", "république tchèque",
			"tschechische republik", "czechia", "czechoslovakia", "checz",
			"chekia", "česko",
		},
		"SK": {
			"slovakia", "slovensko", "slovaquie", "slowakei", "slovak republic",
		},
		"JP": {
			"japan", "nippon", "nihon", "japon", "jappan", "japn",
		},
		"IN": {
			"india", "bharat", "hindustan", "indea", "inida",
		},
		"AF": {
			"afghanistan", "afganistan", "afganisthan", "afgahnistan", "aghanistan",
		},
		"TR": {
			"turkey", "türkiye", "turkiye", "turky", "turkie",
		},
		"PH": {
			"philippines", "philipines", "philipinnes", "phillippines",
			"the phillipines", "pilipinas", "pinoyland",
		},
		"IR": {
			"iran", "persia", "iraan", "irun", "irá", "irák",
		},
		"IQ": {
			"iraq", "irak", "irac", "irack",
		},
		"ZA": {
			"south africa", "south afrika", "soth africa", "azania", "sa", "za",
		},
		"AU": {
			"australia", "aussie", "oz", "straya", "down under", "austraila", "austalia",
		},
		"CA": {
			"canada", "canda", "cannada", "canadia", "the great white north",
		},
		"NZ": {
			"new zealand", "new zeeland", "new zeland", "nz", "kiwiland", "aotearoa",
		},
		"MM": {
			"myanmar", "burma", "mianmar", "myanmer", "mayanmar", "myannmar",
		},
		"MX": {
			"mexico", "méxico", "méjico", "mexcio", "mecsico",
		},
		"EG": {
			"egypt", "misr", "egpyt", "egyt", "eygpt",
		},
		"NG": {
			"nigeria", "nieria", "nigeira", "naija",
		},
		"AR": {
			"argentina", "argetina", "argentia",
		},
		"CO": {
			"colombia", "columbia", "colobia", "collombia",
		},
		"TH": {
			"thailand", "siam", "tailand", "thialand", "thiland",
		},
		"VN": {
			"viet nam", "vietnam", "veitnam", "vietnem",
		},
		"ID": {
			"indonesia", "indonisia", "indonezia",
		},
		"PT": {
			"portugal", "portugol", "portugual", "portgual",
		},
		"BE": {
			"belgium", "belgique", "belgië", "belgien", "belguim", "beljum",
		},
		"CI": {
			"côte d'ivoire", "cote d'ivoire", "ivory coast", "ivore coast",
			"ivri coast", "ivorycoast",
		},
		"AT": {"austria", "österreich", "autriche"},
		"BG": {"bulgaria", "българия", "bulgarie", "bulgarien"},
		"HR": {"croatia", "hrvatska", "croatie", "kroatien"},
		"CY": {"cyprus", "κύπρος", "chypre", "zypern"},
		"DK": {"denmark", "danmark", "danemark", "dänemark"},
		"EE": {"estonia", "eesti", "estonie", "estland"},
		"FI": {"finland", "suomi", "finlande", "finnland"},
		"LV": {"latvia", "latvija", "lettonie", "lettland"},
		"LT": {"lithuania", "lietuva", "lituanie", "litauen"},
		"LU": {"luxembourg", "luxemburg"},
		"MT": {"malta", "malte"},
		"SI": {"slovenia", "slovenija", "slovénie", "slowenien"},
		"BA": {"bosnia", "bosnia and herzegovina"},
		"MK": {"north macedonia", "macedonia"},
		"MD": {"moldova", "republic of moldova"},
		"IE": {"ireland", "éire", "irlande", "irland", "éireann"},
	}
}
