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

// newCountry is a helper to create a country with the new structure
func newCountry(iso2, iso3, name string) domain.Country {
	return domain.Country{
		ISO2:    iso2,
		ISO3:    iso3,
		Names:   map[string]string{"en": name},
		Aliases: []string{},
	}
}

// getCountryData returns hardcoded country data
func getCountryData() []domain.Country {
	return []domain.Country{
		newCountry("AF", "AFG", "Afghanistan"), newCountry("AL", "ALB", "Albania"),
		newCountry("DZ", "DZA", "Algeria"), newCountry("AD", "AND", "Andorra"),
		newCountry("AO", "AGO", "Angola"), newCountry("AG", "ATG", "Antigua and Barbuda"),
		newCountry("AR", "ARG", "Argentina"), newCountry("AM", "ARM", "Armenia"),
		newCountry("AU", "AUS", "Australia"), newCountry("AT", "AUT", "Austria"),
		newCountry("AZ", "AZE", "Azerbaijan"), newCountry("BS", "BHS", "Bahamas"),
		newCountry("BH", "BHR", "Bahrain"), newCountry("BD", "BGD", "Bangladesh"),
		newCountry("BB", "BRB", "Barbados"), newCountry("BY", "BLR", "Belarus"),
		newCountry("BE", "BEL", "Belgium"), newCountry("BZ", "BLZ", "Belize"),
		newCountry("BJ", "BEN", "Benin"), newCountry("BT", "BTN", "Bhutan"),
		newCountry("BO", "BOL", "Bolivia (Plurinational State of)"),
		newCountry("BA", "BIH", "Bosnia and Herzegovina"), newCountry("BW", "BWA", "Botswana"),
		newCountry("BR", "BRA", "Brazil"), newCountry("BN", "BRN", "Brunei Darussalam"),
		newCountry("BG", "BGR", "Bulgaria"), newCountry("BF", "BFA", "Burkina Faso"),
		newCountry("BI", "BDI", "Burundi"), newCountry("CV", "CPV", "Cabo Verde"),
		newCountry("KH", "KHM", "Cambodia"), newCountry("CM", "CMR", "Cameroon"),
		newCountry("CA", "CAN", "Canada"), newCountry("CF", "CAF", "Central African Republic"),
		newCountry("TD", "TCD", "Chad"), newCountry("CL", "CHL", "Chile"),
		newCountry("CN", "CHN", "China"), newCountry("CO", "COL", "Colombia"),
		newCountry("KM", "COM", "Comoros"), newCountry("CG", "COG", "Congo"),
		newCountry("CD", "COD", "Congo, Democratic Republic of the"),
		newCountry("CR", "CRI", "Costa Rica"), newCountry("CI", "CIV", "Côte d'Ivoire"),
		newCountry("HR", "HRV", "Croatia"), newCountry("CU", "CUB", "Cuba"),
		newCountry("CY", "CYP", "Cyprus"), newCountry("CZ", "CZE", "Czechia"),
		newCountry("DK", "DNK", "Denmark"), newCountry("DJ", "DJI", "Djibouti"),
		newCountry("DM", "DMA", "Dominica"), newCountry("DO", "DOM", "Dominican Republic"),
		newCountry("EC", "ECU", "Ecuador"), newCountry("EG", "EGY", "Egypt"),
		newCountry("SV", "SLV", "El Salvador"), newCountry("GQ", "GNQ", "Equatorial Guinea"),
		newCountry("ER", "ERI", "Eritrea"), newCountry("EE", "EST", "Estonia"),
		newCountry("SZ", "SWZ", "Eswatini"), newCountry("ET", "ETH", "Ethiopia"),
		newCountry("FI", "FIN", "Finland"), newCountry("FR", "FRA", "France"),
		newCountry("GA", "GAB", "Gabon"), newCountry("GM", "GMB", "Gambia"),
		newCountry("GE", "GEO", "Georgia"), newCountry("DE", "DEU", "Germany"),
		newCountry("GH", "GHA", "Ghana"), newCountry("GR", "GRC", "Greece"),
		newCountry("GD", "GRD", "Grenada"), newCountry("GT", "GTM", "Guatemala"),
		newCountry("GN", "GIN", "Guinea"), newCountry("GW", "GNB", "Guinea-Bissau"),
		newCountry("GY", "GUY", "Guyana"), newCountry("HT", "HTI", "Haiti"),
		newCountry("HN", "HND", "Honduras"), newCountry("HU", "HUN", "Hungary"),
		newCountry("IS", "ISL", "Iceland"), newCountry("IN", "IND", "India"),
		newCountry("ID", "IDN", "Indonesia"), newCountry("IR", "IRN", "Iran (Islamic Republic of)"),
		newCountry("IQ", "IRQ", "Iraq"), newCountry("IE", "IRL", "Ireland"),
		newCountry("IL", "ISR", "Israel"), newCountry("IT", "ITA", "Italy"),
		newCountry("JM", "JAM", "Jamaica"), newCountry("JP", "JPN", "Japan"),
		newCountry("JO", "JOR", "Jordan"), newCountry("KZ", "KAZ", "Kazakhstan"),
		newCountry("KE", "KEN", "Kenya"), newCountry("KW", "KWT", "Kuwait"),
		newCountry("KG", "KGZ", "Kyrgyzstan"), newCountry("LA", "LAO", "Lao People's Democratic Republic"),
		newCountry("LV", "LVA", "Latvia"), newCountry("LB", "LBN", "Lebanon"),
		newCountry("LS", "LSO", "Lesotho"), newCountry("LR", "LBR", "Liberia"),
		newCountry("LY", "LBY", "Libya"), newCountry("LI", "LIE", "Liechtenstein"),
		newCountry("LT", "LTU", "Lithuania"), newCountry("LU", "LUX", "Luxembourg"),
		newCountry("MG", "MDG", "Madagascar"), newCountry("MW", "MWI", "Malawi"),
		newCountry("MY", "MYS", "Malaysia"), newCountry("MV", "MDV", "Maldives"),
		newCountry("ML", "MLI", "Mali"), newCountry("MT", "MLT", "Malta"),
		newCountry("MR", "MRT", "Mauritania"), newCountry("MU", "MUS", "Mauritius"),
		newCountry("MX", "MEX", "Mexico"), newCountry("MD", "MDA", "Moldova, Republic of"),
		newCountry("MC", "MCO", "Monaco"), newCountry("MN", "MNG", "Mongolia"),
		newCountry("ME", "MNE", "Montenegro"), newCountry("MA", "MAR", "Morocco"),
		newCountry("MZ", "MOZ", "Mozambique"), newCountry("MM", "MMR", "Myanmar"),
		newCountry("NA", "NAM", "Namibia"), newCountry("NP", "NPL", "Nepal"),
		newCountry("NL", "NLD", "Netherlands"), newCountry("NZ", "NZL", "New Zealand"),
		newCountry("NI", "NIC", "Nicaragua"), newCountry("NE", "NER", "Niger"),
		newCountry("NG", "NGA", "Nigeria"), newCountry("KP", "PRK", "North Korea"),
		newCountry("MK", "MKD", "North Macedonia"), newCountry("NO", "NOR", "Norway"),
		newCountry("OM", "OMN", "Oman"), newCountry("PK", "PAK", "Pakistan"),
		newCountry("PS", "PSE", "Palestine, State of"), newCountry("PA", "PAN", "Panama"),
		newCountry("PY", "PRY", "Paraguay"), newCountry("PE", "PER", "Peru"),
		newCountry("PH", "PHL", "Philippines"), newCountry("PL", "POL", "Poland"),
		newCountry("PT", "PRT", "Portugal"), newCountry("PR", "PRI", "Puerto Rico"),
		newCountry("QA", "QAT", "Qatar"), newCountry("RO", "ROU", "Romania"),
		newCountry("RU", "RUS", "Russian Federation"), newCountry("RW", "RWA", "Rwanda"),
		newCountry("KN", "KNA", "Saint Kitts and Nevis"), newCountry("LC", "LCA", "Saint Lucia"),
		newCountry("VC", "VCT", "Saint Vincent and the Grenadines"),
		newCountry("SM", "SMR", "San Marino"), newCountry("ST", "STP", "Sao Tome and Principe"),
		newCountry("SA", "SAU", "Saudi Arabia"),
		newCountry("SN", "SEN", "Senegal"), newCountry("RS", "SRB", "Serbia"),
		newCountry("SC", "SYC", "Seychelles"), newCountry("SL", "SLE", "Sierra Leone"),
		newCountry("SG", "SGP", "Singapore"), newCountry("SK", "SVK", "Slovakia"),
		newCountry("SI", "SVN", "Slovenia"), newCountry("SO", "SOM", "Somalia"),
		newCountry("ZA", "ZAF", "South Africa"), newCountry("KR", "KOR", "South Korea"),
		newCountry("SS", "SSD", "South Sudan"), newCountry("ES", "ESP", "Spain"),
		newCountry("LK", "LKA", "Sri Lanka"), newCountry("SD", "SDN", "Sudan"),
		newCountry("SE", "SWE", "Sweden"), newCountry("CH", "CHE", "Switzerland"),
		newCountry("SY", "SYR", "Syrian Arab Republic"), newCountry("TW", "TWN", "Taiwan, Province of China"),
		newCountry("TJ", "TJK", "Tajikistan"),
		newCountry("TZ", "TZA", "Tanzania, United Republic of"), newCountry("TH", "THA", "Thailand"),
		newCountry("TG", "TGO", "Togo"),
		newCountry("TT", "TTO", "Trinidad and Tobago"), newCountry("TN", "TUN", "Tunisia"),
		newCountry("TR", "TUR", "Turkey"), newCountry("TM", "TKM", "Turkmenistan"),
		newCountry("UG", "UGA", "Uganda"), newCountry("UA", "UKR", "Ukraine"),
		newCountry("AE", "ARE", "United Arab Emirates"),
		newCountry("GB", "GBR", "United Kingdom of Great Britain and Northern Ireland"),
		newCountry("US", "USA", "United States of America"),
		newCountry("UY", "URY", "Uruguay"), newCountry("UZ", "UZB", "Uzbekistan"),
		newCountry("VE", "VEN", "Venezuela (Bolivarian Republic of)"),
		newCountry("VN", "VNM", "Viet Nam"), newCountry("EH", "ESH", "Western Sahara"),
		newCountry("YE", "YEM", "Yemen"), newCountry("ZM", "ZMB", "Zambia"),
		newCountry("ZW", "ZWE", "Zimbabwe"),
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
