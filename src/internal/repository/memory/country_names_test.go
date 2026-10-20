package memory

import (
	"encoding/json"
	"os"
	"strings"
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

	owners := make(map[string]map[string]struct{})
	for code, labels := range aliases {
		for _, label := range labels {
			key := textNormalizer.Normalize(label)
			if owners[key] == nil {
				owners[key] = make(map[string]struct{})
			}
			owners[key][code] = struct{}{}
		}
	}

	for key, aliasOwners := range owners {
		if len(aliasOwners) > 1 && canonicalOwners[key] == "" {
			t.Errorf("normalized alias key %q is assigned to multiple countries without a canonical name/code", key)
		}
	}

	for code, labels := range aliases {
		for _, label := range labels {
			key := textNormalizer.Normalize(label)
			want := canonicalOwners[key]
			if want == "" {
				want = code
			}
			if exactCode, ok := repo.exactNameToCode[strings.TrimSpace(label)]; ok {
				want = exactCode
			} else if containsNonASCII(label) {
				if exactCode, ok := repo.accentNameToCode[accentPreservingKey(label)]; ok {
					want = exactCode
				}
			}
			got, err := repo.FindByName(label)
			if err != nil {
				t.Errorf("lookup alias %q: %v, want %s", label, err, want)
				continue
			}
			if got.ISO2 != want {
				t.Errorf("alias %q resolved to %s, want %s", label, got.ISO2, want)
			}
		}
	}
}

func TestAllVietnameseCLDRTerritoryNamesResolve(t *testing.T) {
	loader := data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")
	countries, err := loader.LoadCountries()
	if err != nil {
		t.Fatalf("load bundled countries: %v", err)
	}
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), loader)
	if err != nil {
		t.Fatalf("load country repository: %v", err)
	}

	fixture, err := os.ReadFile("../../../../data/cldr/48.2.3/vi/territories.json")
	if err != nil {
		t.Fatalf("read pinned CLDR Vietnamese territory fixture: %v", err)
	}
	var source struct {
		Main map[string]struct {
			LocaleDisplayNames struct {
				Territories map[string]string `json:"territories"`
			} `json:"localeDisplayNames"`
		} `json:"main"`
	}
	if err := json.Unmarshal(fixture, &source); err != nil {
		t.Fatalf("parse pinned CLDR fixture: %v", err)
	}
	vi, ok := source.Main["vi"]
	if !ok {
		t.Fatal("pinned CLDR fixture has no vi locale")
	}

	assigned := make(map[string]struct{}, len(countries))
	for _, country := range countries {
		assigned[country.ISO2] = struct{}{}
	}
	checked := 0
	for code, label := range vi.LocaleDisplayNames.Territories {
		if _, ok := assigned[code]; !ok {
			continue
		}
		got, err := repo.FindByName(label)
		if err != nil {
			t.Errorf("CLDR vi label for %s: %v", code, err)
			continue
		}
		if got.ISO2 != code {
			t.Errorf("CLDR vi label for %s resolved to %s", code, got.ISO2)
		}
		checked++
	}
	if checked != len(countries) {
		t.Fatalf("checked %d CLDR Vietnamese territory labels, want all %d assigned ISO alpha-2 entries", checked, len(countries))
	}
}

func TestVietnameseNamesDoNotGetCapturedByISOCodeKeys(t *testing.T) {
	repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ query, want string }{
		{"Áo", "AT"},
		{"áo", "AT"},
		{"A\u0301o", "AT"}, // decomposed accent form
		{"Bỉ", "BE"},
		{"bỈ", "BE"},
		{"Nga", "RU"}, // CLDR's title-cased name shares an ASCII spelling with ISO3 NGA.
		{"AO", "AO"},
		{"ao", "AO"},
		{"BI", "BI"},
		{"bi", "BI"},
		{"NGA", "NG"},
		{"nga", "NG"},
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

func TestIndonesianCLDRCountryNameResolves(t *testing.T) {
	// CLDR 48.2.3 lists "Amerika Serikat" as the Indonesian territory name for US.
	for _, tc := range []struct {
		name   string
		loader data.Loader
	}{
		{name: "bundled CSV", loader: data.NewCSVLoader("../../../../data/countries.csv", "../../../../data/aliases.csv")},
		{name: "memory", loader: data.NewMemoryLoader()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, err := NewCountryRepository(normalizer.NewTextNormalizer(), tc.loader)
			if err != nil {
				t.Fatal(err)
			}
			got, err := repo.FindByName("Amerika Serikat")
			if err != nil {
				t.Fatalf("lookup Indonesian CLDR label: %v", err)
			}
			if got.ISO2 != "US" {
				t.Fatalf("lookup Indonesian CLDR label = %s, want US", got.ISO2)
			}
		})
	}
}
