package filter

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		input  string
		wantF  BasicFilter
		wantOk bool
	}{
		// General rules should have the domain "", as that will be used in the script to inject them in all pages
		{"###cookie_alert", BasicFilter{Domains: []string{""}, CSSSelector: "#cookie_alert"}, true},

		// Normal rules that block certain CSS selectors
		{"example.com##.ad", BasicFilter{Domains: []string{"example.com"}, CSSSelector: ".ad"}, true},
		{"example.com###ad", BasicFilter{Domains: []string{"example.com"}, CSSSelector: "#ad"}, true},
		{"a.com, b.com###ad", BasicFilter{Domains: []string{"a.com", "b.com"}, CSSSelector: "#ad"}, true},

		// Invalid rules that should be rejected
		// has() is not supported by chromium
		{"example.com##.ad:has(.child)", BasicFilter{}, false},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			gotF, gotOk := ParseLine(tt.input)
			if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("ParseLine(%q) gotF = %#v, want %#v", tt.input, gotF, tt.wantF)
			}
			if gotOk != tt.wantOk {
				t.Errorf("ParseLine(%q) gotOk = %#v, want %#v", tt.input, gotOk, tt.wantOk)
			}
		})
	}
}
