package filter

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		input  string
		wantF  Rule
		wantOk bool
	}{
		// General rules should have the domain "", as that will be used in the script to inject them in all pages
		{"###cookie_alert", Rule{Domains: []string{""}, CSSSelector: "#cookie_alert"}, true},
		// Wildcard "*" is also supported
		{"*###cookie_alert", Rule{Domains: []string{""}, CSSSelector: "#cookie_alert"}, true},

		// Normal rules that block certain CSS selectors
		{"example.com##.ad", Rule{Domains: []string{"example.com"}, CSSSelector: ".ad"}, true},
		{"example.com###ad", Rule{Domains: []string{"example.com"}, CSSSelector: "#ad"}, true},
		{"a.com, b.com###ad", Rule{Domains: []string{"a.com", "b.com"}, CSSSelector: "#ad"}, true},

		// Rules for injecting CSS styles, usually to fix scrolling issues
		{"ndtv.com##body:style(overflow: auto !important)", Rule{Domains: []string{"ndtv.com"}, InjectedCSS: "body{overflow: auto !important}"}, true},
		{"seb.lt,seb.ee##body,html:style(height: auto !important; overflow: auto !important)", Rule{Domains: []string{"seb.lt", "seb.ee"}, InjectedCSS: "body,html{height: auto !important; overflow: auto !important}"}, true},
		{"example.com#$#body { background-color: #333!important; }", Rule{Domains: []string{"example.com"}, InjectedCSS: "body { background-color: #333!important; }"}, true},

		// Some rules I found in real files
		{`yandex.ru#$#div[class*="_with-url-actualizer_yes"] > div.adv_type_top { display: none !important; }`, Rule{Domains: []string{"yandex.ru"}, InjectedCSS: `div[class*="_with-url-actualizer_yes"] > div.adv_type_top { display: none !important; }`}, true},
		{"##.gdpr-box:not(body):not(html)", Rule{Domains: []string{""}, CSSSelector: ".gdpr-box:not(body):not(html)"}, true},
		{"arm.com##.c-policies", Rule{Domains: []string{"arm.com"}, CSSSelector: ".c-policies"}, true},

		// Invalid rules that should be rejected
		// has() is not supported by chromium
		{"example.com##.ad:has(.child)", Rule{}, false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if tt.wantOk && len(tt.wantF.Domains) == 0 {
				panic("Invalid test case: in case of no expected domain, the domain should be the empty string")
			}
			gotF, gotOk := ParseLine(tt.input)
			if gotOk != tt.wantOk {
				t.Errorf("ParseLine(%q) gotOk = %#v, want %#v", tt.input, gotOk, tt.wantOk)
			} else if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("ParseLine(%q) gotF = %#v, want %#v", tt.input, gotF, tt.wantF)
			}
		})
	}
}
