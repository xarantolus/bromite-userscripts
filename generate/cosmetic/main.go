package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"cosmetic/filter"
	"cosmetic/util"
)

//go:embed script-template.js
var scriptTemplateString string

var scriptTemplate = template.Must(template.New("").Parse(scriptTemplateString))

func main() {
	var (
		inputLists   = flag.String("input", "filter-lists.txt", "Path to file that defines URLs to blocklists")
		scriptTarget = flag.String("output", "cosmetic.user.js", "Path to output file")
	)
	flag.Parse()

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

	var filters []filter.BasicFilter
	for _, fp := range filterOutputFiles {
		ff := util.FiltersFromFile(fp)
		if len(ff) == 0 {
			log.Printf("[Warning] No rules found in file %q\n", fp)
		}
		filters = append(filters, ff...)
	}
	fmt.Printf("Found %d filters in these files\n", len(filters))

	lookupTable := filter.Combine(filters)

	var compiledRules = map[string]string{}
	for k, f := range lookupTable {
		compiledRules[k] = strings.Join(f, ",")
	}

	fmt.Printf("Combined them for %d domains\n", len(compiledRules))

	rules, err := json.Marshal(compiledRules)
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create(*scriptTarget)
	if err != nil {
		log.Fatalf("creating output file: %s\n", err.Error())
	}

	_, err = outputFile.WriteString("// THIS FILE IS AUTO-GENERATED. DO NOT EDIT. See generate/idcac directory for more info\n")
	if err != nil {
		log.Fatalf("could not write auto generated message: %s\n", err.Error())
	}

	err = scriptTemplate.Execute(outputFile, map[string]string{
		"version": time.Now().Format("2006.01.02"),
		"rules":   string(rules),
	})
	if err != nil {
		log.Fatalf("Error generating script text: %s\n", err.Error())
	}

	err = outputFile.Close()
	if err != nil {
		log.Fatalf("could not close output file: %s\n", err.Error())
	}
}
