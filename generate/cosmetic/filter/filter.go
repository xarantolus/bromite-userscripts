package filter

import (
	"strings"
	"unicode"

	"github.com/andybalholm/cascadia"
)

type BasicFilter struct {
	Domains []string

	CSSSelector string
}

func isIncompatibleSelector(s string) bool {
	// We want only valid selectors, so we check if we can parse it
	_, err := cascadia.Parse(s)
	if err != nil {
		return true
	}

	// Chromium doesn't seem to support the :has() selector
	if strings.Contains(s, ":has(") {
		return true
	}

	// We assume that anything else is supported
	return false
}

// See https://help.eyeo.com/en/adblockplus/how-to-write-filters, "Content Filters"
func ParseLine(line string) (f BasicFilter, ok bool) {
	split := strings.SplitN(line, "##", 2)
	if len(split) != 2 {
		return f, false
	}

	// We currently only support very basic filters.
	// This check makes sure the other filters types are filtered out
	if strings.ContainsAny(split[0], "*~#@") {
		return f, false
	}

	// Make sure we only get valid selectors
	if isIncompatibleSelector(split[1]) {
		return f, false
	}

	domains := strings.FieldsFunc(split[0], func(r rune) bool {
		return unicode.IsSpace(r) || r == ','
	})
	if len(domains) == 0 {
		// General rules for all domains need the empty domain to work with the script
		domains = append(domains, "")
	}

	return BasicFilter{
		Domains:     domains,
		CSSSelector: split[1],
	}, true
}
