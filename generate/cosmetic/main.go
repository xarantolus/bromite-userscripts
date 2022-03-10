package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"

	"cosmetic/filter"
	"cosmetic/util"
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

func main() {
	var (
		inputLists   = flag.String("input", "filter-lists.txt", "Path to file that defines URLs to blocklists")
		scriptTarget = flag.String("output", "cosmetic.user.js", "Path to output file")
	)
	flag.Parse()

	scriptTemplateContent, err := ioutil.ReadFile("script-template.js")
	if err != nil {
		log.Fatalf("reading script template file: %s\n", err.Error())
	}
	var scriptTemplate = template.Must(template.New("").Parse(string(scriptTemplateContent)))

	filterURLs, err := util.ReadListFile(*inputLists)
	if err != nil {
		log.Fatalf("cannot load list of filter URLs: %s\n", err.Error())
	}

	tempDir, err := ioutil.TempDir("", "cosmetic-filter-*")
	if err != nil {
		log.Fatalf("creating temp dir for cosmetic filters: %s\n", err.Error())
	}
	defer os.RemoveAll(tempDir)

	filterOutputFiles, err := util.DownloadURLs(filterURLs, tempDir)
	if err != nil {
		log.Fatalf("error downloading filter lists: %s\n", err.Error())
	}
	log.Printf("Downloaded %d filter files\n", len(filterOutputFiles))

	var filters []filter.Rule
	for _, fp := range filterOutputFiles {
		ff := util.FiltersFromFile(fp)
		if len(ff) == 0 {
			log.Printf("[Warning] No rules found in file %q\n", fp)
		}
		filters = append(filters, ff...)
	}
	fmt.Printf("Found %d filters in these files\n", len(filters))

	lookupTable := filter.Combine(filters)

	var duplicateCount = map[string]int{}
	for _, f := range lookupTable {
		joined := joinSorted(f.Selectors, ",")
		duplicateCount[joined] = duplicateCount[joined] + 1
		joined = joinSorted(f.InjectedCSS, "")
		duplicateCount[joined] = duplicateCount[joined] + 1
	}

	var deduplicatedStrings []string
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
	var (
		compiledSelectorRules  = map[string]interface{}{}
		compiledInjectionRules = map[string]interface{}{}
	)
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

	outputFile, err := os.Create(*scriptTarget)
	if err != nil {
		log.Fatalf("creating output file: %s\n", err.Error())
	}

	_, err = outputFile.WriteString("// THIS FILE IS AUTO-GENERATED. DO NOT EDIT. See generate/cosmetic directory for more info\n")
	if err != nil {
		log.Fatalf("could not write auto generated message: %s\n", err.Error())
	}

	err = scriptTemplate.Execute(outputFile, map[string]string{
		"version":             time.Now().Format("2006.01.02"),
		"rules":               toJSObject(compiledSelectorRules),
		"injectionRules":      toJSObject(compiledInjectionRules),
		"deduplicatedStrings": toJSObject(deduplicatedStrings),
	})
	if err != nil {
		log.Fatalf("Error generating script text: %s\n", err.Error())
	}

	err = outputFile.Close()
	if err != nil {
		log.Fatalf("could not close output file: %s\n", err.Error())
	}
}
