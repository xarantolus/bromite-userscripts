package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"idcac/amo"
	"idcac/extract"

	"github.com/xarantolus/jsonextract"
)

func toJSString(obj interface{}) string {
	s, err := json.Marshal(obj)
	if err != nil {
		panic("turning object into JS string: " + err.Error())
	}
	return string(s)

}

func isCommons(m map[string]interface{}) bool {
	for k := range m {
		_, numerr := strconv.Atoi(k)
		if numerr != nil {
			return false
		}
	}

	return true
}

func isRules(m map[string]interface{}) bool {
	for k := range m {
		if !strings.ContainsRune(k, '.') {
			return false
		}
	}
	return true
}

func mapFunctions(fnts map[int]string) string {
	var sb = strings.Builder{}

	sb.WriteString("{")

	var keys []int
	for key := range fnts {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	var counter int
	for _, key := range keys {
		sb.WriteString(toJSString(strconv.Itoa(key)))
		sb.WriteString(":")
		sb.WriteString("(function () {")
		sb.WriteString(fnts[key])
		sb.WriteString("})")
		counter++

		if counter != len(fnts) {
			sb.WriteString(",")
		}
	}

	sb.WriteString("}")

	return sb.String()
}

const extensionURL = "https://addons.mozilla.org/en/firefox/addon/i-dont-care-about-cookies/"

func main() {
	var (
		scriptTarget = flag.String("output", "idcac.user.js", "Path to output file")
	)
	flag.Parse()

	scriptTemplateContent, err := ioutil.ReadFile("script-template.js")
	if err != nil {
		log.Fatalf("reading script template file: %s\n", err.Error())
	}
	var scriptTemplate = template.Must(template.New("").Parse(string(scriptTemplateContent)))

	extensionInfo, err := amo.ScrapeInfo(extensionURL)
	if err != nil {
		log.Fatalf("scraping idcac extension info: %s\n", err.Error())
	}

	extensionBaseDir, err := os.MkdirTemp("", "idcac-extension-*")
	if err != nil {
		log.Fatalf("creating temporary directory for idcac extension: %s\n", err.Error())
	}
	defer os.RemoveAll(extensionBaseDir)

	extStream, cancel, err := amo.DownloadFile(extensionInfo.URL)
	defer cancel()
	if err != nil {
		log.Fatalf("downloading extension zip file: %s\n", err.Error())
	}
	defer extStream.Close()

	err = extract.Zip(extStream, extensionBaseDir)
	if err != nil {
		log.Fatalf("extracting extension zip: %s\n", err.Error())
	}

	f, err := os.Open(filepath.Join(extensionBaseDir, "data/rules.js"))
	if err != nil {
		log.Fatalf("Open rules file from extension: %s\n", err.Error())
	}
	defer f.Close()

	var (
		commons        string
		rules          string
		cookieBlockCSS string

		javascriptFixes = make(map[int]string)
	)

	err = jsonextract.Reader(f, func(b []byte) error {
		var data = make(map[string]interface{})

		err = json.Unmarshal(b, &data)
		if err != nil {
			return err
		}

		if isCommons(data) {
			commons = string(b)
		} else if isRules(data) {
			rules = string(b)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Reading/Converting rules file: %s\n", err.Error())
	}

	cookieBlockCSSBytes, err := os.ReadFile(filepath.Join(extensionBaseDir, "data/css/common.css"))
	if err != nil {
		log.Fatalf("Error reading common css rules: %s\n", err.Error())
	}
	cookieBlockCSS = string(cookieBlockCSSBytes)

	err = filepath.WalkDir(filepath.Join(extensionBaseDir, "data/js"), func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		var base = filepath.Base(path)

		var number int

		_, err = fmt.Sscanf(base, "common%d.js", &number)
		if err != nil && base != "common.js" {
			return nil
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading common javascript fix file: %w", err)
		}

		javascriptFixes[number] = strings.TrimSpace(string(fileContent))

		return nil
	})

	if len(commons) == 0 || len(rules) == 0 || len(javascriptFixes) == 0 {
		log.Fatalf("Unexpected lengths of commons(%d)/rules(%d)/javascriptFixes(%d) -- expected at least one of each", len(commons), len(rules), len(javascriptFixes))
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
		"version":         time.Now().Format("2006.01.02"),
		"commons":         commons,
		"rules":           rules,
		"javascriptFixes": mapFunctions(javascriptFixes),
		"cookieBlockCSS":  toJSString(cookieBlockCSS),
	})
	if err != nil {
		log.Fatalf("Error generating script text: %s\n", err.Error())
	}

	err = outputFile.Close()
	if err != nil {
		log.Fatalf("could not close output file: %s\n", err.Error())
	}
}
