package filter

import (
	"strings"

	"github.com/andybalholm/cascadia"
)

type BasicFilter struct {
	Domains []string

	CSSSelector string
}

// See https://help.eyeo.com/en/adblockplus/how-to-write-filters, "Content Filters"
func ParseLine(line string) (f BasicFilter, ok bool) {
	split := strings.SplitN(line, "##", 2)
	if len(split) != 2 {
		return f, false
	}

	_, err := cascadia.Parse(split[1])
	if err != nil {
		return f, false
	}

	if strings.ContainsAny(split[0], "*~") || strings.Contains(split[0], "#@") {
		return f, false
	}

	return BasicFilter{
		Domains:     strings.Split(split[0], ","),
		CSSSelector: split[1],
	}, true
}
