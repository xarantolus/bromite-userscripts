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

		// Normal rules that block certain CSS selectors
		{"example.com##.ad", Rule{Domains: []string{"example.com"}, CSSSelector: ".ad"}, true},
		{"example.com###ad", Rule{Domains: []string{"example.com"}, CSSSelector: "#ad"}, true},
		{"a.com, b.com###ad", Rule{Domains: []string{"a.com", "b.com"}, CSSSelector: "#ad"}, true},

		// Rules for injecting CSS styles, usually to fix scrolling issues
		{"ndtv.com##body:style(overflow: auto !important)", Rule{Domains: []string{"ndtv.com"}, InjectedCSS: "body{overflow: auto !important}"}, true},
		{"seb.lt,seb.ee##body,html:style(height: auto !important; overflow: auto !important)", Rule{Domains: []string{"seb.lt", "seb.ee"}, InjectedCSS: "body,html{height: auto !important; overflow: auto !important}"}, true},

		// Invalid rules that should be rejected
		// has() is not supported by chromium
		{"example.com##.ad:has(.child)", Rule{}, false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			gotF, gotOk := ParseLine(tt.input)
			if gotOk != tt.wantOk {
				t.Errorf("ParseLine(%q) gotOk = %#v, want %#v", tt.input, gotOk, tt.wantOk)
			} else if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("ParseLine(%q) gotF = %#v, want %#v", tt.input, gotF, tt.wantF)
			}
		})
	}
}
