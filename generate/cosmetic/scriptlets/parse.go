package scriptlets

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type scriptlet struct {
	Names []string `json:"names"`
	JS    string   `json:"scriptlet"`
}

func Parse(fpath string) (fnLookupTable map[string]string, functionDefinitions string, err error) {
	f, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer f.Close()

	type scriptletFile struct {
		Version    string      `json:"version"`
		Scriptlets []scriptlet `json:"scriptlets"`
	}
	scFile := new(scriptletFile)
	err = json.NewDecoder(f).Decode(scFile)
	if err != nil {
		return
	}

	var defs strings.Builder
	fnLookupTable = make(map[string]string)

	for _, scriptlet := range scFile.Scriptlets {
		fname, ok := parseFunctionName(scriptlet.JS)
		if !ok {
			err = fmt.Errorf("cannot parse function name in scriptlet %q", scriptlet.JS)
			return
		}

		for _, n := range scriptlet.Names {
			fnLookupTable[n] = fname
		}

		defs.WriteString(scriptlet.JS)
		defs.WriteByte('\n')
	}

	return fnLookupTable, defs.String(), nil
}

func toJSObject(x interface{}) string {
	b, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func TableToJS(fnLookupTable map[string]string) string {
	var b strings.Builder

	b.WriteByte('{')

	for alias, fn := range fnLookupTable {
		b.WriteString(toJSObject(alias))
		b.WriteByte(':')
		b.WriteString(fn)
		b.WriteByte(',')
	}
	b.WriteByte('}')

	return b.String()
}

var functionNameRegex = regexp.MustCompile(`function\s+(.*?)\s*\(`)

// Input: e.g. "function abortOnPropertyRead(source,args){" or "function jsonPrune(source,args){"
func parseFunctionName(js string) (fname string, ok bool) {
	results := functionNameRegex.FindAllStringSubmatch(js, 1)
	if len(results) == 0 || len(results[0]) < 2 {
		return
	}

	return results[0][1], true
}
