package filter

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/andybalholm/cascadia"
)

type Rule struct {
	Domains []string

	CSSSelector string

	InjectedCSS string
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

var (
	injectedStyleRegex = regexp.MustCompile(`(.*?)\:style\((.*?)\)`)
)

// See https://help.eyeo.com/en/adblockplus/how-to-write-filters, "Content Filters"
func ParseLine(line string) (f Rule, ok bool) {
	var isCSSInjection bool

	var split []string
	// Check which type of rule we got in this line
	if split = strings.SplitN(line, "##", 2); len(split) == 2 {
		// Element hiding rule
		// Example:   domain1.com,domain2.com##.blocked-element
		isCSSInjection = false
	} else if split = strings.SplitN(line, "#$#", 2); len(split) == 2 {
		// A CSS injection
		// Example:   domain1.com,domain2.com#$#.cookie { display: none!important; }
		isCSSInjection = true
	} else {
		// The statement in this line is not recognized, ignore it
		return f, false
	}

	// We currently only support very basic filters.
	// This check makes sure the other filters types are filtered out
	if split[0] != "*" && strings.ContainsAny(split[0], "*~#@") {
		return f, false
	}

	var (
		injectedStyle string
		selector      = split[1]
	)
	if strings.Contains(split[1], ":style") {
		matches := injectedStyleRegex.FindStringSubmatch(split[1])
		if len(matches) != 3 {
			return f, false
		}
		selector = ""
		injectedStyle = matches[1] + "{" + matches[2] + "}"
	} else if isCSSInjection {
		selector = ""
		injectedStyle = split[1]
	} else {
		// Make sure we only get valid selectors
		if isIncompatibleSelector(selector) {
			return f, false
		}
	}

	domains := strings.FieldsFunc(split[0], func(r rune) bool {
		return unicode.IsSpace(r) || r == ','
	})
	if len(domains) == 0 {
		// General rules for all domains need the empty domain to work with the script
		domains = append(domains, "")
	}
	if split[0] == "*" {
		domains = []string{""}
	}

	for i := 0; i < len(domains); i++ {
		domains[i] = strings.TrimPrefix(domains[i], "www.")
	}

	return Rule{
		Domains:     domains,
		CSSSelector: selector,
		InjectedCSS: injectedStyle,
	}, true
}
