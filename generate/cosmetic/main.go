package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"

	_ "embed"

	"cosmetic/filter"
	"cosmetic/filterlists"
	"cosmetic/topdomains"
	"cosmetic/util"

	"golang.org/x/sync/errgroup"
)

func joinSorted(f []string, comma string) string {
	sort.Strings(f)
	return strings.Join(f, comma)
}

func toJSObject(x interface{}) string {
	b, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}

var (
	//go:embed script-template.js
	scriptTemplateRaw string
	scriptTemplate    = template.Must(template.New("").Parse(string(scriptTemplateRaw)))
)

func readFilterLists(filterListFiles []string, topDomains *topdomains.TopDomainStorage) (compiledSelectorRules map[string]interface{}, compiledInjectionRules map[string]interface{}, deduplicatedStrings []string) {
	var filters []filter.Rule
	for _, fp := range filterListFiles {
		ff := util.FiltersFromFile(fp)
		if len(ff) == 0 {
			log.Printf("[Warning] No rules found in file %q\n", fp)
		}
		filters = append(filters, ff...)
	}
	fmt.Printf("Found %d filters in these files\n", len(filters))

	lookupTable := filter.Combine(filters)

	if topDomains != nil {
		// Now only keep the filters for top/important domains
		topDomainLookupTable := make(map[string]filter.CombineResult)
		for domain, filter := range lookupTable {
			if topDomains.Contains(domain) {
				topDomainLookupTable[domain] = filter
			}
		}
		fmt.Printf("Selected %d top domains from %d domains with available filters\n", len(topDomainLookupTable), len(lookupTable))

		// Also keep the default/general rules for all pages
		topDomainLookupTable[""] = lookupTable[""]

		lookupTable = topDomainLookupTable
	}

	var duplicateCount = map[string]int{}
	for _, f := range lookupTable {
		joined := joinSorted(f.Selectors, ",")
		duplicateCount[joined] = duplicateCount[joined] + 1
		joined = joinSorted(f.InjectedCSS, "")
		duplicateCount[joined] = duplicateCount[joined] + 1
	}

	for f, count := range duplicateCount {
		if count > 1 {
			deduplicatedStrings = append(deduplicatedStrings, f)
		}
	}
	sort.Strings(deduplicatedStrings)

	var deduplicatedIndexMapping = map[string]int{}
	for i, r := range deduplicatedStrings {
		deduplicatedIndexMapping[r] = i
	}

	// The compiled rules are either
	// - a string, which is a css selector (usually selecting many elements)
	// - an int, which is the index of a common rule (that was present more than once)
	compiledSelectorRules = map[string]interface{}{}
	compiledInjectionRules = map[string]interface{}{}
	for domain, filter := range lookupTable {
		if len(filter.Selectors) > 0 {
			joined := joinSorted(filter.Selectors, ",")
			if duplicateCount[joined] > 1 {
				compiledSelectorRules[domain] = deduplicatedIndexMapping[joined]
			} else {
				compiledSelectorRules[domain] = joined
			}
		}

		if len(filter.InjectedCSS) > 0 {
			joined := joinSorted(filter.InjectedCSS, "")
			if duplicateCount[joined] > 1 {
				compiledInjectionRules[domain] = deduplicatedIndexMapping[joined]
			} else {
				compiledInjectionRules[domain] = joined
			}
		}
	}
	fmt.Printf("Combined them for %d domains\n", len(compiledSelectorRules))
	return
}

// if topDomains is nil, all domains are included - otherwise we generate the "lite" list
func generateListForCountry(outputDir string, listURLs []string, countryName, countryCode string, topDomainsFilter *topdomains.TopDomainStorage) (scriptPath string, statsLine string, err error) {
	var filename = fmt.Sprintf("cosmetic-%s.user.js", countryCode)
	if topDomainsFilter != nil {
		filename = fmt.Sprintf("cosmetic-%s-lite.user.js", countryCode)
	}

	compiledSelectorRules, compiledInjectionRules, deduplicatedStrings := readFilterLists(listURLs, topDomainsFilter)

	scriptPath = path.Join(outputDir, filename)
	outputFile, err := os.Create(scriptPath)
	if err != nil {
		err = fmt.Errorf("creating output file: %w", err)
		return
	}

	_, err = outputFile.WriteString("// THIS FILE IS AUTO-GENERATED. DO NOT EDIT. See generate/cosmetic directory for more info\n")
	if err != nil {
		err = fmt.Errorf("could not write auto generated message: %w", err)
		return
	}

	statsLine = fmt.Sprintf("blockers for %d domains, injected CSS rules for %d domains", len(compiledSelectorRules), len(compiledInjectionRules))
	err = scriptTemplate.Execute(outputFile, map[string]interface{}{
		"version":             time.Now().Format("2006.01.02"),
		"rules":               toJSObject(compiledSelectorRules),
		"injectionRules":      toJSObject(compiledInjectionRules),
		"deduplicatedStrings": toJSObject(deduplicatedStrings),
		"statistics":          statsLine,
		"countryName":         countryName,
		"countryCode":         countryCode,
		"filename":            filename,
		"isLite":              topDomainsFilter != nil,
	})
	if err != nil {
		err = fmt.Errorf("error generating script text: %w", err)
		return
	}

	err = outputFile.Close()
	if err != nil {
		err = fmt.Errorf("could not close output file: %w", err)
		return
	}

	return
}

func processLanguage(lang filterlists.Language, filterLists filterlists.FilterLists, extraLists []string, urlTempDir string, outputDir string, topDomainsFilter topdomains.TopDomainStorage) (fullScriptInfo, liteScriptInfo ScriptInfo, err error) {
	log.Printf("Processing language %q\n", lang.Name)

	filterListsForLang := filterLists.ForLanguages([]filterlists.Language{lang})

	var filterListURLs []string
	for _, fl := range filterListsForLang {
		filterListURLs = append(filterListURLs, fl.PrimaryViewURL)
	}

	filterListURLs = append(filterListURLs, extraLists...)

	codeAdditionsPath := path.Join("additions", fmt.Sprintf("%s.txt", lang.Iso6391))
	countrySpecificAdditions, err := util.ReadListFile(codeAdditionsPath)
	if err != nil && !os.IsNotExist(err) {
		err = fmt.Errorf("error reading code additions file: %w", err)
		return ScriptInfo{}, ScriptInfo{}, err
	}
	filterListURLs = append(filterListURLs, countrySpecificAdditions...)

	filterFiles, err := util.DownloadURLs(filterListURLs, urlTempDir)
	if err != nil {
		// Ignore - prefer generating an empty script over having none
		log.Printf("[Warning] Error downloading filter lists for language %q: %s\n", lang.Name, err.Error())
	}

	fullScriptPath, fullStatsLine, err := generateListForCountry(outputDir, filterFiles, lang.Name, lang.Iso6391, nil)
	if err != nil {
		err = fmt.Errorf("error generating full script for language %q: %w", lang.Name, err)
		return
	}

	// Now with filtering of top domains
	liteScriptPath, liteStatsLine, err := generateListForCountry(outputDir, filterFiles, lang.Name, lang.Iso6391, &topDomainsFilter)
	if err != nil {
		err = fmt.Errorf("error generating lite script for language %q: %w", lang.Name, err)
		return
	}

	// Get size of both
	fi, err := os.Stat(fullScriptPath)
	if err != nil {
		err = fmt.Errorf("error getting size of full script: %w", err)
		return
	}
	li, err := os.Stat(liteScriptPath)
	if err != nil {
		err = fmt.Errorf("error getting size of lite script: %w", err)
		return
	}

	return ScriptInfo{
			ScriptBaseName: path.Base(fullScriptPath),
			LanguageName:   lang.Name,
			LanguageCode:   lang.Iso6391,
			Stats:          fullStatsLine,
			Size:           fi.Size(),
		}, ScriptInfo{
			ScriptBaseName: path.Base(liteScriptPath),
			LanguageName:   lang.Name,
			LanguageCode:   lang.Iso6391,
			Stats:          liteStatsLine,
			Size:           li.Size(),
			IsLite:         true,
		},
		nil
}

type ScriptInfo struct {
	ScriptBaseName string `json:"script_name"`
	LanguageName   string `json:"language_name"`
	LanguageCode   string `json:"language_code"`
	Stats          string `json:"stats"`
	IsLite         bool   `json:"is_lite"`
	Size           int64  `json:"size"`
}

const outputDir = "cosmetic-outputs"

func main() {
	var (
		baseURLList = flag.String("baseList", "base-filters.txt", "Path to base URL list")
	)
	flag.Parse()

	baseFilterURLs, err := util.ReadListFile(*baseURLList)
	if err != nil {
		log.Fatalf("cannot load list of filter URLs: %s\n", err.Error())
	}

	urlTempDir, err := os.MkdirTemp("", "cosmetic-urls-*")
	if err != nil {
		log.Fatalf("creating temp dir for cosmetic filters: %s\n", err.Error())
	}
	defer os.RemoveAll(urlTempDir)

	td, err := topdomains.FromFile("top1m/top-1m.csv", 1_000_000)
	if err != nil {
		log.Fatalf("error loading top domains: %s\n", err.Error())
	}

	// Delete and recreate output directory
	err = os.RemoveAll(outputDir)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("error deleting output directory: %s\n", err.Error())
	}
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("error creating output directory: %s\n", err.Error())
	}

	languages, err := filterlists.FetchLanguages()
	if err != nil {
		log.Fatalf("error fetching languages: %s\n", err.Error())
	}

	filterLists, err := filterlists.FetchFilterLists()
	if err != nil {
		log.Fatalf("error fetching filter lists: %s\n", err.Error())
	}

	var (
		scriptInfos     []ScriptInfo
		scriptInfosLock sync.Mutex
	)

	var baseLanguage = filterlists.Language{
		ID:      -1337,
		Iso6391: "base",
		Name:    "Base",
	}

	var eg errgroup.Group
	eg.Go(func() error {
		fullScriptInfo, liteScriptInfo, err := processLanguage(baseLanguage, filterLists, baseFilterURLs, urlTempDir, outputDir, td)
		if err != nil {
			return err
		}

		scriptInfosLock.Lock()
		scriptInfos = append(scriptInfos, fullScriptInfo)
		scriptInfos = append(scriptInfos, liteScriptInfo)
		scriptInfosLock.Unlock()
		return nil
	})

	for _, lang := range languages {
		lcpy := lang
		eg.Go(func() error {
			fullScriptInfo, liteScriptInfo, err := processLanguage(lcpy, filterLists, nil, urlTempDir, outputDir, td)
			if err != nil {
				return err
			}
			scriptInfosLock.Lock()
			scriptInfos = append(scriptInfos, fullScriptInfo)
			scriptInfos = append(scriptInfos, liteScriptInfo)
			scriptInfosLock.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		log.Fatalf("error processing languages: %s\n", err.Error())
	}

	sort.Slice(scriptInfos, func(i, j int) bool {
		return scriptInfos[i].ScriptBaseName < scriptInfos[j].ScriptBaseName
	})

	scriptInfosJSON, err := json.Marshal(scriptInfos)
	if err != nil {
		log.Fatalf("error generating script infos: %s\n", err.Error())
	}

	err = os.WriteFile(path.Join(outputDir, "cosmetic_info.json"), scriptInfosJSON, 0644)
	if err != nil {
		log.Fatalf("error writing script infos: %s\n", err.Error())
	}

	// now the same thing with jsonp
	scriptInfosJSONP := append([]byte("jsonpCosmeticInfo("), append(scriptInfosJSON, []byte(");")...)...)
	err = os.WriteFile(path.Join(outputDir, "cosmetic_info.jsonp"), scriptInfosJSONP, 0644)
	if err != nil {
		log.Fatalf("error writing script infos: %s\n", err.Error())
	}
}
