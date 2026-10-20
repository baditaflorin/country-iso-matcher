package memory

import (
	"testing"

	"country-iso-matcher/src/internal/data"
	"country-iso-matcher/src/pkg/normalizer"
)

func TestCommaCountryNamesAndCongoDisambiguation(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ query, want string }{
		{"Congo", "CG"},
		{"Congo, Democratic Republic of the", "CD"},
		{"Moldova, Republic of", "MD"},
		{"Palestine, State of", "PS"},
		{"Taiwan, Province of China", "TW"},
		{"Tanzania, United Republic of", "TZ"},
	} {
		got, err := repo.FindByName(tc.query)
		if err != nil {
			t.Errorf("lookup %q: %v, want %s", tc.query, err, tc.want)
			continue
		}
		if got.ISO2 != tc.want {
			t.Errorf("lookup %q = %s, want %s", tc.query, got.ISO2, tc.want)
		}
	}
}

func TestCLDRLocalizedTerritoryNamesResolveAcrossScripts(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}

	// Representative discovery and held-out labels from Unicode CLDR 48.2.3.
	for _, tc := range []struct{ locale, label, want string }{
		{"de", "Vereinigte Arabische Emirate", "AE"},
		{"es", "Antigua y Barbuda", "AG"},
		{"fr", "Émirats arabes unis", "AE"},
		{"pt", "Emirados Árabes Unidos", "AE"},
		{"tr", "Arnavutluk", "AL"},
		{"ar", "أندورا", "AD"},
		{"hi", "एंडोरा", "AD"},
		{"ja", "アンドラ", "AD"},
		{"ru", "Андорра", "AD"},
		{"uk", "Обʼєднані Арабські Емірати", "AE"},
		{"zh", "安道尔", "AD"},
		{"it", "Emirati Arabi Uniti", "AE"},
		{"nl", "Albanië", "AL"},
	} {
		got, err := repo.FindByName(tc.label)
		if err != nil {
			t.Errorf("CLDR %s label %q: %v, want %s", tc.locale, tc.label, err, tc.want)
			continue
		}
		if got.ISO2 != tc.want {
			t.Errorf("CLDR %s label %q = %s, want %s", tc.locale, tc.label, got.ISO2, tc.want)
		}
	}
}

func TestPortugueseGuianaResolvesToGuyana(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}
	got, err := repo.FindByName("Guiana")
	if err != nil {
		t.Fatal(err)
	}
	if got.ISO2 != "GY" {
		t.Fatalf("Portuguese CLDR label Guiana resolved to %s, want GY", got.ISO2)
	}
}

func TestJapaneseCountryAliasesResolveAcrossWidthVariants(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}

	for _, query := range []string{"ＵＳＡ", "ｱﾒﾘｶ合衆国"} {
		got, err := repo.FindByName(query)
		if err != nil {
			t.Errorf("lookup width variant %q: %v, want US", query, err)
			continue
		}
		if got.ISO2 != "US" {
			t.Errorf("lookup width variant %q = %s, want US", query, got.ISO2)
		}
	}
	if _, err := repo.FindByName("ｱﾒﾘｶ"); err == nil {
		t.Fatal("width folding should not turn a partial label into a country alias")
	}
}

func TestUNM49CountryLabelsStayDistinctAcrossLanguages(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}

	// UN M49 country/area labels found missing from the bundled catalog.
	// Paired names ensure that shared Congo/Korea roots do not cross-map.
	for _, tc := range []struct{ language, label, want string }{
		{"en", "Democratic Republic of the Congo", "CD"},
		{"en", "Congo", "CG"},
		{"en", "Netherlands (Kingdom of the)", "NL"},
		{"en", "Democratic People's Republic of Korea", "KP"},
		{"en", "Republic of Korea", "KR"},
		{"zh", "刚果民主共和国", "CD"},
		{"zh", "刚果", "CG"},
		{"zh", "大韩民国", "KR"},
		{"zh", "朝鲜民主主义人民共和国", "KP"},
		{"ru", "Демократическая Республика Конго", "CD"},
		{"ru", "Конго", "CG"},
		{"ru", "Республика Корея", "KR"},
		{"ru", "Корейская Народно-Демократическая Республика", "KP"},
		{"ru", "Южная Африка", "ZA"},
		{"ru", "Южная Джорджия и Южные Сандвичевы острова", "GS"},
		{"fr", "République démocratique du Congo", "CD"},
		{"fr", "Congo", "CG"},
		{"fr", "République de Corée", "KR"},
		{"fr", "République populaire démocratique de Corée", "KP"},
		{"es", "República Democrática del Congo", "CD"},
		{"es", "Congo", "CG"},
		{"es", "República de Corea", "KR"},
		{"es", "República Popular Democrática de Corea", "KP"},
		{"ar", "جمهورية الكونغو الديمقراطية", "CD"},
		{"ar", "الكونغو", "CG"},
		{"ar", "جمهورية كوريا", "KR"},
		{"ar", "جمهورية كوريا الشعبية الديمقراطية", "KP"},
	} {
		got, err := repo.FindByName(tc.label)
		if err != nil {
			t.Errorf("UN M49 %s label %q: %v, want %s", tc.language, tc.label, err, tc.want)
			continue
		}
		if got.ISO2 != tc.want {
			t.Errorf("UN M49 %s label %q = %s, want %s", tc.language, tc.label, got.ISO2, tc.want)
		}
	}
}

func TestArabicM49LabelsUseValidNomenclature(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}

	// The current M49 Arabic table contains malformed strings for NL, IO, TF,
	// and UM. Keep the corrected/established Arabic labels resolvable and do not
	// introduce the malformed source forms as aliases.
	for _, tc := range []struct{ label, want string }{
		{"مملكة هولندا", "NL"},
		{"الإقليم البريطاني في المحيط الهندي", "IO"},
		{"الأقاليم الجنوبية الفرنسية", "TF"},
		{"جزر الولايات المتحدة النائية", "UM"},
	} {
		got, err := repo.FindByName(tc.label)
		if err != nil {
			t.Errorf("valid Arabic label %q: %v, want %s", tc.label, err, tc.want)
			continue
		}
		if got.ISO2 != tc.want {
			t.Errorf("valid Arabic label %q = %s, want %s", tc.label, got.ISO2, tc.want)
		}
	}

	for _, label := range []string{
		"هولندا (مملكة _)",
		"المحيط الهندي الإقليم البريطاني في",
		"الأراضي الفرنسية الجنوبية الجنوبية",
		"نائية التابعة للولايات المتحدة",
	} {
		if got, err := repo.FindByName(label); err == nil {
			t.Errorf("malformed Arabic M49 label %q unexpectedly resolved to %s", label, got.ISO2)
		}
	}
}

func TestEveryAliasResolvesWithoutCrossCountryAliasCollisions(t *testing.T) {
	textNormalizer := normalizer.NewTextNormalizer()
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	countries, err := loader.LoadCountries()
	if err != nil {
		t.Fatal(err)
	}
	aliases, err := loader.LoadAliases()
	if err != nil {
		t.Fatal(err)
	}
	repo, err := NewCountryRepository(textNormalizer, loader)
	if err != nil {
		t.Fatal(err)
	}

	canonicalOwners := make(map[string]string)
	for _, country := range countries {
		for _, name := range country.Names {
			canonicalOwners[textNormalizer.Normalize(name)] = country.ISO2
		}
		canonicalOwners[textNormalizer.Normalize(country.ISO2)] = country.ISO2
		canonicalOwners[textNormalizer.Normalize(country.ISO3)] = country.ISO2
	}

	type aliasExample struct {
		code  string
		label string
	}
	owners := make(map[string]map[string]struct{})
	examples := make(map[string]aliasExample)
	for code, labels := range aliases {
		for _, label := range labels {
			key := textNormalizer.Normalize(label)
			if owners[key] == nil {
				owners[key] = make(map[string]struct{})
			}
			owners[key][code] = struct{}{}
			examples[key] = aliasExample{code: code, label: label}
		}
	}

	for key, aliasOwners := range owners {
		canonicalOwner := canonicalOwners[key]
		if len(aliasOwners) > 1 && canonicalOwner == "" {
			t.Errorf("normalized alias %q is assigned to multiple countries without a canonical name/code", examples[key].label)
			continue
		}

		want := canonicalOwner
		if want == "" {
			for code := range aliasOwners {
				want = code
			}
		}
		got, err := repo.FindByName(examples[key].label)
		if err != nil {
			t.Errorf("lookup alias %q: %v, want %s", examples[key].label, err, want)
			continue
		}
		if got.ISO2 != want {
			t.Errorf("alias %q resolved to %s, want %s", examples[key].label, got.ISO2, want)
		}
	}
}

func TestProductionOfficeCountryLabelsMatchAcrossLocales(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}

	// These country-only labels were observed in bounded, read-only production
	// office-location samples. Russian/Ukrainian labels are the discovery set;
	// Chinese/Japanese/Korean labels are the disjoint script holdout.
	for _, tc := range []struct{ label, want string }{
		{"Россия", "RU"},
		{"Україна", "UA"},
		{"中国", "CN"},
		{"日本", "JP"},
		{"한국", "KR"},
		{"Magyarország", "HU"},
		{"Česká republika", "CZ"},
		{"Ελλάδα", "GR"},
	} {
		got, err := repo.FindByName(tc.label)
		if err != nil {
			t.Errorf("lookup %q: %v, want %s", tc.label, err, tc.want)
			continue
		}
		if got.ISO2 != tc.want {
			t.Errorf("lookup %q = %s, want %s", tc.label, got.ISO2, tc.want)
		}
	}

	// Nearby negative checks ensure the two Cyrillic endonyms stay attached to
	// their own countries rather than becoming cross-country aliases.
	for _, tc := range []struct{ label, wrong string }{
		{"Россия", "UA"},
		{"Україна", "RU"},
	} {
		got, err := repo.FindByName(tc.label)
		if err == nil && got.ISO2 == tc.wrong {
			t.Errorf("lookup %q incorrectly resolved to %s", tc.label, tc.wrong)
		}
	}
}
