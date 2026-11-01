package util

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"cosmetic/filter"
)

func ParseFilterList(f io.Reader) (filters []filter.Rule) {
	scan := bufio.NewScanner(f)
	scan.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

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
	if err := scan.Err(); err != nil {
		log.Printf("[Warning] Error while parsing filter list, it may be incomplete: %s\n", err.Error())
	}

	return filters
}

func FiltersFromFile(filepath string) (filters []filter.Rule) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return ParseFilterList(f)
}
