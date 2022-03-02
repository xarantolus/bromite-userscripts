package util

import (
	"bufio"
	"io"
	"os"
	"strings"

	"cosmetic/filter"
)

func ParseFilterList(f io.Reader) (filters []filter.Rule) {
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

func FiltersFromFile(filepath string) (filters []filter.Rule) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return ParseFilterList(f)
}
