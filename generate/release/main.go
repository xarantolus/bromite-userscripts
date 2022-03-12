package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func env(name string, defaultVal string) string {
	v := os.Getenv(name)
	if strings.TrimSpace(v) == "" {
		return defaultVal
	}
	return v
}

const outputTemplate = "This release contains all scripts provided by this repository.\n\nPlease see [the main project page](https://github.com/{{.repo}}) for a description of the scripts.{{if .stats}}\n\n**Stats**:{{range .stats}}\n* `{{.ScriptName}}`: {{.StatsLine}}{{end}}{{end}}"

var statMarker = []byte("/// @stats")

type stats struct {
	ScriptName string
	StatsLine  string
}

var errNoMarker = fmt.Errorf("no stats line (a line starting with %q) found", string(statMarker))

func getStats(filename string) (s stats, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	var line string
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		by := bytes.TrimSpace(scan.Bytes())

		if bytes.HasPrefix(by, statMarker) {
			line = string(bytes.TrimSpace(by[len(statMarker):]))
			break
		}
	}
	if err = scan.Err(); err != nil {
		return
	}
	if len(line) == 0 {
		err = errNoMarker
		return
	}

	return stats{
		ScriptName: filepath.Base(filename),
		StatsLine:  line,
	}, nil
}

func main() {
	var (
		outputFilename = flag.String("output", "release.md", "Path to release file path")

		repo = env("GITHUB_REPOSITORY", "xarantolus/bromite-userscripts")

		tmpl = template.Must(template.New("").Parse(outputTemplate))
	)
	flag.Parse()

	var stats []stats
	for _, inputFile := range flag.Args() {
		s, err := getStats(inputFile)
		if err != nil {
			if errors.Is(err, errNoMarker) {
				log.Printf("Warning in %s: %s\n", filepath.Base(inputFile), err.Error())
			} else {
				log.Printf("Error in %s: %s\n", filepath.Base(inputFile), err.Error())
			}
			continue
		}

		stats = append(stats, s)
	}

	f, err := os.Create(*outputFilename)
	if err != nil {
		log.Fatalf("creating output file: %s\n", err.Error())
	}

	err = tmpl.Execute(f, map[string]interface{}{
		"repo":  repo,
		"stats": stats,
	})
	if err != nil {
		log.Fatalf("executing template: %s\n", err.Error())
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("closing file: %s\n", err.Error())
	}
}
