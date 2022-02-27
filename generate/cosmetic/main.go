package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"cosmetic/filter"
)

func parseFilterList(f io.Reader) (filters []filter.BasicFilter) {
	scan := bufio.NewScanner(f)

	for scan.Scan() {
		txt := strings.TrimSpace(scan.Text())

		if len(txt) == 0 || strings.HasPrefix(txt, "!") {
			continue
		}

		filter, ok := filter.ParseLine(txt)
		if ok {
			filters = append(filters, filter)
		}
	}

	return filters
}

func filtersFromFile(filepath string) (filters []filter.BasicFilter) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return parseFilterList(f)
}

//go:embed script-template.js
var scriptTemplateString string

var scriptTemplate = template.Must(template.New("").Parse(scriptTemplateString))

func main() {
	var filters []filter.BasicFilter

	err := filepath.Walk("t", func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		filters = append(filters, filtersFromFile(path)...)
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d filters\n", len(filters))

	lookup := filter.Combine(filters)

	var compiledRules = map[string]string{}
	for k, f := range lookup {
		compiledRules[k] = strings.Join(f, ",")
	}

	fmt.Printf("Combined them for %d domains", len(compiledRules))

	rules, err := json.Marshal(compiledRules)
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create("out.js")
	if err != nil {
		log.Fatalf("creating output file: %s\n", err.Error())
	}

	_, err = outputFile.WriteString("// THIS FILE IS AUTO-GENERATED. DO NOT EDIT. See generate/idcac directory for more info\n")
	if err != nil {
		log.Fatalf("could not write auto generated message: %s\n", err.Error())
	}

	err = scriptTemplate.Execute(outputFile, map[string]string{
		"version": "1.0.0",
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
